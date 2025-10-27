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
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/config"
	httphandler "github.com/youngermaster/distributed-bookstore/catalog-service/internal/handler/http"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/repository"
	"github.com/youngermaster/distributed-bookstore/catalog-service/internal/service"
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
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed database with sample data (only in development)
	if cfg.Env == "development" {
		if err := config.SeedDatabase(db); err != nil {
			log.Printf("Warning: Failed to seed database: %v", err)
		}
	}

	// Initialize repositories
	bookRepo := repository.NewBookRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	publisherRepo := repository.NewPublisherRepository(db)

	// Initialize service
	catalogService := service.NewCatalogService(bookRepo, authorRepo, categoryRepo, publisherRepo)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "catalog-service",
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
	httphandler.SetupRoutes(app, catalogService)

	// Get port from environment or use default
	httpPort := cfg.HTTPPort

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", httpPort)
		log.Printf("🚀 catalog-service HTTP server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// TODO: Start gRPC server in separate goroutine

	log.Printf("✅ catalog-service started successfully")
	log.Printf("📡 HTTP Server: http://localhost:%s", httpPort)
	log.Printf("📡 gRPC Server: localhost:%s", cfg.GRPCPort)
	log.Printf("🏥 Health Check: http://localhost:%s/health", httpPort)
	log.Printf("📚 API Base: http://localhost:%s/api/v1", httpPort)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down catalog-service...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Error during HTTP server shutdown: %v", err)
	}

	// TODO: Stop gRPC server
	// TODO: Close Redis connections
	// TODO: Close RabbitMQ connections

	log.Println("✅ catalog-service stopped gracefully")
}

func createHealthCheckHandler(db interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "catalog-service",
		})
	}
}

func createReadinessHandler(db interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Check database connection
		// TODO: Check Redis connection
		// TODO: Check RabbitMQ connection

		return c.JSON(fiber.Map{
			"status": "ready",
			"checks": fiber.Map{
				"database": "ok",
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
		"error": err.Error(),
		"code":  code,
	})
}
