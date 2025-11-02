package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/cart-service/internal/service"
)

func SetupRoutes(app *fiber.App, cartService service.CartService) {
	handler := NewCartHandler(cartService)

	// API v1 routes
	api := app.Group("/api/v1")

	// Cart routes
	cart := api.Group("/cart")
	cart.Get("/:cartId", handler.GetCart)
	cart.Post("/:cartId/items", handler.AddItem)
	cart.Put("/:cartId/items/:bookId", handler.UpdateItem)
	cart.Delete("/:cartId/items/:bookId", handler.RemoveItem)
	cart.Delete("/:cartId", handler.ClearCart)
}
