export interface Cart {
  id: string;
  user_id?: string;
  items: CartItem[];
  total: number;
  item_count: number;
  created_at: string;
  updated_at: string;
}

export interface CartItem {
  book_id: string;
  quantity: number;
  price: number;
  subtotal: number;
  added_at: string;
}

export interface AddToCartRequest {
  book_id: string;
  quantity: number;
  price: number;
}

export interface UpdateCartItemRequest {
  quantity: number;
}
