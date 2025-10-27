package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/service"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(app *fiber.App, catalogService service.CatalogService) {
	// Create handlers
	bookHandler := NewBookHandler(catalogService)

	// API v1 group
	api := app.Group("/api/v1")

	// Book routes
	books := api.Group("/books")
	books.Get("/", bookHandler.ListBooks)
	books.Get("/search", bookHandler.SearchBooks)
	books.Get("/:id", bookHandler.GetBook)
	books.Post("/", bookHandler.CreateBook)
	books.Put("/:id", bookHandler.UpdateBook)
	books.Delete("/:id", bookHandler.DeleteBook)
	books.Patch("/:id/stock", bookHandler.UpdateBookStock)

	// Author routes
	authors := api.Group("/authors")
	authors.Get("/", createGetAuthorsHandler(catalogService))
	authors.Get("/:id", createGetAuthorHandler(catalogService))
	authors.Post("/", createCreateAuthorHandler(catalogService))
	authors.Put("/:id", createUpdateAuthorHandler(catalogService))
	authors.Delete("/:id", createDeleteAuthorHandler(catalogService))

	// Category routes
	categories := api.Group("/categories")
	categories.Get("/", createGetCategoriesHandler(catalogService))
	categories.Get("/:id", createGetCategoryHandler(catalogService))
	categories.Post("/", createCreateCategoryHandler(catalogService))
	categories.Put("/:id", createUpdateCategoryHandler(catalogService))
	categories.Delete("/:id", createDeleteCategoryHandler(catalogService))

	// Publisher routes
	publishers := api.Group("/publishers")
	publishers.Get("/", createGetPublishersHandler(catalogService))
	publishers.Get("/:id", createGetPublisherHandler(catalogService))
	publishers.Post("/", createCreatePublisherHandler(catalogService))
	publishers.Put("/:id", createUpdatePublisherHandler(catalogService))
	publishers.Delete("/:id", createDeletePublisherHandler(catalogService))
}

// Author handlers (inline for simplicity)

func createGetAuthorsHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 20)

		result, err := svc.ListAuthors(c.Context(), page, pageSize)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to list authors",
				Message: err.Error(),
			})
		}

		return c.JSON(result)
	}
}

func createGetAuthorHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid author ID",
				Message: err.Error(),
			})
		}

		author, err := svc.GetAuthor(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Author not found",
				Message: err.Error(),
			})
		}

		return c.JSON(author)
	}
}

func createCreateAuthorHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req service.CreateAuthorRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid request body",
				Message: err.Error(),
			})
		}

		author, err := svc.CreateAuthor(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to create author",
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(author)
	}
}

func createUpdateAuthorHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid author ID",
				Message: err.Error(),
			})
		}

		var req service.UpdateAuthorRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid request body",
				Message: err.Error(),
			})
		}

		author, err := svc.UpdateAuthor(c.Context(), id, req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to update author",
				Message: err.Error(),
			})
		}

		return c.JSON(author)
	}
}

func createDeleteAuthorHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid author ID",
				Message: err.Error(),
			})
		}

		if err := svc.DeleteAuthor(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to delete author",
				Message: err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Category handlers

func createGetCategoriesHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hierarchical := c.QueryBool("hierarchical", false)

		categories, err := svc.ListCategories(c.Context(), hierarchical)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to list categories",
				Message: err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"categories": categories,
			"total":      len(categories),
		})
	}
}

func createGetCategoryHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid category ID",
				Message: err.Error(),
			})
		}

		category, err := svc.GetCategory(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Category not found",
				Message: err.Error(),
			})
		}

		return c.JSON(category)
	}
}

func createCreateCategoryHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req service.CreateCategoryRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid request body",
				Message: err.Error(),
			})
		}

		category, err := svc.CreateCategory(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to create category",
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(category)
	}
}

func createUpdateCategoryHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid category ID",
				Message: err.Error(),
			})
		}

		var req service.UpdateCategoryRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid request body",
				Message: err.Error(),
			})
		}

		category, err := svc.UpdateCategory(c.Context(), id, req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to update category",
				Message: err.Error(),
			})
		}

		return c.JSON(category)
	}
}

func createDeleteCategoryHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid category ID",
				Message: err.Error(),
			})
		}

		if err := svc.DeleteCategory(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to delete category",
				Message: err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Publisher handlers

func createGetPublishersHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 20)

		result, err := svc.ListPublishers(c.Context(), page, pageSize)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to list publishers",
				Message: err.Error(),
			})
		}

		return c.JSON(result)
	}
}

func createGetPublisherHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid publisher ID",
				Message: err.Error(),
			})
		}

		publisher, err := svc.GetPublisher(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Publisher not found",
				Message: err.Error(),
			})
		}

		return c.JSON(publisher)
	}
}

func createCreatePublisherHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req service.CreatePublisherRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid request body",
				Message: err.Error(),
			})
		}

		publisher, err := svc.CreatePublisher(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to create publisher",
				Message: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(publisher)
	}
}

func createUpdatePublisherHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid publisher ID",
				Message: err.Error(),
			})
		}

		var req service.UpdatePublisherRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid request body",
				Message: err.Error(),
			})
		}

		publisher, err := svc.UpdatePublisher(c.Context(), id, req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to update publisher",
				Message: err.Error(),
			})
		}

		return c.JSON(publisher)
	}
}

func createDeletePublisherHandler(svc service.CatalogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := parseUUID(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error:   "Invalid publisher ID",
				Message: err.Error(),
			})
		}

		if err := svc.DeletePublisher(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error:   "Failed to delete publisher",
				Message: err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
