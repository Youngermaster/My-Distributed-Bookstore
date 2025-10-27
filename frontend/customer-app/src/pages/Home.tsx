import { useQuery } from "@tanstack/react-query";
import { booksAPI, categoriesAPI } from "@/lib/api";
import { useNavigate } from "@tanstack/react-router";
import BookGrid from "@/components/BookGrid";
import GenreCard from "@/components/GenreCard";
import SearchBar from "@/components/SearchBar";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { type Book } from "@/types/book";

export default function Home() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");

  // Fetch featured books (first 10)
  const { data: booksData, isLoading: booksLoading } = useQuery({
    queryKey: ["books", "featured"],
    queryFn: async () => {
      const response = await booksAPI.list({ page: 1, page_size: 10 });
      return response.data;
    },
  });

  // Fetch all categories
  const { data: categoriesData, isLoading: categoriesLoading } = useQuery({
    queryKey: ["categories"],
    queryFn: async () => {
      const response = await categoriesAPI.list();
      return response.data;
    },
  });

  // Search books
  const { data: searchData, isLoading: searchLoading } = useQuery({
    queryKey: ["books", "search", searchQuery],
    queryFn: async () => {
      if (!searchQuery) return null;
      const response = await booksAPI.search(searchQuery);
      return response.data;
    },
    enabled: !!searchQuery,
  });

  const handleBookClick = (book: Book) => {
    navigate({ to: "/books/$id", params: { id: book.id } });
  };

  const handleSearch = (query: string) => {
    setSearchQuery(query);
  };

  const displayBooks =
    searchQuery && searchData ? searchData.books : booksData?.books || [];
  const isLoading = searchQuery ? searchLoading : booksLoading;

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <div className="bg-linear-to-r from-primary to-primary/80 text-primary-foreground py-16">
        <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-4xl md:text-5xl font-bold mb-4 text-center">
            Welcome to Ohara Bookstore
          </h1>
          <p className="text-xl mb-8 text-center opacity-90">
            Discover your next favorite book
          </p>
          <div className="flex justify-center">
            <SearchBar
              onSearch={handleSearch}
              placeholder="Search for books, authors, or topics..."
            />
          </div>
        </div>
      </div>

      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Search Results or Featured Books */}
        <section className="mb-16">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-3xl font-bold text-foreground">
              {searchQuery ? "Search Results" : "Featured Books"}
            </h2>
            {!searchQuery && (
              <Button
                variant="outline"
                onClick={() => navigate({ to: "/books" })}
              >
                View All Books
              </Button>
            )}
          </div>
          <BookGrid
            books={displayBooks}
            isLoading={isLoading}
            onBookClick={handleBookClick}
          />
        </section>

        {/* Categories Section - Only show when not searching */}
        {!searchQuery && (
          <section>
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-3xl font-bold text-foreground">
                Browse by Genre
              </h2>
              <Button
                variant="outline"
                onClick={() => navigate({ to: "/genres" })}
              >
                View All Genres
              </Button>
            </div>
            {categoriesLoading ? (
              <div className="text-center py-10">
                <p className="text-muted-foreground">Loading categories...</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                {categoriesData?.categories.slice(0, 8).map((category) => (
                  <GenreCard
                    key={category.id}
                    category={category}
                    onClick={() =>
                      navigate({
                        to: "/genres/$slug",
                        params: { slug: category.slug },
                      })
                    }
                  />
                ))}
              </div>
            )}
          </section>
        )}
      </div>
    </div>
  );
}
