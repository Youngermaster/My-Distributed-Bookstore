import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { booksAPI, wishlistAPI, cartAPI, recommendationsAPI } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import { useCartStore } from "@/store/cartStore";
import { createId } from "@/lib/id";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Heart, ShoppingCart } from "lucide-react";
import BookGrid from "@/components/BookGrid";
import { ReviewStats, ReviewList, ReviewForm } from "@/components/reviews";
import { useReviewService } from "@/services";
import { type Book } from "@/types/book";
import { toast } from "sonner";

export default function BookDetail() {
  const { id } = useParams({ from: "/books/$id" });
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { isAuthenticated, user } = useAuthStore();
  const { cartId, setCartId } = useCartStore();
  const { useBookReviews, useBookReviewStats, useVoteReview, useCreateReview } =
    useReviewService();
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  // Initialize cart ID if not exists
  useEffect(() => {
    if (!cartId) {
      const newCartId = createId();
      setCartId(newCartId);
    }
  }, [cartId, setCartId]);

  // Track book view when component mounts
  useEffect(() => {
    if (id && isAuthenticated) {
      recommendationsAPI
        .trackInteraction({
          book_id: id,
          interaction_type: "view",
        })
        .catch((error) => {
          console.error("Failed to track view:", error);
        });
    }
  }, [id, isAuthenticated]);

  const { data: bookData, isLoading } = useQuery({
    queryKey: ["book", id],
    queryFn: () => booksAPI.get(id!).then((res) => res.data),
    enabled: !!id,
  });

  // Fetch similar books
  const { data: similarBooksData, isLoading: similarBooksLoading } = useQuery({
    queryKey: ["recommendations", "similar", id],
    queryFn: () =>
      recommendationsAPI.getSimilar(id!, { limit: 6 }).then((res) => res.data),
    enabled: !!id,
  });

  // Helper function to fetch books by IDs
  const fetchBooksByIds = async (bookIds: string[]): Promise<Book[]> => {
    const bookPromises = bookIds.map((bookId) =>
      booksAPI.get(bookId).then((res) => res.data)
    );
    const books = await Promise.all(bookPromises);
    return books.filter((book): book is Book => book !== null);
  };

  // Fetch full book details for similar books
  const recommendationIds =
    similarBooksData?.recommendations?.map(
      (recommendation) => recommendation.book_id
    ) ?? [];

  const { data: similarBooks } = useQuery({
    queryKey: ["books", "similar", recommendationIds],
    queryFn: () => fetchBooksByIds(recommendationIds),
    enabled: recommendationIds.length > 0,
  });

  // Reviews
  const bookId = id ?? "";
  const { data: reviewListData, isLoading: reviewsLoading } = useBookReviews(
    bookId,
    1,
    10
  );
  const { data: reviewStatsData, isLoading: reviewStatsLoading } =
    useBookReviewStats(bookId);
  const voteReviewMutation = useVoteReview();
  const createReviewMutation = useCreateReview();

  const reviews = reviewListData?.reviews ?? [];

  const addToWishlistMutation = useMutation({
    mutationFn: (bookId: string) => wishlistAPI.add(bookId),
    onSuccess: (_, bookId) => {
      queryClient.invalidateQueries({ queryKey: ["wishlist"] });
      setMessage({ type: "success", text: "Added to wishlist!" });
      setTimeout(() => setMessage(null), 3000);
      toast.success("Added to wishlist.");

      // Track wishlist interaction
      if (isAuthenticated) {
        recommendationsAPI
          .trackInteraction({
            book_id: bookId,
            interaction_type: "wishlist",
          })
          .catch((error) => {
            console.error("Failed to track wishlist interaction:", error);
          });
      }
    },
    onError: (error: any) => {
      const errorMsg =
        error.response?.data?.message || "Failed to add to wishlist";
      setMessage({ type: "error", text: errorMsg });
      setTimeout(() => setMessage(null), 3000);
      toast.error(errorMsg);
    },
  });

  const addToCartMutation = useMutation({
    mutationFn: ({ bookId, price }: { bookId: string; price: number }) =>
      cartAPI.addItem(cartId!, {
        book_id: bookId,
        quantity: 1,
        price,
      }),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
      setMessage({ type: "success", text: "Added to cart!" });
      setTimeout(() => setMessage(null), 3000);
      toast.success("Added to cart.");

      // Track add to cart interaction
      if (isAuthenticated) {
        recommendationsAPI
          .trackInteraction({
            book_id: variables.bookId,
            interaction_type: "add_to_cart",
          })
          .catch((error) => {
            console.error("Failed to track cart interaction:", error);
          });
      }
    },
    onError: (error: any) => {
      const errorMsg = error.response?.data?.message || "Failed to add to cart";
      setMessage({ type: "error", text: errorMsg });
      setTimeout(() => setMessage(null), 3000);
      toast.error(errorMsg);
    },
  });

  const handleAddToWishlist = () => {
    if (!isAuthenticated) {
      toast.info("Please log in to add items to your wishlist.");
      navigate({ to: "/login" });
      return;
    }
    if (id) {
      addToWishlistMutation.mutate(id);
    }
  };

  const handleAddToCart = () => {
    if (id && bookData) {
      addToCartMutation.mutate({ bookId: id, price: bookData.price });
    }
  };

  const handleReviewHelpful = (reviewId: string) => {
    if (!id) {
      return;
    }

    if (!isAuthenticated || !user?.id) {
      setMessage({ type: "error", text: "Please log in to vote on reviews." });
      setTimeout(() => setMessage(null), 3000);
      toast.error("Please log in to vote on reviews.");
      if (!isAuthenticated) {
        navigate({ to: "/login" });
      }
      return;
    }

    if (voteReviewMutation.isPending) {
      return;
    }

    voteReviewMutation.mutate(
      {
        reviewId,
        bookId: id,
        payload: {
          user_id: user.id,
          is_helpful: true,
        },
      },
      {
        onSuccess: () => {
          setMessage({ type: "success", text: "Thanks for your feedback!" });
          setTimeout(() => setMessage(null), 3000);
          toast.success("Thanks for your feedback!");
        },
        onError: () => {
          setMessage({
            type: "error",
            text: "Unable to record your vote right now.",
          });
          setTimeout(() => setMessage(null), 3000);
          toast.error("Unable to record your vote right now.");
        },
      }
    );
  };

  const handleReviewSubmit = async (data: {
    rating: number;
    title: string;
    content: string;
  }) => {
    if (!id || !user?.id) {
      if (!isAuthenticated) {
        navigate({ to: "/login" });
        return;
      }

      setMessage({ type: "error", text: "We couldn't verify your account." });
      setTimeout(() => setMessage(null), 3000);
      toast.error("We couldn't verify your account.");
      return;
    }

    await createReviewMutation.mutateAsync(
      {
        book_id: id,
        user_id: user.id,
        rating: data.rating,
        title: data.title,
        content: data.content,
      },
      {
        onSuccess: () => {
          setMessage({
            type: "success",
            text: "Review submitted successfully.",
          });
          setTimeout(() => setMessage(null), 3000);
          toast.success("Review submitted successfully.");
        },
        onError: () => {
          setMessage({
            type: "error",
            text: "We couldn't submit your review right now. Please try again shortly.",
          });
          setTimeout(() => setMessage(null), 3000);
          toast.error(
            "We couldn't submit your review right now. Please try again shortly."
          );
        },
      }
    );
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-600">Loading book details...</p>
      </div>
    );
  }

  if (!bookData) {
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

  const book = bookData;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        <Link to="/" className="inline-block mb-6">
          <Button variant="outline">← Back to Home</Button>
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
                    onClick={handleAddToCart}
                    disabled={
                      book.stock_quantity === 0 || addToCartMutation.isPending
                    }
                    className="flex-1"
                  >
                    <ShoppingCart className="mr-2 h-4 w-4" />
                    {addToCartMutation.isPending
                      ? "Adding..."
                      : book.stock_quantity === 0
                      ? "Out of Stock"
                      : "Add to Cart"}
                  </Button>
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

        {/* Reviews Section */}
        <div className="mt-12 space-y-6">
          <h2 className="text-3xl font-bold text-gray-900">Reviews</h2>

          {/* Review Stats */}
          {reviewStatsLoading ? (
            <div className="text-muted-foreground">
              Loading review insights...
            </div>
          ) : reviewStatsData ? (
            <ReviewStats stats={reviewStatsData} />
          ) : (
            <div className="text-muted-foreground">
              No review data available yet.
            </div>
          )}

          {/* Review Form */}
          <ReviewForm
            bookId={id!}
            bookTitle={book.title}
            isAuthenticated={isAuthenticated}
            isSubmitting={createReviewMutation.isPending}
            onSubmit={handleReviewSubmit}
          />

          {/* Review List */}
          {reviewsLoading ? (
            <div className="text-muted-foreground">Loading reviews...</div>
          ) : (
            <ReviewList
              reviews={reviews}
              onHelpful={handleReviewHelpful}
              isVoting={voteReviewMutation.isPending}
            />
          )}
        </div>

        {/* Similar Books Section */}
        {similarBooks && similarBooks.length > 0 && (
          <div className="mt-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-6">
              Similar Books You Might Like
            </h2>
            <BookGrid
              books={similarBooks}
              isLoading={similarBooksLoading}
              onBookClick={(book) =>
                navigate({ to: "/books/$id", params: { id: book.id } })
              }
            />
          </div>
        )}
      </div>
    </div>
  );
}
