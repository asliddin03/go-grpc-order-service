package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	inventoryv1 "github.com/asliddin03/go-grpc-order-service/inventory-service/gen/inventory/v1"
	grpcHandler "github.com/asliddin03/go-grpc-order-service/inventory-service/internal/handler/grpc"
	"github.com/asliddin03/go-grpc-order-service/inventory-service/internal/service"
)

func main() {
	inventoryService := service.NewInventoryService()
	inventoryHandler := grpcHandler.NewInventoryHandler(inventoryService)

	grpcServer := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryHandler)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("inventory-service listening on :50053")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
