package grpc

import (
	"time"

	orderv1 "github.com/asliddin03/go-grpc-order-service/order-service/gen/order/v1"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
)

func toProtoOrder(order *domain.Order) *orderv1.Order {
	if order == nil {
		return nil
	}

	items := make([]*orderv1.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &orderv1.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &orderv1.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		CreatedAt:  formatTime(order.CreatedAt),
		UpdatedAt:  formatTime(order.UpdatedAt),
		Items:      items,
	}
}

func toProtoOrders(orders []domain.Order) []*orderv1.Order {
	result := make([]*orderv1.Order, 0, len(orders))

	for i := range orders {
		result = append(result, toProtoOrder(&orders[i]))
	}

	return result
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.RFC3339)
}
