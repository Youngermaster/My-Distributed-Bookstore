import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { orderAPI } from "@/lib/api";
import type {
  CreateOrderRequest,
  UpdateOrderStatusRequest,
} from "@/types/order";
import { toast } from "@/lib/toast";

/**
 * Hook for order operations
 * Provides order creation, listing, and management
 */
export function useOrderService() {
  const queryClient = useQueryClient();

  return {
    // Get single order
    useOrder: (id: string) =>
      useQuery({
        queryKey: ["order", id],
        queryFn: async () => {
          const response = await orderAPI.get(id);
          return response.data;
        },
        enabled: !!id,
      }),

    // List all orders (paginated)
    useOrders: (page = 1, pageSize = 20) =>
      useQuery({
        queryKey: ["orders", page, pageSize],
        queryFn: async () => {
          const response = await orderAPI.list(page, pageSize);
          return response.data;
        },
      }),

    // Get user's orders
    useUserOrders: (userId: string, page = 1, pageSize = 20) =>
      useQuery({
        queryKey: ["orders", "user", userId, page, pageSize],
        queryFn: async () => {
          const response = await orderAPI.getUserOrders(userId, page, pageSize);
          return response.data;
        },
        enabled: !!userId,
      }),

    // Create order
    useCreateOrder: () =>
      useMutation({
        mutationFn: (data: CreateOrderRequest) => orderAPI.create(data),
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["orders"] });
          toast.success("Order created successfully");
        },
        onError: () => {
          toast.error("Failed to create order");
        },
      }),

    // Update order status (admin only)
    useUpdateOrderStatus: () =>
      useMutation({
        mutationFn: ({
          id,
          data,
        }: {
          id: string;
          data: UpdateOrderStatusRequest;
        }) => orderAPI.updateStatus(id, data),
        onSuccess: (_, variables) => {
          queryClient.invalidateQueries({ queryKey: ["orders"] });
          queryClient.invalidateQueries({ queryKey: ["order", variables.id] });
          toast.success("Order status updated");
        },
        onError: () => {
          toast.error("Failed to update order status");
        },
      }),

    // Cancel order
    useCancelOrder: () =>
      useMutation({
        mutationFn: (id: string) => orderAPI.cancel(id),
        onSuccess: (_, id) => {
          queryClient.invalidateQueries({ queryKey: ["orders"] });
          queryClient.invalidateQueries({ queryKey: ["order", id] });
          toast.success("Order cancelled");
        },
        onError: () => {
          toast.error("Failed to cancel order");
        },
      }),
  };
}
