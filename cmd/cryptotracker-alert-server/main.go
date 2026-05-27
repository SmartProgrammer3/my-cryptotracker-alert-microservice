package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	configPckg "cryptotracker/alert/internal/config"
	pb "cryptotracker/alert/api/proto/v1"
	serverPckg "cryptotracker/alert/internal/server"
)

func main() {
	fmt.Println("A iniciar o meu Microsserviço CryptoTracker - Alert!!")

	// 1. Carregar as configurações e verificar se deu erro.
	cfg, err := configPckg.LoadConfig()
	if err != nil {
		log.Fatalf("Erro crítico nas configurações: %v", err)
	}

	// 2. Criar um Listener TCP para abrir a porta de rede.
	endereco := fmt.Sprintf(":%s", cfg.Port)
	listener, err := net.Listen("tcp", endereco)
	if err != nil {
		log.Fatalf("Falha ao abrir a porta %s: %v", cfg.Port, err)
	}

	// 3. Criar uma nova instância do servidor gRPC da Google
	grpcServer := grpc.NewServer()

	// 4. Instanciar o servidor de alertas
	alertaService := &serverPckg.AlertaServer{}

	// 5. Registar o teu serviço no servidor gRPC da Google
	pb.RegisterAlertServiceServer(grpcServer, alertaService)

	reflection.Register(grpcServer)

    log.Printf("Servidor gRPC a correr na porta %s...\n", cfg.Port)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Falha ao arrancar o servidor gRPC: %v", err)
    }
}