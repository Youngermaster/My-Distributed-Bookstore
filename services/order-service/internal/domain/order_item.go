package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	BookID    uuid.UUID      `gorm:"type:uuid;not null" json:"book_id"`
	Quantity  int            `gorm:"type:integer;not null" json:"quantity"`
	UnitPrice float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal  float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the table name for GORM
func (OrderItem) TableName() string {
	return "order_items"
}
