package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

type CartHandler struct {
	cartClient *proxy.CartClient
}

func NewCartHandler(cartClient *proxy.CartClient) *CartHandler {
	return &CartHandler{
		cartClient: cartClient,
	}
}

// ProxyCartRequest is a generic handler that proxies requests to the cart service
func (h *CartHandler) ProxyCartRequest(c *fiber.Ctx) error {
	// Extract path after /api/v1/cart
	path := strings.TrimPrefix(c.Path(), "/api/v1/cart")

	// Build full path
	fullPath := "/api/v1/cart" + path

	// Add query parameters if any
	if c.Request().URI().QueryString() != nil {
		fullPath += "?" + string(c.Request().URI().QueryString())
	}

	// Read request body
	var body []byte
	if c.Body() != nil {
		body = c.Body()
	}

	// Copy relevant headers
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if keyStr == "Content-Type" || keyStr == "Authorization" {
			headers[keyStr] = string(value)
		}
	})

	// Make proxied request
	responseBody, statusCode, err := h.cartClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Failed to communicate with cart service",
			"details": err.Error(),
		})
	}

	// Set response content type
	c.Set("Content-Type", "application/json")

	return c.Status(statusCode).Send(responseBody)
}
