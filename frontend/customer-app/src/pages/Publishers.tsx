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
import { Search, Building2 } from "lucide-react";

export default function Publishers() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const { usePublishers } = useCatalogService();
  const { data, isLoading } = usePublishers(page, 20);

  // Filter publishers by search query
  const filteredPublishers = data?.publishers
    ? data.publishers.filter((publisher: any) =>
        publisher.name.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : [];

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Publishers
          </h1>
          <p className="text-muted-foreground">
            Explore books from leading publishers
          </p>
        </div>

        {/* Search */}
        <div className="mb-8">
          <div className="relative max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-5 h-5" />
            <Input
              type="text"
              placeholder="Search publishers..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {/* Publishers Grid */}
        {isLoading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Loading publishers...</p>
          </div>
        ) : filteredPublishers.length === 0 ? (
          <Card>
            <CardContent className="pt-6 text-center">
              <Building2 className="w-16 h-16 mx-auto mb-4 text-muted-foreground" />
              <p className="text-lg text-muted-foreground">
                {searchQuery
                  ? "No publishers found"
                  : "No publishers available"}
              </p>
            </CardContent>
          </Card>
        ) : (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {filteredPublishers.map((publisher: any) => (
                <PublisherCard
                  key={publisher.id}
                  publisher={publisher}
                  onClick={() =>
                    navigate({
                      to: "/publishers/$id",
                      params: { id: publisher.id },
                    })
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

interface PublisherCardProps {
  publisher: {
    id: string;
    name: string;
    country?: string;
    books_count?: number;
  };
  onClick: () => void;
}

function PublisherCard({ publisher, onClick }: PublisherCardProps) {
  return (
    <Card
      className="cursor-pointer hover:shadow-lg transition-shadow"
      onClick={onClick}
    >
      <CardHeader>
        <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
          <Building2 className="w-8 h-8 text-primary" />
        </div>
        <CardTitle className="text-center text-lg">{publisher.name}</CardTitle>
        {publisher.country && (
          <CardDescription className="text-center">
            {publisher.country}
          </CardDescription>
        )}
        {publisher.books_count && (
          <CardDescription className="text-center">
            {publisher.books_count}{" "}
            {publisher.books_count === 1 ? "book" : "books"}
          </CardDescription>
        )}
      </CardHeader>
    </Card>
  );
}
