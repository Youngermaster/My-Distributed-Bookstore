package server

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/config"
	"github.com/youngermaster/distributed-bookstore/notification-service/internal/notification"
)

type Server struct {
	app     *fiber.App
	cfg     *config.Config
	history *notification.History
	logger  zerolog.Logger
}

func New(cfg *config.Config, history *notification.History, logger zerolog.Logger) *Server {
	app := fiber.New(fiber.Config{
		AppName: "notification-service",
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	s := &Server{
		app:     app,
		cfg:     cfg,
		history: history,
		logger:  logger,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	s.app.Get("/health", s.healthHandler)
	s.app.Get("/ready", s.readyHandler)

	api := s.app.Group("/api/v1")
	api.Get("/notifications/recent", s.recentNotifications)
}

func (s *Server) healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": s.cfg.ServiceName,
	})
}

func (s *Server) readyHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ready",
		"service": s.cfg.ServiceName,
	})
}

func (s *Server) recentNotifications(c *fiber.Ctx) error {
	limitParam := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 20
	}

	result := s.history.Recent(limit)
	return c.JSON(fiber.Map{
		"notifications": result,
		"count":         len(result),
	})
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
