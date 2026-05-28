package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server 		 ServerConfig
	DB     		 DBConfig
	BinanceWSURL string
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		BinanceWSURL: os.Getenv("BINANCE_WS_URL"),
	}

	required := map[string]string{
		"SERVER_PORT": 	  cfg.Server.Port,
		"DB_HOST":     	  cfg.DB.Host,
		"DB_PORT":     	  cfg.DB.Port,
		"DB_USER":     	  cfg.DB.User,
		"DB_PASSWORD": 	  cfg.DB.Password,
		"DB_NAME":     	  cfg.DB.Name,
		"BINANCE_WS_URL": cfg.BinanceWSURL,
	}

	for key, val := range required {
		if val == "" {
			return nil, fmt.Errorf("a variável '%s' é obrigatória, mas não está definida", key)
		}
	}

	return cfg, nil
}