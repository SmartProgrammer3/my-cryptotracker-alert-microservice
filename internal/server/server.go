package server

import (
	"context"
	"fmt"
	"log"

	pb "cryptotracker/alert/api/proto/v1"
)

// AlertaServer vai implementar a interface gerada pelo gRPC
type AlertaServer struct {
	pb.UnimplementedAlertServiceServer 
}

// CreateAlert trata do pedido de criação de um novo alerta
func (s *AlertaServer) CreateAlert(ctx context.Context, req *pb.CreateAlertRequest) (*pb.CreateAlertResponse, error) {
	log.Printf("[gRPC] Pedido para criar alerta para o símbolo %s a %.2f\n", req.Symbol, req.TargetPrice)

	// Simulação: Fingimos que guardámos na Base de Dados e gerámos um ID
	idFalso := "alert_uuid_12345"

	return &pb.CreateAlertResponse{
		AlertId: idFalso,
		Success: true,
		Message: fmt.Sprintf("Alerta ativo para %s no preço alvo de %.2f!", req.Symbol, req.TargetPrice),
	}, nil
}