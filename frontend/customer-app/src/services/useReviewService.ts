import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { reviewAPI } from "@/lib/api";
import type { CreateReviewRequest, ReviewVoteRequest } from "@/types/review";

export function useReviewService() {
  const queryClient = useQueryClient();

  return {
    useBookReviews: (bookId: string, page = 1, pageSize = 10) =>
      useQuery({
        queryKey: ["reviews", bookId, page, pageSize],
        queryFn: async () => {
          const response = await reviewAPI.listForBook(bookId, page, pageSize);
          return response.data;
        },
        enabled: !!bookId,
      }),

    useBookReviewStats: (bookId: string) =>
      useQuery({
        queryKey: ["reviews", "stats", bookId],
        queryFn: async () => {
          const response = await reviewAPI.getStatsForBook(bookId);
          return response.data;
        },
        enabled: !!bookId,
      }),

    useCreateReview: () =>
      useMutation({
        mutationFn: (payload: CreateReviewRequest) => reviewAPI.create(payload),
        onSuccess: (_, variables) => {
          queryClient.invalidateQueries({
            queryKey: ["reviews", variables.book_id],
          });
          queryClient.invalidateQueries({
            queryKey: ["reviews", "stats", variables.book_id],
          });
        },
      }),

    useVoteReview: () =>
      useMutation({
        mutationFn: (variables: {
          reviewId: string;
          payload: ReviewVoteRequest;
          bookId: string;
        }) => reviewAPI.vote(variables.reviewId, variables.payload),
        onSuccess: (_, variables) => {
          queryClient.invalidateQueries({
            queryKey: ["reviews", variables.bookId],
          });
          queryClient.invalidateQueries({
            queryKey: ["reviews", "stats", variables.bookId],
          });
        },
      }),
  };
}
