package db

import (
	"database/sql"
	"time"

	"cryptotracker/alert/internal/alert"
)

// InsertAlert insere um novo alerta na DB com estado PENDING.
func InsertAlert(db *sql.DB, symbol string, targetPrice float64) (int, time.Time, error) {
	createdAt := time.Now().UTC()

	result, err := db.Exec(
		"INSERT INTO alerts (symbol, target_price, status, created_at) VALUES (?, ?, ?, ?)",
		symbol, targetPrice, alert.StatusPending, createdAt,
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

// GetPendingAlerts devolve todos os alertas com estado PENDING.
func GetPendingAlerts(db *sql.DB) ([]alert.Alert, error) {
	rows, err := db.Query(
		"SELECT id, symbol, target_price, created_at FROM alerts WHERE status = ?",
		alert.StatusPending,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []alert.Alert
	for rows.Next() {
		var a alert.Alert
		if err := rows.Scan(&a.ID, &a.Symbol, &a.TargetPrice, &a.CreatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
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