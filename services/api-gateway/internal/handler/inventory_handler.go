package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

// InventoryHandler proxies requests to the inventory service.
type InventoryHandler struct {
	inventoryClient *proxy.InventoryClient
}

// NewInventoryHandler constructs an InventoryHandler.
func NewInventoryHandler(inventoryClient *proxy.InventoryClient) *InventoryHandler {
	return &InventoryHandler{inventoryClient: inventoryClient}
}

// ProxyInventoryRequest forwards requests under /api/v1/inventory to the inventory service.
func (h *InventoryHandler) ProxyInventoryRequest(c *fiber.Ctx) error {
	path := strings.TrimPrefix(c.Path(), "/api/v1/inventory")
	fullPath := "/api/v1/inventory" + path

	if query := c.Request().URI().QueryString(); len(query) > 0 {
		fullPath += "?" + string(query)
	}

	var body []byte
	if data := c.Body(); len(data) > 0 {
		body = data
	}

	headers := map[string]string{}
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		if keyStr == "Content-Type" || keyStr == "Authorization" {
			headers[keyStr] = string(value)
		}
	})

	responseBody, statusCode, err := h.inventoryClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Inventory service unavailable",
			"message": err.Error(),
		})
	}

	c.Status(statusCode)
	if len(responseBody) > 0 {
		c.Set("Content-Type", "application/json")
		return c.Send(responseBody)
	}

	return c.SendStatus(statusCode)
}
