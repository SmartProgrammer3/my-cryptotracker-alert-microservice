package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "cryptotracker/alert/api/proto/v1"

	"cryptotracker/alert/internal/alert"
	"cryptotracker/alert/internal/config"
	"cryptotracker/alert/internal/db"
	"cryptotracker/alert/internal/server"
	"cryptotracker/alert/internal/motor/scout"
	"cryptotracker/alert/internal/motor/sniper"
)

func main() {
	fmt.Println("Starting CryptoTracker - Alert")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config failed to load: %v", err)
	}

	database, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("db failed to connect: %v", err)
	}
	defer database.Close()

	sniperEngine := sniper.New(database)
	scoutEngine := scout.New(sniperEngine, cfg.BinanceWSURL)

	// Arranque, carregar alertas com status PENDING da DB
	status := alert.StatusPending
	pendingAlerts, err := db.GetAlerts(database, db.AlertFilter{Status: &status})
	if err != nil {
		log.Fatalf("load: %v", err)
	}

	sniperEngine.Load(pendingAlerts)
	for _, a := range pendingAlerts {
		scoutEngine.Subscribe(a.Symbol)
	}

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlertServiceServer(grpcServer, server.NewServer(database, scoutEngine, sniperEngine))
	reflection.Register(grpcServer)

	log.Printf("CryptoTracker - Alert listening on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server: %v", err)
	}
}