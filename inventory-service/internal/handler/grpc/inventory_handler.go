package grpc

import (
	"context"

	inventoryv1 "github.com/asliddin03/go-grpc-order-service/inventory-service/gen/inventory/v1"
	"github.com/asliddin03/go-grpc-order-service/inventory-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	inventoryService *service.InventoryService
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) GetProducts(ctx context.Context, req *inventoryv1.GetProductsRequest) (*inventoryv1.GetProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	products := h.inventoryService.GetProducts(req.GetProductIds())

	result := make([]*inventoryv1.Product, 0, len(products))
	for _, product := range products {
		result = append(result, &inventoryv1.Product{
			ProductId: product.ProductID,
			Price:     product.Price,
			Available: product.Available,
		})
	}

	return &inventoryv1.GetProductsResponse{
		Products: result,
	}, nil
}
