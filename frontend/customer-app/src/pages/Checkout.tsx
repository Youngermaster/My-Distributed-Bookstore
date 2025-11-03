import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  useAuthService,
  useCartService,
  useOrderService,
  useRecommendationService,
} from "@/services";
import { booksAPI, recommendationsAPI } from "@/lib/api";
import { useCartStore } from "@/store/cartStore";
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
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { toast } from "@/lib/toast";
import { ShoppingCart, CreditCard, MapPin, Package } from "lucide-react";

export default function CheckoutComplete() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user, isAuthenticated } = useAuthService();
  const { cartId } = useCartStore();
  const { useCart } = useCartService();
  const { useCreateOrder } = useOrderService();
  const { useTrackInteraction } = useRecommendationService();

  const { data: cart, isLoading: cartLoading } = useCart();
  const createOrderMutation = useCreateOrder();
  const trackInteractionMutation = useTrackInteraction();

  const [step, setStep] = useState<"shipping" | "review" | "processing">(
    "shipping"
  );
  const [shippingInfo, setShippingInfo] = useState({
    full_name: user?.full_name || "",
    phone: "",
    address_line1: "",
    address_line2: "",
    city: "",
    state: "",
    postal_code: "",
    country: "",
  });

  // Redirect if not authenticated
  if (!isAuthenticated) {
    navigate({ to: "/login" });
    return null;
  }

  // Redirect if cart is empty
  if (!cart || cart.items.length === 0) {
    navigate({ to: "/cart" });
    return null;
  }

  const handleShippingSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setStep("review");
  };

  const handlePlaceOrder = async () => {
    setStep("processing");

    try {
      // Create order
      const orderData = {
        cart_id: cartId!,
        shipping_address: shippingInfo,
        payment_method: "card", // Placeholder since we don't have payment service
      };

      const order = await createOrderMutation.mutateAsync(orderData);

      // Track purchase interactions for all items
      if (isAuthenticated && cart.items) {
        for (const item of cart.items) {
          try {
            await trackInteractionMutation.mutateAsync({
              book_id: item.book_id,
              interaction_type: "purchase",
            });
          } catch (error) {
            console.error("Failed to track purchase interaction:", error);
          }
        }
      }

      // Clear cart
      queryClient.invalidateQueries({ queryKey: ["cart", cartId] });

      // Navigate to order confirmation
      navigate({ to: "/orders/$orderId", params: { orderId: order.id } });
      toast.success("Order placed successfully!");
    } catch (error: any) {
      console.error("Order creation failed:", error);
      toast.error(error.response?.data?.message || "Failed to create order");
      setStep("review");
    }
  };

  if (cartLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <p className="text-muted-foreground">Loading checkout...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">Checkout</h1>
          <p className="text-muted-foreground">Complete your purchase</p>
        </div>

        {/* Progress Steps */}
        <div className="flex items-center justify-center mb-8">
          <div className="flex items-center">
            <StepIndicator
              icon={<MapPin className="w-5 h-5" />}
              label="Shipping"
              active={step === "shipping"}
              completed={step === "review" || step === "processing"}
            />
            <div className="w-24 h-1 bg-border mx-2" />
            <StepIndicator
              icon={<Package className="w-5 h-5" />}
              label="Review"
              active={step === "review"}
              completed={step === "processing"}
            />
            <div className="w-24 h-1 bg-border mx-2" />
            <StepIndicator
              icon={<CreditCard className="w-5 h-5" />}
              label="Payment"
              active={step === "processing"}
              completed={false}
            />
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Content */}
          <div className="lg:col-span-2">
            {step === "shipping" && (
              <Card>
                <CardHeader>
                  <CardTitle>Shipping Information</CardTitle>
                  <CardDescription>Enter your delivery address</CardDescription>
                </CardHeader>
                <form onSubmit={handleShippingSubmit}>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="full_name">Full Name</Label>
                        <Input
                          id="full_name"
                          value={shippingInfo.full_name}
                          onChange={(e) =>
                            setShippingInfo({
                              ...shippingInfo,
                              full_name: e.target.value,
                            })
                          }
                          required
                        />
                      </div>
                      <div>
                        <Label htmlFor="phone">Phone Number</Label>
                        <Input
                          id="phone"
                          type="tel"
                          value={shippingInfo.phone}
                          onChange={(e) =>
                            setShippingInfo({
                              ...shippingInfo,
                              phone: e.target.value,
                            })
                          }
                          required
                        />
                      </div>
                    </div>

                    <div>
                      <Label htmlFor="address_line1">Address Line 1</Label>
                      <Input
                        id="address_line1"
                        value={shippingInfo.address_line1}
                        onChange={(e) =>
                          setShippingInfo({
                            ...shippingInfo,
                            address_line1: e.target.value,
                          })
                        }
                        required
                      />
                    </div>

                    <div>
                      <Label htmlFor="address_line2">
                        Address Line 2 (Optional)
                      </Label>
                      <Input
                        id="address_line2"
                        value={shippingInfo.address_line2}
                        onChange={(e) =>
                          setShippingInfo({
                            ...shippingInfo,
                            address_line2: e.target.value,
                          })
                        }
                      />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div>
                        <Label htmlFor="city">City</Label>
                        <Input
                          id="city"
                          value={shippingInfo.city}
                          onChange={(e) =>
                            setShippingInfo({
                              ...shippingInfo,
                              city: e.target.value,
                            })
                          }
                          required
                        />
                      </div>
                      <div>
                        <Label htmlFor="state">State/Province</Label>
                        <Input
                          id="state"
                          value={shippingInfo.state}
                          onChange={(e) =>
                            setShippingInfo({
                              ...shippingInfo,
                              state: e.target.value,
                            })
                          }
                          required
                        />
                      </div>
                      <div>
                        <Label htmlFor="postal_code">Postal Code</Label>
                        <Input
                          id="postal_code"
                          value={shippingInfo.postal_code}
                          onChange={(e) =>
                            setShippingInfo({
                              ...shippingInfo,
                              postal_code: e.target.value,
                            })
                          }
                          required
                        />
                      </div>
                    </div>

                    <div>
                      <Label htmlFor="country">Country</Label>
                      <Input
                        id="country"
                        value={shippingInfo.country}
                        onChange={(e) =>
                          setShippingInfo({
                            ...shippingInfo,
                            country: e.target.value,
                          })
                        }
                        required
                      />
                    </div>
                  </CardContent>
                  <CardFooter className="flex justify-between">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => navigate({ to: "/cart" })}
                    >
                      Back to Cart
                    </Button>
                    <Button type="submit">Continue to Review</Button>
                  </CardFooter>
                </form>
              </Card>
            )}

            {step === "review" && (
              <div className="space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Review Order</CardTitle>
                    <CardDescription>
                      Confirm your order details
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <h3 className="font-semibold mb-2">Shipping Address</h3>
                      <p className="text-sm text-muted-foreground">
                        {shippingInfo.full_name}
                        <br />
                        {shippingInfo.address_line1}
                        {shippingInfo.address_line2 && (
                          <>
                            <br />
                            {shippingInfo.address_line2}
                          </>
                        )}
                        <br />
                        {shippingInfo.city}, {shippingInfo.state}{" "}
                        {shippingInfo.postal_code}
                        <br />
                        {shippingInfo.country}
                        <br />
                        Phone: {shippingInfo.phone}
                      </p>
                      <Button
                        variant="link"
                        className="px-0"
                        onClick={() => setStep("shipping")}
                      >
                        Edit Address
                      </Button>
                    </div>

                    <Separator />

                    <div>
                      <h3 className="font-semibold mb-4">Order Items</h3>
                      <div className="space-y-3">
                        {cart.items.map((item) => (
                          <OrderItemRow key={item.book_id} item={item} />
                        ))}
                      </div>
                    </div>
                  </CardContent>
                  <CardFooter className="flex justify-between">
                    <Button
                      variant="outline"
                      onClick={() => setStep("shipping")}
                    >
                      Back
                    </Button>
                    <Button onClick={handlePlaceOrder}>Place Order</Button>
                  </CardFooter>
                </Card>
              </div>
            )}

            {step === "processing" && (
              <Card>
                <CardContent className="pt-6 text-center">
                  <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-primary mx-auto mb-4" />
                  <p className="text-lg font-semibold">
                    Processing your order...
                  </p>
                  <p className="text-sm text-muted-foreground mt-2">
                    Please wait while we create your order
                  </p>
                </CardContent>
              </Card>
            )}
          </div>

          {/* Order Summary Sidebar */}
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
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Tax</span>
                  <span className="font-medium">$0.00</span>
                </div>
                <Separator />
                <div className="flex justify-between">
                  <span className="text-lg font-semibold">Total</span>
                  <span className="text-lg font-bold">
                    ${cart.total.toFixed(2)}
                  </span>
                </div>
                <div className="text-xs text-muted-foreground">
                  <p>
                    {cart.item_count} {cart.item_count === 1 ? "item" : "items"}{" "}
                    in your order
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}

