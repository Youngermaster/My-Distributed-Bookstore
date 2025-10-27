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

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080",
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
            "http://localhost:8082/api/v1/auth/refresh",
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

// Auth API
export const authAPI = {
  login: (data: LoginRequest) =>
    api.post<AuthResponse>("/api/v1/auth/login", data),

  register: (data: RegisterRequest) =>
    api.post<AuthResponse>("/api/v1/auth/register", data),

  logout: () => api.post("/api/v1/auth/logout"),

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

// Wishlist API
export const wishlistAPI = {
  list: () => api.get<WishlistResponse>("/api/v1/users/me/wishlist"),

  add: (book_id: string) =>
    api.post<{ data: WishlistItem }>("/api/v1/users/me/wishlist", { book_id }),

  remove: (book_id: string) =>
    api.delete(`/api/v1/users/me/wishlist/${book_id}`),
};

export default api;
