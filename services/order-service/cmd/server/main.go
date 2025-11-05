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
	"github.com/youngermaster/distributed-bookstore/order-service/internal/config"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/domain"
	httphandler "github.com/youngermaster/distributed-bookstore/order-service/internal/handler/http"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/repository"
	"github.com/youngermaster/distributed-bookstore/order-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	cfg.Print()

	// Initialize database connection
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Get underlying SQL database for connection management
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Run database migrations
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Initialize repository
	orderRepo := repository.NewOrderRepository(db)

	// Initialize service
	orderService := service.NewOrderService(orderRepo, cfg.DefaultPageSize, cfg.MaxPageSize)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "order-service",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoints
	app.Get("/health", createHealthCheckHandler(db))
	app.Get("/ready", createReadinessHandler(db))

	// Setup HTTP routes
	httphandler.SetupRoutes(app, orderService)

	// Get port from environment or use default
	httpPort := cfg.HTTPPort

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", httpPort)
		log.Printf("🚀 order-service HTTP server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down order-service...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Error during HTTP server shutdown: %v", err)
	}

	log.Println("order-service stopped gracefully")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   "Internal Server Error",
		"message": err.Error(),
	})
}

func createHealthCheckHandler(db interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "order-service",
		})
	}
}

func createReadinessHandler(db interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check database connection
		if gormDB, ok := db.(interface{ DB() (interface{}, error) }); ok {
			if sqlDB, err := gormDB.DB(); err == nil {
				if pinger, ok := sqlDB.(interface{ Ping() error }); ok {
					if err := pinger.Ping(); err != nil {
						return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
							"status":   "not ready",
							"database": "unhealthy",
						})
					}
				}
			}
		}

		return c.JSON(fiber.Map{
			"status":   "ready",
			"database": "healthy",
		})
	}
}
