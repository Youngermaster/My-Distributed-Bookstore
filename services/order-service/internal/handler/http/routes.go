package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/service"
)

func SetupRoutes(app *fiber.App, orderService service.OrderService) {
	handler := NewOrderHandler(orderService)

	// API v1 routes
	api := app.Group("/api/v1")

	// Order routes
	orders := api.Group("/orders")
	orders.Post("/", handler.CreateOrder)
	orders.Get("/", handler.ListOrders)
	orders.Get("/:id", handler.GetOrder)
	orders.Patch("/:id/status", handler.UpdateOrderStatus)
	orders.Post("/:id/cancel", handler.CancelOrder)

	// User orders routes
	api.Get("/users/:userId/orders", handler.GetUserOrders)
}
