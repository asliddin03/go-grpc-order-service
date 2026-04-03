package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	orderv1 "github.com/asliddin03/go-grpc-order-service/order-service/gen/order/v1"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/config"
	grpcHandler "github.com/asliddin03/go-grpc-order-service/order-service/internal/handler/grpc"
	postgresrepo "github.com/asliddin03/go-grpc-order-service/order-service/internal/repository/postgres"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/service"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	config       *config.Config
	postgres     *pgxpool.Pool
	orderRepo    *postgresrepo.OrderRepository
	orderService *service.OrderService
	orderHandler *grpcHandler.OrderHandler
	grpcServer   *grpc.Server
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
	a.postgres = postgresPool
	defer a.postgres.Close()

	orderRepo := postgresrepo.NewOrderRepository(postgresPool)
	a.orderRepo = orderRepo

	orderService := service.NewOrderService(
		nil,
		nil,
		orderRepo,
	)
	a.orderService = orderService

	orderHandler := grpcHandler.NewOrderHandler(orderService)
	a.orderHandler = orderHandler

	grpcServer := grpc.NewServer()
	a.grpcServer = grpcServer

	orderv1.RegisterOrderServiceServer(grpcServer, orderHandler)

	addr := ":" + cfg.GRPCPort
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s\n", addr)
		errCh <- grpcServer.Serve(listener)
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stopCh)

	select {
	case err := <-errCh:
		return err
	case sig := <-stopCh:
		log.Printf("received signal %s\n", sig)
		log.Println("shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")
		return nil
	}

	log.Println("order-service started")
	log.Println("postgres connected")
	log.Printf("order repository initialized: %T\n", a.orderRepo)
	log.Printf("order service initialized: %T\n", a.orderService)
	return nil
}
