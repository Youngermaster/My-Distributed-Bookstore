import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { booksAPI, categoriesAPI, authorsAPI, publishersAPI } from "@/lib/api";
import type { BookFilters } from "@/types/book";

/**
 * Hook for catalog/books operations
 * Provides books, categories, authors, publishers management
 */
export function useCatalogService() {
  const queryClient = useQueryClient();

  return {
    // Books
    useBooks: (params?: BookFilters) =>
      useQuery({
        queryKey: ["books", params],
        queryFn: async () => {
          const response = await booksAPI.list(params);
          return response.data;
        },
      }),

    useBook: (id: string) =>
      useQuery({
        queryKey: ["book", id],
        queryFn: async () => {
          const response = await booksAPI.get(id);
          return response.data;
        },
        enabled: !!id,
      }),

    useSearchBooks: (query: string, params?: BookFilters) =>
      useQuery({
        queryKey: ["books", "search", query, params],
        queryFn: async () => {
          const response = await booksAPI.search(query, params);
          return response.data;
        },
        enabled: !!query,
      }),

    // Categories
    useCategories: () =>
      useQuery({
        queryKey: ["categories"],
        queryFn: async () => {
          const response = await categoriesAPI.list();
          return response.data;
        },
      }),

    useCategory: (id: string) =>
      useQuery({
        queryKey: ["category", id],
        queryFn: async () => {
          const response = await categoriesAPI.get(id);
          return response.data;
        },
        enabled: !!id,
      }),

    // Authors
    useAuthors: (page = 1, pageSize = 20) =>
      useQuery({
        queryKey: ["authors", page, pageSize],
        queryFn: async () => {
          const response = await authorsAPI.list(page, pageSize);
          return response.data;
        },
      }),

    useAuthor: (id: string) =>
      useQuery({
        queryKey: ["author", id],
        queryFn: async () => {
          const response = await authorsAPI.get(id);
          return response.data;
        },
        enabled: !!id,
      }),

    // Publishers
    usePublishers: (page = 1, pageSize = 20) =>
      useQuery({
        queryKey: ["publishers", page, pageSize],
        queryFn: async () => {
          const response = await publishersAPI.list(page, pageSize);
          return response.data;
        },
      }),

    usePublisher: (id: string) =>
      useQuery({
        queryKey: ["publisher", id],
        queryFn: async () => {
          const response = await publishersAPI.get(id);
          return response.data;
        },
        enabled: !!id,
      }),

    // Admin operations (create, update, delete books)
    useCreateBook: () =>
      useMutation({
        mutationFn: booksAPI.create,
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["books"] });
        },
      }),

    useUpdateBook: () =>
      useMutation({
        mutationFn: ({ id, data }: { id: string; data: any }) =>
          booksAPI.update(id, data),
        onSuccess: (_, variables) => {
          queryClient.invalidateQueries({ queryKey: ["books"] });
          queryClient.invalidateQueries({ queryKey: ["book", variables.id] });
        },
      }),

    useDeleteBook: () =>
      useMutation({
        mutationFn: booksAPI.delete,
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["books"] });
        },
      }),
  };
}
