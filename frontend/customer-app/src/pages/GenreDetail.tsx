import { useQuery } from "@tanstack/react-query";
import { booksAPI, categoriesAPI } from "@/lib/api";
import { useNavigate, useParams, Link } from "react-router-dom";
import BookGrid from "@/components/BookGrid";
import { Button } from "@/components/ui/button";
import { ChevronLeft, Loader2 } from "lucide-react";
import { type Book } from "@/types/book";

export default function GenreDetail() {
  const { slug } = useParams<{ slug: string }>();
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
    navigate(`/books/${book.id}`);
  };

  if (categoriesLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Loader2 className="h-12 w-12 animate-spin text-gray-400" />
      </div>
    );
  }

  if (!category) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">❓</div>
          <h2 className="text-2xl font-semibold text-gray-700 mb-2">
            Genre not found
          </h2>
          <p className="text-gray-500 mb-6">
            The genre "{slug}" doesn't exist or has been removed
          </p>
          <Button onClick={() => navigate("/genres")}>Browse All Genres</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Breadcrumb */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <nav className="flex items-center text-sm text-gray-600">
            <Link to="/" className="hover:text-gray-900">
              Home
            </Link>
            <span className="mx-2">/</span>
            <Link to="/genres" className="hover:text-gray-900">
              Genres
            </Link>
            <span className="mx-2">/</span>
            <span className="text-gray-900 font-medium">{category.name}</span>
          </nav>
        </div>
      </div>

      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white py-12">
        <div className="max-w-7xl mx-auto px-4">
          <Button
            variant="ghost"
            className="mb-4 text-white hover:text-white hover:bg-white/20"
            onClick={() => navigate("/genres")}
          >
            <ChevronLeft className="h-4 w-4 mr-2" />
            Back to Genres
          </Button>
          <h1 className="text-4xl font-bold mb-2">{category.name}</h1>
          {booksData && (
            <p className="text-lg text-blue-100">
              {booksData.books.length}{" "}
              {booksData.books.length === 1 ? "book" : "books"} in this genre
            </p>
          )}
        </div>
      </div>

      {/* Books Grid */}
      <div className="max-w-7xl mx-auto px-4 py-12">
        <BookGrid
          books={booksData?.books || []}
          isLoading={booksLoading}
          onBookClick={handleBookClick}
        />
      </div>
    </div>
  );
}
