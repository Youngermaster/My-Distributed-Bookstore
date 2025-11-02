package main

import (
"context"
"fmt"
"log"
"os"
"os/signal"
"syscall"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/cors"
"github.com/gofiber/fiber/v2/middleware/logger"
"github.com/gofiber/fiber/v2/middleware/recover"
"github.com/youngermaster/distributed-bookstore/cart-service/internal/config"
httphandler "github.com/youngermaster/distributed-bookstore/cart-service/internal/handler/http"
"github.com/youngermaster/distributed-bookstore/cart-service/internal/repository"
"github.com/youngermaster/distributed-bookstore/cart-service/internal/service"
)

func main() {
// Load configuration
cfg, err := config.Load()
if err != nil {
log.Fatalf("Failed to load configuration: %v", err)
}
cfg.Print()

// Initialize Redis connection
redisClient, err := config.InitRedis(cfg)
if err != nil {
log.Fatalf("Failed to initialize Redis: %v", err)
}
defer redisClient.Close()

// Initialize repository
cartRepo := repository.NewCartRepository(redisClient)

// Initialize service
cartService := service.NewCartService(
cartRepo,
cfg.CartTTLHours,
cfg.MaxItemsPerCart,
cfg.MaxQuantityPerItem,
)

// Create Fiber app
app := fiber.New(fiber.Config{
AppName:      "cart-service",
ErrorHandler: customErrorHandler,
})

// Global middleware
app.Use(recover.New())
app.Use(logger.New(logger.Config{
Format:     "--build{time} | --build{status} | --build{latency} | --build{method} --build{path}\n",
TimeFormat: "2006-01-02 15:04:05",
}))
app.Use(cors.New(cors.Config{
AllowOrigins: "*",
AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))

// Health check endpoints
app.Get("/health", createHealthCheckHandler(redisClient))
app.Get("/ready", createReadinessHandler(redisClient))

// Setup HTTP routes
httphandler.SetupRoutes(app, cartService)

// Get port from environment or use default
httpPort := cfg.HTTPPort

// Start HTTP server in goroutine
go func() {
addr := fmt.Sprintf(":%s", httpPort)
log.Printf("🚀 cart-service HTTP server starting on %s", addr)
if err := app.Listen(addr); err != nil {
log.Fatalf("Failed to start HTTP server: %v", err)
}
}()

// Graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("🛑 Shutting down cart-service...")

if err := app.Shutdown(); err != nil {
log.Printf("Error during HTTP server shutdown: %v", err)
}

log.Println("✅ cart-service stopped gracefully")
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

func createHealthCheckHandler(redisClient interface{}) fiber.Handler {
return func(c *fiber.Ctx) error {
return c.JSON(fiber.Map{
"status":  "healthy",
"service": "cart-service",
})
}
}

func createReadinessHandler(redisClient interface{}) fiber.Handler {
return func(c *fiber.Ctx) error {
// Check Redis connection
if client, ok := redisClient.(interface{ Ping(context.Context) error }); ok {
if err := client.Ping(c.Context()); err != nil {
return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
"status": "not ready",
"redis":  "unhealthy",
})
}
}

return c.JSON(fiber.Map{
"status": "ready",
"redis":  "healthy",
})
}
}
