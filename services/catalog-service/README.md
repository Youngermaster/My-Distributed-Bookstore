# Catalog Service

## Overview

The Catalog Service manages the book catalog, including book information, authors, categories, publishers, and search functionality. It's a core service in the distributed bookstore system built with Go, Fiber, GORM, and PostgreSQL.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Caching**: Redis 7 (optional)
- **Messaging**: RabbitMQ (optional)
- **Protocol**: HTTP/REST + gRPC
- **Ports**:
  - HTTP: 8081
  - gRPC: 50051

## Features

✅ **Implemented:**

- Complete CRUD operations for Books, Authors, Categories, and Publishers
- Advanced search and filtering
- Pagination support
- Hierarchical category management
- Many-to-many relationships (Books ↔ Authors, Books ↔ Categories)
- Auto-migrations with GORM
- Sample data seeding
- Docker support with Docker Compose
- Health check endpoints
- Comprehensive error handling

🚧 **Future:**

- gRPC implementation
- Redis caching
- RabbitMQ event publishing
- Full-text search with PostgreSQL
- Distributed tracing with Jaeger

## Quick Start

### Option 1: Docker Compose (Recommended)

The easiest way to run the service with all dependencies:

```bash
cd services/catalog-service

# Start all services (PostgreSQL + Catalog Service)
docker-compose up -d

# View logs
docker-compose logs -f catalog-service

# Stop all services
docker-compose down

# Reset database (removes all data)
docker-compose down -v
docker-compose up -d
```

**Service URLs:**

- API: http://localhost:8081/api/v1
- Health: http://localhost:8081/health
- PostgreSQL: localhost:5432

### Option 2: Local Development

1. **Install dependencies:**

```bash
go mod download
```

2. **Start PostgreSQL** (if not using Docker):

```bash
# Using Docker for PostgreSQL only
docker run --name catalog-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=catalog_db -p 5432:5432 -d postgres:15-alpine
```

3. **Set environment variables:**

```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=password
export DATABASE_NAME=catalog_db
export HTTP_PORT=8081
export ENV=development
```

4. **Run the service:**

```bash
go run cmd/server/main.go
```

The service will automatically:

- ✅ Connect to PostgreSQL
- ✅ Run database migrations
- ✅ Seed sample data (in development mode)
- ✅ Start HTTP server on port 8081

## Database Schema

### Books

```sql
id              UUID PRIMARY KEY
isbn            VARCHAR(13) UNIQUE
title           VARCHAR(500)
description     TEXT
price           DECIMAL(10,2)
stock_quantity  INTEGER
publisher_id    UUID (FK -> publishers)
cover_image_url TEXT
publication_date DATE
language        VARCHAR(50)
page_count      INTEGER
created_at      TIMESTAMP
updated_at      TIMESTAMP
```

### Authors

```sql
id         UUID PRIMARY KEY
name       VARCHAR(255)
bio        TEXT
birth_date DATE
country    VARCHAR(100)
image_url  TEXT
```

### Categories

```sql
id          UUID PRIMARY KEY
name        VARCHAR(100)
slug        VARCHAR(100) UNIQUE
description TEXT
parent_id   UUID (FK -> categories) -- For hierarchy
```

### Publishers

```sql
id          UUID PRIMARY KEY
name        VARCHAR(255)
country     VARCHAR(100)
website     VARCHAR(500)
description TEXT
```

### Join Tables

- `book_authors` (book_id, author_id, author_order)
- `book_categories` (book_id, category_id)

## API Documentation

Base URL: `http://localhost:8081/api/v1`

### 📚 Books API

#### List Books

```bash
GET /api/v1/books?page=1&page_size=20

# With filters
GET /api/v1/books?min_price=10&max_price=50&in_stock=true&sort_by=price&sort_order=asc
```

**Response:**

