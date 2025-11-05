package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

type RecommendationHandler struct {
	recommendationClient *proxy.RecommendationClient
}

func NewRecommendationHandler(recommendationClient *proxy.RecommendationClient) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationClient: recommendationClient,
	}
}

// ProxyRecommendationRequest is a generic handler that proxies requests to the recommendation service
func (h *RecommendationHandler) ProxyRecommendationRequest(c *fiber.Ctx) error {
	// Extract path after /api/v1/recommendations
	path := strings.TrimPrefix(c.Path(), "/api/v1/recommendations")

	// Build full path
	fullPath := "/api/v1/recommendations" + path

	// Add query parameters if any
	if c.Request().URI().QueryString() != nil {
		fullPath += "?" + string(c.Request().URI().QueryString())
	}

	// Read request body
	var body []byte
	if c.Body() != nil {
		body = c.Body()
	}

	// Copy relevant headers (including X-User-Id for auth)
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		// Copy authentication and content headers
		if keyStr == "Content-Type" || keyStr == "Authorization" || keyStr == "X-User-Id" {
			headers[keyStr] = string(value)
		}
	})

	// Make proxied request
	responseBody, statusCode, err := h.recommendationClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Failed to communicate with recommendation service",
			"details": err.Error(),
		})
	}

	// Set response content type
	c.Set("Content-Type", "application/json")

	return c.Status(statusCode).Send(responseBody)
}

// GetMyRecommendations handles personalized recommendations
func (h *RecommendationHandler) GetMyRecommendations(c *fiber.Ctx) error {
	// Extract user ID from header (set by auth middleware)
	userID := c.Get("X-User-Id")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in request",
		})
	}

	// Get limit from query params
	limit := c.QueryInt("limit", 10)

	// Call recommendation service
	result, err := h.recommendationClient.GetRecommendations(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get recommendations",
			"details": err.Error(),
		})
	}

	return c.JSON(result)
}

// GetSimilarBooks handles similar books recommendations
func (h *RecommendationHandler) GetSimilarBooks(c *fiber.Ctx) error {
	bookID := c.Params("bookId")
	if bookID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book ID is required",
		})
	}

	limit := c.QueryInt("limit", 10)

	result, err := h.recommendationClient.GetSimilarBooks(bookID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get similar books",
			"details": err.Error(),
		})
	}

	return c.JSON(result)
}

// GetTrendingBooks handles trending books
func (h *RecommendationHandler) GetTrendingBooks(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	days := c.QueryInt("days", 7)

	result, err := h.recommendationClient.GetTrendingBooks(limit, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get trending books",
			"details": err.Error(),
		})
	}

	return c.JSON(result)
}
