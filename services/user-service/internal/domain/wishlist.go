package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wishlist represents a user's wishlist item (bookmarked book)
type Wishlist struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_wishlist_user_book,priority:1" json:"user_id"`
	BookID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_wishlist_user_book,priority:2" json:"book_id"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name for GORM
func (Wishlist) TableName() string {
	return "wishlists"
}

// BeforeCreate hook to generate UUID
func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
