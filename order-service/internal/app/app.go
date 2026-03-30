package app

import (
	"context"
	"log"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/config"
	postgresrepo "github.com/asliddin03/go-grpc-order-service/order-service/internal/repository/postgres"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/service"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	config       *config.Config
	postgres     *pgxpool.Pool
	orderRepo    *postgresrepo.OrderRepository
	orderService *service.OrderService
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	a.config = cfg

	postgresPool, err := storage.NewPostgresPool(ctx, cfg.PostgresDSN)
	if err != nil {
		return err
	}
	defer postgresPool.Close()

	a.postgres = postgresPool

	orderRepo := postgresrepo.NewOrderRepository(postgresPool)
	a.orderRepo = orderRepo

	orderService := service.NewOrderService(
		nil,
		nil,
		orderRepo,
	)
	a.orderService = orderService

	log.Println("order-service started")
	log.Println("postgres connected")
	log.Printf("order repository initialized: %T\n", a.orderRepo)
	log.Printf("order service initialized: %T\n", a.orderService)
	return nil
}
