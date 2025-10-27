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

// Category represents a book category/genre
type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	ParentID    *uuid.UUID     `gorm:"type:uuid" json:"parent_id,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Books       []Book         `gorm:"many2many:book_categories;" json:"books,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

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

// BookAuthor represents the many-to-many relationship with ordering
type BookAuthor struct {
	BookID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	AuthorID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"author_id"`
	AuthorOrder int       `gorm:"type:integer;default:0" json:"author_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// BookCategory represents the many-to-many relationship
type BookCategory struct {
	BookID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey" json:"category_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName overrides
func (Book) TableName() string {
	return "books"
}

func (Author) TableName() string {
	return "authors"
}

func (Category) TableName() string {
	return "categories"
}

func (Publisher) TableName() string {
	return "publishers"
}

func (BookAuthor) TableName() string {
	return "book_authors"
}

func (BookCategory) TableName() string {
	return "book_categories"
}
