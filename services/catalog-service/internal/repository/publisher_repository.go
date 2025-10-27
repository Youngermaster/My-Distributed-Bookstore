package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"gorm.io/gorm"
)

type PublisherRepository interface {
	Create(ctx context.Context, publisher *domain.Publisher) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Publisher, error)
	List(ctx context.Context, page, pageSize int) ([]domain.Publisher, int64, error)
	Update(ctx context.Context, publisher *domain.Publisher) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type publisherRepository struct {
	db *gorm.DB
}

func NewPublisherRepository(db *gorm.DB) PublisherRepository {
	return &publisherRepository{db: db}
}

func (r *publisherRepository) Create(ctx context.Context, publisher *domain.Publisher) error {
	if err := r.db.WithContext(ctx).Create(publisher).Error; err != nil {
		return fmt.Errorf("failed to create publisher: %w", err)
	}
	return nil
}

func (r *publisherRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Publisher, error) {
	var publisher domain.Publisher

	err := r.db.WithContext(ctx).
		Preload("Books", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10) // Limit books to prevent loading too much data
		}).
		First(&publisher, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get publisher: %w", err)
	}

	return &publisher, nil
}

func (r *publisherRepository) List(ctx context.Context, page, pageSize int) ([]domain.Publisher, int64, error) {
	var publishers []domain.Publisher
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
	if err := r.db.WithContext(ctx).Model(&domain.Publisher{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count publishers: %w", err)
	}

	// Get publishers
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&publishers).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list publishers: %w", err)
	}

	return publishers, total, nil
}

func (r *publisherRepository) Update(ctx context.Context, publisher *domain.Publisher) error {
	if err := r.db.WithContext(ctx).Save(publisher).Error; err != nil {
		return fmt.Errorf("failed to update publisher: %w", err)
	}
	return nil
}

func (r *publisherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Publisher{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete publisher: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("publisher not found")
	}

	return nil
}
