package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "cryptotracker/alert/api/proto/v1"

	"cryptotracker/alert/internal/alert"
	"cryptotracker/alert/internal/db"
	"cryptotracker/alert/internal/motor/scout"
	"cryptotracker/alert/internal/motor/sniper"
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

// CreateAlert trata do pedido de criação de um novo alerta.
func (s *Server) CreateAlert(ctx context.Context, req *pb.CreateAlertRequest) (*pb.CreateAlertResponse, error) {
	log.Printf("[gRPC] Pedido para criar alerta — símbolo: %s, preço-alvo: %.2f, direção: %s", req.Symbol, req.TargetPrice, req.Direction)

	direction := alert.Direction(req.Direction)
	alertId, createdAt, err := db.InsertAlert(s.db, req.Symbol, req.TargetPrice, direction)
	if err != nil {
		return nil, fmt.Errorf("erro ao guardar alerta: %w", err)
	}

	alert := alert.NewAlert(alertId, req.Symbol, req.TargetPrice, direction, createdAt)

	s.sniper.Add(alert)
	s.scout.Subscribe(req.Symbol)

	return &pb.CreateAlertResponse{
		AlertId: fmt.Sprintf("%d", alert.ID),
		Success: true,
	}, nil
}

// CancelAlert trata do pedido de cancelamento de um novo alerta.
func (s *Server) CancelAlert(ctx context.Context, req *pb.CancelAlertRequest) (*pb.CancelAlertResponse, error) {
	log.Printf("[gRPC] Pedido para cancelar alerta %s", req.AlertId)

	id, err := strconv.Atoi(req.AlertId)
	if err != nil {
		return nil, fmt.Errorf("alert_id inválido: %w", err)
	}

	if err := db.UpdateAlertStatus(s.db, id, alert.StatusCancelled); err != nil {
		return nil, fmt.Errorf("erro ao cancelar alerta: %w", err)
	}

	if symbol, empty := s.sniper.Remove(id); empty {
		s.scout.Unsubscribe(symbol)
	}

	return &pb.CancelAlertResponse{
		Success: true,
	}, nil
}

// ListAlerts trata do pedido de listagem de alertas.
func (s *Server) ListAlerts(ctx context.Context, req *pb.ListAlertsRequest) (*pb.ListAlertsResponse, error) {
	log.Printf("[gRPC] Pedido para listar alertas")

	filter := db.AlertFilter{}

	if req.Status != nil {
		status := alert.Status(req.GetStatus())
		filter.Status = &status
	}

	if req.Symbol != nil {
		symbol := req.GetSymbol()
		filter.Symbol = &symbol
	}

	alerts, err := db.GetAlerts(s.db, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar alertas: %w", err)
	}

	var messages []*pb.AlertMessage
	for _, alert := range alerts {
		messages = append(messages, &pb.AlertMessage{
			AlertId:     fmt.Sprintf("%d", alert.ID),
			Symbol:      alert.Symbol,
			TargetPrice: alert.TargetPrice,
			Direction:   string(alert.Direction),
			Status:      string(alert.Status),
			CreatedAt:   alert.CreatedAt.Format(time.DateTime),
		})
	}

	return &pb.ListAlertsResponse{
		Alerts: messages,
	}, nil
}