package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Publisher represents a book publisher
type Publisher struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Country     string         `gorm:"type:varchar(100)" json:"country"`
	Website     string         `gorm:"type:varchar(500)" json:"website"`
	Description string         `gorm:"type:text" json:"description"`
	Books       []Book         `gorm:"foreignKey:PublisherID" json:"books,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the table name for GORM
func (Publisher) TableName() string {
	return "publishers"
}
