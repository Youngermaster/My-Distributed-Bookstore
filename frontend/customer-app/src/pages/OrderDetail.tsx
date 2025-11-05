import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams, useNavigate, Link } from "@tanstack/react-router";
import { orderAPI, booksAPI } from "@/lib/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { ArrowLeft } from "lucide-react";
import { toast } from "@/lib/toast";
import type { OrderStatus } from "@/types/order";

export default function OrderDetail() {
  const { orderId } = useParams({ from: "/orders/$orderId" });
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: order, isLoading } = useQuery({
    queryKey: ["order", orderId],
    queryFn: () => orderAPI.get(orderId).then((res) => res.data),
  });

  const cancelOrderMutation = useMutation({
    mutationFn: () => orderAPI.cancel(orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["order", orderId] });
      queryClient.invalidateQueries({ queryKey: ["orders"] });
      toast.success("Order cancelled successfully");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Failed to cancel order");
    },
  });

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-muted-foreground">Loading order...</p>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Card>
          <CardContent className="pt-6 text-center">
            <p className="text-muted-foreground mb-4">Order not found</p>
            <Button asChild>
              <Link to="/orders">Back to Orders</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const canCancel = order.status === "pending" || order.status === "confirmed";

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Button
          variant="ghost"
          onClick={() => navigate({ to: "/orders" })}
          className="mb-6"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Orders
        </Button>

        {/* Order Header */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex justify-between items-start">
              <div>
                <CardTitle className="text-2xl">
                  Order #{order.id.substring(0, 8)}
                </CardTitle>
                <CardDescription>
                  Placed on{" "}
                  {new Date(order.created_at).toLocaleDateString("en-US", {
                    year: "numeric",
                    month: "long",
                    day: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </CardDescription>
              </div>
              <OrderStatusBadge status={order.status} />
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {order.shipping_address && (
                <div>
                  <h4 className="font-semibold text-sm mb-1">
                    Shipping Address
                  </h4>
                  <p className="text-sm text-muted-foreground">
                    {order.shipping_address}
                  </p>
                </div>
              )}
              {order.payment_method && (
                <div>
                  <h4 className="font-semibold text-sm mb-1">Payment Method</h4>
                  <p className="text-sm text-muted-foreground capitalize">
                    {order.payment_method.replace("_", " ")}
                  </p>
                </div>
              )}
              <div>
                <h4 className="font-semibold text-sm mb-1">Total Amount</h4>
                <p className="text-2xl font-bold text-primary">
                  ${order.total_amount.toFixed(2)}
                </p>
              </div>
            </div>

            {canCancel && (
              <>
                <Separator className="my-4" />
                <Button
                  variant="destructive"
                  onClick={() => cancelOrderMutation.mutate()}
                  disabled={cancelOrderMutation.isPending}
                >
                  {cancelOrderMutation.isPending
                    ? "Cancelling..."
                    : "Cancel Order"}
                </Button>
              </>
            )}
          </CardContent>
        </Card>

        {/* Order Items */}
        <Card>
          <CardHeader>
            <CardTitle>Order Items ({order.items.length})</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {order.items.map((item, index) => (
              <div key={item.id}>
                {index > 0 && <Separator className="my-4" />}
                <OrderItemCard item={item} />
              </div>
            ))}
          </CardContent>
        </Card>

        {order.notes && (
          <Card className="mt-6">
            <CardHeader>
              <CardTitle>Order Notes</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">{order.notes}</p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}

function OrderStatusBadge({ status }: { status: OrderStatus }) {
  const statusConfig: Record<
    OrderStatus,
    {
      label: string;
      variant: "default" | "secondary" | "destructive" | "outline";
    }
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

function OrderItemCard({
  item,
}: {
  item: {
    id: string;
    book_id: string;
    quantity: number;
    unit_price: number;
    subtotal: number;
  };
}) {
  const { data: book } = useQuery({
    queryKey: ["book", item.book_id],
    queryFn: () => booksAPI.get(item.book_id).then((res) => res.data),
  });

  return (
    <div className="flex gap-4">
      {book?.cover_image_url && (
        <img
          src={book.cover_image_url}
          alt={book.title}
          className="w-20 h-28 object-cover rounded"
        />
      )}
      <div className="flex-1">
        <h4 className="font-semibold mb-1">{book?.title || "Loading..."}</h4>
        {book?.authors && book.authors.length > 0 && (
          <p className="text-sm text-muted-foreground mb-2">
            by {book.authors.map((a) => a.name).join(", ")}
          </p>
        )}
        <div className="flex justify-between items-end">
          <div>
            <p className="text-sm text-muted-foreground">
              Quantity: {item.quantity}
            </p>
            <p className="text-sm text-muted-foreground">
              Unit Price: ${item.unit_price.toFixed(2)}
            </p>
          </div>
          <p className="text-lg font-bold">${item.subtotal.toFixed(2)}</p>
        </div>
      </div>
    </div>
  );
}
