import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { wishlistAPI } from "@/lib/api";
import { toast } from "@/lib/toast";

/**
 * Hook for wishlist operations
 * Provides wishlist management for authenticated users
 */
export function useWishlistService() {
  const queryClient = useQueryClient();

  return {
    // Get wishlist
    useWishlist: (enabled = true) =>
      useQuery({
        queryKey: ["wishlist"],
        queryFn: async () => {
          const response = await wishlistAPI.list();
          return response.data;
        },
        enabled,
      }),

    // Check if book is in wishlist
    useIsInWishlist: (bookId: string) =>
      useQuery({
        queryKey: ["wishlist", "check", bookId],
        queryFn: async () => {
          const response = await wishlistAPI.check(bookId);
          return response.data.in_wishlist;
        },
        enabled: !!bookId,
      }),

    // Add to wishlist
    useAddToWishlist: () =>
      useMutation({
        mutationFn: (bookId: string) => wishlistAPI.add(bookId),
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["wishlist"] });
          toast.success("Added to wishlist");
        },
        onError: () => {
          toast.error("Failed to add to wishlist");
        },
      }),

    // Remove from wishlist
    useRemoveFromWishlist: () =>
      useMutation({
        mutationFn: (bookId: string) => wishlistAPI.remove(bookId),
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["wishlist"] });
          toast.success("Removed from wishlist");
        },
        onError: () => {
          toast.error("Failed to remove from wishlist");
        },
      }),

    // Clear wishlist
    useClearWishlist: () =>
      useMutation({
        mutationFn: () => wishlistAPI.clear(),
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["wishlist"] });
          toast.success("Wishlist cleared");
        },
        onError: () => {
          toast.error("Failed to clear wishlist");
        },
      }),
  };
}
