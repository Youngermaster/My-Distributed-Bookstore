import { useQuery } from "@tanstack/react-query";
import { categoriesAPI } from "@/lib/api";
import { useNavigate } from "react-router-dom";
import GenreCard from "@/components/GenreCard";
import { Loader2 } from "lucide-react";

export default function Genres() {
  const navigate = useNavigate();

  const { data, isLoading, error } = useQuery({
    queryKey: ["categories"],
    queryFn: async () => {
      const response = await categoriesAPI.list();
      return response.data;
    },
  });

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-semibold text-gray-700 mb-2">
            Failed to load genres
          </h2>
          <p className="text-gray-500">
            {error instanceof Error ? error.message : "An error occurred"}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            Browse by Genre
          </h1>
          <p className="text-lg text-gray-600">
            Explore our collection of books organized by genre
          </p>
        </div>
      </div>

      {/* Genres Grid */}
      <div className="max-w-7xl mx-auto px-4 py-12">
        {isLoading ? (
          <div className="flex items-center justify-center py-20">
            <Loader2 className="h-12 w-12 animate-spin text-gray-400" />
          </div>
        ) : data?.categories && data.categories.length > 0 ? (
          <>
            <div className="mb-6">
              <p className="text-gray-600">
                Found {data.total} {data.total === 1 ? "genre" : "genres"}
              </p>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {data.categories.map((category) => (
                <GenreCard
                  key={category.id}
                  category={category}
                  onClick={() => navigate(`/genres/${category.slug}`)}
                />
              ))}
            </div>
          </>
        ) : (
          <div className="text-center py-20">
            <div className="text-6xl mb-4">📚</div>
            <h3 className="text-xl font-semibold text-gray-700 mb-2">
              No genres found
            </h3>
            <p className="text-gray-500">
              Check back later for new genres
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
