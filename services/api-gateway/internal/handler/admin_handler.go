package handler

import (
	"strings"

	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminClient *proxy.AdminClient
}

func NewAdminHandler(adminClient *proxy.AdminClient) *AdminHandler {
	return &AdminHandler{
		adminClient: adminClient,
	}
}

// ProxyAdminRequest forwards all admin requests to admin service
func (h *AdminHandler) ProxyAdminRequest(c *fiber.Ctx) error {
	// Extract path after /api/v1/admin
	path := strings.TrimPrefix(c.Path(), "/api/v1/admin")
	fullPath := "/api/v1/admin" + path

	// Get request body
	body := c.Body()

	// Copy headers
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	// Forward request
	responseBody, statusCode, err := h.adminClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Failed to communicate with admin service",
		})
	}

	// Set response headers
	c.Set("Content-Type", "application/json")

	return c.Status(statusCode).Send(responseBody)
}
