import { Link } from "@tanstack/react-router";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function Checkout() {
  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <Card className="max-w-lg w-full mx-4">
        <CardHeader>
          <CardTitle>Checkout</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4 text-muted-foreground">
          <p>
            Checkout flows are coming soon. In the meantime, you can review your
            cart or track recent purchases from the orders section.
          </p>
          <div className="flex flex-wrap gap-2">
            <Button asChild variant="outline">
              <Link to="/cart">Back to Cart</Link>
            </Button>
            <Button asChild>
              <Link to="/orders">View Orders</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
