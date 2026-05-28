package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb      "cryptotracker/alert/api/proto/v1"

	"cryptotracker/alert/internal/alert"
	"cryptotracker/alert/internal/db"
	"cryptotracker/alert/internal/motor/sniper"
	"cryptotracker/alert/internal/motor/scout"
)

type Server struct {
	pb.UnimplementedAlertServiceServer
	db     *sql.DB
	scout  *scout.Scout
	sniper *sniper.Sniper
}

func NewServer(db *sql.DB, scout *scout.Scout, sniper *sniper.Sniper) *Server {
	return &Server{db: db, scout: scout, sniper: sniper}
}

// CreateAlert trata do pedido de criação de um novo alerta
func (s *Server) CreateAlert(ctx context.Context, req *pb.CreateAlertRequest) (*pb.CreateAlertResponse, error) {
	log.Printf("[gRPC] Pedido para criar alerta — símbolo: %s, preço-alvo: %.2f", req.Symbol, req.TargetPrice)

	alertId, createdAt, err := db.InsertAlert(s.db, req.Symbol, req.TargetPrice)
	if err != nil {
		return nil, fmt.Errorf("erro ao guardar alerta: %w", err)
	}

	alert := alert.NewAlert(alertId, req.Symbol, req.TargetPrice, createdAt)

	s.sniper.Add(alert)
	s.scout.Subscribe(req.Symbol)

	return &pb.CreateAlertResponse{
		AlertId: fmt.Sprintf("%d", alert.ID),
		Success: true,
		Message: fmt.Sprintf("Alerta activo para %s no preço-alvo de %.2f!", req.Symbol, req.TargetPrice),
	}, nil
}