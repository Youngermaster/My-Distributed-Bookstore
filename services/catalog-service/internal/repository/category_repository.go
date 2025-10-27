package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	ListHierarchical(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var category domain.Category

	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var category domain.Category

	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&category, "slug = ?", slug).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category

	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) ListHierarchical(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category

	// Get only root categories (no parent)
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Preload("Children").
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list hierarchical categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Category{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
