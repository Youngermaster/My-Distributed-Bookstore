# API Gateway & Frontend Integration Guide

## What Was Implemented

### API Gateway - Complete Integration

The API Gateway now routes to ALL backend microservices with clean RESTful patterns:

#### Service Routes Implemented:

```
http://localhost:8080/api/v1/
├── /catalog/*           → Catalog Service (port 8081)
│   ├── /books
│   ├── /authors
│   ├── /categories
│   └── /publishers
│
├── /users/*             → User Service (port 8082)  [NEW]
│   ├── /auth/login
│   ├── /auth/register
│   ├── /auth/logout
│   ├── /me
│   └── /me/wishlist
│
├── /cart/*              → Cart Service (port 8083)  [NEW]
│   ├── /{cartId}
│   ├── /{cartId}/items
│   └── /{cartId}/items/{bookId}
│
├── /orders/*            → Order Service (port 8084)  [NEW]
│   ├── /
│   ├── /{orderId}
│   └── /{orderId}/cancel
│
└── /recommendations/*   → Recommendation Service (port 8089)  [NEW]
    ├── /me
    ├── /similar/{bookId}
    ├── /trending
    ├── /popular
    └── /interactions
```

#### Health Checks - Now Monitor All Services:

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
    "recommendation": { "healthy": true }
  }
}

# Readiness check
GET http://localhost:8080/ready
```

---

## Architecture Overview

```
┌─────────────┐
│   Frontend  │
│  (React)    │
│  :3000      │
└──────┬──────┘
       │
       │ HTTP/REST
       ▼
┌─────────────────────────────────────────┐
│         API Gateway (:8080)             │
│  ┌────────────────────────────────┐    │
│  │  1. Validate Request            │    │
│  │  2. Run Middleware (JWT, Rate)  │    │
│  │  3. Route to Service            │    │
│  │  4. Transform Response          │    │
│  └────────────────────────────────┘    │
└───┬───┬────┬────┬──────────┬───────────┘
    │   │    │    │          │
    ▼   ▼    ▼    ▼          ▼
  ┌──┐┌──┐┌───┐┌───┐┌──────────────┐
  │C ││U ││Ca ││O  ││Recommendation│
  │a ││s ││rt ││rd ││   :8089      │
  │t ││e │   ││er ││              │
  │a ││r │:  ││:  ││              │
  │l │  ││83 ││84 ││              │
  │o │:82││   ││   ││              │
  │g │  ││   ││   ││              │
  │  │  ││   ││   ││              │
  │:81││  ││   ││   ││              │
  └──┘└──┘└───┘└───┘└──────────────┘
```

---

## Files Created/Modified

### API Gateway Changes:

#### Configuration (`internal/config/config.go`)
```go
type Config struct {
    // Service URLs
    CatalogServiceURL        string  // http://localhost:8081
    UserServiceURL           string  // http://localhost:8082  [NEW]
    CartServiceURL           string  // http://localhost:8083  [NEW]
    OrderServiceURL          string  // http://localhost:8084  [NEW]
    RecommendationServiceURL string  // http://localhost:8089  [NEW]
}
```

#### Proxy Clients Created:
- `internal/proxy/user_client.go` - User service proxy [NEW]
- `internal/proxy/cart_client.go` - Cart service proxy [NEW]
- `internal/proxy/order_client.go` - Order service proxy [NEW]
- `internal/proxy/recommendation_client.go` - Recommendation proxy [NEW]

#### Handlers Created:
- `internal/handler/user_handler.go` - Routes /users/* [NEW]
- `internal/handler/cart_handler.go` - Routes /cart/* [NEW]
- `internal/handler/order_handler.go` - Routes /orders/* [NEW]
- `internal/handler/recommendation_handler.go` - Routes /recommendations/* [NEW]

#### Main Application (`cmd/server/main.go`)
- Initialized all service clients
- Created all handlers
- Registered all routes
- Updated health checks
- Enhanced logging

---

## Frontend Integration - What's Needed

### Current Frontend State:
The frontend (`frontend/customer-app/`) already has:
- Pages: Login, Register, Cart, Orders, Wishlist, BookDetail
- Types: auth.ts, book.ts, cart.ts, order.ts, user.ts, wishlist.ts
- Stores: authStore.ts, cartStore.ts, themeStore.ts
- API Client: lib/api.ts (needs updating)
- Routes: TanStack Router setup

### Frontend Updates Required:

#### 1. Update API Client (`src/lib/api.ts`)

Currently the API client has mixed routing. Update to route everything through API Gateway:

**After (through API Gateway):**
```typescript
// All through API Gateway on port 8080
export const authAPI = {
  login: (data: LoginRequest) =>
    api.post<AuthResponse>("/api/v1/users/auth/login", data),

  register: (data: RegisterRequest) =>
    api.post<AuthResponse>("/api/v1/users/auth/register", data),

  logout: () =>
    api.post("/api/v1/users/auth/logout"),

  me: () =>
    api.get<{ data: User }>("/api/v1/users/me"),
};

