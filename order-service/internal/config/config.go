package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	PostgresDSN             string
	GRPCPort                string
	AuthServiceAddress      string
	InventoryServiceAddress string
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

	if _, err := strconv.Atoi(grpcPort); err != nil {
		return nil, errors.New("GRPC_PORT must be a number")
	}

	authServiceAddress := os.Getenv("AUTH_SERVICE_ADDRESS")
	if authServiceAddress == "" {
		authServiceAddress = "localhost:50052"
	}

	inventoryServiceAddress := os.Getenv("INVENTORY_SERVICE_ADDRESS")
	if inventoryServiceAddress == "" {
		inventoryServiceAddress = "localhost:50053"
	}

	return &Config{
		PostgresDSN:             postgresDSN,
		GRPCPort:                grpcPort,
		AuthServiceAddress:      authServiceAddress,
		InventoryServiceAddress: inventoryServiceAddress,
	}, nil
}
