package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/events"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/repository"
)

type CartService interface {
	GetCart(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error)
	AddItem(ctx context.Context, cartID uuid.UUID, req AddItemRequest) (*domain.Cart, error)
	UpdateItem(ctx context.Context, cartID uuid.UUID, bookID uuid.UUID, quantity int) (*domain.Cart, error)
	RemoveItem(ctx context.Context, cartID uuid.UUID, bookID uuid.UUID) (*domain.Cart, error)
	ClearCart(ctx context.Context, cartID uuid.UUID) error
}

type cartService struct {
	repo     repository.CartRepository
	cartTTL  time.Duration
	maxItems int
	maxQty   int
	eventPub *events.Publisher
}

func NewCartService(repo repository.CartRepository, cartTTLHours, maxItems, maxQty int, eventPublisher *events.Publisher) CartService {
	return &cartService{
		repo:     repo,
		cartTTL:  time.Duration(cartTTLHours) * time.Hour,
		maxItems: maxItems,
		maxQty:   maxQty,
		eventPub: eventPublisher,
	}
}

func (s *cartService) GetCart(ctx context.Context, cartID uuid.UUID) (*domain.Cart, error) {
	return s.repo.GetCart(ctx, cartID)
}

func (s *cartService) AddItem(ctx context.Context, cartID uuid.UUID, req AddItemRequest) (*domain.Cart, error) {
	// Validate quantity
	if req.Quantity > s.maxQty {
		return nil, fmt.Errorf("quantity exceeds maximum allowed (%d)", s.maxQty)
	}

	// Get current cart
	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists
	found := false
	for i := range cart.Items {
		if cart.Items[i].BookID == req.BookID {
			newQty := cart.Items[i].Quantity + req.Quantity
			if newQty > s.maxQty {
				return nil, fmt.Errorf("total quantity would exceed maximum allowed (%d)", s.maxQty)
			}
			cart.Items[i].Quantity = newQty
			cart.Items[i].Price = req.Price // Update price
			found = true
			break
		}
	}

	if !found {
		// Check max items limit
		if len(cart.Items) >= s.maxItems {
			return nil, fmt.Errorf("cart has reached maximum items limit (%d)", s.maxItems)
		}

		// Add new item
		cart.Items = append(cart.Items, domain.CartItem{
			BookID:   req.BookID,
			Quantity: req.Quantity,
			Price:    req.Price,
			AddedAt:  time.Now(),
		})
	}

	// Recalculate totals
	cart.CalculateTotals()

	// Save cart
	if err := s.repo.SaveCart(ctx, cart, s.cartTTL); err != nil {
		return nil, err
	}

	s.publishEvent(events.NewNotificationEvent("cart.item_added", "cart-service").WithPayload(map[string]interface{}{
		"cart_id":     cartID.String(),
		"book_id":     req.BookID.String(),
		"quantity":    req.Quantity,
		"total_items": len(cart.Items),
		"total_price": cart.Total,
	}))

	return cart, nil
}

func (s *cartService) UpdateItem(ctx context.Context, cartID uuid.UUID, bookID uuid.UUID, quantity int) (*domain.Cart, error) {
	// Validate quantity
	if quantity > s.maxQty {
		return nil, fmt.Errorf("quantity exceeds maximum allowed (%d)", s.maxQty)
	}

	// Get current cart
	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	// Find and update item
	found := false
	for i := range cart.Items {
		if cart.Items[i].BookID == bookID {
			cart.Items[i].Quantity = quantity
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("item not found in cart")
	}

	// Recalculate totals
	cart.CalculateTotals()

	// Save cart
	if err := s.repo.SaveCart(ctx, cart, s.cartTTL); err != nil {
		return nil, err
	}

	s.publishEvent(events.NewNotificationEvent("cart.item_updated", "cart-service").WithPayload(map[string]interface{}{
		"cart_id":  cartID.String(),
		"book_id":  bookID.String(),
		"quantity": quantity,
	}))

	return cart, nil
}

func (s *cartService) RemoveItem(ctx context.Context, cartID uuid.UUID, bookID uuid.UUID) (*domain.Cart, error) {
	// Get current cart
	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	// Find and remove item
	newItems := []domain.CartItem{}
	for _, item := range cart.Items {
		if item.BookID != bookID {
			newItems = append(newItems, item)
		}
	}

	cart.Items = newItems

	// Recalculate totals
	cart.CalculateTotals()

	// Save cart
	if err := s.repo.SaveCart(ctx, cart, s.cartTTL); err != nil {
		return nil, err
	}

	s.publishEvent(events.NewNotificationEvent("cart.item_removed", "cart-service").WithPayload(map[string]interface{}{
		"cart_id": cartID.String(),
		"book_id": bookID.String(),
	}))

	return cart, nil
}

func (s *cartService) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	if err := s.repo.DeleteCart(ctx, cartID); err != nil {
		return err
	}

	s.publishEvent(events.NewNotificationEvent("cart.cleared", "cart-service").WithPayload(map[string]interface{}{
		"cart_id": cartID.String(),
	}))

	return nil
}

func (s *cartService) publishEvent(evt events.NotificationEvent) {
	if s.eventPub == nil {
		return
	}

	if err := s.eventPub.Publish(context.Background(), evt.EventType, evt); err != nil {
		log.Printf("failed to publish cart event %s: %v", evt.EventType, err)
	}
}
