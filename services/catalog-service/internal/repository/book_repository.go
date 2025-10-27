package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	GetByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	List(ctx context.Context, filter BookFilter) ([]domain.Book, int64, error)
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, filter BookFilter) ([]domain.Book, int64, error)
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

type BookFilter struct {
	Page         int
	PageSize     int
	CategoryID   *uuid.UUID
	AuthorID     *uuid.UUID
	PublisherID  *uuid.UUID
	MinPrice     *float64
	MaxPrice     *float64
	InStock      *bool
	SortBy       string // title, price, created_at
	SortOrder    string // asc, desc
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	if err := r.db.WithContext(ctx).Create(book).Error; err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	// Load associations
	return r.loadAssociations(ctx, book)
}

func (r *bookRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	var book domain.Book

	err := r.db.WithContext(ctx).
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		First(&book, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

func (r *bookRepository) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	var book domain.Book

	err := r.db.WithContext(ctx).
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		First(&book, "isbn = ?", isbn).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}

	return &book, nil
}

func (r *bookRepository) List(ctx context.Context, filter BookFilter) ([]domain.Book, int64, error) {
	var books []domain.Book
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Book{})

	// Apply filters
	query = r.applyFilters(query, filter)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	// Apply pagination and sorting
	query = r.applySorting(query, filter)
	query = r.applyPagination(query, filter)

	// Execute query with preloads
	err := query.
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		Find(&books).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list books: %w", err)
	}

	return books, total, nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	// Update book fields
	if err := r.db.WithContext(ctx).Save(book).Error; err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	// Load associations
	return r.loadAssociations(ctx, book)
}

func (r *bookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Book{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete book: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

func (r *bookRepository) Search(ctx context.Context, query string, filter BookFilter) ([]domain.Book, int64, error) {
	var books []domain.Book
	var total int64

	searchQuery := r.db.WithContext(ctx).Model(&domain.Book{})

	// Search in title, description, and ISBN
	if query != "" {
		searchPattern := "%" + query + "%"
		searchQuery = searchQuery.Where(
			"title ILIKE ? OR description ILIKE ? OR isbn ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Apply additional filters
	searchQuery = r.applyFilters(searchQuery, filter)

	// Get total count
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Apply pagination and sorting
	searchQuery = r.applySorting(searchQuery, filter)
	searchQuery = r.applyPagination(searchQuery, filter)

	// Execute query with preloads
	err := searchQuery.
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		Find(&books).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	return books, total, nil
}

func (r *bookRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Book{}).
		Where("id = ?", id).
		Update("stock_quantity", quantity)

	if result.Error != nil {
		return fmt.Errorf("failed to update stock: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Helper methods

func (r *bookRepository) applyFilters(query *gorm.DB, filter BookFilter) *gorm.DB {
	if filter.CategoryID != nil {
		query = query.Joins("JOIN book_categories ON books.id = book_categories.book_id").
			Where("book_categories.category_id = ?", *filter.CategoryID)
	}

	if filter.AuthorID != nil {
		query = query.Joins("JOIN book_authors ON books.id = book_authors.book_id").
			Where("book_authors.author_id = ?", *filter.AuthorID)
	}

	if filter.PublisherID != nil {
		query = query.Where("publisher_id = ?", *filter.PublisherID)
	}

	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	if filter.InStock != nil && *filter.InStock {
		query = query.Where("stock_quantity > 0")
	}

	return query
}

func (r *bookRepository) applySorting(query *gorm.DB, filter BookFilter) *gorm.DB {
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)
	return query.Order(orderClause)
}

func (r *bookRepository) applyPagination(query *gorm.DB, filter BookFilter) *gorm.DB {
	page := filter.Page
	if page < 1 {
		page = 1
	}

	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return query.Offset(offset).Limit(pageSize)
}

func (r *bookRepository) loadAssociations(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).
		Preload("Authors").
		Preload("Categories").
		Preload("Publisher").
		First(book, "id = ?", book.ID).Error
}
