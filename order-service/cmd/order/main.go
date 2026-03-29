package main

import (
	"log"

	"github.com/asliddin03/go-grpc-order-service/tree/main/order-service/internal/app"
)

func main() {
	application := app.New()

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
