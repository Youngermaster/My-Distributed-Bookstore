import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { orderAPI } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Package, ChevronRight } from "lucide-react";
import type { OrderStatus } from "@/types/order";

export default function Orders() {
  const { user } = useAuthStore();
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { data: ordersData, isLoading } = useQuery({
    queryKey: ["orders", user?.id, page],
    queryFn: () =>
      user
        ? orderAPI.getUserOrders(user.id, page, pageSize).then((res) => res.data)
        : Promise.resolve(null),
    enabled: !!user,
  });

  if (!user) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card>
          <CardContent className="pt-6 text-center">
            <p className="text-muted-foreground mb-4">
              Please login to view your orders
            </p>
            <Button asChild>
              <Link to="/login">Login</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-muted-foreground">Loading orders...</p>
      </div>
    );
  }

  const orders = ordersData?.orders || [];
  const isEmpty = orders.length === 0;

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">My Orders</h1>
          <p className="text-muted-foreground">
            {isEmpty
              ? "You haven't placed any orders yet"
              : `${ordersData?.total || 0} total ${ordersData?.total === 1 ? "order" : "orders"}`}
          </p>
        </div>

        {isEmpty ? (
          <Card>
            <CardContent className="pt-6 text-center">
              <Package className="w-16 h-16 mx-auto mb-4 text-muted-foreground" />
              <p className="text-lg text-muted-foreground mb-4">
                No orders found
              </p>
              <Button asChild>
                <Link to="/books">Start Shopping</Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <>
            <div className="space-y-4">
              {orders.map((order) => (
                <Card key={order.id}>
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle className="text-lg">
                          Order #{order.id.substring(0, 8)}
                        </CardTitle>
                        <CardDescription>
                          Placed on{" "}
                          {new Date(order.created_at).toLocaleDateString("en-US", {
                            year: "numeric",
                            month: "long",
                            day: "numeric",
                          })}
                        </CardDescription>
                      </div>
                      <OrderStatusBadge status={order.status} />
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="flex justify-between items-center">
                      <div className="space-y-1">
                        <p className="text-sm text-muted-foreground">
                          {order.item_count} {order.item_count === 1 ? "item" : "items"}
                        </p>
                        <p className="text-lg font-bold">
                          ${order.total_amount.toFixed(2)}
                        </p>
                      </div>
                      <Button variant="outline" asChild>
                        <Link to={`/orders/${order.id}`}>
                          View Details
                          <ChevronRight className="w-4 h-4 ml-2" />
                        </Link>
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>

            {/* Pagination */}
            {ordersData && ordersData.total_pages > 1 && (
              <div className="flex justify-center items-center gap-4 mt-8">
                <Button
                  variant="outline"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                >
                  Previous
                </Button>
                <span className="text-sm text-muted-foreground">
                  Page {page} of {ordersData.total_pages}
                </span>
                <Button
                  variant="outline"
                  onClick={() => setPage((p) => p + 1)}
                  disabled={page >= ordersData.total_pages}
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

function OrderStatusBadge({ status }: { status: OrderStatus }) {
  const statusConfig: Record<
    OrderStatus,
    { label: string; variant: "default" | "secondary" | "destructive" | "outline" }
  > = {
    pending: { label: "Pending", variant: "outline" },
    confirmed: { label: "Confirmed", variant: "secondary" },
    processing: { label: "Processing", variant: "default" },
    shipped: { label: "Shipped", variant: "default" },
    delivered: { label: "Delivered", variant: "secondary" },
    cancelled: { label: "Cancelled", variant: "destructive" },
  };

  const config = statusConfig[status];
  return <Badge variant={config.variant}>{config.label}</Badge>;
}
