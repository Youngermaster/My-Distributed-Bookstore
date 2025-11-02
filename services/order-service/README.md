# Order Service

Simple order management service for the distributed bookstore.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Port**: 8084

## Architecture

Simple clean architecture following catalog-service pattern:

```
order-service/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/                     # Configuration
│   ├── domain/                     # Domain entities (Order, OrderItem)
│   ├── repository/                 # PostgreSQL storage
│   ├── service/                    # Business logic
│   └── handler/http/               # HTTP endpoints
```

## API Endpoints

```http
GET    /health                         # Health check
GET    /ready                          # Readiness check
POST   /api/v1/orders                  # Create order
GET    /api/v1/orders                  # List all orders (paginated)
GET    /api/v1/orders/:id              # Get order by ID
PATCH  /api/v1/orders/:id/status       # Update order status
POST   /api/v1/orders/:id/cancel       # Cancel order
GET    /api/v1/users/:userId/orders    # Get user's orders (paginated)
```

## Order Statuses

- `pending` - Order created, awaiting confirmation
- `confirmed` - Order confirmed
- `processing` - Order being processed
- `shipped` - Order shipped
- `delivered` - Order delivered
- `cancelled` - Order cancelled

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
# Create order
curl -X POST http://localhost:8084/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "items": [
      {
        "book_id": "456e7890-e89b-12d3-a456-426614174000",
        "quantity": 2,
        "unit_price": 29.99
      }
    ],
    "shipping_address": "123 Main St, City, Country",
    "payment_method": "credit_card"
  }'

# Get order
curl http://localhost:8084/api/v1/orders/{order_id}

# List all orders (paginated)
curl http://localhost:8084/api/v1/orders?page=1&page_size=20

# Get user orders
curl http://localhost:8084/api/v1/users/{user_id}/orders?page=1&page_size=10

# Update order status
curl -X PATCH http://localhost:8084/api/v1/orders/{order_id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "confirmed"}'

# Cancel order
curl -X POST http://localhost:8084/api/v1/orders/{order_id}/cancel
```

## Features

✅ PostgreSQL storage with GORM  
✅ Order lifecycle management  
✅ Status transitions with validation  
✅ User order history  
✅ Pagination support  
✅ Automatic total calculation  
✅ Health & readiness checks

## Configuration

Key environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=order_db
DEFAULT_PAGE_SIZE=20           # Default pagination size
MAX_PAGE_SIZE=100              # Maximum items per page
```

---

**Simple, working, following catalog-service pattern** ✅
