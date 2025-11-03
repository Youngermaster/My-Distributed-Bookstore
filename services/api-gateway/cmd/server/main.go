package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/config"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/handler"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/middleware"
	"github.com/youngermaster/distributed-bookstore/api-gateway/internal/proxy"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg.Print()

	// Initialize service clients
	catalogClient := proxy.NewCatalogClient(cfg.CatalogServiceURL, cfg.RequestTimeout)
	userClient := proxy.NewUserClient(cfg.UserServiceURL, cfg.RequestTimeout)
	cartClient := proxy.NewCartClient(cfg.CartServiceURL, cfg.RequestTimeout)
	orderClient := proxy.NewOrderClient(cfg.OrderServiceURL, cfg.RequestTimeout)
	recommendationClient := proxy.NewRecommendationClient(cfg.RecommendationServiceURL, cfg.RequestTimeout)
	adminClient := proxy.NewAdminClient(cfg.AdminServiceURL, cfg.RequestTimeout)

	// Initialize handlers
	catalogHandler := handler.NewCatalogHandler(catalogClient)
	userHandler := handler.NewUserHandler(userClient)
	cartHandler := handler.NewCartHandler(cartClient)
	orderHandler := handler.NewOrderHandler(orderClient)
	recommendationHandler := handler.NewRecommendationHandler(recommendationClient)
	adminHandler := handler.NewAdminHandler(adminClient)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "API Gateway",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Rate limiting middleware
	if cfg.RateLimitEnabled {
		rateLimiter := middleware.NewRateLimiter(cfg.RateLimitMax, cfg.RateLimitWindow)
		app.Use(rateLimiter.Middleware())
		log.Printf("Rate limiting enabled: %d requests per %v", cfg.RateLimitMax, cfg.RateLimitWindow)
	}

	// Health check endpoints
	app.Get("/health", createHealthCheckHandler(catalogClient, userClient, cartClient, orderClient, recommendationClient, adminClient))
	app.Get("/ready", createReadinessHandler(catalogClient, userClient, cartClient, orderClient, recommendationClient, adminClient))

	// API routes
	api := app.Group("/api/v1")

	// Catalog routes (proxy to catalog-service)
	catalog := api.Group("/catalog")

	// Books
	books := catalog.Group("/books")
	books.Get("/", catalogHandler.GetBooks)
	books.Get("/search", catalogHandler.SearchBooks)
	books.Get("/:id", catalogHandler.GetBookByID)
	books.Post("/", catalogHandler.CreateBook)
	books.Put("/:id", catalogHandler.UpdateBook)
	books.Delete("/:id", catalogHandler.DeleteBook)
	books.Patch("/:id/stock", catalogHandler.ProxyCatalogRequest)

	// Authors
	catalog.All("/authors*", catalogHandler.ProxyAuthors)

	// Categories
	catalog.All("/categories*", catalogHandler.ProxyCategories)

	// Publishers
	catalog.All("/publishers*", catalogHandler.ProxyPublishers)

	// User routes (proxy to user-service)
	users := api.Group("/users")
	users.All("*", userHandler.ProxyUserRequest)

	// Cart routes (proxy to cart-service)
	cart := api.Group("/cart")
	cart.All("*", cartHandler.ProxyCartRequest)

	// Order routes (proxy to order-service)
	orders := api.Group("/orders")
	orders.All("*", orderHandler.ProxyOrderRequest)

	// Recommendation routes (proxy to recommendation-service)
	recommendations := api.Group("/recommendations")
	recommendations.Get("/me", recommendationHandler.GetMyRecommendations)
	recommendations.Get("/similar/:bookId", recommendationHandler.GetSimilarBooks)
	recommendations.Get("/trending", recommendationHandler.GetTrendingBooks)
	recommendations.All("*", recommendationHandler.ProxyRecommendationRequest)

	// Admin routes (proxy to admin-service) - requires authentication and admin role
	admin := api.Group("/admin")
	admin.All("*", adminHandler.ProxyAdminRequest)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "API Gateway",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": fiber.Map{
				"health":             "/health",
				"ready":              "/ready",
				"catalog":            "/api/v1/catalog",
				"books":              "/api/v1/catalog/books",
				"authors":            "/api/v1/catalog/authors",
				"categories":         "/api/v1/catalog/categories",
				"publishers":         "/api/v1/catalog/publishers",
				"users":              "/api/v1/users",
				"auth":               "/api/v1/users/auth",
				"wishlist":           "/api/v1/users/me/wishlist",
				"cart":               "/api/v1/cart",
				"orders":             "/api/v1/orders",
				"recommendations":    "/api/v1/recommendations",
				"my_recommendations": "/api/v1/recommendations/me",
				"trending":           "/api/v1/recommendations/trending",
			},
		})
	})

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("🚀 API Gateway starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("API Gateway started successfully")
	log.Printf("Server: http://localhost:%s", cfg.Port)
	log.Printf("Health: http://localhost:%s/health", cfg.Port)
	log.Printf("\nAvailable Routes:")
	log.Printf("  Catalog: http://localhost:%s/api/v1/catalog", cfg.Port)
	log.Printf("  Users: http://localhost:%s/api/v1/users", cfg.Port)
	log.Printf("  Cart: http://localhost:%s/api/v1/cart", cfg.Port)
	log.Printf("  Orders: http://localhost:%s/api/v1/orders", cfg.Port)
	log.Printf("  Recommendations: http://localhost:%s/api/v1/recommendations", cfg.Port)
	log.Printf("  Admin: http://localhost:%s/api/v1/admin", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")
	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	}

	log.Println("API Gateway stopped gracefully")
}

