package sniper

import (
	"cryptotracker/alert/internal/alert"
	"database/sql"
)

func New(db *sql.DB) *Sniper {
	return &Sniper{
		alerts: make(map[string][]alert.Alert),
		db:     db,
	}
}