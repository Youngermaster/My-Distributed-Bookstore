import axios, { AxiosError } from "axios";
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  RefreshTokenResponse,
} from "@/types/auth";
import type { Book, BookFilters, BooksResponse, Category } from "@/types/book";
import type { User } from "@/types/user";
import type { WishlistItem, WishlistResponse } from "@/types/wishlist";
import type {
  Cart,
  AddToCartRequest,
  UpdateCartItemRequest,
} from "@/types/cart";
import type {
  Order,
  CreateOrderRequest,
  UpdateOrderStatusRequest,
  OrderListResponse,
} from "@/types/order";
import type {
  RecommendationResponse,
  InteractionRequest,
  InteractionResponse,
  TrendingBooksParams,
  PopularBooksParams,
  SimilarBooksParams,
} from "@/types/recommendation";
import type {
  DashboardStats,
  SalesAnalytics,
  InventoryReport,
  UserGrowthReport,
  SalesAnalyticsParams,
  InventoryReportParams,
  UserGrowthParams,
  TopBooksParams,
  TopBook,
} from "@/types/admin";

const api = axios.create({
  baseURL:
    import.meta.env.VITE_API_URL ||
    (typeof window !== "undefined" && window.location.port === "3000"
      ? "http://localhost:8080"
      : ""),
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor to add JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && originalRequest) {
      const refreshToken = localStorage.getItem("refresh_token");

      if (refreshToken) {
        try {
          const { data } = await axios.post<RefreshTokenResponse>(
            `${api.defaults.baseURL}/api/v1/users/auth/refresh`,
            { refresh_token: refreshToken }
          );

          localStorage.setItem("token", data.token);
          localStorage.setItem("refresh_token", data.refresh_token);

          originalRequest.headers.Authorization = `Bearer ${data.token}`;
          return api(originalRequest);
        } catch (refreshError) {
          localStorage.removeItem("token");
          localStorage.removeItem("refresh_token");
          window.location.href = "/login";
        }
      }
    }

    return Promise.reject(error);
  }
);

// Auth API (through API Gateway -> User Service)
export const authAPI = {
  login: (data: LoginRequest) =>
    api.post<AuthResponse>("/api/v1/users/auth/login", data),

  register: (data: RegisterRequest) =>
    api.post<AuthResponse>("/api/v1/users/auth/register", data),

  logout: () => api.post("/api/v1/users/auth/logout"),

  refresh: (refresh_token: string) =>
    api.post<RefreshTokenResponse>("/api/v1/users/auth/refresh", {
      refresh_token,
    }),

  me: () => api.get<{ data: User }>("/api/v1/users/me"),
};

// Books API (through API Gateway -> Catalog Service)
export const booksAPI = {
  list: (params?: BookFilters) =>
    api.get<BooksResponse>("/api/v1/catalog/books", { params }),

  get: (id: string) => api.get<Book>(`/api/v1/catalog/books/${id}`),

  search: (query: string, params?: BookFilters) =>
    api.get<BooksResponse>("/api/v1/catalog/books/search", {
      params: { q: query, ...params },
    }),

  create: (book: Partial<Book>) =>
    api.post<Book>("/api/v1/catalog/books", book),

  update: (id: string, book: Partial<Book>) =>
    api.put<Book>(`/api/v1/catalog/books/${id}`, book),

  delete: (id: string) => api.delete(`/api/v1/catalog/books/${id}`),
};

// Categories API (through API Gateway -> Catalog Service)
export const categoriesAPI = {
  list: () =>
    api.get<{ categories: Category[]; total: number }>(
      "/api/v1/catalog/categories"
    ),

  get: (id: string) => api.get<Category>(`/api/v1/catalog/categories/${id}`),
};

// Authors API (through API Gateway -> Catalog Service)
export const authorsAPI = {
  list: (page = 1, pageSize = 20) =>
    api.get("/api/v1/catalog/authors", {
      params: { page, page_size: pageSize },
    }),

  get: (id: string) => api.get(`/api/v1/catalog/authors/${id}`),
};

// Publishers API (through API Gateway -> Catalog Service)
export const publishersAPI = {
  list: (page = 1, pageSize = 20) =>
    api.get("/api/v1/catalog/publishers", {
      params: { page, page_size: pageSize },
    }),

  get: (id: string) => api.get(`/api/v1/catalog/publishers/${id}`),
};

