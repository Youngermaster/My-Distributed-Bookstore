package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents a customer order
type Order struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Status          OrderStatus    `gorm:"type:varchar(50);not null;index" json:"status"`
	Items           []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
	TotalAmount     float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	ItemCount       int            `gorm:"type:integer;not null" json:"item_count"`
	ShippingAddress string         `gorm:"type:text" json:"shipping_address,omitempty"`
	PaymentMethod   string         `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the table name for GORM
func (Order) TableName() string {
	return "orders"
}

// CalculateTotal calculates the total amount for the order
func (o *Order) CalculateTotal() {
	o.TotalAmount = 0
	o.ItemCount = 0
	for _, item := range o.Items {
		o.TotalAmount += item.Subtotal
		o.ItemCount += item.Quantity
	}
}
