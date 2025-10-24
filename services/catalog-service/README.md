# Catalog Service

## Overview

The Catalog Service manages the book catalog, including book information, authors, categories, publishers, and search functionality. It's a core service in the distributed bookstore system.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Caching**: Redis 7
- **Messaging**: RabbitMQ
- **Protocol**: HTTP/REST + gRPC
- **Ports**:
  - HTTP: 8081
  - gRPC: 50051

## Responsibilities

- Book catalog management (CRUD operations)
- Advanced search and filtering
- Category hierarchy management
- Author and publisher management
- Book metadata handling
- Stock quantity tracking
- Price management
- Image URL management

## Database Schema

### Books
- id (UUID, PK)
- isbn (VARCHAR(13), UNIQUE)
- title (VARCHAR(500))
- description (TEXT)
- price (DECIMAL(10,2))
- stock_quantity (INTEGER)
- publisher_id (UUID, FK)
- cover_image_url (TEXT)
- metadata (JSONB)
- created_at, updated_at (TIMESTAMP)

### Authors
- id (UUID, PK)
- name (VARCHAR(255))
- bio (TEXT)
- birth_date (DATE)

### Categories
- id (UUID, PK)
- name (VARCHAR(100))
- slug (VARCHAR(100), UNIQUE)
- parent_id (UUID, FK) - for hierarchy

### Publishers
- id (UUID, PK)
- name (VARCHAR(255))
- country (VARCHAR(100))

### Join Tables
- book_authors (book_id, author_id, author_order)
- book_categories (book_id, category_id)

## REST API Endpoints

```
GET    /api/v1/books              # List books (with pagination)
GET    /api/v1/books/:id          # Get book details
GET    /api/v1/books/search?q=    # Search books
POST   /api/v1/books              # Create book (admin)
PUT    /api/v1/books/:id          # Update book (admin)
DELETE /api/v1/books/:id          # Delete book (admin)
GET    /api/v1/categories         # List categories (hierarchical)
GET    /api/v1/authors            # List authors
GET    /api/v1/publishers         # List publishers
```

## gRPC Methods

```protobuf
rpc GetBook(GetBookRequest) returns (GetBookResponse);
rpc SearchBooks(SearchBooksRequest) returns (SearchBooksResponse);
rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
rpc UpdateStock(UpdateStockRequest) returns (UpdateStockResponse);
rpc GetBooksByIds(GetBooksByIdsRequest) returns (GetBooksByIdsResponse);
```

## Events Published

- `catalog.book_created` - New book added to catalog
- `catalog.stock_updated` - Stock level changed
- `catalog.price_updated` - Price changed

## Events Consumed

- None (this is a core service)

## Environment Variables

```bash
# Server
HTTP_PORT=8081
GRPC_PORT=50051
ENV=development

# Database
DATABASE_URL=postgresql://bookstore:password@postgres:5432/catalog_db?sslmode=disable

# Redis
REDIS_URL=redis:6379
REDIS_PASSWORD=

# RabbitMQ
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/

# Observability
JAEGER_AGENT_HOST=jaeger
JAEGER_AGENT_PORT=6831
LOG_LEVEL=info
```

## Getting Started

### Local Development

```bash
# Run with Docker Compose
docker-compose up

# Run tests
go test ./... -v

# Build
go build -o bin/catalog-service cmd/server/main.go
```

## Next Steps

- [ ] Implement database models and migrations
- [ ] Create PostgreSQL repository layer
- [ ] Implement business logic in service layer
- [ ] Create HTTP REST handlers
- [ ] Implement gRPC server
- [ ] Add Redis caching for frequently accessed books
- [ ] Set up RabbitMQ event publishing
- [ ] Add full-text search functionality
- [ ] Implement pagination and filtering
- [ ] Add comprehensive tests
- [ ] Set up distributed tracing with Jaeger
