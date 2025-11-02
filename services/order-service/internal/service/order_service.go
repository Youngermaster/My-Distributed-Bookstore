package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/domain"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID, page, pageSize int) (*OrderListResponse, error)
	ListOrders(ctx context.Context, page, pageSize int) (*OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) (*domain.Order, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
}

type orderService struct {
	repo            repository.OrderRepository
	defaultPageSize int
	maxPageSize     int
}

func NewOrderService(repo repository.OrderRepository, defaultPageSize, maxPageSize int) OrderService {
	return &orderService{
		repo:            repo,
		defaultPageSize: defaultPageSize,
		maxPageSize:     maxPageSize,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
	// Validate items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	// Create order items
	items := make([]domain.OrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		subtotal := itemReq.UnitPrice * float64(itemReq.Quantity)
		items[i] = domain.OrderItem{
			BookID:    itemReq.BookID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			Subtotal:  subtotal,
		}
	}

	// Create order
	order := &domain.Order{
		UserID:          req.UserID,
		Status:          domain.OrderStatusPending,
		Items:           items,
		ShippingAddress: req.ShippingAddress,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Calculate total
	order.CalculateTotal()

	// Save order
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID uuid.UUID, page, pageSize int) (*OrderListResponse, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = s.defaultPageSize
	}
	if pageSize > s.maxPageSize {
		pageSize = s.maxPageSize
	}

	orders, total, err := s.repo.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) ListOrders(ctx context.Context, page, pageSize int) (*OrderListResponse, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = s.defaultPageSize
	}
	if pageSize > s.maxPageSize {
		pageSize = s.maxPageSize
	}

	orders, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &OrderListResponse{
		Orders:     orders,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) (*domain.Order, error) {
	// Validate status
	validStatuses := map[domain.OrderStatus]bool{
		domain.OrderStatusPending:    true,
		domain.OrderStatusConfirmed:  true,
		domain.OrderStatusProcessing: true,
		domain.OrderStatusShipped:    true,
		domain.OrderStatusDelivered:  true,
		domain.OrderStatusCancelled:  true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid order status: %s", status)
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *orderService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	// Check if order exists and can be cancelled
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	// Can only cancel pending or confirmed orders
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusConfirmed {
		return fmt.Errorf("cannot cancel order with status: %s", order.Status)
	}

	return s.repo.UpdateStatus(ctx, id, domain.OrderStatusCancelled)
}
