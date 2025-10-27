import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { wishlistAPI } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Trash2 } from "lucide-react";

export default function Wishlist() {
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuthStore();

  const { data: wishlistData, isLoading } = useQuery({
    queryKey: ["wishlist"],
    queryFn: () => wishlistAPI.list().then((res) => res.data),
    enabled: isAuthenticated,
  });

  const removeFromWishlistMutation = useMutation({
    mutationFn: (bookId: string) => wishlistAPI.remove(bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["wishlist"] });
    },
  });

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Login Required</CardTitle>
            <CardDescription>
              You need to be logged in to view your wishlist
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link to="/login">
              <Button className="w-full">Login</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <p className="text-muted-foreground">Loading your wishlist...</p>
      </div>
    );
  }

  const items = wishlistData?.data || [];

  return (
    <div className="min-h-screen bg-background py-8">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">
              My Wishlist
            </h1>
            <p className="text-muted-foreground">
              {items.length} {items.length === 1 ? "book" : "books"} saved
            </p>
          </div>
          <Link to="/">
            <Button variant="outline">Back to Books</Button>
          </Link>
        </div>

        {items.length === 0 ? (
          <Card>
            <CardHeader>
              <CardTitle>Your wishlist is empty</CardTitle>
              <CardDescription>
                Start adding books you're interested in to your wishlist!
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Link to="/">
                <Button>Browse Books</Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((item) => {
              const book = item.book;
              if (!book) return null;

              return (
                <Card key={item.id} className="flex flex-col">
                  <CardHeader>
                    {book.cover_image_url && (
                      <div className="w-full h-48 mb-4 bg-gray-200 rounded-md overflow-hidden">
                        <img
                          src={book.cover_image_url}
                          alt={book.title}
                          className="w-full h-full object-cover"
                        />
                      </div>
                    )}
                    <CardTitle className="line-clamp-2">{book.title}</CardTitle>
                    <CardDescription className="line-clamp-1">
                      {book.authors?.map((a) => a.name).join(", ") ||
                        "Unknown Author"}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="flex-1">
                    <p className="text-2xl font-bold text-green-600">
                      ${book.price.toFixed(2)}
                    </p>
                    <p className="text-sm text-gray-500 mt-1">
                      {book.stock_quantity > 0 ? (
                        <span className="text-green-600">
                          In Stock ({book.stock_quantity})
                        </span>
                      ) : (
                        <span className="text-red-600">Out of Stock</span>
                      )}
                    </p>
                    <p className="text-xs text-gray-400 mt-2">
                      Added {new Date(item.created_at).toLocaleDateString()}
                    </p>
                  </CardContent>
                  <CardFooter className="flex gap-2">
                    <Link
                      to="/books/$id"
                      params={{ id: book.id }}
                      className="flex-1"
                    >
                      <Button className="w-full" variant="outline">
                        View Details
                      </Button>
                    </Link>
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => removeFromWishlistMutation.mutate(book.id)}
                      disabled={removeFromWishlistMutation.isPending}
                    >
                      <Trash2 className="h-4 w-4 text-red-600" />
                    </Button>
                  </CardFooter>
                </Card>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
