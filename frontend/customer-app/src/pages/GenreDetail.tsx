import { useQuery } from "@tanstack/react-query";
import { booksAPI, categoriesAPI } from "@/lib/api";
import { useNavigate, useParams, Link } from "@tanstack/react-router";
import BookGrid from "@/components/BookGrid";
import { Button } from "@/components/ui/button";
import { ChevronLeft, Loader2 } from "lucide-react";
import { type Book } from "@/types/book";

export default function GenreDetail() {
  const { slug } = useParams({ from: "/genres/$slug" });
  const navigate = useNavigate();

  // Fetch all categories to find the one matching the slug
  const { data: categoriesData, isLoading: categoriesLoading } = useQuery({
    queryKey: ["categories"],
    queryFn: async () => {
      const response = await categoriesAPI.list();
      return response.data;
    },
  });

  // Find the category by slug
  const category = categoriesData?.categories.find((cat) => cat.slug === slug);

  // Fetch books for this category
  const { data: booksData, isLoading: booksLoading } = useQuery({
    queryKey: ["books", "category", category?.id],
    queryFn: async () => {
      if (!category?.id) return null;
      const response = await booksAPI.list({
        category_id: category.id,
        page: 1,
        page_size: 100,
      });
      return response.data;
    },
    enabled: !!category?.id,
  });

  const handleBookClick = (book: Book) => {
    navigate({ to: "/books/$id", params: { id: book.id } });
  };

  if (categoriesLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-12 w-12 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!category) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">❓</div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">
            Genre not found
          </h2>
          <p className="text-muted-foreground mb-6">
            The genre "{slug}" doesn't exist or has been removed
          </p>
          <Button onClick={() => navigate({ to: "/genres" })}>
            Browse All Genres
          </Button>
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
            <Link to="/genres" className="hover:text-foreground">
              Genres
            </Link>
            <span className="mx-2">/</span>
            <span className="text-foreground font-medium">{category.name}</span>
          </nav>
        </div>
      </div>

      {/* Header */}
      <div className="bg-gradient-to-r from-primary to-primary/80 text-primary-foreground py-12">
        <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8">
          <Button
            variant="ghost"
            className="mb-4 text-primary-foreground hover:text-primary-foreground hover:bg-primary-foreground/20"
            onClick={() => navigate({ to: "/genres" })}
          >
            <ChevronLeft className="h-4 w-4 mr-2" />
            Back to Genres
          </Button>
          <h1 className="text-4xl font-bold mb-2">{category.name}</h1>
          {booksData && (
            <p className="text-lg opacity-90">
              {booksData.books.length}{" "}
              {booksData.books.length === 1 ? "book" : "books"} in this genre
            </p>
          )}
        </div>
      </div>

      {/* Books Grid */}
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <BookGrid
          books={booksData?.books || []}
          isLoading={booksLoading}
          onBookClick={handleBookClick}
        />
      </div>
    </div>
  );
}
