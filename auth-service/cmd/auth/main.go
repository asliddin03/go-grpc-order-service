package main

import (
	"log"
	"net"

	authv1 "github.com/asliddin03/go-grpc-order-service/auth-service/gen/auth/v1"
	grpcHandler "github.com/asliddin03/go-grpc-order-service/auth-service/internal/handler/grpc"
	"github.com/asliddin03/go-grpc-order-service/auth-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	authService := service.NewAuthService()
	authHandler := grpcHandler.NewAuthHandler(authService)

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("auth-service listening on :50052")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
