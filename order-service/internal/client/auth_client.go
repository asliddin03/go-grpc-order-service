package client

import (
	"context"
	"fmt"

	authv1 "github.com/asliddin03/go-grpc-order-service/auth-service/gen/auth/v1"
	"github.com/asliddin03/go-grpc-order-service/order-service/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial auth-service: %w", err)
	}

	client := authv1.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *AuthClient) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *AuthClient) ValidateUser(ctx context.Context, userID int64) error {
	resp, err := c.client.ValidateUser(ctx, &authv1.ValidateUserRequest{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	if !resp.Valid {
		return domain.ErrInvalidUserID
	}

	return nil
}