interface StepIndicatorProps {
  icon: React.ReactNode;
  label: string;
  active: boolean;
  completed: boolean;
}

function StepIndicator({ icon, label, active, completed }: StepIndicatorProps) {
  return (
    <div className="flex flex-col items-center">
      <div
        className={`w-12 h-12 rounded-full flex items-center justify-center ${
          completed
            ? "bg-primary text-primary-foreground"
            : active
            ? "bg-primary text-primary-foreground"
            : "bg-muted text-muted-foreground"
        }`}
      >
        {icon}
      </div>
      <span className="text-xs mt-2 font-medium">{label}</span>
    </div>
  );
}

interface OrderItemRowProps {
  item: {
    book_id: string;
    quantity: number;
    price: number;
    subtotal: number;
  };
}

function OrderItemRow({ item }: OrderItemRowProps) {
  const { data: book } = useQuery({
    queryKey: ["book", item.book_id],
    queryFn: () => booksAPI.get(item.book_id).then((res) => res.data),
  });

  return (
    <div className="flex gap-3 py-2 border-b last:border-b-0">
      {book?.cover_image_url && (
        <img
          src={book.cover_image_url}
          alt={book.title}
          className="w-16 h-20 object-cover rounded"
        />
      )}
      <div className="flex-1">
        <p className="font-semibold text-sm">{book?.title || "Loading..."}</p>
        <p className="text-xs text-muted-foreground">
          Quantity: {item.quantity} × ${item.price.toFixed(2)}
        </p>
      </div>
      <div className="text-right">
        <p className="font-semibold">${item.subtotal.toFixed(2)}</p>
      </div>
    </div>
  );
}
