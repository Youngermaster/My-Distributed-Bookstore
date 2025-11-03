package domain

import (
	"time"

	"github.com/google/uuid"
)

// BookCategory represents the many-to-many relationship between books and categories
type BookCategory struct {
	BookID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey" json:"category_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName overrides the table name for GORM
func (BookCategory) TableName() string {
	return "book_categories"
}
