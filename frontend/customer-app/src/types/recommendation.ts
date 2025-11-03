export interface RecommendationItem {
  book_id: string;
  score: number;
  reason?: string;
}

export interface RecommendationResponse {
  user_id: string;
  recommendations: RecommendationItem[];
  algorithm: string;
  total: number;
}

export interface InteractionRequest {
  book_id: string;
  interaction_type: "view" | "add_to_cart" | "purchase" | "review" | "wishlist";
}

export interface InteractionResponse {
  success: boolean;
  message: string;
}

export interface TrendingBooksParams {
  limit?: number;
  days?: number;
}

export interface PopularBooksParams {
  limit?: number;
}

export interface SimilarBooksParams {
  limit?: number;
}
