export interface Book {
  id: string;
  isbn: string;
  title: string;
  description: string;
  price: number;
  stock_quantity: number;
  publisher_id: string;
  publication_date: string;
  language: string;
  pages: number;
  format: string;
  cover_image_url?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
  authors?: Author[];
  categories?: Category[];
  publisher?: Publisher;
}

export interface Author {
  id: string;
  name: string;
  bio?: string;
  birth_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  parent_id?: string;
  created_at: string;
  updated_at: string;
}

export interface Publisher {
  id: string;
  name: string;
  country?: string;
  website?: string;
  created_at: string;
  updated_at: string;
}

export interface BookFilters {
  title?: string;
  author?: string;
  author_id?: string;
  category_id?: string;
  min_price?: number;
  max_price?: number;
  page?: number;
  page_size?: number;
  limit?: number;
  offset?: number;
}

export interface BooksResponse {
  books: Book[];
  total: number;
  page: number;
  page_size: number;
}
