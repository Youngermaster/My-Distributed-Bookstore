package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

type OrderHandler struct {
	orderClient *proxy.OrderClient
}

func NewOrderHandler(orderClient *proxy.OrderClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}

// ProxyOrderRequest is a generic handler that proxies requests to the order service
func (h *OrderHandler) ProxyOrderRequest(c *fiber.Ctx) error {
	// Extract path after /api/v1/orders
	path := strings.TrimPrefix(c.Path(), "/api/v1/orders")

	// Build full path
	fullPath := "/api/v1/orders" + path

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
	responseBody, statusCode, err := h.orderClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Failed to communicate with order service",
			"details": err.Error(),
		})
	}

	// Set response content type
	c.Set("Content-Type", "application/json")

	return c.Status(statusCode).Send(responseBody)
}
