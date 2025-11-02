import { useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { v4 as uuidv4 } from "uuid";
import { cartAPI, booksAPI } from "@/lib/api";
import { useCartStore } from "@/store/cartStore";
import { useAuthStore } from "@/store/authStore";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Trash2, Plus, Minus, ShoppingCart } from "lucide-react";
import { toast } from "sonner";

export default function Cart() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { cartId, setCartId } = useCartStore();
  const { user } = useAuthStore();

  // Initialize cart ID if not exists
  useEffect(() => {
    if (!cartId) {
      const newCartId = uuidv4();
      setCartId(newCartId);
    }
  }, [cartId, setCartId]);

  // Fetch cart
  const { data: cart, isLoading } = useQuery({
    queryKey: ["cart", cartId],
    queryFn: () => cartAPI.get(cartId!).then((res) => res.data),
    enabled: !!cartId,
  });

  // Remove item mutation
  const removeItemMutation = useMutation({
    mutationFn: ({ bookId }: { bookId: string }) =>
      cartAPI.removeItem(cartId!, bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
      toast.success("Item removed from cart");
    },
    onError: () => {
      toast.error("Failed to remove item");
    },
  });

  // Update quantity mutation
  const updateQuantityMutation = useMutation({
    mutationFn: ({ bookId, quantity }: { bookId: string; quantity: number }) =>
      cartAPI.updateItem(cartId!, bookId, { quantity }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
    },
    onError: () => {
      toast.error("Failed to update quantity");
    },
  });

  // Clear cart mutation
  const clearCartMutation = useMutation({
    mutationFn: () => cartAPI.clear(cartId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
      toast.success("Cart cleared");
    },
    onError: () => {
      toast.error("Failed to clear cart");
    },
  });

  const handleQuantityChange = (bookId: string, currentQty: number, delta: number) => {
    const newQty = currentQty + delta;
    if (newQty < 1) {
      removeItemMutation.mutate({ bookId });
    } else if (newQty <= 99) {
      updateQuantityMutation.mutate({ bookId, quantity: newQty });
    }
  };

  const handleCheckout = () => {
    if (!user) {
      toast.error("Please login to checkout");
      navigate({ to: "/login" });
      return;
    }
    navigate({ to: "/checkout" });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-muted-foreground">Loading cart...</p>
      </div>
    );
  }

  const isEmpty = !cart || cart.items.length === 0;

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Shopping Cart
          </h1>
          <p className="text-muted-foreground">
            {isEmpty
              ? "Your cart is empty"
              : `${cart.item_count} ${cart.item_count === 1 ? "item" : "items"} in your cart`}
          </p>
        </div>

        {isEmpty ? (
          <Card>
            <CardContent className="pt-6 text-center">
              <ShoppingCart className="w-16 h-16 mx-auto mb-4 text-muted-foreground" />
              <p className="text-lg text-muted-foreground mb-4">
                Your cart is empty
              </p>
              <Button onClick={() => navigate({ to: "/books" })}>
                Continue Shopping
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Cart Items */}
            <div className="lg:col-span-2 space-y-4">
              {cart.items.map((item) => (
                <CartItemCard
                  key={item.book_id}
                  item={item}
                  onQuantityChange={(delta) =>
                    handleQuantityChange(item.book_id, item.quantity, delta)
                  }
                  onRemove={() => removeItemMutation.mutate({ bookId: item.book_id })}
                />
              ))}

              <div className="flex justify-between items-center pt-4">
                <Button
                  variant="outline"
                  onClick={() => navigate({ to: "/books" })}
                >
                  Continue Shopping
                </Button>
                <Button
                  variant="destructive"
                  onClick={() => clearCartMutation.mutate()}
                >
                  Clear Cart
                </Button>
              </div>
            </div>

            {/* Order Summary */}
            <div className="lg:col-span-1">
              <Card className="sticky top-4">
                <CardHeader>
                  <CardTitle>Order Summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Subtotal</span>
                    <span className="font-medium">${cart.total.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Shipping</span>
                    <span className="font-medium">FREE</span>
                  </div>
                  <div className="border-t pt-4">
                    <div className="flex justify-between">
                      <span className="text-lg font-semibold">Total</span>
                      <span className="text-lg font-bold">
                        ${cart.total.toFixed(2)}
                      </span>
                    </div>
                  </div>
                </CardContent>
                <CardFooter>
                  <Button className="w-full" onClick={handleCheckout}>
                    Proceed to Checkout
                  </Button>
                </CardFooter>
              </Card>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

interface CartItemCardProps {
  item: {
    book_id: string;
    quantity: number;
    price: number;
    subtotal: number;
  };
  onQuantityChange: (delta: number) => void;
  onRemove: () => void;
}

function CartItemCard({ item, onQuantityChange, onRemove }: CartItemCardProps) {
  const { data: book } = useQuery({
    queryKey: ["book", item.book_id],
    queryFn: () => booksAPI.get(item.book_id).then((res) => res.data),
  });

  return (
    <Card>
      <CardContent className="pt-6">
        <div className="flex gap-4">
          {book?.cover_image_url && (
            <img
              src={book.cover_image_url}
              alt={book.title}
              className="w-24 h-32 object-cover rounded"
            />
          )}
          <div className="flex-1">
            <h3 className="font-semibold text-lg mb-1">
              {book?.title || "Loading..."}
            </h3>
            {book?.authors && book.authors.length > 0 && (
              <p className="text-sm text-muted-foreground mb-2">
                by {book.authors.map((a) => a.name).join(", ")}
              </p>
            )}
            <p className="text-lg font-bold text-primary">
              ${item.price.toFixed(2)}
            </p>
          </div>
          <div className="flex flex-col items-end justify-between">
            <Button
              variant="ghost"
              size="icon"
              onClick={onRemove}
              className="text-destructive"
            >
              <Trash2 className="w-4 h-4" />
            </Button>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="icon"
                onClick={() => onQuantityChange(-1)}
              >
                <Minus className="w-4 h-4" />
              </Button>
              <Input
                type="number"
                value={item.quantity}
                readOnly
                className="w-16 text-center"
              />
              <Button
                variant="outline"
                size="icon"
                onClick={() => onQuantityChange(1)}
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
            <p className="text-lg font-semibold">
              ${item.subtotal.toFixed(2)}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
