package grpc

import (
	"context"

	inventoryv1 "github.com/asliddin03/go-grpc-order-service/inventory-service/gen/inventory/v1"
	"github.com/asliddin03/go-grpc-order-service/inventory-service/internal/service"
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
		return &inventoryv1.GetProductsResponse{
			Products: []*inventoryv1.Product{},
		}, nil
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
