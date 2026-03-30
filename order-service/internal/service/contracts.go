package service

import (
	"context"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
)

type AuthClient interface {
	ValidateUser(ctx context.Context, userID int64) error
}

type InventoryProduct struct {
	ProductID int64
	Price     int64
	Available bool
}

type InventoryClient interface {
	GetProducts(ctx context.Context, productIDs []int64) (map[int64]InventoryProduct, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, orderID int64) (*domain.Order, error)
	ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error)
}
