import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { type Book } from "@/types/book";
import { ShoppingCart, Heart } from "lucide-react";

interface BookCardProps {
  book: Book;
  onClick?: () => void;
  onAddToCart?: () => void;
  onAddToWishlist?: () => void;
}

export default function BookCard({
  book,
  onClick,
  onAddToCart,
  onAddToWishlist,
}: BookCardProps) {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(price);
  };

  const authorsText =
    book.authors?.map((a) => a.name).join(", ") || "Unknown Author";

  return (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow cursor-pointer group">
      <div onClick={onClick}>
        <div className="aspect-[2/3] bg-gray-100 relative overflow-hidden">
          {book.cover_image_url ? (
            <img
              src={book.cover_image_url}
              alt={book.title}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-gray-400">
              <span className="text-4xl">📚</span>
            </div>
          )}
          {book.stock_quantity === 0 && (
            <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
              <span className="text-white font-semibold text-lg">
                Out of Stock
              </span>
            </div>
          )}
        </div>
        <CardContent className="p-4">
          <h3 className="font-semibold text-lg line-clamp-2 mb-1 group-hover:text-blue-600 transition-colors">
            {book.title}
          </h3>
          <p className="text-sm text-gray-600 mb-2 line-clamp-1">
            {authorsText}
          </p>
          <div className="flex items-center gap-2">
            <span className="text-lg font-bold text-gray-900">
              {formatPrice(book.price)}
            </span>
            {book.stock_quantity > 0 && book.stock_quantity < 5 && (
              <span className="text-xs text-orange-600">
                Only {book.stock_quantity} left
              </span>
            )}
          </div>
        </CardContent>
      </div>
      <CardFooter className="p-4 pt-0 flex gap-2">
        <Button
          className="flex-1"
          size="sm"
          onClick={(e) => {
            e.stopPropagation();
            onAddToCart?.();
          }}
          disabled={book.stock_quantity === 0}
        >
          <ShoppingCart className="h-4 w-4 mr-2" />
          Add to Cart
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={(e) => {
            e.stopPropagation();
            onAddToWishlist?.();
          }}
        >
          <Heart className="h-4 w-4" />
        </Button>
      </CardFooter>
    </Card>
  );
}
