import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { booksAPI } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Pencil, Trash2, Plus } from "lucide-react";
import type { Book } from "@/types/book";

export default function ManageBooks() {
  const queryClient = useQueryClient();
  const { isAuthenticated, user } = useAuthStore();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingBook, setEditingBook] = useState<Book | null>(null);

  // Check if user is admin
  const isAdmin = user?.role?.name === "admin";

  const { data: booksData, isLoading } = useQuery({
    queryKey: ["books"],
    queryFn: () =>
      booksAPI.list({ limit: 100, offset: 0 }).then((res) => res.data),
    enabled: isAuthenticated && isAdmin,
  });

  const deleteBookMutation = useMutation({
    mutationFn: (bookId: string) => booksAPI.delete(bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books"] });
    },
  });

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Login Required</CardTitle>
            <CardDescription>
              You need to be logged in as an admin
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link to="/login">
              <Button className="w-full">Login</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!isAdmin) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Access Denied</CardTitle>
            <CardDescription>
              You don't have permission to access this page
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link to="/">
              <Button className="w-full">Back to Home</Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-600">Loading books...</p>
      </div>
    );
  }

  const books = booksData?.books || [];

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold text-gray-900 mb-2">
              Manage Books
            </h1>
            <p className="text-gray-600">
              {books.length} {books.length === 1 ? "book" : "books"} total
            </p>
          </div>
          <div className="flex gap-2">
            <Link to="/">
              <Button variant="outline">Back to Store</Button>
            </Link>
            <Button onClick={() => setShowCreateForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Add Book
            </Button>
          </div>
        </div>

        {showCreateForm && (
          <BookForm
            onClose={() => setShowCreateForm(false)}
            onSuccess={() => {
              setShowCreateForm(false);
              queryClient.invalidateQueries({ queryKey: ["books"] });
            }}
          />
        )}

        {editingBook && (
          <BookForm
            book={editingBook}
            onClose={() => setEditingBook(null)}
            onSuccess={() => {
              setEditingBook(null);
              queryClient.invalidateQueries({ queryKey: ["books"] });
            }}
          />
        )}

        <div className="space-y-4">
          {books.map((book) => (
            <Card key={book.id}>
              <CardContent className="p-6">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900">
                      {book.title}
                    </h3>
                    <p className="text-sm text-gray-600 mt-1">
                      ISBN: {book.isbn} | Price: ${book.price.toFixed(2)} |
                      Stock: {book.stock_quantity}
                    </p>
                    {book.description && (
                      <p className="text-sm text-gray-500 mt-2 line-clamp-2">
                        {book.description}
                      </p>
                    )}
                  </div>
                  <div className="flex gap-2 ml-4">
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => setEditingBook(book)}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => {
                        if (
                          confirm("Are you sure you want to delete this book?")
                        ) {
                          deleteBookMutation.mutate(book.id);
                        }
                      }}
                      disabled={deleteBookMutation.isPending}
                    >
                      <Trash2 className="h-4 w-4 text-red-600" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}

interface BookFormProps {
  book?: Book;
  onClose: () => void;
  onSuccess: () => void;
}

function BookForm({ book, onClose, onSuccess }: BookFormProps) {
  const [formData, setFormData] = useState({
    isbn: book?.isbn || "",
    title: book?.title || "",
    description: book?.description || "",
    price: book?.price?.toString() || "",
    stock_quantity: book?.stock_quantity?.toString() || "",
    language: book?.language || "en",
    pages: book?.pages?.toString() || "",
    format: book?.format || "paperback",
  });

  const createBookMutation = useMutation({
    mutationFn: (data: Partial<Book>) => booksAPI.create(data),
    onSuccess: () => onSuccess(),
  });

  const updateBookMutation = useMutation({
    mutationFn: (data: Partial<Book>) => booksAPI.update(book!.id, data),
    onSuccess: () => onSuccess(),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const data = {
      ...formData,
      price: parseFloat(formData.price),
      stock_quantity: parseInt(formData.stock_quantity),
      pages: formData.pages ? parseInt(formData.pages) : undefined,
    };

    if (book) {
      updateBookMutation.mutate(data);
    } else {
      createBookMutation.mutate(data);
    }
  };

  const isLoading =
    createBookMutation.isPending || updateBookMutation.isPending;

  return (
    <Card className="mb-6">
      <CardHeader>
        <CardTitle>{book ? "Edit Book" : "Add New Book"}</CardTitle>
        <CardDescription>
          {book
            ? "Update book information"
            : "Fill in the details for the new book"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="isbn">ISBN *</Label>
              <Input
                id="isbn"
                value={formData.isbn}
                onChange={(e) =>
                  setFormData({ ...formData, isbn: e.target.value })
                }
                required
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="title">Title *</Label>
              <Input
                id="title"
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
                required
                disabled={isLoading}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <textarea
              id="description"
              className="w-full min-h-[100px] rounded-md border border-input bg-transparent px-3 py-2 text-sm"
              value={formData.description}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
              disabled={isLoading}
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label htmlFor="price">Price *</Label>
              <Input
                id="price"
                type="number"
                step="0.01"
                value={formData.price}
                onChange={(e) =>
                  setFormData({ ...formData, price: e.target.value })
                }
                required
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="stock_quantity">Stock Quantity *</Label>
              <Input
                id="stock_quantity"
                type="number"
                value={formData.stock_quantity}
                onChange={(e) =>
                  setFormData({ ...formData, stock_quantity: e.target.value })
                }
                required
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="pages">Pages</Label>
              <Input
                id="pages"
                type="number"
                value={formData.pages}
                onChange={(e) =>
                  setFormData({ ...formData, pages: e.target.value })
                }
                disabled={isLoading}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="language">Language</Label>
              <Input
                id="language"
                value={formData.language}
                onChange={(e) =>
                  setFormData({ ...formData, language: e.target.value })
                }
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="format">Format</Label>
              <select
                id="format"
                className="w-full h-9 rounded-md border border-input bg-transparent px-3 py-1 text-sm"
                value={formData.format}
                onChange={(e) =>
                  setFormData({ ...formData, format: e.target.value })
                }
                disabled={isLoading}
              >
                <option value="paperback">Paperback</option>
                <option value="hardcover">Hardcover</option>
                <option value="ebook">eBook</option>
              </select>
            </div>
          </div>

          <div className="flex gap-2 justify-end">
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Saving..." : book ? "Update Book" : "Create Book"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
