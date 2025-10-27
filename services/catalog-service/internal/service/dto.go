package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
)

// Book DTOs

type CreateBookRequest struct {
	ISBN            string      `json:"isbn" validate:"required,len=13"`
	Title           string      `json:"title" validate:"required,min=1,max=500"`
	Description     string      `json:"description"`
	Price           float64     `json:"price" validate:"required,gt=0"`
	StockQuantity   int         `json:"stock_quantity" validate:"gte=0"`
	PublisherID     *uuid.UUID  `json:"publisher_id"`
	CoverImageURL   string      `json:"cover_image_url"`
	PublicationDate *time.Time  `json:"publication_date"`
	Language        string      `json:"language"`
	PageCount       int         `json:"page_count" validate:"gte=0"`
	AuthorIDs       []uuid.UUID `json:"author_ids"`
	CategoryIDs     []uuid.UUID `json:"category_ids"`
}

type UpdateBookRequest struct {
	Title           *string     `json:"title,omitempty"`
	Description     *string     `json:"description,omitempty"`
	Price           *float64    `json:"price,omitempty"`
	StockQuantity   *int        `json:"stock_quantity,omitempty"`
	PublisherID     *uuid.UUID  `json:"publisher_id,omitempty"`
	CoverImageURL   *string     `json:"cover_image_url,omitempty"`
	PublicationDate *time.Time  `json:"publication_date,omitempty"`
	Language        *string     `json:"language,omitempty"`
	PageCount       *int        `json:"page_count,omitempty"`
	AuthorIDs       []uuid.UUID `json:"author_ids,omitempty"`
	CategoryIDs     []uuid.UUID `json:"category_ids,omitempty"`
}

type BookListResponse struct {
	Books      []domain.Book `json:"books"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// Author DTOs

type CreateAuthorRequest struct {
	Name      string     `json:"name" validate:"required,min=1,max=255"`
	Bio       string     `json:"bio"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Country   string     `json:"country"`
	ImageURL  string     `json:"image_url"`
}

type UpdateAuthorRequest struct {
	Name      *string    `json:"name,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Country   *string    `json:"country,omitempty"`
	ImageURL  *string    `json:"image_url,omitempty"`
}

type AuthorListResponse struct {
	Authors    []domain.Author `json:"authors"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// Category DTOs

type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=100"`
	Slug        string     `json:"slug" validate:"required,min=1,max=100"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        *string    `json:"name,omitempty"`
	Slug        *string    `json:"slug,omitempty"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

// Publisher DTOs

type CreatePublisherRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Country     string `json:"country"`
	Website     string `json:"website"`
	Description string `json:"description"`
}

type UpdatePublisherRequest struct {
	Name        *string `json:"name,omitempty"`
	Country     *string `json:"country,omitempty"`
	Website     *string `json:"website,omitempty"`
	Description *string `json:"description,omitempty"`
}

type PublisherListResponse struct {
	Publishers []domain.Publisher `json:"publishers"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