func createHealthCheckHandler(
	catalogClient *proxy.CatalogClient,
	userClient *proxy.UserClient,
	cartClient *proxy.CartClient,
	orderClient *proxy.OrderClient,
	recommendationClient *proxy.RecommendationClient,
	adminClient *proxy.AdminClient,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check all services health
		catalogHealthy := catalogClient.HealthCheck() == nil
		userHealthy := userClient.HealthCheck() == nil
		cartHealthy := cartClient.HealthCheck() == nil
		orderHealthy := orderClient.HealthCheck() == nil
		recommendationHealthy := recommendationClient.HealthCheck() == nil
		adminHealthy := adminClient.HealthCheck() == nil

		allHealthy := catalogHealthy && userHealthy && cartHealthy && orderHealthy && recommendationHealthy && adminHealthy
		overallStatus := "ok"
		if !allHealthy {
			overallStatus = "degraded"
		}

		return c.JSON(fiber.Map{
			"status":  overallStatus,
			"service": "api-gateway",
			"services": fiber.Map{
				"catalog":        fiber.Map{"healthy": catalogHealthy},
				"user":           fiber.Map{"healthy": userHealthy},
				"cart":           fiber.Map{"healthy": cartHealthy},
				"order":          fiber.Map{"healthy": orderHealthy},
				"recommendation": fiber.Map{"healthy": recommendationHealthy},
				"admin":          fiber.Map{"healthy": adminHealthy},
			},
		})
	}
}

func createReadinessHandler(
	catalogClient *proxy.CatalogClient,
	userClient *proxy.UserClient,
	cartClient *proxy.CartClient,
	orderClient *proxy.OrderClient,
	recommendationClient *proxy.RecommendationClient,
	adminClient *proxy.AdminClient,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check all critical services are reachable
		notReadyServices := []string{}

		if err := catalogClient.HealthCheck(); err != nil {
			notReadyServices = append(notReadyServices, "catalog")
		}
		if err := userClient.HealthCheck(); err != nil {
			notReadyServices = append(notReadyServices, "user")
		}
		if err := cartClient.HealthCheck(); err != nil {
			notReadyServices = append(notReadyServices, "cart")
		}
		if err := orderClient.HealthCheck(); err != nil {
			notReadyServices = append(notReadyServices, "order")
		}
		if err := recommendationClient.HealthCheck(); err != nil {
			notReadyServices = append(notReadyServices, "recommendation")
		}
		if err := adminClient.HealthCheck(); err != nil {
			notReadyServices = append(notReadyServices, "admin")
		}

		if len(notReadyServices) > 0 {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":            "not ready",
				"unavailable_services": notReadyServices,
			})
		}

		return c.JSON(fiber.Map{
			"status": "ready",
			"services": fiber.Map{
				"catalog":        "ready",
				"user":           "ready",
				"cart":           "ready",
				"order":          "ready",
				"recommendation": "ready",
				"admin":          "ready",
			},
		})
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   "Request failed",
		"message": err.Error(),
		"code":    code,
	})
}
