package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Book represents a book in the catalog
type Book struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ISBN            string         `gorm:"type:varchar(13);uniqueIndex;not null" json:"isbn"`
	Title           string         `gorm:"type:varchar(500);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	Price           float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQuantity   int            `gorm:"type:integer;default:0" json:"stock_quantity"`
	PublisherID     *uuid.UUID     `gorm:"type:uuid" json:"publisher_id,omitempty"`
	Publisher       *Publisher     `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	CoverImageURL   string         `gorm:"type:text" json:"cover_image_url"`
	PublicationDate *time.Time     `gorm:"type:date" json:"publication_date,omitempty"`
	Language        string         `gorm:"type:varchar(50);default:'English'" json:"language"`
	PageCount       int            `gorm:"type:integer" json:"page_count"`
	Authors         []Author       `gorm:"many2many:book_authors;" json:"authors,omitempty"`
	Categories      []Category     `gorm:"many2many:book_categories;" json:"categories,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName overrides the table name for GORM
func (Book) TableName() string {
	return "books"
}
