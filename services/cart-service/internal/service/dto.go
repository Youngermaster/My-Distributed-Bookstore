package service

import (
	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/domain"
)

// AddItemRequest represents the request to add an item to cart
type AddItemRequest struct {
	BookID   uuid.UUID `json:"book_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,min=1"`
	Price    float64   `json:"price" validate:"required,min=0"`
}

// UpdateItemRequest represents the request to update item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// CartResponse represents the cart response
type CartResponse struct {
	*domain.Cart
}
