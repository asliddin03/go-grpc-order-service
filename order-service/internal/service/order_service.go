package service

import (
	"context"
	"time"

	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
)

type OrderService struct {
	authClient      AuthClient
	inventoryClient InventoryClient
	orderRepository OrderRepository
}

func NewOrderService(
	authClient AuthClient,
	inventoryClient InventoryClient,
	orderRepository OrderRepository,
) *OrderService {
	return &OrderService{
		authClient:      authClient,
		inventoryClient: inventoryClient,
		orderRepository: orderRepository,
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	userID int64,
	items []domain.OrderItem,
) (*domain.Order, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	if len(items) == 0 {
		return nil, domain.ErrOrderItemsRequired
	}

	productIDs := make([]int64, 0, len(items))

	for _, item := range items {
		if item.ProductID <= 0 {
			return nil, domain.ErrInvalidProductID
		}

		if item.Quantity <= 0 {
			return nil, domain.ErrInvalidQuantity
		}

		productIDs = append(productIDs, item.ProductID)
	}

	if err := s.authClient.ValidateUser(ctx, userID); err != nil {
		return nil, err
	}

	products, err := s.inventoryClient.GetProducts(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	orderItems := make([]domain.OrderItem, 0, len(items))
	var totalPrice int64

	for _, item := range items {
		product, ok := products[item.ProductID]
		if !ok || !product.Available {
			return nil, domain.ErrProductUnavailable
		}

		orderItem := domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		orderItems = append(orderItems, orderItem)
		totalPrice += product.Price * int64(item.Quantity)
	}

	now := time.Now()

	order := &domain.Order{
		UserID:     userID,
		Status:     domain.OrderStatusPending,
		TotalPrice: totalPrice,
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      orderItems,
	}

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	if orderID <= 0 {
		return nil, domain.ErrOrderNotFound
	}

	return s.orderRepository.GetByID(ctx, orderID)
}

func (s *OrderService) ListOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	return s.orderRepository.ListByUserID(ctx, userID)
}
