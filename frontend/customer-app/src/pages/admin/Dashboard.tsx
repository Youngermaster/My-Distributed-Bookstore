import { ProtectedRoute } from "@/components/ProtectedRoute";
import { useAdminService } from "@/services";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DollarSign,
  ShoppingCart,
  Users,
  Package,
  TrendingUp,
  AlertTriangle,
} from "lucide-react";

export default function AdminDashboard() {
  return (
    <ProtectedRoute requireRole="admin">
      <DashboardContent />
    </ProtectedRoute>
  );
}

function DashboardContent() {
  const { useDashboard } = useAdminService();
  const { data: stats, isLoading, error } = useDashboard();

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-muted-foreground">Loading dashboard...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-destructive">Error loading dashboard</p>
      </div>
    );
  }

  if (!stats) {
    return null;
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Admin Dashboard
          </h1>
          <p className="text-muted-foreground">
            Monitor your bookstore's performance and key metrics
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {/* Total Revenue */}
          <StatCard
            title="Total Revenue"
            value={`$${stats.total_revenue.toFixed(2)}`}
            subtitle={`$${stats.revenue_today.toFixed(2)} today`}
            icon={<DollarSign className="w-5 h-5" />}
            trend={`+$${stats.revenue_week.toFixed(2)} this week`}
          />

          {/* Total Orders */}
          <StatCard
            title="Total Orders"
            value={stats.total_orders.toString()}
            subtitle={`${stats.total_orders_today} today`}
            icon={<ShoppingCart className="w-5 h-5" />}
            trend={`${stats.total_orders_week} this week`}
          />

          {/* Total Users */}
          <StatCard
            title="Total Users"
            value={stats.total_users.toString()}
            subtitle={`${stats.new_users_today} new today`}
            icon={<Users className="w-5 h-5" />}
            trend={`+${stats.new_users_week} this week`}
          />

          {/* Total Books */}
          <StatCard
            title="Total Books"
            value={stats.total_books.toString()}
            subtitle={`${stats.out_of_stock_count} out of stock`}
            icon={<Package className="w-5 h-5" />}
            trend={`${stats.low_stock_count} low stock`}
            alert={stats.low_stock_count > 0}
          />
        </div>

        {/* Additional Metrics */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Average Order Value */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <TrendingUp className="w-5 h-5" />
                Average Order Value
              </CardTitle>
              <CardDescription>Per order metrics</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold text-primary">
                ${stats.average_order_value.toFixed(2)}
              </p>
            </CardContent>
          </Card>

          {/* Inventory Alerts */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <AlertTriangle className="w-5 h-5" />
                Inventory Alerts
              </CardTitle>
              <CardDescription>Stock status overview</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Low Stock:</span>
                  <span className="font-semibold text-warning">
                    {stats.low_stock_count} items
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Out of Stock:</span>
                  <span className="font-semibold text-destructive">
                    {stats.out_of_stock_count} items
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Top Selling Books */}
        <Card>
          <CardHeader>
            <CardTitle>Top Selling Books</CardTitle>
            <CardDescription>Best performing books by revenue</CardDescription>
          </CardHeader>
          <CardContent>
            {stats.top_selling_books && stats.top_selling_books.length > 0 ? (
              <div className="space-y-4">
                {stats.top_selling_books.map((book, index) => (
                  <div
                    key={book.book_id}
                    className="flex items-center justify-between py-2 border-b last:border-b-0"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center font-bold text-primary">
                        {index + 1}
                      </div>
                      <div>
                        <p className="font-semibold">{book.title}</p>
                        <p className="text-sm text-muted-foreground">
                          {book.total_sold} sold
                        </p>
                      </div>
                    </div>
                    <p className="font-bold text-primary">
                      ${book.revenue.toFixed(2)}
                    </p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-muted-foreground text-center py-4">
                No sales data available
              </p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

interface StatCardProps {
  title: string;
  value: string;
  subtitle: string;
  icon: React.ReactNode;
  trend: string;
  alert?: boolean;
}

function StatCard({
  title,
  value,
  subtitle,
  icon,
  trend,
  alert,
}: StatCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        <div className={alert ? "text-destructive" : "text-muted-foreground"}>
          {icon}
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-3xl font-bold mb-1">{value}</div>
        <p className="text-sm text-muted-foreground mb-2">{subtitle}</p>
        <p
          className={`text-xs ${
            alert ? "text-warning" : "text-muted-foreground"
          }`}
        >
          {trend}
        </p>
      </CardContent>
    </Card>
  );
}
