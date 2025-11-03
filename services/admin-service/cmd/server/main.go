package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/config"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/grpc"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/handler"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/middleware"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/repository"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := initDatabase(cfg.DatabaseURL, cfg.DBMaxIdleConns, cfg.DBMaxOpenConns)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database connected successfully")

	// Initialize service clients
	clients := grpc.NewServiceClients(
		cfg.UserServiceURL,
		cfg.CatalogServiceURL,
		cfg.OrderServiceURL,
		cfg.CartServiceURL,
		cfg.InventoryServiceURL,
		cfg.RecommendationServiceURL,
	)

	log.Println("Service clients initialized")

	// Initialize repository
	adminRepo := repository.NewAdminRepository(db)

	// Initialize services
	analyticsService := service.NewAnalyticsService(adminRepo, clients)

	// Initialize handlers
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	healthHandler := handler.NewHealthHandler(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Admin Service",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check routes (no auth required)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// API routes with authentication and admin authorization
	api := app.Group("/api/v1")
	
	// Apply auth middleware to all API routes
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	api.Use(middleware.RequireAdmin())

	// Analytics routes
	admin := api.Group("/admin")
	admin.Get("/dashboard", analyticsHandler.GetDashboard)
	admin.Get("/analytics/sales", analyticsHandler.GetSalesAnalytics)
	admin.Get("/analytics/inventory", analyticsHandler.GetInventoryReport)
	admin.Get("/analytics/users", analyticsHandler.GetUserGrowth)
	admin.Get("/top-books", analyticsHandler.GetTopSellingBooks)

	// Start server
	port := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("Starting Admin Service on port %s", cfg.HTTPPort)
	log.Printf("Environment: %s", cfg.Env)
	log.Printf("CORS allowed origins: %s", cfg.CORSAllowOrigins)
	
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(databaseURL string, maxIdleConns, maxOpenConns int) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
