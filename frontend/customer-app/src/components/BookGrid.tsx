import { Book } from "@/types/book";
import BookCard from "./BookCard";
import { Loader2 } from "lucide-react";

interface BookGridProps {
  books: Book[];
  isLoading?: boolean;
  onBookClick?: (book: Book) => void;
  onAddToCart?: (book: Book) => void;
  onAddToWishlist?: (book: Book) => void;
}

export default function BookGrid({
  books,
  isLoading = false,
  onBookClick,
  onAddToCart,
  onAddToWishlist,
}: BookGridProps) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="h-12 w-12 animate-spin text-gray-400" />
      </div>
    );
  }

  if (!books || books.length === 0) {
    return (
      <div className="text-center py-20">
        <div className="text-6xl mb-4">📚</div>
        <h3 className="text-xl font-semibold text-gray-700 mb-2">
          No books found
        </h3>
        <p className="text-gray-500">
          Try adjusting your search or browse our categories
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
      {books.map((book) => (
        <BookCard
          key={book.id}
          book={book}
          onClick={() => onBookClick?.(book)}
          onAddToCart={() => onAddToCart?.(book)}
          onAddToWishlist={() => onAddToWishlist?.(book)}
        />
      ))}
    </div>
  );
}
