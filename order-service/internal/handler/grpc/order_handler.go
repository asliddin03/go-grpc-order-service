package grpc

import (
	"context"
	"errors"

	orderv1 "github.com/asliddin03/go-grpc-order-service/order-service/gen/order/v1"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	orderv1.UnimplementedOrderServiceServer
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) GetOrder(ctx context.Context,
	req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	order, err := h.orderService.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, mapError(err)
	}

	return &orderv1.GetOrderResponse{
		Order: toProtoOrder(order),
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context,
	req *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	orders, err := h.orderService.ListOrders(ctx, req.GetUserId())
	if err != nil {
		return nil, mapError(err)
	}

	return &orderv1.ListOrdersResponse{
		Orders: toProtoOrders(orders),
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidUserID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
