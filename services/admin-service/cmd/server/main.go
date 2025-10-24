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
)

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "admin-service",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[]  -   - \n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", healthCheckHandler)
	app.Get("/ready", readinessHandler)

	// TODO: Initialize database connection
	// TODO: Run database migrations
	// TODO: Initialize Redis connection
	// TODO: Initialize RabbitMQ connection
	// TODO: Setup gRPC server
	// TODO: Register HTTP routes
	// TODO: Setup event consumers
	// TODO: Initialize Jaeger tracer
	// TODO: Setup Prometheus metrics

	// API routes group
	api := app.Group("/api/v1")
	
	// TODO: Add route handlers here
	_ = api // Remove this line when routes are added

	// Get port from environment or use default
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8090"
	}

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", httpPort)
		log.Printf("admin-service HTTP server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// TODO: Start gRPC server in separate goroutine

	log.Printf("admin-service started successfully")
	log.Printf("HTTP Port: %s", httpPort)
	log.Printf("gRPC Port: 50060")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down admin-service...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during HTTP server shutdown: %v", err)
	}

	// TODO: Stop gRPC server
	// TODO: Close database connections
	// TODO: Close Redis connections
	// TODO: Close RabbitMQ connections

	log.Println("admin-service stopped gracefully")
}

func healthCheckHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "admin-service",
	})
}

func readinessHandler(c *fiber.Ctx) error {
	// TODO: Check database connection
	// TODO: Check Redis connection
	// TODO: Check RabbitMQ connection
	
	return c.JSON(fiber.Map{
		"status": "ready",
		"checks": fiber.Map{
			"database": "ok",
			"redis":    "ok",
			"rabbitmq": "ok",
		},
	})
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
