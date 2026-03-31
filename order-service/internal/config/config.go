package config

import (
	"errors"
	"os"
)

type Config struct {
	PostgresDSN string
	GRPCPort    string
}

func Load() (*Config, error) {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		return nil, errors.New("POSTGRES_DSN is required")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	return &Config{
		PostgresDSN: postgresDSN,
		GRPCPort:    grpcPort,
	}, nil
}
