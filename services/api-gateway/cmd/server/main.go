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

	// Initialize catalog service client
	catalogClient := proxy.NewCatalogClient(cfg.CatalogServiceURL, cfg.RequestTimeout)

	// Initialize recommendation service client
	recommendationClient := proxy.NewRecommendationClient(cfg.RecommendationServiceURL, cfg.RequestTimeout)

	// Initialize handlers
	catalogHandler := handler.NewCatalogHandler(catalogClient)
	recommendationHandler := handler.NewRecommendationHandler(recommendationClient)

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
		log.Printf("✅ Rate limiting enabled: %d requests per %v", cfg.RateLimitMax, cfg.RateLimitWindow)
	}

	// Health check endpoints
	app.Get("/health", createHealthCheckHandler(catalogClient, recommendationClient))
	app.Get("/ready", createReadinessHandler(catalogClient, recommendationClient))

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

	// Recommendation routes (proxy to recommendation-service)
	recommendations := api.Group("/recommendations")
	recommendations.Get("/me", recommendationHandler.GetMyRecommendations)
	recommendations.Get("/similar/:bookId", recommendationHandler.GetSimilarBooks)
	recommendations.Get("/trending", recommendationHandler.GetTrendingBooks)
	recommendations.All("*", recommendationHandler.ProxyRecommendationRequest)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "API Gateway",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": fiber.Map{
				"health":           "/health",
				"ready":            "/ready",
				"catalog":          "/api/v1/catalog",
				"books":            "/api/v1/catalog/books",
				"authors":          "/api/v1/catalog/authors",
				"categories":       "/api/v1/catalog/categories",
				"publishers":       "/api/v1/catalog/publishers",
				"recommendations":  "/api/v1/recommendations",
				"my_recommendations": "/api/v1/recommendations/me",
				"trending":         "/api/v1/recommendations/trending",
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

	log.Printf("✅ API Gateway started successfully")
	log.Printf("📡 Server: http://localhost:%s", cfg.Port)
	log.Printf("🏥 Health: http://localhost:%s/health", cfg.Port)
	log.Printf("📚 Catalog: http://localhost:%s/api/v1/catalog", cfg.Port)
	log.Printf("🎯 Recommendations: http://localhost:%s/api/v1/recommendations", cfg.Port)
	log.Printf("🔗 Catalog Service: %s", cfg.CatalogServiceURL)
	log.Printf("🔗 Recommendation Service: %s", cfg.RecommendationServiceURL)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down API Gateway...")
	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	}

	log.Println("✅ API Gateway stopped gracefully")
}

func createHealthCheckHandler(catalogClient *proxy.CatalogClient, recommendationClient *proxy.RecommendationClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check catalog service health
		catalogHealthy := true
		if err := catalogClient.HealthCheck(); err != nil {
			catalogHealthy = false
		}

		// Check recommendation service health
		recommendationHealthy := true
		if err := recommendationClient.HealthCheck(); err != nil {
			recommendationHealthy = false
		}

		overallStatus := "ok"
		if !catalogHealthy || !recommendationHealthy {
			overallStatus = "degraded"
		}

		return c.JSON(fiber.Map{
			"status":  overallStatus,
			"service": "api-gateway",
			"services": fiber.Map{
				"catalog": fiber.Map{
					"healthy": catalogHealthy,
				},
				"recommendation": fiber.Map{
					"healthy": recommendationHealthy,
				},
			},
		})
	}
}

func createReadinessHandler(catalogClient *proxy.CatalogClient, recommendationClient *proxy.RecommendationClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if catalog service is reachable
		catalogReady := true
		if err := catalogClient.HealthCheck(); err != nil {
			catalogReady = false
		}

		// Check if recommendation service is reachable
		recommendationReady := true
		if err := recommendationClient.HealthCheck(); err != nil {
			recommendationReady = false
		}

		if !catalogReady {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"reason": "catalog service is not available",
			})
		}

		if !recommendationReady {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "not ready",
				"reason": "recommendation service is not available",
			})
		}

		return c.JSON(fiber.Map{
			"status": "ready",
			"services": fiber.Map{
				"catalog": "ready",
				"recommendation": "ready",
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
