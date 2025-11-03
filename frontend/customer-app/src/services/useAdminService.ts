import { useQuery } from "@tanstack/react-query";
import { adminAPI } from "@/lib/api";
import type {
  SalesAnalyticsParams,
  InventoryReportParams,
  UserGrowthParams,
  TopBooksParams,
} from "@/types/admin";

/**
 * Hook for admin operations
 * Provides dashboard stats, analytics, and reports
 * Requires admin role
 */
export function useAdminService() {
  return {
    // Dashboard stats
    useDashboard: (enabled = true) =>
      useQuery({
        queryKey: ["admin", "dashboard"],
        queryFn: async () => {
          const response = await adminAPI.getDashboard();
          return response.data.data;
        },
        enabled,
        staleTime: 1000 * 60, // 1 minute
      }),

    // Sales analytics
    useSalesAnalytics: (params?: SalesAnalyticsParams, enabled = true) =>
      useQuery({
        queryKey: ["admin", "sales", params],
        queryFn: async () => {
          const response = await adminAPI.getSalesAnalytics(params);
          return response.data.data;
        },
        enabled,
        staleTime: 1000 * 60 * 5, // 5 minutes
      }),

    // Inventory report
    useInventoryReport: (params?: InventoryReportParams, enabled = true) =>
      useQuery({
        queryKey: ["admin", "inventory", params],
        queryFn: async () => {
          const response = await adminAPI.getInventoryReport(params);
          return response.data.data;
        },
        enabled,
        staleTime: 1000 * 60 * 5, // 5 minutes
      }),

    // User growth report
    useUserGrowth: (params?: UserGrowthParams, enabled = true) =>
      useQuery({
        queryKey: ["admin", "userGrowth", params],
        queryFn: async () => {
          const response = await adminAPI.getUserGrowth(params);
          return response.data.data;
        },
        enabled,
        staleTime: 1000 * 60 * 5, // 5 minutes
      }),

    // Top books
    useTopBooks: (params?: TopBooksParams, enabled = true) =>
      useQuery({
        queryKey: ["admin", "topBooks", params],
        queryFn: async () => {
          const response = await adminAPI.getTopBooks(params);
          return response.data.data.books;
        },
        enabled,
        staleTime: 1000 * 60 * 5, // 5 minutes
      }),
  };
}
