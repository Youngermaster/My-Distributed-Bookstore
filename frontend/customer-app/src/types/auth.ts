import type { User } from "./user";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  full_name: string;
  phone?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_at: string;
  user: User;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export type RefreshTokenResponse = AuthResponse;