```json
{
  "books": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

#### Get Book by ID

```bash
GET /api/v1/books/{book_id}
```

**Response:**

```json
{
  "id": "uuid",
  "isbn": "9781492032649",
  "title": "Building Microservices",
  "description": "A comprehensive guide...",
  "price": 49.99,
  "stock_quantity": 50,
  "publisher_id": "uuid",
  "publisher": {
    "id": "uuid",
    "name": "O'Reilly Media",
    "country": "USA"
  },
  "authors": [
    {
      "id": "uuid",
      "name": "Martin Fowler",
      "bio": "Software development expert"
    }
  ],
  "categories": [
    {
      "id": "uuid",
      "name": "Software Architecture",
      "slug": "software-architecture"
    }
  ],
  "cover_image_url": "https://...",
  "publication_date": "2018-03-09T00:00:00Z",
  "language": "English",
  "page_count": 280,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Search Books

```bash
GET /api/v1/books/search?q=microservices&page=1&page_size=10
```

#### Create Book

```bash
POST /api/v1/books
Content-Type: application/json

{
  "isbn": "9781234567890",
  "title": "My New Book",
  "description": "An amazing book about...",
  "price": 29.99,
  "stock_quantity": 100,
  "publisher_id": "uuid-of-publisher",
  "cover_image_url": "https://example.com/cover.jpg",
  "publication_date": "2024-01-01T00:00:00Z",
  "language": "English",
  "page_count": 350,
  "author_ids": ["uuid-1", "uuid-2"],
  "category_ids": ["uuid-1", "uuid-2"]
}
```

**Example with curl:**

```bash
curl -X POST http://localhost:8081/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "isbn": "9781234567890",
    "title": "Distributed Systems Design",
    "description": "Learn how to build scalable distributed systems",
    "price": 59.99,
    "stock_quantity": 100,
    "language": "English",
    "page_count": 450
  }'
```

#### Update Book

```bash
PUT /api/v1/books/{book_id}
Content-Type: application/json

{
  "title": "Updated Title",
  "price": 39.99,
  "stock_quantity": 75
}
```

**Example:**

```bash
curl -X PUT http://localhost:8081/api/v1/books/{book-id} \
  -H "Content-Type: application/json" \
  -d '{
    "price": 44.99,
    "stock_quantity": 60
  }'
```

#### Delete Book

```bash
DELETE /api/v1/books/{book_id}
```

**Example:**

```bash
curl -X DELETE http://localhost:8081/api/v1/books/{book-id}
```

#### Update Book Stock

```bash
PATCH /api/v1/books/{book_id}/stock
Content-Type: application/json

{
  "quantity": 150
}
```

### 👤 Authors API

#### List Authors

```bash
GET /api/v1/authors?page=1&page_size=20
```

#### Get Author by ID

```bash
GET /api/v1/authors/{author_id}
```

#### Create Author

```bash
POST /api/v1/authors
Content-Type: application/json

{
  "name": "Jane Doe",
  "bio": "Award-winning author...",
  "birth_date": "1980-05-15T00:00:00Z",
  "country": "USA",
  "image_url": "https://example.com/author.jpg"
}
```

**Example:**

```bash
curl -X POST http://localhost:8081/api/v1/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Robert C. Martin",
    "bio": "Software engineer known as Uncle Bob",
    "country": "USA"
  }'
```

#### Update Author

```bash
PUT /api/v1/authors/{author_id}
Content-Type: application/json

{
  "bio": "Updated biography...",
  "image_url": "https://example.com/new-photo.jpg"
}
```

#### Delete Author

```bash
DELETE /api/v1/authors/{author_id}
```

### 📂 Categories API

#### List Categories

```bash
# Flat list
GET /api/v1/categories

# Hierarchical tree
GET /api/v1/categories?hierarchical=true
```

#### Get Category by ID

```bash
GET /api/v1/categories/{category_id}
```

#### Create Category

```bash
POST /api/v1/categories
Content-Type: application/json

{
  "name": "Science Fiction",
  "slug": "science-fiction",
  "description": "Futuristic and imaginative books",
  "parent_id": null
}
```

**Create subcategory:**

```bash
curl -X POST http://localhost:8081/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cyberpunk",
    "slug": "cyberpunk",
    "description": "High-tech, low-life sci-fi",
    "parent_id": "uuid-of-scifi-category"
  }'
```

#### Update Category

```bash
PUT /api/v1/categories/{category_id}
Content-Type: application/json

{
  "description": "Updated description"
}
```

#### Delete Category

```bash
DELETE /api/v1/categories/{category_id}
```

### 🏢 Publishers API

#### List Publishers

```bash
GET /api/v1/publishers?page=1&page_size=20
```

#### Get Publisher by ID

```bash
GET /api/v1/publishers/{publisher_id}
```

#### Create Publisher

```bash
POST /api/v1/publishers
Content-Type: application/json

{
  "name": "Tech Books Publishing",
  "country": "USA",
  "website": "https://techbooks.com",
  "description": "Leading technology book publisher"
}
```

**Example:**

```bash
curl -X POST http://localhost:8081/api/v1/publishers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Addison-Wesley",
    "country": "USA",
    "website": "https://www.awprofessional.com",
    "description": "Computer science and technology publisher"
  }'
```

#### Update Publisher

```bash
PUT /api/v1/publishers/{publisher_id}
Content-Type: application/json

{
  "website": "https://new-website.com"
}
```

#### Delete Publisher

```bash
DELETE /api/v1/publishers/{publisher_id}
```

## Complete End-to-End Example

Here's a complete workflow to create a book with all relationships:

```bash
# 1. Create a publisher
PUBLISHER_ID=$(curl -s -X POST http://localhost:8081/api/v1/publishers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Press",
    "country": "USA",
    "website": "https://techpress.com"
  }' | jq -r '.id')

# 2. Create authors
AUTHOR1_ID=$(curl -s -X POST http://localhost:8081/api/v1/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "bio": "Expert in distributed systems",
    "country": "USA"
  }' | jq -r '.id')

AUTHOR2_ID=$(curl -s -X POST http://localhost:8081/api/v1/authors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob Smith",
    "bio": "Cloud architecture specialist",
    "country": "UK"
  }' | jq -r '.id')

# 3. Create categories
CATEGORY1_ID=$(curl -s -X POST http://localhost:8081/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cloud Computing",
    "slug": "cloud-computing",
    "description": "Cloud platforms and services"
  }' | jq -r '.id')

CATEGORY2_ID=$(curl -s -X POST http://localhost:8081/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "System Design",
    "slug": "system-design",
    "description": "Designing scalable systems"
  }' | jq -r '.id')

# 4. Create the book with all relationships
curl -X POST http://localhost:8081/api/v1/books \
  -H "Content-Type: application/json" \
  -d "{
    \"isbn\": \"9781234567890\",
    \"title\": \"Mastering Cloud Architecture\",
    \"description\": \"A comprehensive guide to building scalable cloud applications\",
    \"price\": 79.99,
    \"stock_quantity\": 200,
    \"publisher_id\": \"$PUBLISHER_ID\",
    \"language\": \"English\",
    \"page_count\": 520,
    \"author_ids\": [\"$AUTHOR1_ID\", \"$AUTHOR2_ID\"],
    \"category_ids\": [\"$CATEGORY1_ID\", \"$CATEGORY2_ID\"]
  }" | jq

# 5. Search for the book
curl -s "http://localhost:8081/api/v1/books/search?q=Cloud" | jq

# 6. List books with filters
curl -s "http://localhost:8081/api/v1/books?min_price=50&max_price=100&in_stock=true" | jq
```

## Testing with Sample Data

The service automatically seeds sample data in development mode. After starting, you'll have:

- **3 Publishers**: O'Reilly Media, Manning Publications, Addison-Wesley
- **5 Authors**: Martin Fowler, Robert C. Martin, Eric Evans, Andrew S. Tanenbaum, Maarten van Steen
- **5 Categories**: Programming, Distributed Systems, Software Architecture, Databases, Cloud Computing
- **3 Books**:
  - Building Microservices
  - Clean Code
  - Distributed Systems: Principles and Paradigms

Try these queries:

```bash
# List all books
curl http://localhost:8081/api/v1/books | jq

# Search for "distributed"
curl "http://localhost:8081/api/v1/books/search?q=distributed" | jq

# Get all categories
curl http://localhost:8081/api/v1/categories | jq

# Get all authors
curl http://localhost:8081/api/v1/authors | jq
```

## Environment Variables

```bash
# Server Configuration
HTTP_PORT=8081              # HTTP server port
GRPC_PORT=50051            # gRPC server port (future)
ENV=development            # Environment: development, staging, production

# Database Configuration (Option 1: Individual parameters)
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=bookstore
DATABASE_PASSWORD=password
DATABASE_NAME=catalog_db
DATABASE_SSLMODE=disable

# Database Configuration (Option 2: Connection URL)
DATABASE_URL=postgresql://bookstore:password@localhost:5432/catalog_db?sslmode=disable

# Redis (Optional - for future caching)
REDIS_URL=redis:6379
REDIS_PASSWORD=

# RabbitMQ (Optional - for future event publishing)
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/

# Observability (Optional - for future monitoring)
JAEGER_AGENT_HOST=jaeger
JAEGER_AGENT_PORT=6831
LOG_LEVEL=info
```

## Project Structure

```
catalog-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go           # Configuration management
│   │   └── database.go         # Database connection & migrations
│   ├── domain/
│   │   └── models.go           # Domain entities (Book, Author, etc.)
│   ├── repository/
│   │   ├── book_repository.go
│   │   ├── author_repository.go
│   │   ├── category_repository.go
│   │   └── publisher_repository.go
│   ├── service/
│   │   ├── catalog_service.go  # Business logic
│   │   └── dto.go              # Data Transfer Objects
│   └── handler/
│       └── http/
│           ├── book_handler.go
│           ├── routes.go       # Route definitions
│           ├── responses.go
│           └── helpers.go
├── migrations/
│   ├── 001_init_schema.up.sql
│   └── 001_init_schema.down.sql
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Development

### Build

```bash
go build -o bin/catalog-service cmd/server/main.go
```

### Run Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/repository/... -v
```

### Database Migrations

Migrations are automatically applied on startup using GORM AutoMigrate. For manual migration management, SQL files are available in `migrations/`:

```bash
# Apply migrations manually (if needed)
psql $DATABASE_URL < migrations/001_init_schema.up.sql

# Rollback migrations
psql $DATABASE_URL < migrations/001_init_schema.down.sql
```

## Troubleshooting

### Service won't start

1. **Check PostgreSQL is running:**

```bash
docker ps | grep postgres
# or
pg_isready -h localhost -p 5432
```

2. **Check logs:**

```bash
docker-compose logs catalog-service
```

3. **Verify environment variables:**

```bash
# Make sure DATABASE_URL or individual DB parameters are set
echo $DATABASE_URL
```

### Database connection errors

```bash
# Reset the database
docker-compose down -v
docker-compose up -d postgres
# Wait a few seconds
docker-compose up -d catalog-service
```

### Port already in use

```bash
# Find process using port 8081
lsof -i :8081

# Kill the process
kill -9 <PID>

# Or change the port
export HTTP_PORT=8082
```

## Performance Tips

### Database Indexes

The migrations include indexes for:

- ISBN lookups
- Title searches (using trigram)
- Price range queries
- Stock availability
- Foreign key relationships

### Query Optimization

```bash
# Use pagination to limit results
GET /api/v1/books?page=1&page_size=20

# Filter early to reduce data
GET /api/v1/books?in_stock=true&min_price=10

# Sort efficiently
GET /api/v1/books?sort_by=price&sort_order=asc
```

## Next Steps

- [x] ✅ Implement database models and migrations
- [x] ✅ Create PostgreSQL repository layer
- [x] ✅ Implement business logic in service layer
- [x] ✅ Create HTTP REST handlers
- [x] ✅ Implement pagination and filtering
- [ ] Implement gRPC server
- [ ] Add Redis caching for frequently accessed books
- [ ] Set up RabbitMQ event publishing
- [ ] Add full-text search with PostgreSQL
- [ ] Add comprehensive unit tests
- [ ] Add integration tests
- [ ] Set up distributed tracing with Jaeger
- [ ] Add Prometheus metrics
- [ ] Implement rate limiting
- [ ] Add API documentation with Swagger

## License

Part of the Distributed Bookstore project.
