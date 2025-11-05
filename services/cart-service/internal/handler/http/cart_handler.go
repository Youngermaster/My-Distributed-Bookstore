package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/service"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// GetCart returns the current cart
// GET /api/v1/cart/:cartId
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	cartIDStr := c.Params("cartId")
	cartID, err := parseUUID(cartIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid cart ID",
			Message: err.Error(),
		})
	}

	cart, err := h.service.GetCart(c.Context(), cartID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to get cart",
			Message: err.Error(),
		})
	}

	return c.JSON(cart)
}

// AddItem adds an item to the cart
// POST /api/v1/cart/:cartId/items
func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	cartIDStr := c.Params("cartId")
	cartID, err := parseUUID(cartIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid cart ID",
			Message: err.Error(),
		})
	}

	var req service.AddItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	cart, err := h.service.AddItem(c.Context(), cartID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to add item",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(cart)
}

// UpdateItem updates the quantity of an item
// PUT /api/v1/cart/:cartId/items/:bookId
func (h *CartHandler) UpdateItem(c *fiber.Ctx) error {
	cartIDStr := c.Params("cartId")
	cartID, err := parseUUID(cartIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid cart ID",
			Message: err.Error(),
		})
	}

	bookIDStr := c.Params("bookId")
	bookID, err := parseUUID(bookIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid book ID",
			Message: err.Error(),
		})
	}

	var req service.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	cart, err := h.service.UpdateItem(c.Context(), cartID, bookID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to update item",
			Message: err.Error(),
		})
	}

	return c.JSON(cart)
}

// RemoveItem removes an item from the cart
// DELETE /api/v1/cart/:cartId/items/:bookId
func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	cartIDStr := c.Params("cartId")
	cartID, err := parseUUID(cartIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid cart ID",
			Message: err.Error(),
		})
	}

	bookIDStr := c.Params("bookId")
	bookID, err := parseUUID(bookIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid book ID",
			Message: err.Error(),
		})
	}

	cart, err := h.service.RemoveItem(c.Context(), cartID, bookID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Failed to remove item",
			Message: err.Error(),
		})
	}

	return c.JSON(cart)
}

// ClearCart clears all items from the cart
// DELETE /api/v1/cart/:cartId
func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	cartIDStr := c.Params("cartId")
	cartID, err := parseUUID(cartIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid cart ID",
			Message: err.Error(),
		})
	}

	if err := h.service.ClearCart(c.Context(), cartID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to clear cart",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Cart cleared successfully",
	})
}
