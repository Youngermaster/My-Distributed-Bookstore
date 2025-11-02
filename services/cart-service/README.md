# Cart Service

Simple shopping cart service for the distributed bookstore.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Storage**: Redis (in-memory key-value store)
- **Port**: 8083

## Architecture

Simple clean architecture following catalog-service pattern:

```
cart-service/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/                     # Configuration
│   ├── domain/                     # Domain entities (Cart, CartItem)
│   ├── repository/                 # Redis storage
│   ├── service/                    # Business logic
│   └── handler/http/               # HTTP endpoints
```

## API Endpoints

```http
GET    /health                        # Health check
GET    /ready                         # Readiness check
GET    /api/v1/cart/:cartId           # Get cart
POST   /api/v1/cart/:cartId/items     # Add item
PUT    /api/v1/cart/:cartId/items/:bookId   # Update quantity
DELETE /api/v1/cart/:cartId/items/:bookId   # Remove item
DELETE /api/v1/cart/:cartId           # Clear cart
```

## Quick Start

```bash
# Start with Docker
docker-compose up -d

# Or run locally
cp .env.example .env
go run cmd/server/main.go
```

## Testing

```bash
# Create a new cart (use any UUID)
CART_ID=$(uuidgen)

# Get cart
curl http://localhost:8083/api/v1/cart/$CART_ID

# Add item
curl -X POST http://localhost:8083/api/v1/cart/$CART_ID/items \
  -H "Content-Type: application/json" \
  -d '{
    "book_id":"123e4567-e89b-12d3-a456-426614174000",
    "quantity":2,
    "price":29.99
  }'

# Update quantity
curl -X PUT http://localhost:8083/api/v1/cart/$CART_ID/items/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{"quantity":3}'

# Remove item
curl -X DELETE http://localhost:8083/api/v1/cart/$CART_ID/items/123e4567-e89b-12d3-a456-426614174000

# Clear cart
curl -X DELETE http://localhost:8083/api/v1/cart/$CART_ID
```

## Features

✅ Redis-only storage (simple & fast)  
✅ UUID-based cart IDs  
✅ Automatic TTL (7 days default)  
✅ Quantity & item limits  
✅ Real-time total calculation  
✅ Health & readiness checks

## Configuration

Key environment variables:

```env
CART_TTL_HOURS=168           # Cart expiration (7 days)
MAX_ITEMS_PER_CART=50        # Maximum items per cart
MAX_QUANTITY_PER_ITEM=99     # Maximum quantity per item
```

---

**Simple, working, following catalog-service pattern** ✅
