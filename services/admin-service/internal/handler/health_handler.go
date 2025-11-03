package handler

import (
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Health returns service health status
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Admin service is healthy", fiber.Map{
		"status":  "healthy",
		"service": "admin-service",
	})
}

// Ready returns service readiness status
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		return response.Error(c, fiber.StatusServiceUnavailable, "Database connection error")
	}

	if err := sqlDB.Ping(); err != nil {
		return response.Error(c, fiber.StatusServiceUnavailable, "Database not reachable")
	}

	return response.Success(c, fiber.StatusOK, "Admin service is ready", fiber.Map{
		"status":   "ready",
		"service":  "admin-service",
		"database": "connected",
	})
}
