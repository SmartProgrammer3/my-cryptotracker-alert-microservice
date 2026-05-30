package db

import (
	"database/sql"
	"strings"
	"time"

	"cryptotracker/alert/internal/alert"
)

type AlertFilter struct {
    Status *alert.Status
    Symbol *string
}

// InsertAlert insere um novo alerta na DB com estado PENDING.
func InsertAlert(db *sql.DB, symbol string, targetPrice float64, direction alert.Direction) (int, time.Time, error) {
	createdAt := time.Now().UTC()

	result, err := db.Exec(
		"INSERT INTO alerts (symbol, target_price, direction, status, created_at) VALUES (?, ?, ?, ?, ?)",
		symbol, targetPrice, direction, alert.StatusPending, createdAt,
	)
	if err != nil {
		return 0, time.Time{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, time.Time{}, err
	}

	return int(id), createdAt, nil
}

// UpdateAlertStatus atualiza o estado de um alerta. Se foi TRIGGERED ou CANCELED
func UpdateAlertStatus(db *sql.DB, id int, status alert.Status) error {
	if status == alert.StatusTriggered {
		_, err := db.Exec(
			"UPDATE alerts SET status = ?, triggered_at = NOW() WHERE id = ?",
			status, id,
		)
		return err
	}

	_, err := db.Exec(
		"UPDATE alerts SET status = ? WHERE id = ?",
		status, id,
	)
	return err
}

// GetAlerts devolve alertas com suporte a filtros. AlertFilter{} devolve todos os alertas.
func GetAlerts(db *sql.DB, filter AlertFilter) ([]alert.Alert, error) {
	query := "SELECT id, symbol, target_price, direction, status, created_at FROM alerts"
	args := []any{}
	conditions := []string{}

	if filter.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *filter.Status)
	}

	if filter.Symbol != nil {
		conditions = append(conditions, "symbol = ?")
		args = append(args, *filter.Symbol)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []alert.Alert
	for rows.Next() {
		var a alert.Alert
		if err := rows.Scan(&a.ID, &a.Symbol, &a.TargetPrice, &a.Direction, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}