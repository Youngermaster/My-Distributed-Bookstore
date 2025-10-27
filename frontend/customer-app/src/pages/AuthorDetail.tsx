import { useQuery } from "@tanstack/react-query";
import { authorsAPI, booksAPI } from "@/lib/api";
import { useNavigate, useParams, Link } from "react-router-dom";
import BookGrid from "@/components/BookGrid";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ChevronLeft, Loader2, User, Calendar } from "lucide-react";
import { type Book } from "@/types/book";

export default function AuthorDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  // Fetch author details
  const {
    data: authorData,
    isLoading: authorLoading,
    error,
  } = useQuery({
    queryKey: ["author", id],
    queryFn: async () => {
      if (!id) throw new Error("Author ID is required");
      const response = await authorsAPI.get(id);
      return response.data;
    },
    enabled: !!id,
  });

  // Fetch books by this author
  const { data: booksData, isLoading: booksLoading } = useQuery({
    queryKey: ["books", "author", id],
    queryFn: async () => {
      if (!id) return null;
      const response = await booksAPI.list({
        author_id: id,
        page: 1,
        page_size: 100,
      });
      return response.data;
    },
    enabled: !!id,
  });

  const handleBookClick = (book: Book) => {
    navigate(`/books/${book.id}`);
  };

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">❓</div>
          <h2 className="text-2xl font-semibold text-gray-700 mb-2">
            Author not found
          </h2>
          <p className="text-gray-500 mb-6">
            {error instanceof Error ? error.message : "An error occurred"}
          </p>
          <Button onClick={() => navigate("/")}>Back to Home</Button>
        </div>
      </div>
    );
  }

  if (authorLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Loader2 className="h-12 w-12 animate-spin text-gray-400" />
      </div>
    );
  }

  const author = authorData;

  if (!author) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">❓</div>
          <h2 className="text-2xl font-semibold text-gray-700 mb-2">
            Author not found
          </h2>
          <p className="text-gray-500 mb-6">
            This author doesn't exist or has been removed
          </p>
          <Button onClick={() => navigate("/")}>Back to Home</Button>
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
            <Link to="/authors" className="hover:text-gray-900">
              Authors
            </Link>
            <span className="mx-2">/</span>
            <span className="text-gray-900 font-medium">{author.name}</span>
          </nav>
        </div>
      </div>

      {/* Author Info Section */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white py-12">
        <div className="max-w-7xl mx-auto px-4">
          <Button
            variant="ghost"
            className="mb-4 text-white hover:text-white hover:bg-white/20"
            onClick={() => navigate(-1)}
          >
            <ChevronLeft className="h-4 w-4 mr-2" />
            Back
          </Button>
          <div className="flex items-start gap-8">
            <div className="flex-shrink-0">
              <div className="w-32 h-32 bg-white/20 rounded-full flex items-center justify-center">
                <User className="h-16 w-16 text-white" />
              </div>
            </div>
            <div className="flex-1">
              <h1 className="text-4xl font-bold mb-4">{author.name}</h1>
              {author.birth_date && (
                <div className="flex items-center gap-2 text-blue-100 mb-4">
                  <Calendar className="h-4 w-4" />
                  <span>
                    Born: {new Date(author.birth_date).toLocaleDateString()}
                  </span>
                </div>
              )}
              {booksData && (
                <p className="text-lg text-blue-100">
                  {booksData.books.length} published{" "}
                  {booksData.books.length === 1 ? "book" : "books"}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-12">
        {/* Biography */}
        {author.bio && (
          <Card className="mb-12">
            <CardContent className="p-6">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">
                Biography
              </h2>
              <p className="text-gray-700 leading-relaxed whitespace-pre-line">
                {author.bio}
              </p>
            </CardContent>
          </Card>
        )}

        {/* Books by this author */}
        <section>
          <h2 className="text-3xl font-bold text-gray-900 mb-6">
            Books by {author.name}
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
