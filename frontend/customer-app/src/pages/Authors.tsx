import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useCatalogService } from "@/services";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Search, User } from "lucide-react";

export default function Authors() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const { useAuthors } = useCatalogService();
  const { data, isLoading } = useAuthors(page, 20);

  // Filter authors by search query
  const filteredAuthors = data?.authors
    ? data.authors.filter((author: any) =>
        author.name.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : [];

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">Authors</h1>
          <p className="text-muted-foreground">
            Discover books from your favorite authors
          </p>
        </div>

        {/* Search */}
        <div className="mb-8">
          <div className="relative max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-5 h-5" />
            <Input
              type="text"
              placeholder="Search authors..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {/* Authors Grid */}
        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading authors...</p>
          </div>
        ) : filteredAuthors.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center">
              <User className="w-16 h-16 mx-auto mb-4 text-muted-foreground" />
              <p className="text-lg text-muted-foreground">
                {searchQuery ? "No authors found" : "No authors available"}
              </p>
            </CardContent>
          </Card>
        ) : (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {filteredAuthors.map((author: any) => (
                <AuthorCard
                  key={author.id}
                  author={author}
                  onClick={() =>
                    navigate({ to: "/authors/$id", params: { id: author.id } })
                  }
                />
              ))}
            </div>

            {/* Pagination */}
            {data && data.total_pages > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                <Button
                  variant="outline"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                >
                  Previous
                </Button>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">
                    Page {page} of {data.total_pages}
                  </span>
                </div>
                <Button
                  variant="outline"
                  onClick={() =>
                    setPage((p) => Math.min(data.total_pages, p + 1))
                  }
                  disabled={page === data.total_pages}
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

interface AuthorCardProps {
  author: {
    id: string;
    name: string;
    bio?: string;
    books_count?: number;
  };
  onClick: () => void;
}

function AuthorCard({ author, onClick }: AuthorCardProps) {
  return (
    <Card
      className="cursor-pointer hover:shadow-lg transition-shadow"
      onClick={onClick}
    >
      <CardHeader>
        <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
          <User className="w-8 h-8 text-primary" />
        </div>
        <CardTitle className="text-center text-lg">{author.name}</CardTitle>
        {author.books_count && (
          <CardDescription className="text-center">
            {author.books_count} {author.books_count === 1 ? "book" : "books"}
          </CardDescription>
        )}
      </CardHeader>
      {author.bio && (
        <CardContent>
          <p className="text-sm text-muted-foreground line-clamp-3 text-center">
            {author.bio}
          </p>
        </CardContent>
      )}
    </Card>
  );
}
