package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
)

func TestOrderRepository_CreateAndGetByID(t *testing.T) {
	ctx, _, repository := newTestOrderRepository(t)

	order := newTestOrder(
		42,
		domain.OrderStatusPending,
		3500,
		[]domain.OrderItem{
			{ProductID: 1, Quantity: 2, Price: 1000},
			{ProductID: 2, Quantity: 3, Price: 500},
		},
	)

	err := repository.Create(ctx, order)
	require.NoError(t, err)
	require.NotZero(t, order.ID)

	savedOrder, err := repository.GetByID(ctx, order.ID)
	require.NoError(t, err)
	require.NotNil(t, savedOrder)

	assert.Equal(t, order.ID, savedOrder.ID)
	assert.Equal(t, order.UserID, savedOrder.UserID)
	assert.Equal(t, order.Status, savedOrder.Status)
	assert.Equal(t, order.TotalPrice, savedOrder.TotalPrice)

	require.Len(t, savedOrder.Items, 2)
	assert.Equal(t, int64(1), savedOrder.Items[0].ProductID)
	assert.Equal(t, int32(2), savedOrder.Items[0].Quantity)
	assert.Equal(t, int64(1000), savedOrder.Items[0].Price)

	assert.Equal(t, int64(2), savedOrder.Items[1].ProductID)
	assert.Equal(t, int32(3), savedOrder.Items[1].Quantity)
	assert.Equal(t, int64(500), savedOrder.Items[1].Price)
}

func TestOrderRepository_GetByID_NotFound(t *testing.T) {
	ctx, _, repository := newTestOrderRepository(t)

	order, err := repository.GetByID(ctx, 999999)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
	assert.Nil(t, order)
}

func TestOrderRepository_ListByUserID(t *testing.T) {
	ctx, _, repository := newTestOrderRepository(t)

	firstOrder := newTestOrder(
		42,
		domain.OrderStatusPending,
		1000,
		[]domain.OrderItem{
			{ProductID: 10, Quantity: 1, Price: 1000},
		},
	)

	secondOrder := newTestOrder(
		42,
		domain.OrderStatusConfirmed,
		2500,
		[]domain.OrderItem{
			{ProductID: 20, Quantity: 5, Price: 500},
		},
	)

	otherUserOrder := newTestOrder(
		99,
		domain.OrderStatusPending,
		700,
		[]domain.OrderItem{
			{ProductID: 30, Quantity: 1, Price: 700},
		},
	)

	require.NoError(t, repository.Create(ctx, firstOrder))
	require.NoError(t, repository.Create(ctx, secondOrder))
	require.NoError(t, repository.Create(ctx, otherUserOrder))

	orders, err := repository.ListByUserID(ctx, 42)

	require.NoError(t, err)
	require.Len(t, orders, 2)

	assert.Equal(t, int64(42), orders[0].UserID)
	assert.Equal(t, int64(42), orders[1].UserID)

	assert.Equal(t, firstOrder.ID, orders[0].ID)
	assert.Equal(t, secondOrder.ID, orders[1].ID)

	require.Len(t, orders[0].Items, 1)
	require.Len(t, orders[1].Items, 1)

	assert.Equal(t, int64(10), orders[0].Items[0].ProductID)
	assert.Equal(t, int64(20), orders[1].Items[0].ProductID)
}
