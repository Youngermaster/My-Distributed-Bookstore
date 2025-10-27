import { useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { booksAPI, wishlistAPI } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Heart } from "lucide-react";

export default function BookDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuthStore();
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const { data: bookData, isLoading } = useQuery({
    queryKey: ["book", id],
    queryFn: () => booksAPI.get(id!).then((res) => res.data),
    enabled: !!id,
  });

  const addToWishlistMutation = useMutation({
    mutationFn: (bookId: string) => wishlistAPI.add(bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["wishlist"] });
      setMessage({ type: "success", text: "Added to wishlist!" });
      setTimeout(() => setMessage(null), 3000);
    },
    onError: (error: any) => {
      const errorMsg =
        error.response?.data?.message || "Failed to add to wishlist";
      setMessage({ type: "error", text: errorMsg });
      setTimeout(() => setMessage(null), 3000);
    },
  });

  const handleAddToWishlist = () => {
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }
    if (id) {
      addToWishlistMutation.mutate(id);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-600">Loading book details...</p>
      </div>
    );
  }

  if (!bookData?.data) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Book Not Found</CardTitle>
            <CardDescription>
              The requested book could not be found.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link to="/">
              <Button>Back to Books</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  const book = bookData.data;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        <Link to="/" className="inline-block mb-6">
          <Button variant="outline">← Back to Books</Button>
        </Link>

        {message && (
          <div
            className={`mb-6 p-4 rounded-md ${
              message.type === "success"
                ? "bg-green-50 text-green-800 border border-green-200"
                : "bg-red-50 text-red-800 border border-red-200"
            }`}
          >
            {message.text}
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Book Image */}
          <div className="md:col-span-1">
            <Card>
              <CardContent className="p-0">
                {book.cover_image_url ? (
                  <img
                    src={book.cover_image_url}
                    alt={book.title}
                    className="w-full h-auto rounded-t-xl"
                  />
                ) : (
                  <div className="w-full h-96 bg-gray-200 rounded-t-xl flex items-center justify-center">
                    <p className="text-gray-400">No cover image</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Book Details */}
          <div className="md:col-span-2 space-y-6">
            <div>
              <h1 className="text-4xl font-bold text-gray-900 mb-2">
                {book.title}
              </h1>
              {book.authors && book.authors.length > 0 && (
                <p className="text-xl text-gray-600 mb-4">
                  by {book.authors.map((a) => a.name).join(", ")}
                </p>
              )}
              {book.categories && book.categories.length > 0 && (
                <div className="flex flex-wrap gap-2 mb-4">
                  {book.categories.map((category) => (
                    <span
                      key={category.id}
                      className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm"
                    >
                      {category.name}
                    </span>
                  ))}
                </div>
              )}
            </div>

            <Card>
              <CardContent className="pt-6 space-y-4">
                <div className="flex items-baseline gap-4">
                  <span className="text-4xl font-bold text-green-600">
                    ${book.price.toFixed(2)}
                  </span>
                  <span className="text-gray-600">
                    {book.stock_quantity > 0 ? (
                      <span className="text-green-600 font-medium">
                        In Stock ({book.stock_quantity} available)
                      </span>
                    ) : (
                      <span className="text-red-600 font-medium">
                        Out of Stock
                      </span>
                    )}
                  </span>
                </div>

                <div className="flex gap-2">
                  <Button
                    onClick={handleAddToWishlist}
                    variant="outline"
                    disabled={addToWishlistMutation.isPending}
                  >
                    <Heart className="mr-2 h-4 w-4" />
                    {addToWishlistMutation.isPending
                      ? "Adding..."
                      : "Add to Wishlist"}
                  </Button>
                </div>
              </CardContent>
            </Card>

            {book.description && (
              <Card>
                <CardHeader>
                  <CardTitle>Description</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-gray-700 leading-relaxed">
                    {book.description}
                  </p>
                </CardContent>
              </Card>
            )}

            <Card>
              <CardHeader>
                <CardTitle>Details</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-500">ISBN</p>
                    <p className="font-medium">{book.isbn}</p>
                  </div>
                  {book.publisher && (
                    <div>
                      <p className="text-sm text-gray-500">Publisher</p>
                      <p className="font-medium">{book.publisher.name}</p>
                    </div>
                  )}
                  {book.publication_date && (
                    <div>
                      <p className="text-sm text-gray-500">Publication Date</p>
                      <p className="font-medium">
                        {new Date(book.publication_date).toLocaleDateString()}
                      </p>
                    </div>
                  )}
                  {book.pages && (
                    <div>
                      <p className="text-sm text-gray-500">Pages</p>
                      <p className="font-medium">{book.pages}</p>
                    </div>
                  )}
                  {book.language && (
                    <div>
                      <p className="text-sm text-gray-500">Language</p>
                      <p className="font-medium">
                        {book.language.toUpperCase()}
                      </p>
                    </div>
                  )}
                  {book.format && (
                    <div>
                      <p className="text-sm text-gray-500">Format</p>
                      <p className="font-medium capitalize">{book.format}</p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
