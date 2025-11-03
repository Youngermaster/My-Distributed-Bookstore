package domain

import (
	"time"

	"github.com/google/uuid"
)

// CartItem represents an item in the shopping cart
type CartItem struct {
	BookID    uuid.UUID `json:"book_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Subtotal  float64   `json:"subtotal"`
	AddedAt   time.Time `json:"added_at"`
}
