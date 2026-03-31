package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/asliddin03/go-grpc-order-service/auth-service/gen/auth/v1"
	"github.com/asliddin03/go-grpc-order-service/auth-service/internal/service"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) ValidateUser(ctx context.Context, req *authv1.ValidateUserRequest) (*authv1.ValidateUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	user, err := h.authService.ValidateUser(req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.ValidateUserResponse{
		Valid:  true,
		UserId: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}
