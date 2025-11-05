package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

type NotificationHandler struct {
	client *proxy.NotificationClient
}

func NewNotificationHandler(client *proxy.NotificationClient) *NotificationHandler {
	return &NotificationHandler{client: client}
}

func (h *NotificationHandler) ProxyNotificationRequest(c *fiber.Ctx) error {
	path := c.Path()
	if query := c.Request().URI().QueryString(); len(query) > 0 {
		path += "?" + string(query)
	}

	body := c.Body()
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	responseBody, statusCode, err := h.client.ProxyRequest(c.Method(), path, body, headers)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Failed to communicate with notification service",
			"details": err.Error(),
		})
	}

	c.Set("Content-Type", "application/json")
	return c.Status(statusCode).Send(responseBody)
}
