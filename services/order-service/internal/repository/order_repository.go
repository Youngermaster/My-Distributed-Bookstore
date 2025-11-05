package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.Order, int64, error)
	List(ctx context.Context, page, pageSize int) ([]*domain.Order, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.WithContext(ctx).Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get orders with items
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, total, nil
}

func (r *orderRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	offset := (page - 1) * pageSize

	// Count total
	if err := r.db.WithContext(ctx).Model(&domain.Order{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get orders with items
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, total, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	result := r.db.WithContext(ctx).Model(&domain.Order{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Order{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}
