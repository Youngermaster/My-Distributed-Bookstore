package domain

import "time"

// DashboardStats represents the main dashboard statistics
type DashboardStats struct {
	TotalOrders        int64   `json:"total_orders"`
	TotalOrdersToday   int64   `json:"total_orders_today"`
	TotalOrdersWeek    int64   `json:"total_orders_week"`
	TotalOrdersMonth   int64   `json:"total_orders_month"`
	TotalRevenue       float64 `json:"total_revenue"`
	RevenueToday       float64 `json:"revenue_today"`
	RevenueWeek        float64 `json:"revenue_week"`
	RevenueMonth       float64 `json:"revenue_month"`
	TotalUsers         int64   `json:"total_users"`
	NewUsersToday      int64   `json:"new_users_today"`
	NewUsersWeek       int64   `json:"new_users_week"`
	NewUsersMonth      int64   `json:"new_users_month"`
	TotalBooks         int64   `json:"total_books"`
	LowStockCount      int64   `json:"low_stock_count"`
	AverageOrderValue  float64 `json:"average_order_value"`
	TopSellingBooks    []TopBook `json:"top_selling_books"`
	RecentOrders       int64   `json:"recent_orders_24h"`
}

// TopBook represents a top-selling book
type TopBook struct {
	BookID      string  `json:"book_id"`
	BookTitle   string  `json:"book_title"`
	TotalSales  int64   `json:"total_sales"`
	TotalRevenue float64 `json:"total_revenue"`
}

// SalesAnalytics represents sales data over time
type SalesAnalytics struct {
	Period      string        `json:"period"` // daily, weekly, monthly
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date"`
	TotalOrders int64         `json:"total_orders"`
	TotalRevenue float64      `json:"total_revenue"`
	DataPoints  []SalesDataPoint `json:"data_points"`
}

// SalesDataPoint represents a single data point in sales analytics
type SalesDataPoint struct {
	Date        time.Time `json:"date"`
	Orders      int64     `json:"orders"`
	Revenue     float64   `json:"revenue"`
	Label       string    `json:"label"` // e.g., "Jan 1", "Week 1"
}

// InventoryReport represents inventory status report
type InventoryReport struct {
	TotalBooks      int64              `json:"total_books"`
	LowStockBooks   int64              `json:"low_stock_books"`
	OutOfStockBooks int64              `json:"out_of_stock_books"`
	TotalStockValue float64            `json:"total_stock_value"`
	Items           []InventoryItem    `json:"items"`
}

// InventoryItem represents a single inventory item
type InventoryItem struct {
	BookID          string  `json:"book_id"`
	BookTitle       string  `json:"book_title"`
	ISBN            string  `json:"isbn"`
	StockQuantity   int32   `json:"stock_quantity"`
	Price           float64 `json:"price"`
	StockValue      float64 `json:"stock_value"`
	Status          string  `json:"status"` // "in_stock", "low_stock", "out_of_stock"
}

// UserGrowthReport represents user growth analytics
type UserGrowthReport struct {
	Period     string             `json:"period"`
	TotalUsers int64              `json:"total_users"`
	DataPoints []UserGrowthPoint  `json:"data_points"`
}

// UserGrowthPoint represents a point in user growth
type UserGrowthPoint struct {
	Date      time.Time `json:"date"`
	NewUsers  int64     `json:"new_users"`
	TotalUsers int64    `json:"total_users"`
	Label     string    `json:"label"`
}

// OrderStatusBreakdown represents orders by status
type OrderStatusBreakdown struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// RevenueByCategory represents revenue breakdown by book category
type RevenueByCategory struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
}
