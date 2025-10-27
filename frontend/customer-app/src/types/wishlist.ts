import type { Book } from "./book";

export interface WishlistItem {
  id: string;
  user_id: string;
  book_id: string;
  book?: Book;
  created_at: string;
}

export interface WishlistResponse {
  data: WishlistItem[];
  total: number;
}
