# Distributed Bookstore - Architecture & Implementation Guide

**Last Updated:** November 5, 2025  
**Status:** Production Deployment on AWS EKS

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Service Implementation Details](#service-implementation-details)
3. [Go Services Architecture Pattern](#go-services-architecture-pattern)
4. [Python Services Architecture Pattern](#python-services-architecture-pattern)
5. [API Gateway Integration](#api-gateway-integration)
6. [Frontend Architecture](#frontend-architecture)
7. [Database Design](#database-design)
8. [Testing Guide](#testing-guide)
9. [Deployment Architecture](#deployment-architecture)

---

## System Architecture

### Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      Customer Browser                            │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                    ┌───────────▼────────────┐
                    │  Frontend (React 19)   │
                    │  Port: 5173 (dev)      │
                    │  AWS ELB (production)  │
                    └───────────┬────────────┘
                                │
                    ┌───────────▼─────────────┐
                    │  API Gateway (Go)       │
                    │  Port: 8080             │
                    │  Routes to all services │
                    └───┬──────────────┬──────┘
                        │              │
        ┌───────────────┼──────────────┼───────────────────┐
        │               │              │                   │
   ┌────▼────┐    ┌────▼────┐    ┌───▼────┐    ┌────────▼────────┐
   │ Catalog │    │  User   │    │  Cart  │    │  Order          │
   │ :8081   │    │  :8082  │    │  :8083 │    │  :8084          │
   │ Go      │    │  Go     │    │  Go    │    │  Go             │
   └────┬────┘    └────┬────┘    └───┬────┘    └────┬────────────┘
        │              │              │              │
   ┌────▼────┐    ┌────▼────┐    ┌───▼────┐    ┌───▼─────┐
   │catalog  │    │ userdb  │    │ Redis  │    │ orderdb │
   │   _db   │    │         │    │  :6379 │    │         │
   └─────────┘    └─────────┘    └────────┘    └─────────┘

   ┌──────────────────────────────────────────────────────┐
   │           Python Services (FastAPI)                  │
   ├──────────────┬──────────────┬──────────────────────┬─┤
   │  Inventory   │   Review     │  Recommendation      │ │
   │  :8086       │   :8088      │  :8089               │ │
   │  Python      │   Python     │  Python              │ │
   └──────┬───────┴──────┬───────┴──────────┬───────────┘
          │              │                  │
   ┌──────▼──────┐┌──────▼──────┐   ┌──────▼──────┐
   │inventory_db ││ reviews_db  │   │   (Uses     │
   │             ││             │   │  others)    │
   └─────────────┘└─────────────┘   └─────────────┘

   ┌──────────────────────────────────────────────────────┐
   │      TypeScript Services (Express/Node.js)           │
   ├──────────────┬────────────────────────────────────┬──┤
   │  Payment     │       Notification                 │  │
   │  :8085       │       :8087                        │  │
   │  TypeScript  │       TypeScript                   │  │
   └──────────────┴────────────────────────────────────┴──┘
```

![Diagrama Arquitectura](image.png)



### Service Inventory

| Service | Port | Language | Status | Database | Purpose |
|---------|------|----------|--------|----------|---------|
| Frontend | 5173 | React 19 + TS | ✅ Deployed | - | Customer UI |
| API Gateway | 8080 | Go | ✅ Deployed | - | Unified entry point |
| Catalog | 8081 | Go | ✅ Deployed | catalog_db | Books, authors, categories |
| User | 8082 | Go | ✅ Deployed | userdb | Auth, profiles, wishlist |
| Cart | 8083 | Go | ✅ Deployed | Redis | Shopping cart |
| Order | 8084 | Go | ✅ Deployed | orderdb | Order processing |
| Admin | 8090 | Go | ✅ Deployed | bookstore | Analytics, dashboard |
| Inventory | 8086 | Python | ✅ Deployed | inventory_db | Stock management |
| Review | 8088 | Python | ✅ Deployed | reviews_db | ML sentiment analysis |
| Recommendation | 8089 | Python | ✅ Deployed | - | Personalized recommendations |
| Payment | 8085 | TypeScript | ⚠️ Scaffold | - | Stripe integration |
| Notification | 8087 | TypeScript | ⚠️ Scaffold | - | Email/SMS notifications |

**Production Status:** 10/12 services fully implemented and deployed on AWS EKS

---

## Service Implementation Details

### Catalog Service (Go) - Port 8081

**Purpose:** Complete book catalog management

**Features:**
- ✅ Full CRUD for Books, Authors, Categories, Publishers
- ✅ Advanced filtering (by category, author, publisher, price range)
- ✅ Search functionality with partial matching
- ✅ Stock management
- ✅ Pagination support
- ✅ Many-to-many relationships (books ↔ authors, books ↔ categories)
- ✅ Database seeding for development

**API Endpoints:**
```
Books:
  GET    /api/v1/books              # List with filters (page, limit, category, author)
  GET    /api/v1/books/search       # Search books by title/ISBN
  GET    /api/v1/books/:id          # Get by ID
  POST   /api/v1/books              # Create book
  PUT    /api/v1/books/:id          # Update book
  DELETE /api/v1/books/:id          # Delete book
  PATCH  /api/v1/books/:id/stock    # Update stock quantity

Authors:
  GET    /api/v1/authors
  GET    /api/v1/authors/:id
  POST   /api/v1/authors
  PUT    /api/v1/authors/:id
  DELETE /api/v1/authors/:id

Categories, Publishers: Similar CRUD patterns
```

**Database Schema:**
- `books` - Core book information with ISBN, title, price, stock
- `authors` - Author information with biography
- `categories` - Book categories/genres
- `publishers` - Publisher information
- `book_authors` - Many-to-many join table
- `book_categories` - Many-to-many join table

---

### User Service (Go) - Port 8082

**Purpose:** Authentication and user management

**Features:**
- ✅ User registration with email validation
- ✅ JWT-based authentication (access + refresh tokens)
- ✅ Password hashing with bcrypt
- ✅ Role-based access control (RBAC)
- ✅ User profiles and addresses
- ✅ Wishlist management
- ✅ Session tracking

**API Endpoints:**
```
Authentication:
  POST   /api/v1/auth/register      # Register new user
  POST   /api/v1/auth/login         # Login (returns JWT)
  POST   /api/v1/auth/refresh       # Refresh JWT token
  POST   /api/v1/auth/logout        # Logout (protected)

User Profile:
  GET    /api/v1/users/me           # Get current user (protected)
  PUT    /api/v1/users/me           # Update profile (protected)

Wishlist:
  GET    /api/v1/users/me/wishlist           # Get wishlist
  POST   /api/v1/users/me/wishlist           # Add to wishlist
  DELETE /api/v1/users/me/wishlist/:bookId   # Remove from wishlist
  DELETE /api/v1/users/me/wishlist           # Clear wishlist
  GET    /api/v1/users/me/wishlist/check/:bookId  # Check if in wishlist
```

**Security Features:**
- Password strength requirements (8+ chars, uppercase, lowercase, number, special)
- JWT secret key: "JWT_SECRET" (configurable via environment)
- Token expiration: 24 hours (access), 7 days (refresh)
- Protected routes via middleware

---

### Cart Service (Go) - Port 8083

**Purpose:** Shopping cart management with Redis caching

**Features:**
- ✅ Redis-based cart storage (fast, ephemeral)
- ✅ Session-based carts for anonymous users
- ✅ User-based carts for authenticated users
- ✅ Cart item management (add, update quantity, remove)
- ✅ Automatic cart totals calculation
- ✅ Cart merging on login

**Storage Strategy:**
- Anonymous users: Redis with session ID (TTL: 7 days)
- Authenticated users: Redis with user ID
- On login: Merge anonymous cart → user cart

---

### Order Service (Go) - Port 8084

**Purpose:** Order processing and management

**Features:**
- ✅ Order creation from cart
- ✅ Order status tracking (pending, confirmed, shipped, delivered, cancelled)
- ✅ Order history
- ✅ Order cancellation
- ✅ Integration with inventory and payment services

**Database Schema:**
- `orders` - Order header (user, total, status, shipping address)
- `order_items` - Order line items (book, quantity, price)
- `order_status_history` - Audit trail of status changes

---

### Admin Service (Go) - Port 8090

**Purpose:** Analytics and administrative dashboard

**Features:**
- ✅ Sales analytics (total revenue, order count)
- ✅ Top-selling books
- ✅ User statistics
- ✅ Inventory alerts (low stock)
- ✅ JWT authentication required

**API Endpoints:**
```
Analytics:
  GET    /api/v1/admin/analytics/sales           # Sales summary
  GET    /api/v1/admin/analytics/top-books       # Best sellers
  GET    /api/v1/admin/analytics/users           # User stats
  GET    /api/v1/admin/inventory/low-stock       # Stock alerts
```

**Database:** `bookstore` (single database for analytics queries)

---

### Inventory Service (Python/FastAPI) - Port 8086

**Purpose:** Stock tracking and reservation management

**Features:**
- ✅ Stock level tracking per book
- ✅ Reservation system with automatic expiry
- ✅ Background task for reservation cleanup (runs every 60s)
- ✅ Stock adjustment operations (add, subtract, set)
- ✅ Low stock alerts with configurable thresholds
- ✅ Complete audit trail via stock movements

**API Endpoints:**
```http
POST   /api/v1/inventory                     # Create inventory for book
GET    /api/v1/inventory/{book_id}           # Get stock level
POST   /api/v1/inventory/{book_id}/adjust    # Adjust stock (add/subtract)
GET    /api/v1/inventory/low-stock           # Low stock items
POST   /api/v1/inventory/reserve             # Reserve stock for order
POST   /api/v1/inventory/release/{order_id}  # Release reservation
POST   /api/v1/inventory/commit/{order_id}   # Commit reservation (after payment)
GET    /api/v1/inventory/{book_id}/movements # Stock movement history
GET    /health                                # Health check
```

**Database Schema:**
- `inventory` - Stock levels (book_id, quantity, reorder_level)
- `reservations` - Active reservations with expiry (TTL: 15 minutes)
- `stock_movements` - Audit trail (operation, quantity, reason)

**Background Tasks:**
- Reservation expiry checker (every 60 seconds) - automatically releases expired reservations

---

### Review Service (Python/FastAPI) - Port 8088

**Purpose:** Book reviews with ML-powered sentiment analysis

**Features:**
- ✅ Full CRUD operations for book reviews
- ✅ **Automatic ML sentiment analysis** (TextBlob + NLTK)
- ✅ Review voting system (helpful/not helpful)
- ✅ Review statistics and aggregations
- ✅ Rating distribution analytics
- ✅ Async database operations (SQLAlchemy 2.0)

**API Endpoints:**
```http
POST   /api/v1/reviews                      # Create review (auto-sentiment)
GET    /api/v1/reviews/{review_id}          # Get review
PUT    /api/v1/reviews/{review_id}          # Update review
DELETE /api/v1/reviews/{review_id}          # Delete review
GET    /api/v1/reviews/book/{book_id}       # List reviews (paginated)
GET    /api/v1/reviews/book/{book_id}/stats # Get statistics (avg rating, count)
POST   /api/v1/reviews/{review_id}/vote     # Vote on review (helpful/not)
GET    /health                              # Health check
```

**ML Features:**
- **Sentiment Classification:** positive/neutral/negative
- **Sentiment Score:** -1.0 (very negative) to 1.0 (very positive)
- **Technology:** TextBlob for polarity analysis, NLTK for NLP
- **Auto-download:** NLTK data packages on first startup

**Database Schema:**
- `reviews` - Review content with sentiment (unique constraint: book_id + user_id)
- `review_votes` - User votes on review helpfulness

---

### Recommendation Service (Python/FastAPI) - Port 8089

**Purpose:** Personalized book recommendations using multiple strategies

**Features:**
- ✅ Multiple recommendation algorithms
- ✅ User interaction tracking
- ✅ Trending and popular book lists
- ✅ Similar book suggestions

**API Endpoints:**
```http
GET    /api/v1/recommendations/me                # Personalized recommendations
GET    /api/v1/recommendations/similar/{bookId}  # Books similar to this one
GET    /api/v1/recommendations/trending          # Trending books (7 days default)
GET    /api/v1/recommendations/popular           # Popular books (all time)
POST   /api/v1/recommendations/interactions      # Track user interaction
GET    /health                                   # Health check
```

**Recommendation Strategies:**
1. **Tag-based:** Recommends books with similar categories/genres
2. **Collaborative Filtering:** "Users who bought this also bought..."
3. **Popularity-based:** Trending and best-sellers
4. **Hybrid:** Combines multiple strategies

**Tracked Interactions:**
- `view` - User viewed book details
- `add_to_cart` - Added to shopping cart
- `purchase` - Completed purchase
- `review` - Left a review
- `wishlist` - Added to wishlist

---

## Go Services Architecture Pattern

All Go services follow Clean Architecture with this standard structure:

```
service-name/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point, wires dependencies
├── internal/
│   ├── config/                        # Configuration management
│   │   ├── config.go                  # Load env vars, settings
│   │   └── database.go                # DB initialization, migrations, seeding
│   ├── domain/                        # Domain entities (GORM models)
│   │   └── models.go                  # All database models
│   ├── repository/                    # Data access layer (one file per entity)
│   │   ├── book_repository.go
│   │   ├── author_repository.go
│   │   └── category_repository.go
│   ├── service/                       # Business logic layer
│   │   ├── catalog_service.go         # Service implementation
│   │   └── dto.go                     # Request/Response DTOs
│   ├── handler/                       # HTTP handlers
│   │   └── http/
│   │       ├── routes.go              # Route definitions
│   │       ├── book_handler.go        # Handler per entity
│   │       ├── responses.go           # Common response structures
│   │       └── helpers.go             # Helper functions
│   └── middleware/                    # Custom middleware
│       └── auth.go                    # JWT authentication
├── pkg/                               # Shared utilities (reusable)
│   ├── jwt/                           # JWT utilities
│   └── password/                      # Password hashing
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### Layer Responsibilities

#### 1. Domain Layer (`internal/domain/`)
```go
type Book struct {
    ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    ISBN          string     `gorm:"unique;not null"`
    Title         string     `gorm:"not null"`
    Price         float64    `gorm:"not null"`
    StockQuantity int        `gorm:"default:0"`
    Authors       []Author   `gorm:"many2many:book_authors;"`
    Categories    []Category `gorm:"many2many:book_categories;"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

#### 2. Repository Layer (`internal/repository/`)
```go
type BookRepository interface {
    Create(ctx context.Context, book *domain.Book) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
    List(ctx context.Context, filters map[string]interface{}) ([]*domain.Book, error)
    Update(ctx context.Context, book *domain.Book) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type bookRepository struct {
    db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
    return &bookRepository{db: db}
}
```

#### 3. Service Layer (`internal/service/`)
```go
type CatalogService interface {
    CreateBook(ctx context.Context, req CreateBookRequest) (*domain.Book, error)
    GetBook(ctx context.Context, id uuid.UUID) (*domain.Book, error)
    ListBooks(ctx context.Context, filters ListBooksRequest) ([]*domain.Book, error)
}

type catalogService struct {
    bookRepo      repository.BookRepository
    authorRepo    repository.AuthorRepository
    categoryRepo  repository.CategoryRepository
}

func NewCatalogService(bookRepo repository.BookRepository, ...) CatalogService {
    return &catalogService{
        bookRepo:     bookRepo,
        authorRepo:   authorRepo,
        categoryRepo: categoryRepo,
    }
}
```

#### 4. Handler Layer (`internal/handler/http/`)
```go
type BookHandler struct {
    service service.CatalogService
}

func NewBookHandler(service service.CatalogService) *BookHandler {
    return &BookHandler{service: service}
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
    var req dto.CreateBookRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{Message: "Invalid request"})
    }
    
    book, err := h.service.CreateBook(c.Context(), req)
    if err != nil {
        return c.Status(500).JSON(ErrorResponse{Message: err.Error()})
    }
    
    return c.Status(201).JSON(SuccessResponse{Data: book})
}
```

### Tech Stack (Go Services)

- **Go:** 1.21+
- **Web Framework:** Fiber v2 (Express-like API)
- **ORM:** GORM v2
- **Database:** PostgreSQL 15 with asyncpg driver
- **Authentication:** JWT (golang-jwt/jwt)
- **Password:** bcrypt
- **Logging:** zerolog
- **Validation:** go-playground/validator

---

## Python Services Architecture Pattern

All Python services follow FastAPI best practices with async/await:

```
service-name/
├── app/
│   ├── main.py                        # FastAPI application entry point
│   ├── api/
│   │   └── v1/
│   │       └── endpoints/
│   │           └── resource.py        # API endpoints (one file per resource)
│   ├── core/
│   │   └── config.py                  # Pydantic Settings
│   ├── db/
│   │   └── base.py                    # Async SQLAlchemy setup
│   ├── models/
│   │   └── resource.py                # SQLAlchemy models
│   ├── schemas/
│   │   └── resource.py                # Pydantic request/response schemas
│   ├── services/
│   │   └── resource_service.py        # Business logic
│   ├── tasks/                         # Background tasks (optional)
│   │   └── reservation_expiry.py
│   └── ml/                            # ML modules (optional)
│       └── sentiment.py
├── requirements.txt                   # Python dependencies
├── Dockerfile                         # Multi-stage build
├── docker-compose.yml
└── README.md
```

### Layer Responsibilities

#### 1. Models Layer (`app/models/`)
```python
from sqlalchemy import Column, String, Integer, Float, DateTime, Text
from sqlalchemy.dialects.postgresql import UUID
from app.db.base import Base
import uuid

class Review(Base):
    __tablename__ = "reviews"
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    book_id = Column(UUID(as_uuid=True), nullable=False)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    rating = Column(Integer, nullable=False)  # 1-5
    title = Column(String(200))
    content = Column(Text)
    sentiment_score = Column(Float)  # -1.0 to 1.0
    sentiment_label = Column(String(20))  # positive/neutral/negative
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, onupdate=func.now())
```

#### 2. Schemas Layer (`app/schemas/`)
```python
from pydantic import BaseModel, Field, validator
from uuid import UUID
from datetime import datetime

class ReviewCreateRequest(BaseModel):
    book_id: UUID
    user_id: UUID
    rating: int = Field(..., ge=1, le=5)
    title: str = Field(..., max_length=200)
    content: str
    verified_purchase: bool = False

class ReviewResponse(BaseModel):
    id: UUID
    book_id: UUID
    user_id: UUID
    rating: int
    title: str
    content: str
    sentiment_score: float
    sentiment_label: str
    helpful_votes: int
    created_at: datetime
    updated_at: datetime | None
    
    class Config:
        from_attributes = True
```

#### 3. Service Layer (`app/services/`)
```python
from typing import List, Optional
from sqlalchemy.ext.asyncio import AsyncSession
from app.models.review import Review
from app.schemas.review import ReviewCreateRequest
from app.ml.sentiment import analyze_sentiment

class ReviewService:
    def __init__(self, db: AsyncSession):
        self.db = db
    
    async def create_review(self, data: ReviewCreateRequest) -> Review:
        # Analyze sentiment
        sentiment_result = analyze_sentiment(data.content)
        
        review = Review(
            book_id=data.book_id,
            user_id=data.user_id,
            rating=data.rating,
            title=data.title,
            content=data.content,
            sentiment_score=sentiment_result["score"],
            sentiment_label=sentiment_result["label"],
            verified_purchase=data.verified_purchase
        )
        
        self.db.add(review)
        await self.db.commit()
        await self.db.refresh(review)
        return review
```

#### 4. API Endpoints Layer (`app/api/v1/endpoints/`)
```python
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from typing import Annotated
from app.db.base import get_db
from app.services.review_service import ReviewService
from app.schemas.review import ReviewCreateRequest, ReviewResponse

router = APIRouter(prefix="/reviews", tags=["reviews"])

@router.post("/", response_model=ReviewResponse, status_code=201)
async def create_review(
    review_data: ReviewCreateRequest,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    service = ReviewService(db)
    review = await service.create_review(review_data)
    return review
```

### Tech Stack (Python Services)

- **Python:** 3.11+
- **Framework:** FastAPI 0.115+ with [standard] dependencies
- **ASGI Server:** Uvicorn with uvloop (performance)
- **ORM:** SQLAlchemy 2.0 (async)
- **Database:** PostgreSQL 15 with asyncpg driver
- **Validation:** Pydantic v2
- **ML Libraries:** NLTK, TextBlob (for sentiment analysis)
- **Background Tasks:** FastAPI BackgroundTasks + asyncio

---

## API Gateway Integration

### Purpose
Single entry point for all client requests, routing to appropriate backend services.

### Request Flow
```
Browser/Client
    ↓
API Gateway (Port 8080)
    ↓ Route based on path prefix
    ├─ /api/v1/catalog/*        → Catalog Service (8081)
    ├─ /api/v1/users/*          → User Service (8082)
    ├─ /api/v1/cart/*           → Cart Service (8083)
    ├─ /api/v1/orders/*         → Order Service (8084)
    ├─ /api/v1/inventory/*      → Inventory Service (8086)
    ├─ /api/v1/reviews/*        → Review Service (8088)
    ├─ /api/v1/recommendations/*→ Recommendation Service (8089)
    └─ /api/v1/admin/*          → Admin Service (8090)
```

### Health Checks
```bash
# Overall health status
GET http://localhost:8080/health

Response:
{
  "status": "ok",  // or "degraded"
  "service": "api-gateway",
  "services": {
    "catalog": { "healthy": true },
    "user": { "healthy": true },
    "cart": { "healthy": true },
    "order": { "healthy": true },
    "inventory": { "healthy": true },
    "review": { "healthy": true },
    "recommendation": { "healthy": true },
    "admin": { "healthy": true }
  }
}
```

### Features
- ✅ HTTP proxy to all services
- ✅ Rate limiting (100 req/min configurable)
- ✅ CORS handling (configurable origins)
- ✅ Request logging with structured logs
- ✅ Health check aggregation
- ✅ JWT forwarding to protected services

### Configuration
```go
type Config struct {
    Port                     string  // 8080
    CatalogServiceURL        string  // http://catalog-service:8081
    UserServiceURL           string  // http://user-service:8082
    CartServiceURL           string  // http://cart-service:8083
    OrderServiceURL          string  // http://order-service:8084
    InventoryServiceURL      string  // http://inventory-service:8086
    ReviewServiceURL         string  // http://review-service:8088
    RecommendationServiceURL string  // http://recommendation-service:8089
    AdminServiceURL          string  // http://admin-service:8090
    RateLimitEnabled         bool
    RateLimitMax             int
    RateLimitWindow          string
}
```

---

## Frontend Architecture

### Tech Stack

- **Framework:** React 19
- **Language:** TypeScript 5+
- **Build Tool:** Vite
- **Router:** TanStack Router (file-based routing)
- **State Management:**
  - **Server State:** TanStack Query (React Query)
  - **Client State:** Zustand
- **UI Framework:** ShadcnUI + Tailwind CSS
- **HTTP Client:** Axios

### Project Structure

```
frontend/customer-app/
├── src/
│   ├── main.tsx                       # App entry point
│   ├── App.tsx                        # Root component
│   ├── routes/                        # TanStack Router routes
│   │   ├── __root.tsx
│   │   ├── index.tsx                  # Home page
│   │   ├── books/
│   │   │   ├── index.tsx              # Book list
│   │   │   └── $id.tsx                # Book details
│   │   ├── genres/
│   │   │   ├── index.tsx              # Genre list
│   │   │   └── $slug.tsx              # Genre books
│   │   └── authors/
│   │       └── $id.tsx                # Author page
│   ├── components/                    # React components
│   │   ├── BookCard.tsx
│   │   ├── GenreCard.tsx
│   │   ├── SearchBar.tsx
│   │   ├── BookGrid.tsx
│   │   └── ui/                        # ShadcnUI components
│   ├── lib/
│   │   ├── api.ts                     # API client (Axios)
│   │   └── utils.ts                   # Utilities
│   ├── store/                         # Zustand stores
│   │   ├── authStore.ts
│   │   ├── cartStore.ts
│   │   └── themeStore.ts
│   └── types/                         # TypeScript types
│       ├── auth.ts
│       ├── book.ts
│       ├── cart.ts
│       └── user.ts
├── public/
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── package.json
```

### API Client Configuration

```typescript
// src/lib/api.ts
import axios from 'axios';

// All requests go through API Gateway
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '',  // Empty = relative URLs
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add JWT token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const booksAPI = {
  list: (params) => api.get('/api/v1/catalog/books', { params }),
  search: (q) => api.get('/api/v1/catalog/books/search', { params: { q } }),
  get: (id) => api.get(`/api/v1/catalog/books/${id}`),
};

export const authAPI = {
  login: (data) => api.post('/api/v1/users/auth/login', data),
  register: (data) => api.post('/api/v1/users/auth/register', data),
  logout: () => api.post('/api/v1/users/auth/logout'),
  me: () => api.get('/api/v1/users/me'),
};

export const recommendationsAPI = {
  getPersonalized: (limit = 10) => 
    api.get(`/api/v1/recommendations/me?limit=${limit}`),
  getSimilar: (bookId, limit = 10) => 
    api.get(`/api/v1/recommendations/similar/${bookId}?limit=${limit}`),
  trackInteraction: (data) => 
    api.post('/api/v1/recommendations/interactions', data),
};
```

### Production Configuration (Kubernetes/EKS)

**Frontend Deployment:**
- Nginx serves static files
- Nginx proxies `/api/*` requests to API Gateway (internal Kubernetes service)
- No CORS issues (same-origin requests)
- One LoadBalancer for frontend (API Gateway remains internal)

**Request Flow:**
```
User Browser
    ↓
Frontend LoadBalancer (AWS ELB)
    ↓
Frontend Pod (Nginx)
    ↓
/api/* → Nginx Proxy → API Gateway Service (ClusterIP)
    ↓
Backend Services (all internal)
```

**Nginx Configuration:**
```nginx
location /api/ {
    proxy_pass http://api-gateway.bookstore-dev.svc.cluster.local:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

---

## Database Design

### Database-per-Service Pattern

Each service has its own database for data isolation and independent scaling.

| Database | Service | Tables | Purpose |
|----------|---------|--------|---------|
| `catalog_db` | Catalog | books, authors, categories, publishers, book_authors, book_categories | Book catalog |
| `userdb` | User | users, roles, addresses, sessions, wishlist | User management |
| `orderdb` | Order | orders, order_items, order_status_history | Order processing |
| `inventory_db` | Inventory | inventory, reservations, stock_movements | Stock tracking |
| `reviews_db` | Review | reviews, review_votes | Book reviews |
| `bookstore` | Admin | (queries other databases) | Analytics |

### PostgreSQL Configuration

**Init Scripts** (`infrastructure/k8s/databases/postgres-initdb-configmap.yaml`):

```sql
-- 01-create-databases.sql
CREATE DATABASE catalog_db;
CREATE DATABASE userdb;
CREATE DATABASE orderdb;
CREATE DATABASE inventory_db;
CREATE DATABASE reviews_db;
CREATE DATABASE bookstore;

-- 02-create-users.sql
CREATE USER userservice WITH PASSWORD 'dev_password';
CREATE USER orderservice WITH PASSWORD 'dev_password';
CREATE USER orderuser WITH PASSWORD 'dev_password';

-- 03-grant-permissions.sql
GRANT ALL PRIVILEGES ON DATABASE userdb TO userservice;
GRANT ALL PRIVILEGES ON DATABASE orderdb TO orderservice;
GRANT ALL PRIVILEGES ON DATABASE orderdb TO orderuser;

-- Schema-level permissions for table creation
\c userdb
GRANT ALL ON SCHEMA public TO userservice;
ALTER DEFAULT PRIVILEGES FOR USER userservice IN SCHEMA public GRANT ALL ON TABLES TO userservice;
ALTER DEFAULT PRIVILEGES FOR USER userservice IN SCHEMA public GRANT ALL ON SEQUENCES TO userservice;
```

**Connection Patterns:**

```go
// User Service (DATABASE_URL env var)
dsn := os.Getenv("DATABASE_URL")
// postgresql://userservice:dev_password@postgres:5432/userdb?sslmode=disable

// Order Service (individual env vars)
dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),     // postgres
    os.Getenv("DB_PORT"),     // 5432
    os.Getenv("DB_USER"),     // orderservice
    os.Getenv("DB_PASSWORD"), // dev_password
    os.Getenv("DB_NAME"))     // orderdb
```

### Auto-Migration (GORM)

```go
// Catalog Service
db.AutoMigrate(&domain.Book{}, &domain.Author{}, &domain.Category{}, &domain.Publisher{})

// User Service
db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.Address{}, &domain.Session{})
```

### Indexing Strategy

- **Primary Keys:** UUID (better for distributed systems than auto-increment)
- **Foreign Keys:** Indexed automatically by GORM
- **Search Columns:** Indexed (e.g., `title`, `email`, `isbn`)
- **Join Tables:** Composite indexes on both columns

---

## Testing Guide

### Local Development Testing

#### 1. Start Backend with Docker Compose

```bash
cd services/api-gateway
docker compose up -d

# Wait 30 seconds for database seeding
sleep 30

# Check all services are running
docker compose ps
```

#### 2. Start Frontend

```bash
cd frontend/customer-app
npm install
npm run dev
# Open http://localhost:5173
```

#### 3. Test API Endpoints

```bash
# Health check
curl http://localhost:8080/health | jq

# List books
curl http://localhost:8080/api/v1/catalog/books | jq

# Search books
curl "http://localhost:8080/api/v1/catalog/books/search?q=distributed" | jq

# Register user
curl -X POST http://localhost:8080/api/v1/users/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "full_name": "Test User"
  }' | jq

# Login
curl -X POST http://localhost:8080/api/v1/users/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!"
  }' | jq

# Get personalized recommendations (with JWT)
curl http://localhost:8080/api/v1/recommendations/me?limit=10 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" | jq

# Create review
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": "BOOK_UUID",
    "user_id": "USER_UUID",
    "rating": 5,
    "title": "Excellent book!",
    "content": "This book was amazing! Highly recommended.",
    "verified_purchase": true
  }' | jq
```

### Frontend Testing Checklist

- [ ] Home page loads with featured books and genres
- [ ] Search bar works with debouncing (300ms)
- [ ] Genre navigation works (`/genres` → `/genres/:slug`)
- [ ] Book details page shows correct information
- [ ] Author page shows author bio and books
- [ ] Add to cart works
- [ ] Wishlist functionality works (requires login)
- [ ] User registration and login works
- [ ] Protected routes redirect to login

### E2E User Flow Test

1. **Browse Catalog:**
   - Visit home page
   - Click "Programming" genre
   - See filtered books

2. **Search:**
   - Type "distributed" in search bar
   - See "Distributed Systems" book

3. **View Details:**
   - Click on book card
   - See full details, reviews, recommendations

4. **Track Interaction:**
   - View book (auto-tracked)
   - Add to cart (tracked)
   - See personalized recommendations update

5. **Checkout Flow:**
   - View cart
   - Proceed to checkout
   - Create order
   - See order confirmation

---

## Deployment Architecture

### AWS EKS Production Deployment

**Cluster Configuration:**
- **Name:** my-bookstore
- **Region:** us-east-1
- **Node Type:** t3.medium
- **Node Count:** 3 (auto-scaling 2-4)
- **Kubernetes Version:** 1.28+
- **Namespace:** bookstore-dev

**Container Registry:**
- **ECR Repository:** 905418472239.dkr.ecr.us-east-1.amazonaws.com
- **Images:** All 10 services + frontend
- **Tag Strategy:** `latest` (production would use semantic versioning)

**Service Endpoints:**
- **Frontend:** AWS LoadBalancer (external)
  ```
  http://ab1a1c3c5b1ca49768c26f26e92ca780-844836377.us-east-1.elb.amazonaws.com
  ```
- **API Gateway:** ClusterIP (internal only, proxied via frontend Nginx)
- **All Backend Services:** ClusterIP (internal only)

**Database:**
- **Type:** PostgreSQL 15 StatefulSet
- **Storage:** AWS EBS (gp2, 10Gi)
- **Persistence:** Init scripts in ConfigMap
- **Replicas:** 1 (production would use managed RDS)

**Cache:**
- **Type:** Redis 7 Deployment
- **Persistence:** EmptyDir (ephemeral, production would use ElastiCache)
- **Replicas:** 1

**Deployment Status (Production):**
```
NAME                                  READY   STATUS    RESTARTS
admin-service-xxxx                    1/1     Running   0
api-gateway-xxxx                      2/2     Running   0
cart-service-xxxx                     2/2     Running   0
catalog-service-xxxx                  2/2     Running   0
frontend-xxxx                         2/2     Running   0
inventory-service-xxxx                1/1     Running   0
order-service-xxxx                    2/2     Running   0
postgres-0                            1/1     Running   0
recommendation-service-xxxx           1/1     Running   0
redis-xxxx                            1/1     Running   0
review-service-xxxx                   1/1     Running   0
user-service-xxxx                     2/2     Running   0
```

**Total:** 12/12 services deployed and healthy

### Kubernetes Resources

**Namespace:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: bookstore-dev
```

**Secrets:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
type: Opaque
data:
  postgres-password: ZGV2X3Bhc3N3b3Jk  # base64: dev_password
```

**Services (examples):**
```yaml
# Frontend - LoadBalancer (public)
apiVersion: v1
kind: Service
metadata:
  name: frontend
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 80

# API Gateway - ClusterIP (internal)
apiVersion: v1
kind: Service
metadata:
  name: api-gateway
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: 8080
```

### Deployment Best Practices

✅ **Implemented:**
- Database-per-service pattern
- Init scripts for database setup (persistent via ConfigMap)
- Health checks for all services
- Resource limits and requests
- Multi-replica deployments for high availability
- Secrets for sensitive data
- Namespace isolation
- LoadBalancer for frontend only (API Gateway internal)

🔜 **Production Enhancements:**
- Managed databases (AWS RDS for PostgreSQL)
- Managed cache (AWS ElastiCache for Redis)
- Horizontal Pod Autoscaler (HPA)
- Ingress controller with TLS/SSL
- Service mesh (Istio/Linkerd)
- Monitoring (Prometheus + Grafana)
- Logging (ELK/EFK stack)
- Distributed tracing (Jaeger)
- CI/CD pipelines (GitHub Actions)

---

## Development Workflow

### Adding a New Feature

1. **Backend (Go Service):**
   ```bash
   # Add domain model
   vim services/catalog-service/internal/domain/models.go
   
   # Add repository method
   vim services/catalog-service/internal/repository/book_repository.go
   
   # Add service method
   vim services/catalog-service/internal/service/catalog_service.go
   
   # Add HTTP handler
   vim services/catalog-service/internal/handler/http/book_handler.go
   
   # Add route
   vim services/catalog-service/internal/handler/http/routes.go
   
   # Test locally
   go run cmd/server/main.go
   ```

2. **Backend (Python Service):**
   ```bash
   # Add model
   vim app/models/review.py
   
   # Add schema
   vim app/schemas/review.py
   
   # Add service method
   vim app/services/review_service.py
   
   # Add endpoint
   vim app/api/v1/endpoints/reviews.py
   
   # Test locally
   uvicorn app.main:app --reload --port 8088
   ```

3. **Frontend:**
   ```bash
   # Add API client method
   vim src/lib/api.ts
   
   # Add TypeScript types
   vim src/types/book.ts
   
   # Add component
   vim src/components/BookRecommendations.tsx
   
   # Add route
   vim src/routes/books/$id.tsx
   
   # Test locally
   npm run dev
   ```

4. **Deploy to Production:**
   ```bash
   # Build and push Docker image
   docker build -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest .
   docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest
   
   # Update Kubernetes deployment
   kubectl rollout restart deployment/catalog-service -n bookstore-dev
   
   # Verify
   kubectl get pods -n bookstore-dev
   kubectl logs -f deployment/catalog-service -n bookstore-dev
   ```

---

## Conclusion

This distributed bookstore system demonstrates:

✅ **Microservices Architecture** - 12 independent services  
✅ **Polyglot Persistence** - PostgreSQL (relational) + Redis (cache)  
✅ **Clean Architecture** - Clear separation of concerns  
✅ **Cloud-Native** - Docker containers, Kubernetes orchestration  
✅ **Modern Tech Stack** - Go, Python, TypeScript, React 19  
✅ **Production Deployment** - Running on AWS EKS  
✅ **ML Integration** - Sentiment analysis, recommendations  
✅ **API Gateway Pattern** - Single entry point  
✅ **Database-per-Service** - Data isolation  
✅ **Security** - JWT authentication, RBAC  

**Ready for:**
- [ ] Service mesh (Istio/Linkerd)
- [ ] Event-driven architecture (RabbitMQ/Kafka)
- [ ] CQRS and Event Sourcing
- [ ] Distributed tracing (Jaeger)
- [ ] Monitoring (Prometheus/Grafana)
- [ ] CI/CD automation
- [ ] Multi-region deployment

---

**For detailed deployment instructions, see:** [EKS_DEPLOYMENT_GUIDE.md](../EKS_DEPLOYMENT_GUIDE.md)  
**For project overview, see:** [README.md](../README.md)

