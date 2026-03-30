package config

import (
	"errors"
	"os"
)

type Config struct {
	PostgresDSN string
}

func Load() (*Config, error) {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required")
	}

	return &Config{
		PostgresDSN: postgresDSN,
	}, nil
}
