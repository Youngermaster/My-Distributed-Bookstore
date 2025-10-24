# Cart Service

## Overview

Manages shopping cart functionality with session-based carts for anonymous users and persistent carts for authenticated users.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Storage**: Redis (primary), PostgreSQL (backup)
- **Ports**: HTTP: 8083, gRPC: 50053

## Responsibilities

- Shopping cart management
- Session-based cart for anonymous users
- Persistent cart for authenticated users
- Cart synchronization (anonymous → authenticated)
- Real-time price calculation
- Cart item validation

## Data Model (Redis)

```json
{
  "cart:{user_id}": {
    "items": [
      {
        "book_id": "uuid",
        "quantity": 2,
        "unit_price": 29.99,
        "added_at": "timestamp"
      }
    ],
    "total": 59.98,
    "updated_at": "timestamp"
  }
}
```

## REST API Endpoints

```
GET    /api/v1/cart               # Get current cart
POST   /api/v1/cart/items         # Add item to cart
PUT    /api/v1/cart/items/:bookId # Update quantity
DELETE /api/v1/cart/items/:bookId # Remove from cart
DELETE /api/v1/cart               # Clear cart
POST   /api/v1/cart/sync          # Sync anonymous cart to user cart
```

## gRPC Methods

```protobuf
rpc GetCart(GetCartRequest) returns (GetCartResponse);
rpc AddItem(AddItemRequest) returns (CartResponse);
rpc ClearCart(ClearCartRequest) returns (ClearCartResponse);
```

## Events Consumed

- `catalog.price_updated` - Update cart prices
- `catalog.stock_updated` - Validate cart availability

## Environment Variables

```bash
HTTP_PORT=8083
GRPC_PORT=50053
REDIS_URL=redis:6379
CATALOG_SERVICE_URL=catalog-service:50051
CART_TTL=7d  # Cart expiration for anonymous users
```

## Next Steps

- [ ] Implement Redis operations
- [ ] Create cart handlers
- [ ] Integrate with Catalog Service via gRPC
- [ ] Add cart validation logic
- [ ] Implement cart synchronization
- [ ] Add tests
