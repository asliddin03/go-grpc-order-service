package service

import (
	"context"
	"errors"
	"testing"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder_Success(t *testing.T) {
	t.Parallel()

	authClient := &fakeAuthClient{}

	inventoryClient := &fakeInventoryClient{
		getProductsFunc: func(ctx context.Context, productIDs []int64) (map[int64]InventoryProduct, error) {
			return map[int64]InventoryProduct{
				1: {
					ProductID: 1,
					Price:     1000,
					Available: true,
				},
				2: {
					ProductID: 2,
					Price:     500,
					Available: true,
				},
			}, nil
		},
	}

	repository := &fakeOrderRepository{
		createFunc: func(ctx context.Context, order *domain.Order) error {
			order.ID = 101
			return nil
		},
	}

	orderService := NewOrderService(authClient, inventoryClient, repository)

	items := []domain.OrderItem{
		{
			ProductID: 1,
			Quantity:  2,
		},
		{
			ProductID: 2,
			Quantity:  3,
		},
	}

	order, err := orderService.CreateOrder(context.Background(), 42, items)

	require.NoError(t, err)
	require.NotNil(t, order)

	assert.Equal(t, int64(101), order.ID)
	assert.Equal(t, int64(42), order.UserID)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	assert.Equal(t, int64(3500), order.TotalPrice)

	require.Len(t, order.Items, 2)
	assert.Equal(t, int64(1000), order.Items[0].Price)
	assert.Equal(t, int64(500), order.Items[1].Price)
	assert.False(t, order.CreatedAt.IsZero())
	assert.False(t, order.UpdatedAt.IsZero())
}

func TestOrderService_CreateOrder_Validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		userID  int64
		items   []domain.OrderItem
		wantErr error
	}{
		{
			name:    "invalid user id",
			userID:  0,
			items:   []domain.OrderItem{{ProductID: 1, Quantity: 1}},
			wantErr: domain.ErrInvalidUserID,
		},
		{
			name:    "empty items",
			userID:  42,
			items:   nil,
			wantErr: domain.ErrOrderItemsRequired,
		},
		{
			name:    "invalid product id",
			userID:  42,
			items:   []domain.OrderItem{{ProductID: 0, Quantity: 1}},
			wantErr: domain.ErrInvalidProductID,
		},
		{
			name:    "invalid quantity",
			userID:  42,
			items:   []domain.OrderItem{{ProductID: 1, Quantity: 0}},
			wantErr: domain.ErrInvalidQuantity,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderService := NewOrderService(
				&fakeAuthClient{},
				&fakeInventoryClient{},
				&fakeOrderRepository{},
			)

			order, err := orderService.CreateOrder(context.Background(), tc.userID, tc.items)

			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Nil(t, order)
		})
	}
}

func TestOrderService_CreateOrder_DependencyError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("inventory unavailable")

	inventoryClient := &fakeInventoryClient{
		getProductsFunc: func(ctx context.Context, productIDs []int64) (map[int64]InventoryProduct, error) {
			return nil, expectedErr
		},
	}

	orderService := NewOrderService(
		&fakeAuthClient{},
		inventoryClient,
		&fakeOrderRepository{},
	)

	items := []domain.OrderItem{
		{
			ProductID: 1,
			Quantity:  1,
		},
	}

	order, err := orderService.CreateOrder(context.Background(), 42, items)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, order)
}

func TestOrderService_GetOrder(t *testing.T) {
	t.Parallel()

	expectedOrder := &domain.Order{
		ID:         1,
		UserID:     42,
		Status:     domain.OrderStatusPending,
		TotalPrice: 1500,
		Items: []domain.OrderItem{
			{
				ProductID: 1,
				Quantity:  1,
				Price:     1500,
			},
		},
	}

	repository := &fakeOrderRepository{
		getByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return expectedOrder, nil
		},
	}

	orderService := NewOrderService(
		&fakeAuthClient{},
		&fakeInventoryClient{},
		repository,
	)

	order, err := orderService.GetOrder(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, expectedOrder, order)
}

func TestOrderService_ListOrders(t *testing.T) {
	t.Parallel()

	expectedOrders := []domain.Order{
		{
			ID:         1,
			UserID:     42,
			Status:     domain.OrderStatusPending,
			TotalPrice: 1000,
		},
		{
			ID:         2,
			UserID:     42,
			Status:     domain.OrderStatusConfirmed,
			TotalPrice: 2500,
		},
	}

	repository := &fakeOrderRepository{
		listByUserIDFunc: func(ctx context.Context, userID int64) ([]domain.Order, error) {
			return expectedOrders, nil
		},
	}

	orderService := NewOrderService(
		&fakeAuthClient{},
		&fakeInventoryClient{},
		repository,
	)

	orders, err := orderService.ListOrders(context.Background(), 42)

	require.NoError(t, err)
	require.Len(t, orders, 2)
	assert.Equal(t, expectedOrders, orders)
}

func TestOrderService_GetOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()

	orderService := NewOrderService(
		&fakeAuthClient{},
		&fakeInventoryClient{},
		&fakeOrderRepository{},
	)

	order, err := orderService.GetOrder(context.Background(), 0)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidOrderID)
	assert.Nil(t, order)
}

func TestOrderService_CreateOrder_InvalidProductPrice(t *testing.T) {
	t.Parallel()

	inventoryClient := &fakeInventoryClient{
		getProductsFunc: func(ctx context.Context, productIDs []int64) (map[int64]InventoryProduct, error) {
			return map[int64]InventoryProduct{
				1: {
					ProductID: 1,
					Price:     -100,
					Available: true,
				},
			}, nil
		},
	}

	orderService := NewOrderService(
		&fakeAuthClient{},
		inventoryClient,
		&fakeOrderRepository{},
	)

	items := []domain.OrderItem{
		{
			ProductID: 1,
			Quantity:  2,
		},
	}

	order, err := orderService.CreateOrder(context.Background(), 42, items)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidProductPrice)
	assert.Nil(t, order)
}
