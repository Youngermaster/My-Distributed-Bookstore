package service

import (
	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/domain"
)

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	UserID          uuid.UUID       `json:"user_id" validate:"required"`
	Items           []OrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingAddress string          `json:"shipping_address,omitempty"`
	PaymentMethod   string          `json:"payment_method,omitempty"`
	Notes           string          `json:"notes,omitempty"`
}

// OrderItemRequest represents an item in the order request
type OrderItemRequest struct {
	BookID    uuid.UUID `json:"book_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	UnitPrice float64   `json:"unit_price" validate:"required,min=0"`
}

// UpdateStatusRequest represents the request to update order status
type UpdateStatusRequest struct {
	Status domain.OrderStatus `json:"status" validate:"required"`
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	Orders     []*domain.Order `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
