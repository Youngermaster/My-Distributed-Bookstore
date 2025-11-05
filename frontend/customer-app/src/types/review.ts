export interface Review {
  id: string;
  book_id: string;
  user_id: string;
  rating: number;
  title: string;
  content: string;
  helpful_votes: number;
  verified_purchase: boolean;
  sentiment_score?: number | null;
  sentiment_label?: string | null;
  created_at: string;
  updated_at: string;
  user_name?: string | null;
}

export interface ReviewListResponse {
  reviews: Review[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ReviewStatsResponse {
  book_id: string;
  total_reviews: number;
  average_rating: number;
  rating_distribution: Record<number, number>;
  sentiment_distribution: Record<string, number>;
}

export interface CreateReviewRequest {
  book_id: string;
  user_id: string;
  rating: number;
  title: string;
  content: string;
  verified_purchase?: boolean;
}

export interface ReviewVoteRequest {
  user_id: string;
  is_helpful: boolean;
}
