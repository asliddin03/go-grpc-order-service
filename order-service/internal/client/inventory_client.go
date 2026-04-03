package client

import (
	"context"
	"fmt"

	inventoryv1 "github.com/asliddin03/go-grpc-order-service/inventory-service/gen/inventory/v1"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventoryv1.InventoryServiceClient
}

func NewInventoryClient(address string) (*InventoryClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial inventory-service: %w", err)
	}

	client := inventoryv1.NewInventoryServiceClient(conn)

	return &InventoryClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *InventoryClient) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *InventoryClient) GetProducts(ctx context.Context,
	productIDs []int64) (map[int64]service.InventoryProduct, error) {
	resp, err := c.client.GetProducts(ctx, &inventoryv1.GetProductsRequest{
		ProductIds: productIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("inventory-service get products: %w", err)
	}

	result := make(map[int64]service.InventoryProduct, len(resp.GetProducts()))

	for _, product := range resp.GetProducts() {
		result[product.GetProductId()] = service.InventoryProduct{
			ProductID: product.GetProductId(),
			Price:     product.GetPrice(),
			Available: product.GetAvailable(),
		}
	}

	return result, nil
}
