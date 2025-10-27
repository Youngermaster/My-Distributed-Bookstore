package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

type CatalogHandler struct {
	catalogClient *proxy.CatalogClient
}

func NewCatalogHandler(catalogClient *proxy.CatalogClient) *CatalogHandler {
	return &CatalogHandler{
		catalogClient: catalogClient,
	}
}

// ProxyCatalogRequest is a generic handler that proxies requests to the catalog service
func (h *CatalogHandler) ProxyCatalogRequest(c *fiber.Ctx) error {
	// Extract path after /api/v1/catalog
	path := strings.TrimPrefix(c.Path(), "/api/v1/catalog")

	// Build full path
	fullPath := "/api/v1" + path

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
		// Only copy specific headers
		if keyStr == "Content-Type" || keyStr == "Authorization" {
			headers[keyStr] = string(value)
		}
	})

	// Make proxied request
	responseBody, statusCode, err := h.catalogClient.ProxyRequest(
		c.Method(),
		fullPath,
		body,
		headers,
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Service unavailable",
			"message": err.Error(),
		})
	}

	// Set response status and body
	c.Status(statusCode)
	c.Set("Content-Type", "application/json")

	return c.Send(responseBody)
}

// GetBooks handles GET /api/v1/catalog/books
func (h *CatalogHandler) GetBooks(c *fiber.Ctx) error {
	queryString := string(c.Request().URI().QueryString())

	result, err := h.catalogClient.GetBooks(queryString)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "Failed to fetch books",
			"message": err.Error(),
		})
	}

	return c.JSON(result)
}

// GetBookByID handles GET /api/v1/catalog/books/:id
func (h *CatalogHandler) GetBookByID(c *fiber.Ctx) error {
	bookID := c.Params("id")
	if bookID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request",
			"message": "Book ID is required",
		})
	}

	result, err := h.catalogClient.GetBookByID(bookID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Book not found",
			"message": err.Error(),
		})
	}

	return c.JSON(result)
}

// SearchBooks handles GET /api/v1/catalog/books/search
func (h *CatalogHandler) SearchBooks(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request",
			"message": "Search query is required",
		})
	}

	result, err := h.catalogClient.SearchBooks(query)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "Search failed",
			"message": err.Error(),
		})
	}

	return c.JSON(result)
}

// CreateBook handles POST /api/v1/catalog/books
func (h *CatalogHandler) CreateBook(c *fiber.Ctx) error {
	body := c.Body()

	responseBody, statusCode, err := h.catalogClient.ProxyRequest(
		"POST",
		"/api/v1/books",
		body,
		map[string]string{"Content-Type": "application/json"},
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Failed to create book",
			"message": err.Error(),
		})
	}

	c.Status(statusCode)
	c.Set("Content-Type", "application/json")
	return c.Send(responseBody)
}

// UpdateBook handles PUT /api/v1/catalog/books/:id
func (h *CatalogHandler) UpdateBook(c *fiber.Ctx) error {
	bookID := c.Params("id")
	body := c.Body()

	responseBody, statusCode, err := h.catalogClient.ProxyRequest(
		"PUT",
		"/api/v1/books/"+bookID,
		body,
		map[string]string{"Content-Type": "application/json"},
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Failed to update book",
			"message": err.Error(),
		})
	}

	c.Status(statusCode)
	c.Set("Content-Type", "application/json")
	return c.Send(responseBody)
}

// DeleteBook handles DELETE /api/v1/catalog/books/:id
func (h *CatalogHandler) DeleteBook(c *fiber.Ctx) error {
	bookID := c.Params("id")

	responseBody, statusCode, err := h.catalogClient.ProxyRequest(
		"DELETE",
		"/api/v1/books/"+bookID,
		nil,
		nil,
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Failed to delete book",
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

// Generic proxy handlers for other resources

func (h *CatalogHandler) ProxyAuthors(c *fiber.Ctx) error {
	return h.proxyToPath(c, "/authors")
}

func (h *CatalogHandler) ProxyCategories(c *fiber.Ctx) error {
	return h.proxyToPath(c, "/categories")
}

func (h *CatalogHandler) ProxyPublishers(c *fiber.Ctx) error {
	return h.proxyToPath(c, "/publishers")
}

func (h *CatalogHandler) proxyToPath(c *fiber.Ctx, basePath string) error {
	// Build path
	path := "/api/v1" + basePath
	if c.Params("*") != "" {
		path += "/" + c.Params("*")
	}

	// Add query parameters
	if c.Request().URI().QueryString() != nil {
		path += "?" + string(c.Request().URI().QueryString())
	}

	// Read body
	var body []byte
	if c.Body() != nil && len(c.Body()) > 0 {
		body = c.Body()
	}

	// Proxy request
	responseBody, statusCode, err := h.catalogClient.ProxyRequest(
		c.Method(),
		path,
		body,
		map[string]string{"Content-Type": "application/json"},
	)

	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"error":   "Service request failed",
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
