package config

import (
	"github.com/caarlos0/env/v10"
)

// Config represents config's variables
type Config struct {
	ServerPort  string `env:"SERVER_PORT" envDefault:"0.0.0.0:10011"`
	PostgresURL string `env:"POSTGRES_URL" envDefault:"postgresql://postgres:testpassword@localhost:6432/filmlibrary?sslmode=disable""`
}

// NewConfig loads and parses config file from given paths
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
