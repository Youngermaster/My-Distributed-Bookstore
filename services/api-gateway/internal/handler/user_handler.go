package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

type UserHandler struct {
	userClient *proxy.UserClient
}

func NewUserHandler(userClient *proxy.UserClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// ProxyUserRequest is a generic handler that proxies requests to the user service
func (h *UserHandler) ProxyUserRequest(c *fiber.Ctx) error {
	// Extract path after /api/v1/users
	path := strings.TrimPrefix(c.Path(), "/api/v1/users")

	// Build full path
	fullPath := "/api/v1/users" + path

	// Add query parameters if any
	if c.Request().URI().QueryString() != nil {
		fullPath += "?" + string(c.Request().URI().QueryString())
	}

	// Read request body
	var body []byte
	if c.Body() != nil {
		body = c.Body()
	}

	// Copy relevant headers (Authorization is important here)
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		// Copy auth and content headers
		if keyStr == "Content-Type" || keyStr == "Authorization" {
			headers[keyStr] = string(value)
		}
	})

	// Make proxied request
	responseBody, statusCode, err := h.userClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Failed to communicate with user service",
			"details": err.Error(),
		})
	}

	// Set response content type
	c.Set("Content-Type", "application/json")

	return c.Status(statusCode).Send(responseBody)
}
