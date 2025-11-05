package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/domain"
)

type CartRepository interface {
	GetCart(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error)
	SaveCart(ctx context.Context, cart *domain.Cart, ttl time.Duration) error
	DeleteCart(ctx context.Context, cartID uuid.UUID) error
}

type cartRepository struct {
	redis *redis.Client
}

func NewCartRepository(redis *redis.Client) CartRepository {
	return &cartRepository{redis: redis}
}

func (r *cartRepository) GetCart(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error) {
	key := fmt.Sprintf("cart:%s", cartID.String())
	
	data, err := r.redis.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Cart doesn't exist, create new one
		return &domain.Cart{
			ID:        cartID,
			Items:     []domain.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cart from Redis: %w", err)
	}

	var cart domain.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	return &cart, nil
}

func (r *cartRepository) SaveCart(ctx context.Context, cart *domain.Cart, ttl time.Duration) error {
	key := fmt.Sprintf("cart:%s", cart.ID.String())
	
	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	if err := r.redis.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save cart to Redis: %w", err)
	}

	return nil
}

func (r *cartRepository) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	key := fmt.Sprintf("cart:%s", cartID.String())
	
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cart from Redis: %w", err)
	}

	return nil
}
