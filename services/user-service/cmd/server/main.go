package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/config"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/database"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/dto"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/middleware"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/repository"
	"github.com/youngermaster/distributed-bookstore/user-service/internal/service"
	"github.com/youngermaster/distributed-bookstore/user-service/pkg/jwt"
	"github.com/youngermaster/distributed-bookstore/user-service/pkg/password"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed default roles
	if err := database.SeedDefaultRoles(db); err != nil {
		log.Printf("Warning: Failed to seed default roles: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Initialize utilities
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTAccessDuration, cfg.JWTRefreshDuration)
	passwordSvc := password.NewService(cfg.BcryptCost)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, sessionRepo, jwtService, passwordSvc)
	userService := service.NewUserService(userRepo, addressRepo, passwordSvc)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.CORSAllowOrigins, ","),
		AllowCredentials: cfg.CORSAllowCredentials,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "user-service",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", func(c *fiber.Ctx) error {
		return registerHandler(c, authService)
	})
	auth.Post("/login", func(c *fiber.Ctx) error {
		return loginHandler(c, authService)
	})
	auth.Post("/refresh", func(c *fiber.Ctx) error {
		return refreshTokenHandler(c, authService)
	})

	// User routes (protected)
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtService))
	
	users.Get("/me", func(c *fiber.Ctx) error {
		return getMeHandler(c, userService)
	})
	users.Put("/me", func(c *fiber.Ctx) error {
		return updateProfileHandler(c, userService)
	})
	users.Post("/me/password", func(c *fiber.Ctx) error {
		return changePasswordHandler(c, userService)
	})

	// Address routes (protected)
	users.Get("/me/addresses", func(c *fiber.Ctx) error {
		return getAddressesHandler(c, userService)
	})
	users.Post("/me/addresses", func(c *fiber.Ctx) error {
		return createAddressHandler(c, userService)
	})
	users.Put("/me/addresses/:id", func(c *fiber.Ctx) error {
		return updateAddressHandler(c, userService)
	})
	users.Delete("/me/addresses/:id", func(c *fiber.Ctx) error {
		return deleteAddressHandler(c, userService)
	})

	// Start server
	go func() {
		log.Printf("Starting HTTP server on :%s", cfg.HTTPPort)
		if err := app.Listen(":" + cfg.HTTPPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

// Handlers

func registerHandler(c *fiber.Ctx, authService *service.AuthService) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "validation failed",
			Message: err.Error(),
		})
	}

	user, err := authService.Register(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func loginHandler(c *fiber.Ctx, authService *service.AuthService) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "validation failed",
		})
	}

	response, err := authService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(response)
}

func refreshTokenHandler(c *fiber.Ctx, authService *service.AuthService) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	response, err := authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "invalid refresh token",
		})
	}

	return c.JSON(response)
}

func getMeHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)

	user, err := userService.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: "user not found",
		})
	}

	return c.JSON(user)
}

func updateProfileHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	user, err := userService.UpdateProfile(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(user)
}

func changePasswordHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "validation failed",
		})
	}

	if err := userService.ChangePassword(userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "password changed successfully",
	})
}

func getAddressesHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)

	addresses, err := userService.GetAddresses(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "failed to retrieve addresses",
		})
	}

	return c.JSON(addresses)
}

func createAddressHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.AddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "validation failed",
		})
	}

	address, err := userService.CreateAddress(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(address)
}

func updateAddressHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)
	addressID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid address ID",
		})
	}

	var req dto.AddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid request body",
		})
	}

	address, err := userService.UpdateAddress(userID, addressID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(address)
}

func deleteAddressHandler(c *fiber.Ctx, userService *service.UserService) error {
	userID := c.Locals("user_id").(uuid.UUID)
	addressID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid address ID",
		})
	}

	if err := userService.DeleteAddress(userID, addressID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Success: true,
		Message: "address deleted successfully",
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(dto.ErrorResponse{
		Error: err.Error(),
	})
}

func joinStrings(strings []string, sep string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += sep + strings[i]
	}
	return result
}
