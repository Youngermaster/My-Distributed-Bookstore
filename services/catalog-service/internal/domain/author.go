package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Author represents a book author
type Author struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Bio       string         `gorm:"type:text" json:"bio"`
	BirthDate *time.Time     `gorm:"type:date" json:"birth_date,omitempty"`
	Country   string         `gorm:"type:varchar(100)" json:"country"`
	ImageURL  string         `gorm:"type:text" json:"image_url"`
	Books     []Book         `gorm:"many2many:book_authors;" json:"books,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the table name for GORM
func (Author) TableName() string {
	return "authors"
}
