package service

import (
	"context"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
)

type fakeAuthClient struct {
	validateUserFunc func(ctx context.Context, userID int64) error
}

func (f *fakeAuthClient) ValidateUser(ctx context.Context, userID int64) error {
	if f.validateUserFunc != nil {
		return f.validateUserFunc(ctx, userID)
	}

	return nil
}

type fakeInventoryClient struct {
	getProductsFunc func(ctx context.Context, productIDs []int64) (map[int64]InventoryProduct, error)
}

func (f *fakeInventoryClient) GetProducts(ctx context.Context, productIDs []int64) (map[int64]InventoryProduct, error) {
	if f.getProductsFunc != nil {
		return f.getProductsFunc(ctx, productIDs)
	}

	return map[int64]InventoryProduct{}, nil
}

type fakeOrderRepository struct {
	createFunc       func(ctx context.Context, order *domain.Order) error
	getByIDFunc      func(ctx context.Context, orderID int64) (*domain.Order, error)
	listByUserIDFunc func(ctx context.Context, userID int64) ([]domain.Order, error)
}

func (f *fakeOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if f.createFunc != nil {
		return f.createFunc(ctx, order)
	}

	return nil
}

func (f *fakeOrderRepository) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	if f.getByIDFunc != nil {
		return f.getByIDFunc(ctx, orderID)
	}

	return nil, nil
}

func (f *fakeOrderRepository) ListByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	if f.listByUserIDFunc != nil {
		return f.listByUserIDFunc(ctx, userID)
	}

	return nil, nil
}
