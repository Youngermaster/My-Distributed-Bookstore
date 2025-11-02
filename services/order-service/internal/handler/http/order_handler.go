package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/service"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder creates a new order
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req service.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	order, err := h.service.CreateOrder(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to create order",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetOrder gets an order by ID
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := parseUUID(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
	}

	order, err := h.service.GetOrder(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Order not found",
			Message: err.Error(),
		})
	}

	return c.JSON(order)
}

// ListOrders lists all orders with pagination
// GET /api/v1/orders
func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	response, err := h.service.ListOrders(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to list orders",
			Message: err.Error(),
		})
	}

	return c.JSON(response)
}

// GetUserOrders gets orders for a specific user
// GET /api/v1/users/:userId/orders
func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := parseUUID(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	response, err := h.service.GetUserOrders(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to get user orders",
			Message: err.Error(),
		})
	}

	return c.JSON(response)
}

// UpdateOrderStatus updates the status of an order
// PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := parseUUID(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
	}

	var req service.UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	order, err := h.service.UpdateOrderStatus(c.Context(), id, req.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to update order status",
			Message: err.Error(),
		})
	}

	return c.JSON(order)
}

// CancelOrder cancels an order
// POST /api/v1/orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := parseUUID(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid order ID",
			Message: err.Error(),
		})
	}

	if err := h.service.CancelOrder(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to cancel order",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Order cancelled successfully",
	})
}
