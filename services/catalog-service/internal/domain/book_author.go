package domain

import (
	"time"

	"github.com/google/uuid"
)

// BookAuthor represents the many-to-many relationship between books and authors with ordering
type BookAuthor struct {
	BookID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	AuthorID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"author_id"`
	AuthorOrder int       `gorm:"type:integer;default:0" json:"author_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName overrides the table name for GORM
func (BookAuthor) TableName() string {
	return "book_authors"
}
