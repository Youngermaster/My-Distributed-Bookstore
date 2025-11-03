package handler

import (
	"time"

	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/service"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboard returns dashboard statistics
func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	stats, err := h.analyticsService.GetDashboardStats()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get dashboard stats")
	}

	return response.Success(c, fiber.StatusOK, "Dashboard stats retrieved successfully", stats)
}

// GetSalesAnalytics returns sales analytics
func (h *AnalyticsHandler) GetSalesAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "daily")
	
	// Parse date range
	startDateStr := c.Query("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.Query("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid start_date format")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid end_date format")
	}

	analytics, err := h.analyticsService.GetSalesAnalytics(period, startDate, endDate)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get sales analytics")
	}

	return response.Success(c, fiber.StatusOK, "Sales analytics retrieved successfully", analytics)
}

// GetInventoryReport returns inventory report
func (h *AnalyticsHandler) GetInventoryReport(c *fiber.Ctx) error {
	lowStockThreshold := c.QueryInt("low_stock_threshold", 10)

	report, err := h.analyticsService.GetInventoryReport(lowStockThreshold)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get inventory report")
	}

	return response.Success(c, fiber.StatusOK, "Inventory report retrieved successfully", report)
}

// GetUserGrowth returns user growth analytics
func (h *AnalyticsHandler) GetUserGrowth(c *fiber.Ctx) error {
	period := c.Query("period", "monthly")

	report, err := h.analyticsService.GetUserGrowthReport(period)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get user growth report")
	}

	return response.Success(c, fiber.StatusOK, "User growth report retrieved successfully", report)
}

// GetTopSellingBooks returns top selling books
func (h *AnalyticsHandler) GetTopSellingBooks(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)

	books, err := h.analyticsService.GetTopSellingBooks(limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get top selling books")
	}

	return response.Success(c, fiber.StatusOK, "Top selling books retrieved successfully", books)
}
