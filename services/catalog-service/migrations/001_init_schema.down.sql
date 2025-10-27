-- Drop triggers
DROP TRIGGER IF EXISTS update_books_updated_at ON books;
DROP TRIGGER IF EXISTS update_authors_updated_at ON authors;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_publishers_updated_at ON publishers;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_books_isbn;
DROP INDEX IF EXISTS idx_books_title;
DROP INDEX IF EXISTS idx_books_price;
DROP INDEX IF EXISTS idx_books_stock;
DROP INDEX IF EXISTS idx_books_publisher_id;
DROP INDEX IF EXISTS idx_books_created_at;
DROP INDEX IF EXISTS idx_books_deleted_at;

DROP INDEX IF EXISTS idx_authors_name;
DROP INDEX IF EXISTS idx_authors_deleted_at;

DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_categories_deleted_at;

DROP INDEX IF EXISTS idx_publishers_name;
DROP INDEX IF EXISTS idx_publishers_deleted_at;

DROP INDEX IF EXISTS idx_book_authors_book_id;
DROP INDEX IF EXISTS idx_book_authors_author_id;

DROP INDEX IF EXISTS idx_book_categories_book_id;
DROP INDEX IF EXISTS idx_book_categories_category_id;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS book_categories;
DROP TABLE IF EXISTS book_authors;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS publishers;

-- Drop extensions
DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "uuid-ossp";
