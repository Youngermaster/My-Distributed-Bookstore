# API Gateway Service

## Overview

The API Gateway serves as the single entry point for all client requests in the distributed bookstore system. It handles request routing, authentication, rate limiting, and provides a unified interface for the frontend.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Protocol**: HTTP/REST
- **Port**: 8080 (HTTP)

## Responsibilities

- Single entry point for all client requests
- Request routing to appropriate microservices
- JWT token validation
- Rate limiting (Redis-backed)
- Request/response logging with correlation IDs
- CORS handling
- API composition (aggregate multiple service calls)
- Health check aggregation from all services
- Circuit breaker for downstream services
- Request timeout management
- Metrics collection (Prometheus)

## Key Features

### Authentication
- JWT token validation for protected routes
- Token refresh mechanism
- Role-based access control (RBAC)

### Routing
Routes requests to appropriate services:
- `/api/v1/books/*` → Catalog Service
- `/api/v1/auth/*` → User Service
- `/api/v1/cart/*` → Cart Service
- `/api/v1/orders/*` → Order Service
- `/api/v1/payments/*` → Payment Service
- `/api/v1/inventory/*` → Inventory Service
- `/api/v1/reviews/*` → Review Service
- `/api/v1/recommendations/*` → Recommendation Service
- `/api/v1/admin/*` → Admin Service

### Resilience
- Circuit breaker pattern for downstream services
- Request timeout management
- Automatic retry logic
- Graceful degradation

### Observability
- Distributed tracing with Jaeger
- Prometheus metrics
- Structured logging
- Request correlation IDs

## Project Structure

```
api-gateway/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── handler/
│   │   ├── http/
│   │   │   └── handler.go      # HTTP route handlers
│   │   └── grpc/
│   │       └── client.go       # gRPC client setup
│   ├── middleware/
│   │   ├── auth.go             # JWT validation middleware
│   │   ├── logger.go           # Logging middleware
│   │   ├── cors.go             # CORS middleware
│   │   ├── ratelimit.go        # Rate limiting middleware
│   │   └── metrics.go          # Prometheus metrics middleware
│   ├── routing/
│   │   └── router.go           # Route definitions
│   └── client/
│       └── grpc.go             # gRPC client connections
├── pkg/
│   ├── jwt/                    # JWT utilities
│   ├── validator/              # Input validation
│   ├── errors/                 # Error handling
│   └── response/               # HTTP response helpers
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Environment Variables

```bash
# Server Configuration
PORT=8080
ENV=development

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Redis Configuration (for rate limiting)
REDIS_URL=redis:6379
REDIS_PASSWORD=

# Downstream Services (gRPC)
CATALOG_SERVICE_URL=catalog-service:50051
USER_SERVICE_URL=user-service:50052
CART_SERVICE_URL=cart-service:50053
ORDER_SERVICE_URL=order-service:50054
PAYMENT_SERVICE_URL=payment-service:50055
INVENTORY_SERVICE_URL=inventory-service:50056
REVIEW_SERVICE_URL=review-service:50058
RECOMMENDATION_SERVICE_URL=recommendation-service:50059
ADMIN_SERVICE_URL=admin-service:50060

# Observability
JAEGER_AGENT_HOST=jaeger
JAEGER_AGENT_PORT=6831
LOG_LEVEL=info

# Rate Limiting
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m

# Timeouts
REQUEST_TIMEOUT=30s
CIRCUIT_BREAKER_THRESHOLD=5
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Redis (for rate limiting)

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run locally with Docker Compose:
```bash
docker-compose up
```

3. Run without Docker:
```bash
go run cmd/server/main.go
```

### Testing

```bash
# Run unit tests
go test ./... -v

# Run tests with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### Building

```bash
# Build binary
go build -o bin/api-gateway cmd/server/main.go

# Build Docker image
docker build -t api-gateway:latest .
```

## API Endpoints

### Health Check
```
GET /health
Response: { "status": "ok", "services": {...} }
```

### Catalog
```
GET    /api/v1/books              # List books
GET    /api/v1/books/:id          # Get book details
GET    /api/v1/books/search       # Search books
POST   /api/v1/books              # Create book (admin)
PUT    /api/v1/books/:id          # Update book (admin)
GET    /api/v1/categories         # List categories
```

### Authentication
```
POST   /api/v1/auth/register      # User registration
POST   /api/v1/auth/login         # Login
POST   /api/v1/auth/logout        # Logout
POST   /api/v1/auth/refresh       # Refresh token
```

### User
```
GET    /api/v1/users/me           # Get current user
PUT    /api/v1/users/me           # Update profile
GET    /api/v1/users/me/addresses # List addresses
POST   /api/v1/users/me/addresses # Add address
```

### Cart
```
GET    /api/v1/cart               # Get cart
POST   /api/v1/cart/items         # Add item
PUT    /api/v1/cart/items/:id     # Update quantity
DELETE /api/v1/cart/items/:id     # Remove item
DELETE /api/v1/cart               # Clear cart
```

### Orders
```
POST   /api/v1/orders             # Create order
GET    /api/v1/orders             # List orders
GET    /api/v1/orders/:id         # Get order details
POST   /api/v1/orders/:id/cancel  # Cancel order
```

## Development Guidelines

1. **Code Style**: Follow Go best practices and use `gofmt`
2. **Error Handling**: Always return appropriate HTTP status codes
3. **Logging**: Use structured logging with correlation IDs
4. **Testing**: Maintain >80% code coverage
5. **Documentation**: Keep API documentation up to date

## Troubleshooting

### Common Issues

**Cannot connect to downstream services:**
- Check gRPC URLs in environment variables
- Verify services are running
- Check network connectivity

**Rate limiting too aggressive:**
- Adjust `RATE_LIMIT_MAX` and `RATE_LIMIT_WINDOW`
- Check Redis connection

**High latency:**
- Check Jaeger traces for bottlenecks
- Review timeout configurations
- Consider enabling request caching

## Next Steps

1. Implement main.go with Fiber setup
2. Create JWT middleware
3. Implement gRPC client connections
4. Add rate limiting middleware
5. Implement health check aggregation
6. Add circuit breaker logic
7. Set up distributed tracing
8. Add comprehensive tests
