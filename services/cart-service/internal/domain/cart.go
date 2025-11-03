package domain

import (
	"time"

	"github.com/google/uuid"
)

// Cart represents a shopping cart (Redis-only storage)
type Cart struct {
	ID        uuid.UUID   `json:"id"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"` // nil for anonymous users
	Items     []CartItem  `json:"items"`
	Total     float64     `json:"total"`
	ItemCount int         `json:"item_count"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// CalculateTotals calculates the total and item count for the cart
func (c *Cart) CalculateTotals() {
	c.Total = 0
	c.ItemCount = 0
	for i := range c.Items {
		c.Items[i].Subtotal = c.Items[i].Price * float64(c.Items[i].Quantity)
		c.Total += c.Items[i].Subtotal
		c.ItemCount += c.Items[i].Quantity
	}
	c.UpdatedAt = time.Now()
}
