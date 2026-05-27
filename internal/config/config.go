package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port string
}

// LoadConfig tenta carregar as configurações e retorna um erro se faltarem variáveis obrigatórias
func LoadConfig() (*Config, error) {
	port := os.Getenv("ALERT_MICROSERVICE_PORT")
	
	// Se a variável estiver vazia, lançamos um erro claro
	if port == "" {
		return nil, fmt.Errorf("a variável de ambiente obrigatória 'ALERT_MICROSERVICE_PORT' não está definida")
	}

	return &Config{
		Port: port,
	}, nil
}