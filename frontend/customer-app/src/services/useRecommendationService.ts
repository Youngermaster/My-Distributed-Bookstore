import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { recommendationsAPI } from "@/lib/api";
import type {
  TrendingBooksParams,
  PopularBooksParams,
  SimilarBooksParams,
  InteractionRequest,
} from "@/types/recommendation";

/**
 * Hook for recommendation operations
 * Provides personalized, trending, popular, and similar book recommendations
 */
export function useRecommendationService() {
  const queryClient = useQueryClient();

  return {
    // Get personalized recommendations
    usePersonalizedRecommendations: (
      params?: { limit?: number },
      enabled = true
    ) =>
      useQuery({
        queryKey: ["recommendations", "personalized", params],
        queryFn: async () => {
          const response = await recommendationsAPI.getPersonalized(params);
          return response.data;
        },
        enabled,
      }),

    // Get similar books
    useSimilarBooks: (bookId: string, params?: SimilarBooksParams) =>
      useQuery({
        queryKey: ["recommendations", "similar", bookId, params],
        queryFn: async () => {
          const response = await recommendationsAPI.getSimilar(bookId, params);
          return response.data;
        },
        enabled: !!bookId,
      }),

    // Get trending books
    useTrendingBooks: (params?: TrendingBooksParams) =>
      useQuery({
        queryKey: ["recommendations", "trending", params],
        queryFn: async () => {
          const response = await recommendationsAPI.getTrending(params);
          return response.data;
        },
      }),

    // Get popular books
    usePopularBooks: (params?: PopularBooksParams) =>
      useQuery({
        queryKey: ["recommendations", "popular", params],
        queryFn: async () => {
          const response = await recommendationsAPI.getPopular(params);
          return response.data;
        },
      }),

    // Track interaction
    useTrackInteraction: () =>
      useMutation({
        mutationFn: (data: InteractionRequest) =>
          recommendationsAPI.trackInteraction(data),
        onSuccess: () => {
          // Invalidate recommendation queries to get updated suggestions
          queryClient.invalidateQueries({ queryKey: ["recommendations"] });
        },
      }),
  };
}
