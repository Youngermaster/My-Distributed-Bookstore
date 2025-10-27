package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/repository"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/service"
)

type BookHandler struct {
	service service.CatalogService
}

func NewBookHandler(service service.CatalogService) *BookHandler {
	return &BookHandler{service: service}
}

// CreateBook godoc
// @Summary Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body service.CreateBookRequest true "Book to create"
// @Success 201 {object} domain.Book
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books [post]
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var req service.CreateBookRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	book, err := h.service.CreateBook(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to create book",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

// GetBook godoc
// @Summary Get book by ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/books/{id} [get]
func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid book ID",
			Message: err.Error(),
		})
	}

	book, err := h.service.GetBook(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Book not found",
			Message: err.Error(),
		})
	}

	return c.JSON(book)
}

// ListBooks godoc
// @Summary List books with filters and pagination
// @Tags books
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param category_id query string false "Filter by category ID"
// @Param author_id query string false "Filter by author ID"
// @Param publisher_id query string false "Filter by publisher ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param in_stock query boolean false "Only in-stock books"
// @Param sort_by query string false "Sort by field" Enums(title, price, created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Success 200 {object} service.BookListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books [get]
func (h *BookHandler) ListBooks(c *fiber.Ctx) error {
	filter := repository.BookFilter{
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("page_size", 20),
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	// Parse optional UUID filters
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filter.CategoryID = &id
		}
	}

	if authorID := c.Query("author_id"); authorID != "" {
		if id, err := uuid.Parse(authorID); err == nil {
			filter.AuthorID = &id
		}
	}

	if publisherID := c.Query("publisher_id"); publisherID != "" {
		if id, err := uuid.Parse(publisherID); err == nil {
			filter.PublisherID = &id
		}
	}

	// Parse price filters
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = &price
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = &price
		}
	}

	// Parse boolean filters
	if inStock := c.Query("in_stock"); inStock != "" {
		if stock, err := strconv.ParseBool(inStock); err == nil {
			filter.InStock = &stock
		}
	}

	result, err := h.service.ListBooks(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to list books",
			Message: err.Error(),
		})
	}

	return c.JSON(result)
}

// SearchBooks godoc
// @Summary Search books by query
// @Tags books
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} service.BookListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/search [get]
func (h *BookHandler) SearchBooks(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Search query is required",
			Message: "Please provide a search query using the 'q' parameter",
		})
	}

	filter := repository.BookFilter{
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("page_size", 20),
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	result, err := h.service.SearchBooks(c.Context(), query, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to search books",
			Message: err.Error(),
		})
	}

	return c.JSON(result)
}

// UpdateBook godoc
// @Summary Update a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body service.UpdateBookRequest true "Book updates"
// @Success 200 {object} domain.Book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/{id} [put]
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid book ID",
			Message: err.Error(),
		})
	}

	var req service.UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	book, err := h.service.UpdateBook(c.Context(), id, req)
	if err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Book not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to update book",
			Message: err.Error(),
		})
	}

	return c.JSON(book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Tags books
// @Param id path string true "Book ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/{id} [delete]
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid book ID",
			Message: err.Error(),
		})
	}

	if err := h.service.DeleteBook(c.Context(), id); err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Book not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to delete book",
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateBookStock godoc
// @Summary Update book stock quantity
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param stock body StockUpdateRequest true "Stock update"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/books/{id}/stock [patch]
func (h *BookHandler) UpdateBookStock(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid book ID",
			Message: err.Error(),
		})
	}

	var req StockUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	if err := h.service.UpdateBookStock(c.Context(), id, req.Quantity); err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Book not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to update stock",
			Message: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Stock updated successfully",
	})
}

type StockUpdateRequest struct {
	Quantity int `json:"quantity" validate:"gte=0"`
}
