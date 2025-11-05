// Admin Service Types
// Based on the backend admin service implementation

export interface DashboardStats {
  total_orders: number;
  total_orders_today: number;
  total_orders_week: number;
  total_orders_month: number;
  total_revenue: number;
  revenue_today: number;
  revenue_week: number;
  revenue_month: number;
  total_users: number;
  new_users_today: number;
  new_users_week: number;
  new_users_month: number;
  total_books: number;
  low_stock_count: number;
  out_of_stock_count: number;
  average_order_value: number;
  top_selling_books: TopBook[];
}

export interface TopBook {
  book_id: string;
  title: string;
  total_sold: number;
  revenue: number;
}

export interface SalesAnalytics {
  period: string;
  start_date: string;
  end_date: string;
  total_orders: number;
  total_revenue: number;
  average_order_value: number;
  data_points: SalesDataPoint[];
  revenue_by_category?: RevenueByCategory[];
}

export interface SalesDataPoint {
  date: string;
  orders: number;
  revenue: number;
}

export interface RevenueByCategory {
  category: string;
  revenue: number;
  orders: number;
}

export interface InventoryReport {
  total_items: number;
  total_value: number;
  low_stock_items: number;
  out_of_stock_items: number;
  items: InventoryItem[];
}

export interface InventoryItem {
  book_id: string;
  title: string;
  sku?: string;
  available_quantity: number;
  reserved_quantity: number;
  reorder_level: number;
  unit_price: number;
  total_value: number;
  status: "in_stock" | "low_stock" | "out_of_stock";
}

export interface UserGrowthReport {
  period: string;
  start_date: string;
  end_date: string;
  total_new_users: number;
  data_points: UserGrowthPoint[];
}

export interface UserGrowthPoint {
  date: string;
  new_users: number;
  cumulative_users: number;
}

export interface OrderStatusBreakdown {
  status: string;
  count: number;
  percentage: number;
}

// Request params
export interface SalesAnalyticsParams {
  period?: "day" | "week" | "month" | "year";
  start_date?: string;
  end_date?: string;
}

export interface InventoryReportParams {
  low_stock_only?: boolean;
  out_of_stock_only?: boolean;
}

export interface UserGrowthParams {
  period?: "day" | "week" | "month" | "year";
  start_date?: string;
  end_date?: string;
}

export interface TopBooksParams {
  limit?: number;
  period?: "day" | "week" | "month" | "all";
}
