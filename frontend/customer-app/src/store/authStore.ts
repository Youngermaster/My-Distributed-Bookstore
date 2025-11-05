import { create } from "zustand";
import { authAPI } from "@/lib/api";
import type { User } from "@/types/user";
import type { LoginRequest, RegisterRequest } from "@/types/auth";

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;

  login: (credentials: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  loadUser: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: localStorage.getItem("token"),
  isAuthenticated: !!localStorage.getItem("token"),
  isLoading: false,
  error: null,

  login: async (credentials: LoginRequest) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authAPI.login(credentials);
      const { user, access_token, refresh_token } = response.data;

      localStorage.setItem("token", access_token);
      localStorage.setItem("refresh_token", refresh_token);

      set({
        user,
        token: access_token,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || "Login failed";
      set({ isLoading: false, error: errorMessage });
      throw error;
    }
  },

  register: async (data: RegisterRequest) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authAPI.register(data);
      const { user, access_token, refresh_token } = response.data;

      localStorage.setItem("token", access_token);
      localStorage.setItem("refresh_token", refresh_token);

      set({
        user,
        token: access_token,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.message || "Registration failed";
      set({ isLoading: false, error: errorMessage });
      throw error;
    }
  },

  logout: async () => {
    try {
      await authAPI.logout();
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      set({
        user: null,
        token: null,
        isAuthenticated: false,
        error: null,
      });
    }
  },

  loadUser: async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      set({ isAuthenticated: false, user: null });
      return;
    }

    set({ isLoading: true });
    try {
      const response = await authAPI.me();
      set({
        user: response.data,
        isAuthenticated: true,
        isLoading: false,
        error: null,
      });
    } catch (error) {
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      set({
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
      });
    }
  },

  clearError: () => set({ error: null }),
}));