// Wishlist API (through API Gateway -> User Service)
export const wishlistAPI = {
  list: () => api.get<WishlistResponse>("/api/v1/users/me/wishlist"),

  add: (book_id: string) =>
    api.post<{ data: WishlistItem }>("/api/v1/users/me/wishlist", { book_id }),

  remove: (book_id: string) =>
    api.delete(`/api/v1/users/me/wishlist/${book_id}`),

  clear: () => api.delete("/api/v1/users/me/wishlist"),

  check: (book_id: string) =>
    api.get<{ in_wishlist: boolean }>(
      `/api/v1/users/me/wishlist/check/${book_id}`
    ),
};

// Cart API (through API Gateway -> Cart Service)
export const cartAPI = {
  get: (cartId: string) => api.get<Cart>(`/api/v1/cart/${cartId}`),

  addItem: (cartId: string, data: AddToCartRequest) =>
    api.post<Cart>(`/api/v1/cart/${cartId}/items`, data),

  updateItem: (cartId: string, bookId: string, data: UpdateCartItemRequest) =>
    api.put<Cart>(`/api/v1/cart/${cartId}/items/${bookId}`, data),

  removeItem: (cartId: string, bookId: string) =>
    api.delete<Cart>(`/api/v1/cart/${cartId}/items/${bookId}`),

  clear: (cartId: string) =>
    api.delete<{ success: boolean; message: string }>(`/api/v1/cart/${cartId}`),
};

// Order API (through API Gateway -> Order Service)
export const orderAPI = {
  create: (data: CreateOrderRequest) => api.post<Order>("/api/v1/orders", data),

  get: (id: string) => api.get<Order>(`/api/v1/orders/${id}`),

  list: (page = 1, pageSize = 20) =>
    api.get<OrderListResponse>("/api/v1/orders", {
      params: { page, page_size: pageSize },
    }),

  getUserOrders: (userId: string, page = 1, pageSize = 20) =>
    api.get<OrderListResponse>(`/api/v1/users/${userId}/orders`, {
      params: { page, page_size: pageSize },
    }),

  updateStatus: (id: string, data: UpdateOrderStatusRequest) =>
    api.patch<Order>(`/api/v1/orders/${id}/status`, data),

  cancel: (id: string) =>
    api.post<{ success: boolean; message: string }>(
      `/api/v1/orders/${id}/cancel`
    ),
};

// Recommendations API (through API Gateway -> Recommendation Service)
export const recommendationsAPI = {
  getPersonalized: (params?: { limit?: number }) =>
    api.get<RecommendationResponse>("/api/v1/recommendations/me", { params }),

  getSimilar: (bookId: string, params?: SimilarBooksParams) =>
    api.get<RecommendationResponse>(
      `/api/v1/recommendations/similar/${bookId}`,
      { params }
    ),

  getTrending: (params?: TrendingBooksParams) =>
    api.get<RecommendationResponse>("/api/v1/recommendations/trending", {
      params,
    }),

  getPopular: (params?: PopularBooksParams) =>
    api.get<RecommendationResponse>("/api/v1/recommendations/popular", {
      params,
    }),

  trackInteraction: (data: InteractionRequest) =>
    api.post<InteractionResponse>("/api/v1/recommendations/interactions", data),
};

// Admin API (through API Gateway -> Admin Service)
export const adminAPI = {
  getDashboard: () =>
    api.get<{ success: boolean; data: DashboardStats }>(
      "/api/v1/admin/dashboard"
    ),

  getSalesAnalytics: (params?: SalesAnalyticsParams) =>
    api.get<{ success: boolean; data: SalesAnalytics }>(
      "/api/v1/admin/analytics/sales",
      { params }
    ),

  getInventoryReport: (params?: InventoryReportParams) =>
    api.get<{ success: boolean; data: InventoryReport }>(
      "/api/v1/admin/analytics/inventory",
      { params }
    ),

  getUserGrowth: (params?: UserGrowthParams) =>
    api.get<{ success: boolean; data: UserGrowthReport }>(
      "/api/v1/admin/analytics/users",
      { params }
    ),

  getTopBooks: (params?: TopBooksParams) =>
    api.get<{ success: boolean; data: { books: TopBook[] } }>(
      "/api/v1/admin/top-books",
      { params }
    ),
};

export default api;
