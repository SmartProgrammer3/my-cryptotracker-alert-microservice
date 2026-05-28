package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb        "cryptotracker/alert/api/proto/v1"
	configPkg "cryptotracker/alert/internal/config"
	dbPkg     "cryptotracker/alert/internal/db"
	serverPkg "cryptotracker/alert/internal/server"

	"cryptotracker/alert/internal/motor/sniper"
	"cryptotracker/alert/internal/motor/scout"
)

func main() {
	fmt.Println("Starting CryptoTracker - Alert")

	cfg, err := configPkg.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := dbPkg.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	sniperEngine := sniper.New(database)
	scoutEngine := scout.New(sniperEngine, cfg.BinanceWSURL)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlertServiceServer(grpcServer, serverPkg.NewServer(database, scoutEngine, sniperEngine))
	reflection.Register(grpcServer)

	log.Printf("gRPC server a escutar na porta %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server: %v", err)
	}
}