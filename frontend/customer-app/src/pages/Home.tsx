import { useQuery } from "@tanstack/react-query";
import { booksAPI, categoriesAPI, recommendationsAPI } from "@/lib/api";
import { useNavigate } from "@tanstack/react-router";
import BookGrid from "@/components/BookGrid";
import GenreCard from "@/components/GenreCard";
import SearchBar from "@/components/SearchBar";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { type Book } from "@/types/book";
import { useAuthStore } from "@/store/authStore";

export default function Home() {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");
  const { isAuthenticated } = useAuthStore();

  // Fetch featured books (first 10)
  const { data: booksData, isLoading: booksLoading } = useQuery({
    queryKey: ["books", "featured"],
    queryFn: async () => {
      const response = await booksAPI.list({ page: 1, page_size: 10 });
      return response.data;
    },
  });

  // Fetch personalized recommendations (only for authenticated users)
  const { data: recommendationsData } = useQuery({
    queryKey: ["recommendations", "personalized"],
    queryFn: async () => {
      const response = await recommendationsAPI.getPersonalized({ limit: 10 });
      return response.data;
    },
    enabled: isAuthenticated,
  });

  // Fetch trending books
  const { data: trendingData } = useQuery({
    queryKey: ["recommendations", "trending"],
    queryFn: async () => {
      const response = await recommendationsAPI.getTrending({
        limit: 10,
        days: 7,
      });
      return response.data;
    },
  });

  // Fetch popular books
  const { data: popularData } = useQuery({
    queryKey: ["recommendations", "popular"],
    queryFn: async () => {
      const response = await recommendationsAPI.getPopular({ limit: 10 });
      return response.data;
    },
  });

  // Helper function to fetch books by IDs
  const fetchBooksByIds = async (bookIds: string[]): Promise<Book[]> => {
    const bookPromises = bookIds.map((id) =>
      booksAPI.get(id).then((res) => res.data)
    );
    const books = await Promise.all(bookPromises);
    return books.filter((book): book is Book => book !== null);
  };

  // Fetch full book details for recommendations
  const { data: recommendedBooks, isLoading: recommendedBooksLoading } =
    useQuery({
      queryKey: [
        "books",
        "recommended",
        recommendationsData?.recommendations.map((r) => r.book_id),
      ],
      queryFn: () =>
        fetchBooksByIds(
          recommendationsData?.recommendations.map((r) => r.book_id) || []
        ),
      enabled:
        isAuthenticated &&
        !!recommendationsData &&
        recommendationsData.recommendations.length > 0,
    });

  // Fetch full book details for trending books
  const { data: trendingBooks, isLoading: trendingBooksLoading } = useQuery({
    queryKey: [
      "books",
      "trending",
      trendingData?.recommendations.map((r) => r.book_id),
    ],
    queryFn: () =>
      fetchBooksByIds(
        trendingData?.recommendations.map((r) => r.book_id) || []
      ),
    enabled: !!trendingData && trendingData.recommendations.length > 0,
  });

  // Fetch full book details for popular books
  const { data: popularBooks, isLoading: popularBooksLoading } = useQuery({
    queryKey: [
      "books",
      "popular",
      popularData?.recommendations.map((r) => r.book_id),
    ],
    queryFn: () =>
      fetchBooksByIds(popularData?.recommendations.map((r) => r.book_id) || []),
    enabled: !!popularData && popularData.recommendations.length > 0,
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

        {/* Personalized Recommendations - Only show for authenticated users when not searching */}
        {!searchQuery &&
          isAuthenticated &&
          recommendedBooks &&
          recommendedBooks.length > 0 && (
            <section className="mb-16">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-3xl font-bold text-foreground">
                  Recommended for You
                </h2>
              </div>
              <BookGrid
                books={recommendedBooks}
                isLoading={recommendedBooksLoading}
                onBookClick={handleBookClick}
              />
            </section>
          )}

        {/* Trending Books - Only show when not searching */}
        {!searchQuery && trendingBooks && trendingBooks.length > 0 && (
          <section className="mb-16">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-3xl font-bold text-foreground">
                Trending This Week
              </h2>
            </div>
            <BookGrid
              books={trendingBooks}
              isLoading={trendingBooksLoading}
              onBookClick={handleBookClick}
            />
          </section>
        )}

        {/* Popular Books - Only show when not searching */}
        {!searchQuery && popularBooks && popularBooks.length > 0 && (
          <section className="mb-16">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-3xl font-bold text-foreground">
                Most Popular
              </h2>
            </div>
            <BookGrid
              books={popularBooks}
              isLoading={popularBooksLoading}
              onBookClick={handleBookClick}
            />
          </section>
        )}

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
