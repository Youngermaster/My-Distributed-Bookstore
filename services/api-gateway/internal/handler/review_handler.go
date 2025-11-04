package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

// ReviewHandler proxies requests to the review service.
type ReviewHandler struct {
	reviewClient *proxy.ReviewClient
}

// NewReviewHandler constructs a ReviewHandler.
func NewReviewHandler(reviewClient *proxy.ReviewClient) *ReviewHandler {
	return &ReviewHandler{reviewClient: reviewClient}
}

// ProxyReviewRequest forwards the incoming request to the review service.
func (h *ReviewHandler) ProxyReviewRequest(c *fiber.Ctx) error {
	path := strings.TrimPrefix(c.Path(), "/api/v1/reviews")
	fullPath := "/api/v1/reviews" + path

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

	responseBody, statusCode, err := h.reviewClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Review service unavailable",
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
