import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { booksAPI, categoriesAPI } from "@/lib/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { BookFilters } from "@/types/book";

export default function BookList() {
  const [filters, setFilters] = useState<BookFilters>({
    limit: 20,
    offset: 0,
  });
  const [searchTitle, setSearchTitle] = useState("");

  const { data: booksData, isLoading: booksLoading } = useQuery({
    queryKey: ["books", filters],
    queryFn: () => booksAPI.list(filters).then((res) => res.data),
  });

  const { data: categoriesData } = useQuery({
    queryKey: ["categories"],
    queryFn: () => categoriesAPI.list().then((res) => res.data),
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setFilters((prev) => ({ ...prev, title: searchTitle, offset: 0 }));
  };

  const handleCategoryFilter = (categoryId: string) => {
    setFilters((prev) => ({
      ...prev,
      category_id: prev.category_id === categoryId ? undefined : categoryId,
      offset: 0,
    }));
  };

  const handleNextPage = () => {
    setFilters((prev) => ({
      ...prev,
      offset: (prev.offset || 0) + (prev.limit || 20),
    }));
  };

  const handlePrevPage = () => {
    setFilters((prev) => ({
      ...prev,
      offset: Math.max(0, (prev.offset || 0) - (prev.limit || 20)),
    }));
  };

  const currentPage =
    Math.floor((filters.offset || 0) / (filters.limit || 20)) + 1;
  const totalPages = Math.ceil((booksData?.total || 0) / (filters.limit || 20));

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Bookstore</h1>
          <p className="text-gray-600">Discover your next favorite book</p>
        </div>

        {/* Search and Filters */}
        <div className="mb-6 space-y-4">
          <form onSubmit={handleSearch} className="flex gap-2">
            <div className="flex-1">
              <Input
                type="text"
                placeholder="Search books by title..."
                value={searchTitle}
                onChange={(e) => setSearchTitle(e.target.value)}
              />
            </div>
            <Button type="submit">Search</Button>
          </form>

          {/* Category Filters */}
          {categoriesData?.data && (
            <div className="flex flex-wrap gap-2">
              <Label className="self-center">Categories:</Label>
              {categoriesData.data.map((category) => (
                <Button
                  key={category.id}
                  variant={
                    filters.category_id === category.id ? "default" : "outline"
                  }
                  size="sm"
                  onClick={() => handleCategoryFilter(category.id)}
                >
                  {category.name}
                </Button>
              ))}
            </div>
          )}
        </div>

        {/* Books Grid */}
        {booksLoading ? (
          <div className="text-center py-12">
            <p className="text-gray-600">Loading books...</p>
          </div>
        ) : !booksData?.data || booksData.data.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600">No books found</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6 mb-8">
              {booksData.data.map((book) => (
                <Card key={book.id} className="flex flex-col">
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
                  </CardContent>
                  <CardFooter>
                    <Link to={`/books/${book.id}`} className="w-full">
                      <Button className="w-full">View Details</Button>
                    </Link>
                  </CardFooter>
                </Card>
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center gap-4">
                <Button
                  variant="outline"
                  onClick={handlePrevPage}
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                <span className="text-sm text-gray-600">
                  Page {currentPage} of {totalPages}
                </span>
                <Button
                  variant="outline"
                  onClick={handleNextPage}
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