// Wishlist routes through user service
export const wishlistAPI = {
  list: () =>
    api.get<WishlistResponse>("/api/v1/users/me/wishlist"),

  add: (book_id: string) =>
    api.post("/api/v1/users/me/wishlist", { book_id }),

  remove: (book_id: string) =>
    api.delete(`/api/v1/users/me/wishlist/${book_id}`),

  clear: () =>
    api.delete("/api/v1/users/me/wishlist"),

  check: (book_id: string) =>
    api.get(`/api/v1/users/me/wishlist/check/${book_id}`),
};

// NEW: Recommendations API
export const recommendationsAPI = {
  getPersonalized: (limit = 10) =>
    api.get(`/api/v1/recommendations/me?limit=${limit}`),

  getSimilar: (bookId: string, limit = 10) =>
    api.get(`/api/v1/recommendations/similar/${bookId}?limit=${limit}`),

  getTrending: (limit = 10, days = 7) =>
    api.get(`/api/v1/recommendations/trending?limit=${limit}&days=${days}`),

  getPopular: (limit = 10) =>
    api.get(`/api/v1/recommendations/popular?limit=${limit}`),

  trackInteraction: (data: { book_id: string; interaction_type: string }) =>
    api.post("/api/v1/recommendations/interactions", data),
};
```

#### 2. Add Recommendation Types (`src/types/recommendation.ts`)

Create new file:
```typescript
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
  interaction_type: 'view' | 'add_to_cart' | 'purchase' | 'review' | 'wishlist';
}
```

#### 3. Update Environment Variables

**`.env` or `.env.local`:**
```bash
# All API calls go through the gateway
VITE_API_URL=http://localhost:8080
```

#### 4. Track Book Views (BookDetail.tsx)

When a user views a book, track it:
```typescript
import { recommendationsAPI } from '@/lib/api';

// In BookDetail component
useEffect(() => {
  if (bookId) {
    // Track view interaction
    recommendationsAPI.trackInteraction({
      book_id: bookId,
      interaction_type: 'view'
    }).catch(console.error);  // Silent fail - don't block UI
  }
}, [bookId]);
```

#### 5. Add "Recommended for You" Section (Home.tsx)

```typescript
import { useQuery } from '@tanstack/react-query';
import { recommendationsAPI } from '@/lib/api';

// In Home component
const { data: recommendations } = useQuery({
  queryKey: ['recommendations', 'personalized'],
  queryFn: () => recommendationsAPI.getPersonalized(10),
  enabled: !!isAuthenticated,
});
```

#### 6. Track Cart Additions

When adding to cart:
```typescript
const addToCart = async (bookId: string) => {
  await cartAPI.addItem({ book_id: bookId, quantity: 1 });

  // Track interaction
  await recommendationsAPI.trackInteraction({
    book_id: bookId,
    interaction_type: 'add_to_cart'
  });
};
```

---

## Testing the Integration

### 1. Test API Gateway Routes:

```bash
# Health check (should show all 5 services)
curl http://localhost:8080/health

# List all endpoints
curl http://localhost:8080/

# Test catalog (existing)
curl http://localhost:8080/api/v1/catalog/books

# Test user service
curl -X POST http://localhost:8080/api/v1/users/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 2. Start All Services:

```bash
# API Gateway
cd services/api-gateway && go run cmd/server/main.go

# Backend services (in separate terminals)
cd services/catalog-service && go run cmd/server/main.go
cd services/user-service && go run cmd/server/main.go
cd services/cart-service && go run cmd/server/main.go
cd services/order-service && go run cmd/server/main.go
cd services/recommendation-service && python -m app.main

# Frontend
cd frontend/customer-app && pnpm dev
```

---

## API Routes Summary

| Route Pattern | Service | Port | Purpose |
|--------------|---------|------|---------|
| `/api/v1/catalog/*` | Catalog | 8081 | Books, authors, categories |
| `/api/v1/users/*` | User | 8082 | Auth, profile, wishlist |
| `/api/v1/cart/*` | Cart | 8083 | Shopping cart |
| `/api/v1/orders/*` | Order | 8084 | Order management |
| `/api/v1/recommendations/*` | Recommendation | 8089 | Personalized recommendations |

---

## Authentication Flow

```
1. User submits login
   Frontend → Gateway → User Service

2. User Service validates & returns JWT
   User Service → Gateway → Frontend

3. Frontend stores JWT in localStorage

4. Subsequent requests include JWT
   Frontend (with JWT) → Gateway (validates) → Service

5. Gateway forwards Authorization header to services
```

---

## Benefits of This Architecture

1. **Single Entry Point**: Frontend only needs to know about the gateway (port 8080)
2. **Service Independence**: Backend services can be developed independently
3. **Clean URLs**: RESTful patterns (`/users/*`, `/orders/*`)
4. **Centralized Auth**: Gateway can validate JWTs before routing
5. **Easy Scaling**: Add load balancing at gateway level
6. **Monitoring**: Single point to monitor all service health
7. **Rate Limiting**: Applied at gateway for all services

---

## Next Steps

1. [DONE] API Gateway fully implemented
2. [TODO] Update frontend `api.ts` with new routes
3. [TODO] Add recommendations API to frontend
4. [TODO] Track user interactions (views, cart, purchases)
5. [TODO] Add "Recommended for You" section
6. [TODO] Add "Trending Books" section
7. [TODO] Add "Similar Books" on book details
8. [TODO] Test end-to-end integration

---

**Status**: API Gateway Complete
**Next**: Frontend Integration
**Date**: 2025-11-02
