import { useQuery } from "@tanstack/react-query";
import { publishersAPI, booksAPI } from "@/lib/api";
import { useNavigate, useParams, Link } from "@tanstack/react-router";
import BookGrid from "@/components/BookGrid";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ChevronLeft, Loader2, Building } from "lucide-react";
import { type Book } from "@/types/book";

export default function PublisherDetail() {
  const { id } = useParams({ from: "/publishers/$id" });
  const navigate = useNavigate();

  // Fetch publisher details
  const {
    data: publisherData,
    isLoading: publisherLoading,
    error,
  } = useQuery({
    queryKey: ["publisher", id],
    queryFn: async () => {
      if (!id) throw new Error("Publisher ID is required");
      const response = await publishersAPI.get(id);
      return response.data;
    },
    enabled: !!id,
  });

  // Fetch all books and filter by publisher client-side
  // TODO: Add publisher_id filter to backend API
  const { data: allBooksData, isLoading: booksLoading } = useQuery({
    queryKey: ["books", "all"],
    queryFn: async () => {
      const response = await booksAPI.list({
        page: 1,
        page_size: 1000, // Get all books
      });
      return response.data;
    },
  });

  // Filter books by this publisher
  const booksData = allBooksData
    ? {
        ...allBooksData,
        books: allBooksData.books.filter(
          (book: any) => book.publisher?.id === id
        ),
      }
    : undefined;

  const handleBookClick = (book: Book) => {
    navigate({ to: "/books/$id", params: { id: book.id } });
  };

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">❓</div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">
            Publisher not found
          </h2>
          <p className="text-muted-foreground mb-6">
            {error instanceof Error ? error.message : "An error occurred"}
          </p>
          <Button onClick={() => navigate({ to: "/" })}>Back to Home</Button>
        </div>
      </div>
    );
  }

  if (publisherLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-12 w-12 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const publisher = publisherData;

  if (!publisher) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">❓</div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">
            Publisher not found
          </h2>
          <p className="text-muted-foreground mb-6">
            This publisher doesn't exist or has been removed
          </p>
          <Button onClick={() => navigate({ to: "/" })}>Back to Home</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Breadcrumb */}
      <div className="bg-card border-b border-border">
        <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <nav className="flex items-center text-sm text-muted-foreground">
            <Link to="/" className="hover:text-foreground">
              Home
            </Link>
            <span className="mx-2">/</span>
            <Link to="/publishers" className="hover:text-foreground">
              Publishers
            </Link>
            <span className="mx-2">/</span>
            <span className="text-foreground font-medium">
              {publisher.name}
            </span>
          </nav>
        </div>
      </div>

      {/* Publisher Info Section */}
      <div className="bg-gradient-to-r from-primary to-primary/80 text-primary-foreground py-12">
        <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8">
          <Button
            variant="ghost"
            className="mb-4 text-primary-foreground hover:text-primary-foreground hover:bg-primary-foreground/20"
            onClick={() => window.history.back()}
          >
            <ChevronLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
          <div className="flex items-start gap-8">
            <div className="flex-shrink-0">
              <div className="w-32 h-32 bg-primary-foreground/20 rounded-lg flex items-center justify-center">
                <Building className="h-16 w-16 text-primary-foreground" />
              </div>
            </div>
            <div className="flex-1">
              <h1 className="text-4xl font-bold mb-4">{publisher.name}</h1>
              {publisher.country && (
                <div className="flex items-center gap-2 opacity-90 mb-4">
                  <span>📍 {publisher.country}</span>
                </div>
              )}
              {booksData && (
                <p className="text-lg opacity-90">
                  {booksData.books.length} published{" "}
                  {booksData.books.length === 1 ? "book" : "books"}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Biography/Description */}
        {publisher.description && (
          <Card className="mb-12">
            <CardContent className="p-6">
              <h2 className="text-2xl font-bold text-foreground mb-4">
                About {publisher.name}
              </h2>
              <p className="text-muted-foreground leading-relaxed whitespace-pre-line">
                {publisher.description}
              </p>
            </CardContent>
          </Card>
        )}

        {/* Books by this publisher */}
        <section>
          <h2 className="text-3xl font-bold text-foreground mb-6">
            Books from {publisher.name}
          </h2>
          <BookGrid
            books={booksData?.books || []}
            isLoading={booksLoading}
            onBookClick={handleBookClick}
          />
        </section>
      </div>
    </div>
  );
}
