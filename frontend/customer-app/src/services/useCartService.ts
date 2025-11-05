import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { cartAPI } from "@/lib/api";
import { useCartStore } from "@/store/cartStore";
import type { AddToCartRequest, UpdateCartItemRequest } from "@/types/cart";
import { toast } from "@/lib/toast";

/**
 * Hook for cart operations
 * Provides cart management with automatic cart ID handling
 */
export function useCartService() {
  const queryClient = useQueryClient();
  const { cartId } = useCartStore();

  return {
    // Get cart
    useCart: () =>
      useQuery({
        queryKey: ["cart", cartId],
        queryFn: async () => {
          if (!cartId) return null;
          const response = await cartAPI.get(cartId);
          return response.data;
        },
        enabled: !!cartId,
      }),

    // Add item to cart
    useAddToCart: () =>
      useMutation({
        mutationFn: (data: AddToCartRequest) => {
          if (!cartId) throw new Error("No cart ID");
          return cartAPI.addItem(cartId, data);
        },
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
          toast.success("Added to cart");
        },
        onError: () => {
          toast.error("Failed to add to cart");
        },
      }),

    // Update cart item
    useUpdateCartItem: () =>
      useMutation({
        mutationFn: ({
          bookId,
          data,
        }: {
          bookId: string;
          data: UpdateCartItemRequest;
        }) => {
          if (!cartId) throw new Error("No cart ID");
          return cartAPI.updateItem(cartId, bookId, data);
        },
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
        },
        onError: () => {
          toast.error("Failed to update cart");
        },
      }),

    // Remove item from cart
    useRemoveFromCart: () =>
      useMutation({
        mutationFn: (bookId: string) => {
          if (!cartId) throw new Error("No cart ID");
          return cartAPI.removeItem(cartId, bookId);
        },
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
          toast.success("Removed from cart");
        },
        onError: () => {
          toast.error("Failed to remove item");
        },
      }),

    // Clear cart
    useClearCart: () =>
      useMutation({
        mutationFn: () => {
          if (!cartId) throw new Error("No cart ID");
          return cartAPI.clear(cartId);
        },
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["cart", cartId] });
          toast.success("Cart cleared");
        },
        onError: () => {
          toast.error("Failed to clear cart");
        },
      }),
  };
}
