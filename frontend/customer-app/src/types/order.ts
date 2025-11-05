export type OrderStatus =
  | "pending"
  | "confirmed"
  | "processing"
  | "shipped"
  | "delivered"
  | "cancelled";

export interface Order {
  id: string;
  user_id: string;
  status: OrderStatus;
  items: OrderItem[];
  total_amount: number;
  item_count: number;
  shipping_address?: string;
  payment_method?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  id: string;
  order_id: string;
  book_id: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  created_at: string;
}

export interface CreateOrderRequest {
  user_id: string;
  items: CreateOrderItemRequest[];
  shipping_address?: string;
  payment_method?: string;
  notes?: string;
}

export interface CreateOrderItemRequest {
  book_id: string;
  quantity: number;
  unit_price: number;
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus;
}

export interface OrderListResponse {
  orders: Order[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}
