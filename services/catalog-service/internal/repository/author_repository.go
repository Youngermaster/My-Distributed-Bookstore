package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	Create(ctx context.Context, author *domain.Author) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Author, error)
	List(ctx context.Context, page, pageSize int) ([]domain.Author, int64, error)
	Update(ctx context.Context, author *domain.Author) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string) ([]domain.Author, error)
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) Create(ctx context.Context, author *domain.Author) error {
	if err := r.db.WithContext(ctx).Create(author).Error; err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}
	return nil
}

func (r *authorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Author, error) {
	var author domain.Author

	err := r.db.WithContext(ctx).
		Preload("Books", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10) // Limit books to prevent loading too much data
		}).
		First(&author, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	return &author, nil
}

func (r *authorRepository) List(ctx context.Context, page, pageSize int) ([]domain.Author, int64, error) {
	var authors []domain.Author
	var total int64

	// Default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.WithContext(ctx).Model(&domain.Author{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count authors: %w", err)
	}

	// Get authors
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&authors).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list authors: %w", err)
	}

	return authors, total, nil
}

func (r *authorRepository) Update(ctx context.Context, author *domain.Author) error {
	if err := r.db.WithContext(ctx).Save(author).Error; err != nil {
		return fmt.Errorf("failed to update author: %w", err)
	}
	return nil
}

func (r *authorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Author{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete author: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("author not found")
	}

	return nil
}

func (r *authorRepository) Search(ctx context.Context, query string) ([]domain.Author, error) {
	var authors []domain.Author

	searchPattern := "%" + query + "%"
	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR bio ILIKE ?", searchPattern, searchPattern).
		Order("name ASC").
		Limit(50).
		Find(&authors).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search authors: %w", err)
	}

	return authors, nil
}
