package domain

import "errors"

var (
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrOrderItemsRequired = errors.New("order items are required")
	ErrInvalidProductID   = errors.New("invalid product id")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrProductUnavailable = errors.New("product unavailable")
	ErrOrderNotFound      = errors.New("order not found")
)
