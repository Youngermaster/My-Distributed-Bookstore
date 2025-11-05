package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Address represents a shipping or billing address
type Address struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	AddressLine1 string         `gorm:"size:255;not null" json:"address_line1"`
	AddressLine2 string         `gorm:"size:255" json:"address_line2,omitempty"`
	City         string         `gorm:"size:100;not null" json:"city"`
	State        string         `gorm:"size:100" json:"state,omitempty"`
	PostalCode   string         `gorm:"size:20;not null" json:"postal_code"`
	Country      string         `gorm:"size:100;not null" json:"country"`
	IsDefault    bool           `gorm:"default:false" json:"is_default"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name for GORM
func (Address) TableName() string {
	return "addresses"
}

// BeforeCreate hook to generate UUID
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
