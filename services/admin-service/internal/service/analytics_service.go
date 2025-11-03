package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/domain"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/grpc"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/internal/repository"
)

type AnalyticsService interface {
	GetDashboardStats() (*domain.DashboardStats, error)
	GetSalesAnalytics(period string, startDate, endDate time.Time) (*domain.SalesAnalytics, error)
	GetInventoryReport(lowStockThreshold int) (*domain.InventoryReport, error)
	GetUserGrowthReport(period string) (*domain.UserGrowthReport, error)
	GetTopSellingBooks(limit int) ([]domain.TopBook, error)
}

type analyticsService struct {
	repo    repository.AdminRepository
	clients *grpc.ServiceClients
}

func NewAnalyticsService(repo repository.AdminRepository, clients *grpc.ServiceClients) AnalyticsService {
	return &analyticsService{
		repo:    repo,
		clients: clients,
	}
}

func (s *analyticsService) GetDashboardStats() (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	// Get user statistics from user service
	userData, err := s.clients.UserService.Get("/api/v1/users/stats")
	if err == nil {
		var userStats struct {
			TotalUsers    int64 `json:"total_users"`
			NewToday      int64 `json:"new_today"`
			NewWeek       int64 `json:"new_week"`
			NewMonth      int64 `json:"new_month"`
		}
		if err := json.Unmarshal(userData, &userStats); err == nil {
			stats.TotalUsers = userStats.TotalUsers
			stats.NewUsersToday = userStats.NewToday
			stats.NewUsersWeek = userStats.NewWeek
			stats.NewUsersMonth = userStats.NewMonth
		}
	}

	// Get catalog statistics
	catalogData, err := s.clients.CatalogService.Get("/api/v1/books/stats")
	if err == nil {
		var catalogStats struct {
			TotalBooks int64 `json:"total_books"`
		}
		if err := json.Unmarshal(catalogData, &catalogStats); err == nil {
			stats.TotalBooks = catalogStats.TotalBooks
		}
	}

	// Get order statistics from order service
	orderData, err := s.clients.OrderService.Get("/api/v1/orders/stats")
	if err == nil {
		var orderStats struct {
			TotalOrders      int64   `json:"total_orders"`
			OrdersToday      int64   `json:"orders_today"`
			OrdersWeek       int64   `json:"orders_week"`
			OrdersMonth      int64   `json:"orders_month"`
			TotalRevenue     float64 `json:"total_revenue"`
			RevenueToday     float64 `json:"revenue_today"`
			RevenueWeek      float64 `json:"revenue_week"`
			RevenueMonth     float64 `json:"revenue_month"`
			AverageOrderValue float64 `json:"average_order_value"`
			RecentOrders     int64   `json:"recent_orders_24h"`
		}
		if err := json.Unmarshal(orderData, &orderStats); err == nil {
			stats.TotalOrders = orderStats.TotalOrders
			stats.TotalOrdersToday = orderStats.OrdersToday
			stats.TotalOrdersWeek = orderStats.OrdersWeek
			stats.TotalOrdersMonth = orderStats.OrdersMonth
			stats.TotalRevenue = orderStats.TotalRevenue
			stats.RevenueToday = orderStats.RevenueToday
			stats.RevenueWeek = orderStats.RevenueWeek
			stats.RevenueMonth = orderStats.RevenueMonth
			stats.AverageOrderValue = orderStats.AverageOrderValue
			stats.RecentOrders = orderStats.RecentOrders
		}
	}

	// Get inventory statistics
	inventoryData, err := s.clients.InventoryService.Get("/api/v1/inventory/stats")
	if err == nil {
		var inventoryStats struct {
			LowStockCount int64 `json:"low_stock_count"`
		}
		if err := json.Unmarshal(inventoryData, &inventoryStats); err == nil {
			stats.LowStockCount = inventoryStats.LowStockCount
		}
	}

	// Get top selling books
	topBooks, _ := s.GetTopSellingBooks(5)
	stats.TopSellingBooks = topBooks

	return stats, nil
}

func (s *analyticsService) GetSalesAnalytics(period string, startDate, endDate time.Time) (*domain.SalesAnalytics, error) {
	// Call order service for sales data
	path := fmt.Sprintf("/api/v1/orders/analytics?period=%s&start=%s&end=%s",
		period,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)

	data, err := s.clients.OrderService.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales analytics: %w", err)
	}

	var analytics domain.SalesAnalytics
	if err := json.Unmarshal(data, &analytics); err != nil {
		return nil, fmt.Errorf("failed to parse sales analytics: %w", err)
	}

	return &analytics, nil
}

func (s *analyticsService) GetInventoryReport(lowStockThreshold int) (*domain.InventoryReport, error) {
	// Get inventory data from inventory service
	path := fmt.Sprintf("/api/v1/inventory/report?low_stock_threshold=%d", lowStockThreshold)
	
	data, err := s.clients.InventoryService.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory report: %w", err)
	}

	var report domain.InventoryReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to parse inventory report: %w", err)
	}

	return &report, nil
}

func (s *analyticsService) GetUserGrowthReport(period string) (*domain.UserGrowthReport, error) {
	// Get user growth data from user service
	path := fmt.Sprintf("/api/v1/users/growth?period=%s", period)
	
	data, err := s.clients.UserService.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get user growth report: %w", err)
	}

	var report domain.UserGrowthReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to parse user growth report: %w", err)
	}

	return &report, nil
}

func (s *analyticsService) GetTopSellingBooks(limit int) ([]domain.TopBook, error) {
	// Get top selling books from order service
	path := fmt.Sprintf("/api/v1/orders/top-books?limit=%d", limit)
	
	data, err := s.clients.OrderService.Get(path)
	if err != nil {
		// Return empty slice if service is unavailable
		return []domain.TopBook{}, nil
	}

	var books []domain.TopBook
	if err := json.Unmarshal(data, &books); err != nil {
		return []domain.TopBook{}, nil
	}

	return books, nil
}
