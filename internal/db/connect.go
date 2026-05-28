package db

import (
	"database/sql"
	"fmt"

	"cryptotracker/alert/internal/config"
	"github.com/go-sql-driver/mysql"
)

func Connect(cfg config.DBConfig) (*sql.DB, error) {
	dsn := mysql.Config{
		User:   cfg.User,
		Passwd: cfg.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DBName: cfg.Name,
	}

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir ligação à DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar à DB: %w", err)
	}

	return db, nil
}
