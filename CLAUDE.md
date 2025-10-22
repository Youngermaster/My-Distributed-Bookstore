# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A production-ready distributed bookstore system built with microservices architecture on AWS EKS, following industry best practices and principles from "Distributed Systems" by Tanenbaum & van Steen. This project demonstrates real-world implementation of distributed systems concepts including event-driven architecture, saga patterns, gRPC communication, and comprehensive observability.

### Tech Stack

**Backend Microservices** (Polyglot Architecture):
- **Go 1.21+**: Primary language for high-performance services
  - Fiber (web framework)
  - GORM (ORM)
  - gRPC with Protocol Buffers
- **Node.js/TypeScript**: For I/O-bound services
  - Express.js
  - Payment gateway integrations
  - WebSocket support
- **Python 3.11+**: For ML/AI services
  - FastAPI
  - scikit-learn (recommendations)
  - NLTK (sentiment analysis)
- **JWT**: Authentication & authorization
- **PostgreSQL 15**: Primary relational database
- **Redis 7**: Caching, sessions, rate limiting
- **RabbitMQ 3.12**: Message broker for async communication

**Frontend**:
- TypeScript + React 18+
- shadcn/ui component library
- TanStack Query v5 (data fetching)
- Zustand (state management)
- React Router v6
- Tailwind CSS

**Infrastructure**:
- Docker & Docker Compose (local development)
- Kubernetes 1.28+ (AWS EKS)
- Nginx Ingress Controller
- PostgreSQL (per-service databases)
- Redis (shared cache cluster)
- RabbitMQ (StatefulSet cluster)
- S3 (static assets, book covers)
- CloudFront (CDN)
- Route53 (DNS)
- ACM (SSL certificates)

**Observability**:
- Jaeger (distributed tracing)
- Prometheus (metrics collection)
- Grafana (visualization dashboards)
- ELK Stack (centralized logging)
  - Fluentd (log shipping)
  - Elasticsearch (log storage)
  - Kibana (log analysis)

## Architecture Principles

### Microservices Design
- **Database per Service**: Each microservice owns its data
- **API Gateway Pattern**: Single entry point for clients
- **Service Discovery**: Kubernetes DNS
- **Event-Driven Architecture**: Asynchronous communication via RabbitMQ
- **Saga Pattern**: Distributed transactions without 2PC
- **CQRS**: Separate read/write models for complex queries
- **Circuit Breaker**: Fault tolerance with automatic recovery
- **Health Checks**: Liveness and readiness probes

### Communication Patterns
- **Synchronous**: REST (client ↔ API Gateway), gRPC (service ↔ service)
- **Asynchronous**: RabbitMQ (events, notifications, long-running tasks)
- **Protocol Buffers**: Strongly-typed service contracts

### Distributed Systems Concepts (Tanenbaum)
- **Consistency**: Eventual consistency for high availability
- **Replication**: Multi-replica deployments for fault tolerance
- **Naming**: Kubernetes service discovery
- **Coordination**: Event-driven saga orchestration
- **Fault Tolerance**: Auto-healing, circuit breakers, retries
- **Scalability**: Horizontal pod autoscaling
- **Transparency**: Users see single coherent system

## Project Structure

```
distributed-bookstore/
├── services/
│   ├── api-gateway/              # API Gateway (Golang)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── catalog-service/          # Book catalog (Golang)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── proto/
│   │   ├── migrations/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── user-service/             # Auth & users (Golang)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── proto/
│   │   ├── migrations/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── cart-service/             # Shopping cart (Golang)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── proto/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── order-service/            # Order processing (Golang)
│   │   ├── cmd/
│   │   ├── internal/
│   │   │   ├── saga/            # Saga orchestrator
│   │   │   └── events/          # Event publishers
│   │   ├── proto/
│   │   ├── migrations/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── payment-service/          # Payments (Node.js/TypeScript)
│   │   ├── src/
│   │   │   ├── controllers/
│   │   │   ├── grpc/
│   │   │   └── stripe/
│   │   ├── proto/
│   │   ├── Dockerfile
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── inventory-service/        # Stock management (Golang)
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── proto/
│   │   ├── migrations/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── notification-service/     # Email/SMS (Node.js/TypeScript)
│   │   ├── src/
│   │   │   ├── consumers/       # RabbitMQ consumers
│   │   │   ├── email/           # SendGrid/SES
│   │   │   └── sms/             # Twilio (optional)
│   │   ├── Dockerfile
│   │   ├── package.json
│   │   └── tsconfig.json
│   ├── review-service/           # Reviews & ratings (Python/FastAPI)
│   │   ├── app/
│   │   │   ├── api/
│   │   │   ├── ml/              # Sentiment analysis
│   │   │   └── repository/
│   │   ├── proto/
│   │   ├── migrations/
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   ├── recommendation-service/   # ML recommendations (Python/FastAPI)
│   │   ├── app/
│   │   │   ├── api/
│   │   │   ├── ml/              # Collaborative filtering
│   │   │   └── repository/
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   └── admin-service/            # Admin & reporting (Golang)
│       ├── cmd/
│       ├── internal/
│       ├── proto/
│       ├── Dockerfile
│       └── go.mod
├── frontend/
│   ├── customer-app/             # Customer React app
│   │   ├── src/
│   │   │   ├── components/
│   │   │   │   ├── ui/          # shadcn/ui components
│   │   │   │   ├── books/
│   │   │   │   ├── cart/
│   │   │   │   └── orders/
│   │   │   ├── pages/
│   │   │   ├── api/             # API client
│   │   │   ├── hooks/
│   │   │   ├── store/           # Zustand stores
│   │   │   └── types/
│   │   ├── public/
│   │   ├── Dockerfile
│   │   ├── package.json
│   │   └── tsconfig.json
│   └── admin-app/                # Admin dashboard (optional)
├── infrastructure/
│   └── k8s/                      # Kubernetes manifests
│       ├── namespaces/
│       │   ├── production.yaml
│       │   └── staging.yaml
│       ├── ingress/
│       │   ├── nginx-controller.yaml
│       │   └── bookstore-ingress.yaml
│       ├── services/
│       │   ├── api-gateway/
│       │   │   ├── deployment.yaml
│       │   │   ├── service.yaml
│       │   │   ├── hpa.yaml
│       │   │   └── configmap.yaml
│       │   ├── catalog-service/
│       │   ├── user-service/
│       │   ├── cart-service/
│       │   ├── order-service/
│       │   ├── payment-service/
│       │   ├── inventory-service/
│       │   ├── notification-service/
│       │   ├── review-service/
│       │   ├── recommendation-service/
│       │   └── admin-service/
│       ├── databases/
│       │   ├── postgres/
│       │   │   ├── statefulset.yaml
│       │   │   ├── service.yaml
│       │   │   └── pvc.yaml
│       │   └── redis/
│       │       ├── deployment.yaml
│       │       └── service.yaml
│       ├── messaging/
│       │   ├── rabbitmq-statefulset.yaml
│       │   ├── rabbitmq-service.yaml
│       │   └── rabbitmq-configmap.yaml
│       ├── observability/
│       │   ├── jaeger/
│       │   ├── prometheus/
│       │   ├── grafana/
│       │   └── elk/
│       ├── network-policies/
│       │   ├── default-deny.yaml
│       │   └── service-specific/
│       └── secrets/
│           ├── db-credentials.yaml
│           └── jwt-secret.yaml
├── proto/                        # Shared protobuf definitions
│   ├── catalog.proto
│   ├── user.proto
│   ├── cart.proto
│   ├── order.proto
│   ├── payment.proto
│   ├── inventory.proto
│   └── common.proto
├── scripts/
│   ├── build-all.sh
│   ├── deploy-local.sh
│   ├── deploy-k8s.sh
│   ├── generate-protos.sh
│   └── init-databases.sh
├── docs/
│   ├── architecture.md
│   ├── api-specs/
│   ├── database-schemas/
│   ├── deployment-guide.md
│   └── troubleshooting.md
├── .github/
│   └── workflows/
│       ├── ci.yaml
│       └── cd.yaml
├── docker-compose.yml            # Local development
├── Makefile
└── README.md
```

## Core Microservices

### 1. API Gateway (Golang)
**Port**: 8080 (HTTP)

**Responsibilities**:
- Single entry point for all client requests
- Request routing to appropriate services
- JWT validation
- Rate limiting (Redis-backed)
- Request/response logging with correlation IDs
- CORS handling
- API composition (aggregate multiple service calls)

**Key Features**:
- Health check aggregation from all services
- Circuit breaker for downstream services
- Request timeout management
- Metrics collection (Prometheus)

**Endpoints**:
- `GET /health` - Aggregated health status
- `GET /api/v1/books/*` - Proxy to catalog service
- `POST /api/v1/auth/*` - Proxy to user service
- `GET /api/v1/cart/*` - Proxy to cart service
- `POST /api/v1/orders/*` - Proxy to order service

### 2. Catalog Service (Golang)
**Port**: 8081 (HTTP), 50051 (gRPC)
**Database**: `catalog_db` (PostgreSQL)

**Responsibilities**:
- Book catalog management
- Search and filtering
- Category management
- Publisher and author information
- Book metadata

**Key Entities**:
```sql
books (
  id UUID PRIMARY KEY,
  isbn VARCHAR(13) UNIQUE,
  title VARCHAR(500),
  description TEXT,
  price DECIMAL(10,2),
  stock_quantity INTEGER,
  publisher_id UUID,
  cover_image_url TEXT,
  metadata JSONB,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)

authors (
  id UUID PRIMARY KEY,
  name VARCHAR(255),
  bio TEXT,
  birth_date DATE
)

book_authors (
  book_id UUID,
  author_id UUID,
  author_order INTEGER,
  PRIMARY KEY (book_id, author_id)
)

categories (
  id UUID PRIMARY KEY,
  name VARCHAR(100),
  slug VARCHAR(100) UNIQUE,
  parent_id UUID,
  FOREIGN KEY (parent_id) REFERENCES categories(id)
)

book_categories (
  book_id UUID,
  category_id UUID,
  PRIMARY KEY (book_id, category_id)
)

publishers (
  id UUID PRIMARY KEY,
  name VARCHAR(255),
  country VARCHAR(100)
)
```

**REST Endpoints**:
- `GET /api/v1/books` - List books (pagination, filters)
- `GET /api/v1/books/:id` - Get book details
- `GET /api/v1/books/search?q=query` - Search books
- `POST /api/v1/books` - Create book (admin only)
- `PUT /api/v1/books/:id` - Update book (admin only)
- `GET /api/v1/categories` - List categories (hierarchical)

**gRPC Methods**:
- `GetBook(id) -> BookDetails`
- `SearchBooks(query) -> BookList`
- `CheckStock(bookId) -> StockInfo`
- `UpdateStock(bookId, quantity) -> bool`

**Events Published**:
- `catalog.book_created` - New book added
- `catalog.stock_updated` - Stock level changed
- `catalog.price_updated` - Price changed

### 3. User Service (Golang)
**Port**: 8082 (HTTP), 50052 (gRPC)
**Database**: `users_db` (PostgreSQL)

**Responsibilities**:
- User registration and authentication
- JWT token generation and validation
- User profile management
- Address management
- Role-based access control (RBAC)

**Key Entities**:
```sql
users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  full_name VARCHAR(255),
  phone VARCHAR(20),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  last_login_at TIMESTAMP
)

roles (
  id UUID PRIMARY KEY,
  name VARCHAR(50) UNIQUE,
  permissions JSONB
)

user_roles (
  user_id UUID,
  role_id UUID,
  PRIMARY KEY (user_id, role_id)
)

addresses (
  id UUID PRIMARY KEY,
  user_id UUID,
  address_line1 VARCHAR(255),
  address_line2 VARCHAR(255),
  city VARCHAR(100),
  state VARCHAR(100),
  postal_code VARCHAR(20),
  country VARCHAR(100),
  is_default BOOLEAN,
  FOREIGN KEY (user_id) REFERENCES users(id)
)

sessions (
  id UUID PRIMARY KEY,
  user_id UUID,
  token_hash VARCHAR(255),
  ip_address INET,
  user_agent TEXT,
  expires_at TIMESTAMP,
  created_at TIMESTAMP
)
```

**REST Endpoints**:
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - Login (returns JWT)
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update profile
- `GET /api/v1/users/me/addresses` - List addresses
- `POST /api/v1/users/me/addresses` - Add address

**gRPC Methods**:
- `ValidateToken(token) -> UserClaims`
- `GetUser(userId) -> UserProfile`
- `CheckPermission(userId, resource) -> bool`

**Events Published**:
- `user.registered` - New user signed up
- `user.logged_in` - User logged in
- `user.profile_updated` - Profile changed

### 4. Cart Service (Golang)
**Port**: 8083 (HTTP), 50053 (gRPC)
**Storage**: Redis (primary), PostgreSQL (persistent backup)

**Responsibilities**:
- Shopping cart management
- Session-based cart for anonymous users
- Persistent cart for authenticated users
- Cart synchronization
- Price calculation

**Data Model** (Redis):
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

**REST Endpoints**:
- `GET /api/v1/cart` - Get current cart
- `POST /api/v1/cart/items` - Add item to cart
- `PUT /api/v1/cart/items/:bookId` - Update quantity
- `DELETE /api/v1/cart/items/:bookId` - Remove from cart
- `DELETE /api/v1/cart` - Clear cart

**gRPC Methods**:
- `GetCart(userId) -> CartDetails`
- `AddItem(userId, bookId, quantity) -> Cart`
- `ClearCart(userId) -> bool`

**Events Consumed**:
- `catalog.price_updated` - Update cart prices
- `catalog.stock_updated` - Validate cart availability

### 5. Order Service (Golang)
**Port**: 8084 (HTTP), 50054 (gRPC)
**Database**: `orders_db` (PostgreSQL)

**Responsibilities**:
- Order creation and management
- Saga orchestration for distributed transactions
- Order status tracking
- Order history
- Coordination with payment and inventory services

**Key Entities**:
```sql
orders (
  id UUID PRIMARY KEY,
  user_id UUID,
  status VARCHAR(50),
  total_amount DECIMAL(10,2),
  shipping_address_id UUID,
  payment_method VARCHAR(50),
  tracking_number VARCHAR(100),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)

order_items (
  id UUID PRIMARY KEY,
  order_id UUID,
  book_id UUID,
  quantity INTEGER,
  unit_price DECIMAL(10,2),
  subtotal DECIMAL(10,2),
  FOREIGN KEY (order_id) REFERENCES orders(id)
)

order_status_history (
  id UUID PRIMARY KEY,
  order_id UUID,
  previous_status VARCHAR(50),
  new_status VARCHAR(50),
  changed_by UUID,
  notes TEXT,
  changed_at TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id)
)
```

**Order Status Flow**:
```
PENDING → PAYMENT_PROCESSING → PAID → PREPARING → SHIPPED → DELIVERED
                ↓                 ↓         ↓
            CANCELLED       CANCELLED  CANCELLED
```

**REST Endpoints**:
- `POST /api/v1/orders` - Create order from cart
- `GET /api/v1/orders` - List user orders
- `GET /api/v1/orders/:id` - Get order details
- `POST /api/v1/orders/:id/cancel` - Cancel order
- `GET /api/v1/orders/:id/tracking` - Get tracking info

**gRPC Methods**:
- `CreateOrder(userId, items) -> OrderId`
- `GetOrder(orderId) -> OrderDetails`
- `UpdateOrderStatus(orderId, status) -> bool`

**Saga Orchestration** (Choreography Pattern):
```
1. Order Service creates order (status: PENDING)
2. Publishes: order.created
3. Payment Service processes payment
   - Success: Publishes payment.completed
   - Failure: Publishes payment.failed
4. Inventory Service reserves stock
   - Success: Publishes inventory.reserved
   - Failure: Publishes inventory.reservation_failed
5. Order Service updates status based on events
6. Notification Service sends confirmation
```

**Compensating Transactions**:
- If payment fails: Cancel order
- If inventory reservation fails: Refund payment, cancel order
- If order cancelled: Refund payment, release inventory

**Events Published**:
- `order.created` - New order placed
- `order.confirmed` - Payment and inventory successful
- `order.cancelled` - Order cancelled
- `order.shipped` - Order shipped
- `order.delivered` - Order delivered

**Events Consumed**:
- `payment.completed` - Payment successful
- `payment.failed` - Payment failed
- `inventory.reserved` - Stock reserved
- `inventory.reservation_failed` - Stock unavailable

### 6. Payment Service (Node.js/TypeScript)
**Port**: 8085 (HTTP), 50055 (gRPC)
**Database**: `payments_db` (PostgreSQL)

**Responsibilities**:
- Payment processing via Stripe API
- Payment method management
- Refund processing
- Payment status tracking
- PCI compliance handling

**Key Entities**:
```sql
payments (
  id UUID PRIMARY KEY,
  order_id UUID,
  user_id UUID,
  amount DECIMAL(10,2),
  currency VARCHAR(3),
  payment_method VARCHAR(50),
  status VARCHAR(50),
  stripe_payment_intent_id VARCHAR(255),
  stripe_charge_id VARCHAR(255),
  failure_reason TEXT,
  metadata JSONB,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)

refunds (
  id UUID PRIMARY KEY,
  payment_id UUID,
  amount DECIMAL(10,2),
  reason TEXT,
  status VARCHAR(50),
  stripe_refund_id VARCHAR(255),
  created_at TIMESTAMP,
  FOREIGN KEY (payment_id) REFERENCES payments(id)
)
```

**REST Endpoints**:
- `POST /api/v1/payments` - Create payment intent
- `POST /api/v1/payments/:id/confirm` - Confirm payment
- `GET /api/v1/payments/:id` - Get payment status
- `POST /api/v1/payments/:id/refund` - Issue refund

**gRPC Methods**:
- `ProcessPayment(orderId, amount, method) -> PaymentResult`
- `GetPaymentStatus(paymentId) -> PaymentStatus`
- `RefundPayment(paymentId, amount) -> RefundResult`

**Events Published**:
- `payment.processing` - Payment initiated
- `payment.completed` - Payment successful
- `payment.failed` - Payment failed
- `payment.refunded` - Refund issued

**Events Consumed**:
- `order.created` - Process payment for order
- `order.cancelled` - Refund payment if applicable

### 7. Inventory Service (Golang)
**Port**: 8086 (HTTP), 50056 (gRPC)
**Database**: `inventory_db` (PostgreSQL)

**Responsibilities**:
- Real-time stock tracking
- Stock reservation for orders
- Stock release on cancellation
- Low stock alerts
- Multi-warehouse support (future)

**Key Entities**:
```sql
inventory (
  id UUID PRIMARY KEY,
  book_id UUID UNIQUE,
  available_quantity INTEGER,
  reserved_quantity INTEGER,
  reorder_level INTEGER,
  last_restocked_at TIMESTAMP,
  updated_at TIMESTAMP
)

stock_movements (
  id UUID PRIMARY KEY,
  book_id UUID,
  movement_type VARCHAR(50),
  quantity INTEGER,
  reference_type VARCHAR(50),
  reference_id UUID,
  notes TEXT,
  created_at TIMESTAMP
)

reservations (
  id UUID PRIMARY KEY,
  book_id UUID,
  order_id UUID,
  quantity INTEGER,
  status VARCHAR(50),
  expires_at TIMESTAMP,
  created_at TIMESTAMP
)
```

**REST Endpoints**:
- `GET /api/v1/inventory/:bookId` - Get stock level
- `POST /api/v1/inventory/:bookId/adjust` - Adjust stock (admin)
- `GET /api/v1/inventory/low-stock` - Get low stock items

**gRPC Methods**:
- `CheckStock(bookId) -> StockLevel`
- `ReserveStock(orderId, items) -> ReservationResult`
- `ReleaseReservation(orderId) -> bool`
- `CommitReservation(orderId) -> bool`

**Events Published**:
- `inventory.updated` - Stock level changed
- `inventory.low_stock` - Stock below threshold
- `inventory.reserved` - Stock reserved for order
- `inventory.reservation_failed` - Insufficient stock

**Events Consumed**:
- `order.created` - Reserve stock
- `order.cancelled` - Release reservation
- `payment.completed` - Commit reservation
- `catalog.book_created` - Initialize inventory

### 8. Notification Service (Node.js/TypeScript)
**Port**: 8087 (HTTP)
**Storage**: PostgreSQL (notification history)

**Responsibilities**:
- Email notifications (SendGrid/SES)
- SMS notifications (Twilio - optional)
- Push notifications (future)
- Notification templates
- Delivery tracking
- RabbitMQ consumer for all notification events

**Key Entities**:
```sql
notifications (
  id UUID PRIMARY KEY,
  user_id UUID,
  type VARCHAR(50),
  channel VARCHAR(20),
  recipient VARCHAR(255),
  subject VARCHAR(500),
  body TEXT,
  status VARCHAR(50),
  sent_at TIMESTAMP,
  error_message TEXT,
  metadata JSONB,
  created_at TIMESTAMP
)

notification_templates (
  id UUID PRIMARY KEY,
  name VARCHAR(100) UNIQUE,
  type VARCHAR(50),
  subject_template TEXT,
  body_template TEXT,
  variables JSONB
)
```

**Notification Types**:
- Order confirmation
- Order shipped
- Order delivered
- Payment receipt
- Password reset
- Welcome email
- Low stock alert (admin)

**Events Consumed** (RabbitMQ):
- `user.registered` - Send welcome email
- `order.created` - Send order confirmation
- `order.shipped` - Send shipping notification
- `order.delivered` - Send delivery confirmation
- `payment.completed` - Send payment receipt
- `inventory.low_stock` - Alert admin

**Email Templates** (Handlebars):
```html
<!-- order-confirmation.hbs -->
<h1>Order Confirmation</h1>
<p>Hi {{userName}},</p>
<p>Your order #{{orderId}} has been confirmed!</p>
<h2>Order Details:</h2>
<ul>
  {{#each items}}
    <li>{{bookTitle}} - Quantity: {{quantity}} - ${{price}}</li>
  {{/each}}
</ul>
<p><strong>Total: ${{total}}</strong></p>
```

### 9. Review Service (Python/FastAPI)
**Port**: 8088 (HTTP), 50058 (gRPC)
**Database**: `reviews_db` (PostgreSQL)

**Responsibilities**:
- Book reviews and ratings
- Sentiment analysis (ML)
- Review moderation
- Rating aggregation
- Helpful votes

**Key Entities**:
```sql
reviews (
  id UUID PRIMARY KEY,
  book_id UUID,
  user_id UUID,
  rating INTEGER CHECK (rating >= 1 AND rating <= 5),
  title VARCHAR(255),
  content TEXT,
  sentiment_score FLOAT,
  sentiment_label VARCHAR(20),
  verified_purchase BOOLEAN,
  helpful_votes INTEGER,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  UNIQUE(book_id, user_id)
)

review_votes (
  review_id UUID,
  user_id UUID,
  is_helpful BOOLEAN,
  created_at TIMESTAMP,
  PRIMARY KEY (review_id, user_id)
)
```

**REST Endpoints**:
- `GET /api/v1/books/:bookId/reviews` - List reviews
- `POST /api/v1/books/:bookId/reviews` - Submit review
- `PUT /api/v1/reviews/:id` - Update review
- `DELETE /api/v1/reviews/:id` - Delete review
- `POST /api/v1/reviews/:id/vote` - Vote helpful/not helpful

**gRPC Methods**:
- `GetBookReviews(bookId) -> ReviewList`
- `GetAverageRating(bookId) -> RatingInfo`

**ML Pipeline**:
```python
# Sentiment analysis using NLTK/TextBlob
def analyze_sentiment(review_text):
    blob = TextBlob(review_text)
    polarity = blob.sentiment.polarity
    
    if polarity > 0.3:
        return "positive", polarity
    elif polarity < -0.3:
        return "negative", polarity
    else:
        return "neutral", polarity
```

**Events Published**:
- `review.submitted` - New review posted
- `review.updated` - Review edited

**Events Consumed**:
- `order.delivered` - Enable review submission

### 10. Recommendation Service (Python/FastAPI)
**Port**: 8089 (HTTP), 50059 (gRPC)
**Database**: `recommendations_db` (PostgreSQL + pgvector)

**Responsibilities**:
- Personalized book recommendations
- Collaborative filtering
- Content-based filtering
- "Customers who bought this also bought" recommendations
- Trending books

**Key Entities**:
```sql
user_interactions (
  id UUID PRIMARY KEY,
  user_id UUID,
  book_id UUID,
  interaction_type VARCHAR(50),
  score FLOAT,
  created_at TIMESTAMP
)

recommendations_cache (
  id UUID PRIMARY KEY,
  user_id UUID,
  book_ids UUID[],
  algorithm VARCHAR(50),
  score FLOAT,
  expires_at TIMESTAMP,
  created_at TIMESTAMP
)

book_embeddings (
  book_id UUID PRIMARY KEY,
  embedding VECTOR(128),
  updated_at TIMESTAMP
)
```

**Interaction Types**:
- View (weight: 1)
- Add to cart (weight: 3)
- Purchase (weight: 5)
- Review (weight: 4)

**REST Endpoints**:
- `GET /api/v1/recommendations` - Get personalized recommendations
- `GET /api/v1/books/:bookId/similar` - Get similar books
- `GET /api/v1/recommendations/trending` - Get trending books

**gRPC Methods**:
- `GetRecommendations(userId, limit) -> BookList`
- `GetSimilarBooks(bookId, limit) -> BookList`

**ML Algorithms**:
1. **Collaborative Filtering** (scikit-learn):
```python
from sklearn.neighbors import NearestNeighbors

def collaborative_filtering(user_id, n_recommendations=10):
    # User-item matrix
    user_item_matrix = create_user_item_matrix()
    
    # KNN model
    model = NearestNeighbors(metric='cosine', algorithm='brute')
    model.fit(user_item_matrix)
    
    # Find similar users
    distances, indices = model.kneighbors(
        user_item_matrix[user_id], 
        n_neighbors=20
    )
    
    # Aggregate recommendations
    return aggregate_recommendations(indices)
```

2. **Content-Based Filtering**:
```python
# Using book embeddings (genres, authors, descriptions)
def content_based_filtering(book_id, n_recommendations=10):
    book_vector = get_book_embedding(book_id)
    
    # Cosine similarity with all books
    similar_books = find_similar_vectors(
        book_vector, 
        n_recommendations
    )
    
    return similar_books
```

**Events Consumed**:
- `order.completed` - Update user interaction history
- `review.submitted` - Update interaction scores
- `user.registered` - Initialize empty recommendation cache

### 11. Admin Service (Golang)
**Port**: 8090 (HTTP), 50060 (gRPC)
**Database**: Aggregates data from all services

**Responsibilities**:
- Admin dashboard backend
- Analytics and reporting
- User management (admin functions)
- Book management (admin functions)
- Order management (admin view)
- System health monitoring

**REST Endpoints**:
- `GET /api/v1/admin/dashboard` - Dashboard stats
- `GET /api/v1/admin/orders` - All orders
- `GET /api/v1/admin/users` - All users
- `GET /api/v1/admin/analytics/sales` - Sales analytics
- `GET /api/v1/admin/analytics/inventory` - Inventory reports
- `POST /api/v1/admin/books` - Bulk book import

**Dashboard Metrics**:
- Total orders today/week/month
- Revenue today/week/month
- Active users
- Low stock items
- Average order value
- Top selling books
- Order status breakdown

## RabbitMQ Event Architecture

### Exchange Configuration

```yaml
Exchanges:
  orders:
    type: fanout
    durable: true
    description: "Order lifecycle events"
  
  payments:
    type: direct
    durable: true
    description: "Payment processing events"
  
  notifications:
    type: topic
    durable: true
    description: "User notifications"
    routing_keys:
      - notification.email.*
      - notification.sms.*
      - notification.push.*
  
  inventory:
    type: topic
    durable: true
    description: "Inventory management"
    routing_keys:
      - inventory.updated
      - inventory.low_stock
      - inventory.reserved
  
  catalog:
    type: topic
    durable: true
    description: "Catalog changes"
    routing_keys:
      - catalog.book_created
      - catalog.price_updated
      - catalog.stock_updated
```

### Queue Bindings

```yaml
Queues:
  # Order events
  order.payment:
    exchange: orders
    durable: true
    consumer: payment-service
  
  order.inventory:
    exchange: orders
    durable: true
    consumer: inventory-service
  
  order.notification:
    exchange: orders
    durable: true
    consumer: notification-service
  
  # Payment events
  payment.order:
    exchange: payments
    routing_key: payment.completed
    consumer: order-service
  
  payment.notification:
    exchange: payments
    routing_key: payment.*
    consumer: notification-service
  
  # Notification events
  notification.email:
    exchange: notifications
    routing_key: notification.email.*
    consumer: notification-service
  
  # Inventory events
  inventory.catalog:
    exchange: inventory
    routing_key: inventory.updated
    consumer: catalog-service
  
  inventory.admin:
    exchange: inventory
    routing_key: inventory.low_stock
    consumer: admin-service
  
  # Catalog events
  catalog.recommendation:
    exchange: catalog
    routing_key: catalog.*
    consumer: recommendation-service
```

### Event Schema Examples

```go
// order.created event
type OrderCreatedEvent struct {
    EventID   string    `json:"event_id"`
    Timestamp time.Time `json:"timestamp"`
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Items     []struct {
        BookID   string  `json:"book_id"`
        Quantity int     `json:"quantity"`
        Price    float64 `json:"price"`
    } `json:"items"`
    TotalAmount       float64 `json:"total_amount"`
    ShippingAddressID string  `json:"shipping_address_id"`
}

// payment.completed event
type PaymentCompletedEvent struct {
    EventID       string    `json:"event_id"`
    Timestamp     time.Time `json:"timestamp"`
    PaymentID     string    `json:"payment_id"`
    OrderID       string    `json:"order_id"`
    Amount        float64   `json:"amount"`
    PaymentMethod string    `json:"payment_method"`
    StripeChargeID string   `json:"stripe_charge_id"`
}

// inventory.reserved event
type InventoryReservedEvent struct {
    EventID     string    `json:"event_id"`
    Timestamp   time.Time `json:"timestamp"`
    OrderID     string    `json:"order_id"`
    Items       []struct {
        BookID   string `json:"book_id"`
        Quantity int    `json:"quantity"`
    } `json:"items"`
    ReservationID string `json:"reservation_id"`
    ExpiresAt     time.Time `json:"expires_at"`
}
```

## gRPC Service Contracts

### Example: Catalog Service Proto

```protobuf
syntax = "proto3";

package catalog;

option go_package = "github.com/yourusername/bookstore/proto/catalog";

service CatalogService {
  rpc GetBook(GetBookRequest) returns (GetBookResponse);
  rpc SearchBooks(SearchBooksRequest) returns (SearchBooksResponse);
  rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
  rpc UpdateStock(UpdateStockRequest) returns (UpdateStockResponse);
  rpc GetBooksByIds(GetBooksByIdsRequest) returns (GetBooksByIdsResponse);
}

message GetBookRequest {
  string book_id = 1;
}

message GetBookResponse {
  string id = 1;
  string isbn = 2;
  string title = 3;
  string description = 4;
  double price = 5;
  int32 stock_quantity = 6;
  string cover_image_url = 7;
  repeated string author_names = 8;
  repeated string category_names = 9;
  string publisher_name = 10;
  double average_rating = 11;
  int32 review_count = 12;
}

message SearchBooksRequest {
  string query = 1;
  int32 page = 2;
  int32 page_size = 3;
  repeated string genres = 4;
  double min_price = 5;
  double max_price = 6;
  string sort_by = 7; // price, title, rating
  string sort_order = 8; // asc, desc
}

message SearchBooksResponse {
  repeated GetBookResponse books = 1;
  int32 total = 2;
  int32 page = 3;
  int32 total_pages = 4;
}

message CheckStockRequest {
  string book_id = 1;
}

message CheckStockResponse {
  string book_id = 1;
  int32 available_quantity = 2;
  int32 reserved_quantity = 3;
  bool in_stock = 4;
}

message UpdateStockRequest {
  string book_id = 1;
  int32 quantity_change = 2;
  string operation_type = 3; // add, subtract, set
}

message UpdateStockResponse {
  bool success = 1;
  int32 new_stock_quantity = 2;
}

message GetBooksByIdsRequest {
  repeated string book_ids = 1;
}

message GetBooksByIdsResponse {
  repeated GetBookResponse books = 1;
}
```

## Service Architecture Pattern (Golang)

Each Go service follows this structure:

```
service-name/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── domain/
│   │   └── models.go              # Domain entities
│   ├── repository/
│   │   ├── interface.go           # Repository interface
│   │   └── postgres/
│   │       └── repository.go      # PostgreSQL implementation
│   ├── service/
│   │   └── service.go             # Business logic
│   ├── handler/
│   │   ├── http/
│   │   │   └── handler.go         # REST handlers
│   │   └── grpc/
│   │       └── server.go          # gRPC handlers
│   ├── middleware/
│   │   ├── auth.go                # JWT validation
│   │   ├── logger.go              # Logging middleware
│   │   ├── cors.go                # CORS handling
│   │   └── metrics.go             # Prometheus metrics
│   ├── events/
│   │   ├── publisher.go           # RabbitMQ publisher
│   │   └── consumer.go            # RabbitMQ consumer
│   └── tracing/
│       └── jaeger.go              # Distributed tracing
├── pkg/
│   ├── jwt/                       # JWT utilities
│   ├── validator/                 # Input validation
│   ├── errors/                    # Error handling
│   └── response/                  # HTTP response helpers
├── proto/
│   └── service.proto              # gRPC definitions
├── migrations/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── Dockerfile
├── go.mod
└── go.sum
```

### Example: main.go Structure

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/yourusername/bookstore/catalog-service/internal/config"
    "github.com/yourusername/bookstore/catalog-service/internal/handler/grpc"
    "github.com/yourusername/bookstore/catalog-service/internal/handler/http"
    "github.com/yourusername/bookstore/catalog-service/internal/repository/postgres"
    "github.com/yourusername/bookstore/catalog-service/internal/service"
    "github.com/yourusername/bookstore/catalog-service/internal/tracing"
    pb "github.com/yourusername/bookstore/proto/catalog"
    grpclib "google.golang.org/grpc"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize Jaeger tracing
    tracer, closer, err := tracing.InitJaeger("catalog-service")
    if err != nil {
        log.Printf("Warning: Failed to init Jaeger: %v", err)
    }
    defer closer.Close()

    // Initialize database
    db, err := postgres.NewDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize repository
    repo := postgres.NewCatalogRepository(db)

    // Initialize service
    svc := service.NewCatalogService(repo)

    // Initialize HTTP server (Fiber)
    app := fiber.New(fiber.Config{
        ErrorHandler: http.CustomErrorHandler,
    })

    // Setup HTTP routes
    http.SetupRoutes(app, svc, tracer)

    // Initialize gRPC server
    grpcServer := grpclib.NewServer()
    pb.RegisterCatalogServiceServer(grpcServer, grpc.NewCatalogServer(svc))

    // Start gRPC server in goroutine
    go func() {
        lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
        if err != nil {
            log.Fatalf("Failed to listen: %v", err)
        }
        log.Printf("gRPC server listening on :%d", cfg.GRPCPort)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Failed to serve gRPC: %v", err)
        }
    }()

    // Start HTTP server in goroutine
    go func() {
        log.Printf("HTTP server listening on :%d", cfg.HTTPPort)
        if err := app.Listen(fmt.Sprintf(":%d", cfg.HTTPPort)); err != nil {
            log.Fatalf("Failed to start HTTP server: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down servers...")

    // Shutdown HTTP server
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := app.ShutdownWithContext(ctx); err != nil {
        log.Printf("HTTP server shutdown error: %v", err)
    }

    // Stop gRPC server
    grpcServer.GracefulStop()

    log.Println("Servers stopped")
}
```

## Docker Configuration

### Multi-stage Dockerfile (Golang)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations (if needed)
COPY --from=builder /app/migrations ./migrations

# Expose ports
EXPOSE 8081 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run
CMD ["./main"]
```

### Docker Compose (Local Development)

```yaml
version: '3.9'

services:
  # Infrastructure
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: bookstore
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-databases.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U bookstore"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    environment:
      RABBITMQ_DEFAULT_USER: bookstore
      RABBITMQ_DEFAULT_PASS: dev_password
    ports:
      - "5672:5672"
      - "15672:15672"
    healthcheck:
      test: ["CMD", "rabbitmqctl", "status"]
      interval: 30s
      timeout: 10s
      retries: 5

  # Microservices
  api-gateway:
    build:
      context: ./services/api-gateway
      dockerfile: Dockerfile
    environment:
      PORT: 8080
      CATALOG_SERVICE_URL: catalog-service:50051
      USER_SERVICE_URL: user-service:50052
      CART_SERVICE_URL: cart-service:50053
      ORDER_SERVICE_URL: order-service:50054
      REDIS_URL: redis:6379
      JWT_SECRET: dev_jwt_secret_change_in_production
      JAEGER_AGENT_HOST: jaeger
      JAEGER_AGENT_PORT: 6831
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  catalog-service:
    build:
      context: ./services/catalog-service
      dockerfile: Dockerfile
    environment:
      HTTP_PORT: 8081
      GRPC_PORT: 50051
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/catalog_db?sslmode=disable
      REDIS_URL: redis:6379
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
      JAEGER_AGENT_HOST: jaeger
      JAEGER_AGENT_PORT: 6831
    ports:
      - "8081:8081"
      - "50051:50051"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  user-service:
    build:
      context: ./services/user-service
      dockerfile: Dockerfile
    environment:
      HTTP_PORT: 8082
      GRPC_PORT: 50052
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/users_db?sslmode=disable
      JWT_SECRET: dev_jwt_secret_change_in_production
      JWT_EXPIRATION: 24h
      JAEGER_AGENT_HOST: jaeger
    ports:
      - "8082:8082"
      - "50052:50052"
    depends_on:
      postgres:
        condition: service_healthy

  cart-service:
    build:
      context: ./services/cart-service
      dockerfile: Dockerfile
    environment:
      HTTP_PORT: 8083
      GRPC_PORT: 50053
      REDIS_URL: redis:6379
      CATALOG_SERVICE_URL: catalog-service:50051
      JAEGER_AGENT_HOST: jaeger
    ports:
      - "8083:8083"
      - "50053:50053"
    depends_on:
      redis:
        condition: service_healthy

  order-service:
    build:
      context: ./services/order-service
      dockerfile: Dockerfile
    environment:
      HTTP_PORT: 8084
      GRPC_PORT: 50054
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/orders_db?sslmode=disable
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
      CATALOG_SERVICE_URL: catalog-service:50051
      USER_SERVICE_URL: user-service:50052
      JAEGER_AGENT_HOST: jaeger
    ports:
      - "8084:8084"
      - "50054:50054"
    depends_on:
      postgres:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  payment-service:
    build:
      context: ./services/payment-service
      dockerfile: Dockerfile
    environment:
      PORT: 8085
      GRPC_PORT: 50055
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/payments_db?sslmode=disable
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
      STRIPE_API_KEY: ${STRIPE_API_KEY:-sk_test_xxx}
      JAEGER_AGENT_HOST: jaeger
    ports:
      - "8085:8085"
      - "50055:50055"
    depends_on:
      postgres:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  inventory-service:
    build:
      context: ./services/inventory-service
      dockerfile: Dockerfile
    environment:
      HTTP_PORT: 8086
      GRPC_PORT: 50056
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/inventory_db?sslmode=disable
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
      JAEGER_AGENT_HOST: jaeger
    ports:
      - "8086:8086"
      - "50056:50056"
    depends_on:
      postgres:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  notification-service:
    build:
      context: ./services/notification-service
      dockerfile: Dockerfile
    environment:
      PORT: 8087
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/notifications_db?sslmode=disable
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
      SENDGRID_API_KEY: ${SENDGRID_API_KEY}
      FROM_EMAIL: noreply@bookstore.com
    ports:
      - "8087:8087"
    depends_on:
      postgres:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy

  review-service:
    build:
      context: ./services/review-service
      dockerfile: Dockerfile
    environment:
      PORT: 8088
      GRPC_PORT: 50058
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/reviews_db?sslmode=disable
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
    ports:
      - "8088:8088"
      - "50058:50058"
    depends_on:
      postgres:
        condition: service_healthy

  recommendation-service:
    build:
      context: ./services/recommendation-service
      dockerfile: Dockerfile
    environment:
      PORT: 8089
      GRPC_PORT: 50059
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/recommendations_db?sslmode=disable
      RABBITMQ_URL: amqp://bookstore:dev_password@rabbitmq:5672/
      CATALOG_SERVICE_URL: catalog-service:50051
    ports:
      - "8089:8089"
      - "50059:50059"
    depends_on:
      postgres:
        condition: service_healthy

  admin-service:
    build:
      context: ./services/admin-service
      dockerfile: Dockerfile
    environment:
      HTTP_PORT: 8090
      GRPC_PORT: 50060
      DATABASE_URL: postgresql://bookstore:dev_password@postgres:5432/admin_db?sslmode=disable
      CATALOG_SERVICE_URL: catalog-service:50051
      ORDER_SERVICE_URL: order-service:50054
      USER_SERVICE_URL: user-service:50052
      INVENTORY_SERVICE_URL: inventory-service:50056
    ports:
      - "8090:8090"
      - "50060:50060"
    depends_on:
      postgres:
        condition: service_healthy

  # Frontend
  frontend:
    build:
      context: ./frontend/customer-app
      dockerfile: Dockerfile
      args:
        REACT_APP_API_URL: http://localhost:8080
    ports:
      - "3000:80"
    depends_on:
      - api-gateway

  # Observability
  jaeger:
    image: jaegertracing/all-in-one:1.51
    environment:
      COLLECTOR_ZIPKIN_HOST_PORT: :9411
    ports:
      - "5775:5775/udp"
      - "6831:6831/udp"
      - "6832:6832/udp"
      - "5778:5778"
      - "16686:16686"
      - "14268:14268"
      - "14250:14250"
      - "9411:9411"

  prometheus:
    image: prom/prometheus:v2.48.0
    volumes:
      - ./infrastructure/observability/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:10.2.2
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
      GF_USERS_ALLOW_SIGN_UP: false
    volumes:
      - ./infrastructure/observability/grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./infrastructure/observability/grafana/datasources:/etc/grafana/provisioning/datasources
      - grafana_data:/var/lib/grafana
    ports:
      - "3001:3000"
    depends_on:
      - prometheus

volumes:
  postgres_data:
  prometheus_data:
  grafana_data:

networks:
  default:
    name: bookstore-network
```

## Kubernetes Deployment

### Example: Catalog Service Deployment

```yaml
# infrastructure/k8s/services/catalog-service/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: catalog-service
  namespace: production
  labels:
    app: catalog-service
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: catalog-service
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: catalog-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8081"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: catalog-service-sa
      
      # Anti-affinity
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  app: catalog-service
              topologyKey: kubernetes.io/hostname
      
      containers:
      - name: catalog-service
        image: <ECR_REGISTRY>/catalog-service:latest
        imagePullPolicy: Always
        
        ports:
        - name: http
          containerPort: 8081
          protocol: TCP
        - name: grpc
          containerPort: 50051
          protocol: TCP
        
        env:
        - name: HTTP_PORT
          value: "8081"
        - name: GRPC_PORT
          value: "50051"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: catalog-db-secret
              key: url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: catalog-config
              key: redis_url
        - name: RABBITMQ_URL
          valueFrom:
            secretKeyRef:
              name: rabbitmq-secret
              key: url
        - name: JAEGER_AGENT_HOST
          value: "jaeger-agent.observability.svc.cluster.local"
        - name: JAEGER_AGENT_PORT
          value: "6831"
        - name: LOG_LEVEL
          value: "info"
        
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
        
        volumeMounts:
        - name: tmp
          mountPath: /tmp
      
      volumes:
      - name: tmp
        emptyDir: {}
      
      imagePullSecrets:
      - name: ecr-registry-secret

---
# infrastructure/k8s/services/catalog-service/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: catalog-service
  namespace: production
  labels:
    app: catalog-service
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    targetPort: 8081
    protocol: TCP
  - name: grpc
    port: 50051
    targetPort: 50051
    protocol: TCP
  selector:
    app: catalog-service

---
# infrastructure/k8s/services/catalog-service/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: catalog-service-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: catalog-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 30
      - type: Pods
        value: 4
        periodSeconds: 30
      selectPolicy: Max
```

### Network Policy Example

```yaml
# infrastructure/k8s/network-policies/catalog-service-netpol.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: catalog-service-netpol
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: catalog-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow from API Gateway
  - from:
    - podSelector:
        matchLabels:
          app: api-gateway
    ports:
    - protocol: TCP
      port: 8081
  # Allow gRPC from other services
  - from:
    - podSelector:
        matchLabels:
          app: order-service
    - podSelector:
        matchLabels:
          app: cart-service
    - podSelector:
        matchLabels:
          app: recommendation-service
    ports:
    - protocol: TCP
      port: 50051
  egress:
  # Allow to PostgreSQL
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  # Allow to Redis
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  # Allow to RabbitMQ
  - to:
    - podSelector:
        matchLabels:
          app: rabbitmq
    ports:
    - protocol: TCP
      port: 5672
  # Allow DNS
  - to:
    - namespaceSelector: {}
      podSelector:
        matchLabels:
          k8s-app: kube-dns
    ports:
    - protocol: UDP
      port: 53
  # Allow to Jaeger
  - to:
    - namespaceSelector:
        matchLabels:
          name: observability
      podSelector:
        matchLabels:
          app: jaeger-agent
    ports:
    - protocol: UDP
      port: 6831
```

## Frontend Architecture (React + shadcn/ui)

### Project Structure

```
frontend/customer-app/
├── src/
│   ├── components/
│   │   ├── ui/                    # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── input.tsx
│   │   │   ├── toast.tsx
│   │   │   └── ...
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   ├── Footer.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── Layout.tsx
│   │   ├── books/
│   │   │   ├── BookCard.tsx
│   │   │   ├── BookGrid.tsx
│   │   │   ├── BookDetails.tsx
│   │   │   ├── BookFilters.tsx
│   │   │   └── BookSearch.tsx
│   │   ├── cart/
│   │   │   ├── CartDrawer.tsx
│   │   │   ├── CartItem.tsx
│   │   │   ├── CartSummary.tsx
│   │   │   └── EmptyCart.tsx
│   │   ├── orders/
│   │   │   ├── OrderCard.tsx
│   │   │   ├── OrderDetails.tsx
│   │   │   ├── OrderStatus.tsx
│   │   │   └── OrderHistory.tsx
│   │   └── auth/
│   │       ├── LoginForm.tsx
│   │       ├── RegisterForm.tsx
│   │       └── ProtectedRoute.tsx
│   ├── pages/
│   │   ├── HomePage.tsx
│   │   ├── BookDetailsPage.tsx
│   │   ├── SearchPage.tsx
│   │   ├── CartPage.tsx
│   │   ├── CheckoutPage.tsx
│   │   ├── OrderHistoryPage.tsx
│   │   ├── OrderDetailsPage.tsx
│   │   ├── ProfilePage.tsx
│   │   └── NotFoundPage.tsx
│   ├── api/
│   │   ├── client.ts              # Axios instance
│   │   ├── books.ts
│   │   ├── auth.ts
│   │   ├── cart.ts
│   │   ├── orders.ts
│   │   └── users.ts
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useCart.ts
│   │   ├── useBooks.ts
│   │   ├── useOrders.ts
│   │   └── useDebounce.ts
│   ├── store/
│   │   ├── authStore.ts
│   │   ├── cartStore.ts
│   │   └── index.ts
│   ├── types/
│   │   ├── book.ts
│   │   ├── user.ts
│   │   ├── order.ts
│   │   ├── cart.ts
│   │   └── api.ts
│   ├── utils/
│   │   ├── format.ts
│   │   ├── validation.ts
│   │   └── constants.ts
│   ├── App.tsx
│   ├── main.tsx
│   └── index.css
├── public/
├── Dockerfile
├── nginx.conf
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── vite.config.ts
```

### Key Features Implementation

#### API Client with Interceptors

```typescript
// src/api/client.ts
import axios from 'axios';
import { useAuthStore } from '../store/authStore';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor (add JWT token)
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor (handle errors)
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token expired, logout user
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

#### Auth Store (Zustand)

```typescript
// src/store/authStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { User } from '../types/user';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
  updateUser: (user: Partial<User>) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user, token) =>
        set({ user, token, isAuthenticated: true }),
      logout: () =>
        set({ user: null, token: null, isAuthenticated: false }),
      updateUser: (updatedUser) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...updatedUser } : null,
        })),
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

#### Book Search with TanStack Query

```typescript
// src/hooks/useBooks.ts
import { useQuery } from '@tanstack/react-query';
import { searchBooks } from '../api/books';
import { BookSearchParams } from '../types/book';

export const useBooks = (params: BookSearchParams) => {
  return useQuery({
    queryKey: ['books', params],
    queryFn: () => searchBooks(params),
    staleTime: 5 * 60 * 1000, // 5 minutes
    keepPreviousData: true,
  });
};

export const useBook = (bookId: string) => {
  return useQuery({
    queryKey: ['book', bookId],
    queryFn: () => getBookById(bookId),
    enabled: !!bookId,
  });
};
```

## Development Workflow

### Initial Setup

```bash
# Clone repository
git clone <repo-url>
cd distributed-bookstore

# Start infrastructure (local development)
docker-compose up -d postgres redis rabbitmq

# Generate protobuf files
make proto-gen

# Initialize databases (run migrations)
./scripts/init-databases.sh

# Start all services
make services-start

# Start frontend
cd frontend/customer-app
npm install
npm run dev
```

### Common Commands (Makefile)

```makefile
# Makefile

.PHONY: proto-gen services-start services-stop docker-build k8s-deploy

# Generate protobuf files for all services
proto-gen:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/*.proto
	@echo "Done!"

# Start all services locally
services-start:
	docker-compose up -d

# Stop all services
services-stop:
	docker-compose down

# Build all Docker images
docker-build:
	@echo "Building all services..."
	./scripts/build-all.sh

# Deploy to Kubernetes
k8s-deploy:
	kubectl apply -f infrastructure/k8s/namespaces/
	kubectl apply -f infrastructure/k8s/databases/
	kubectl apply -f infrastructure/k8s/messaging/
	kubectl apply -f infrastructure/k8s/services/
	kubectl apply -f infrastructure/k8s/ingress/

# Run tests for all Go services
test:
	@echo "Running tests..."
	@for service in services/*/; do \
		if [ -f $$service/go.mod ]; then \
			echo "Testing $$service..."; \
			cd $$service && go test ./... -v -cover; \
		fi \
	done

# Lint all Go code
lint:
	@echo "Linting..."
	@for service in services/*/; do \
		if [ -f $$service/go.mod ]; then \
			echo "Linting $$service..."; \
			cd $$service && golangci-lint run; \
		fi \
	done

# Run database migrations
migrate-up:
	@echo "Running migrations..."
	./scripts/run-migrations.sh up

migrate-down:
	@echo "Reverting migrations..."
	./scripts/run-migrations.sh down

# Clean up
clean:
	docker-compose down -v
	rm -rf */bin */dist
```

## Distributed Systems Observability

### Jaeger Tracing Configuration

```go
// internal/tracing/jaeger.go
package tracing

import (
    "io"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
)

func InitJaeger(serviceName string) (opentracing.Tracer, io.Closer, error) {
    cfg := &config.Configuration{
        ServiceName: serviceName,
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans:           true,
            LocalAgentHostPort: "jaeger-agent:6831",
        },
    }
    
    tracer, closer, err := cfg.NewTracer(
        config.Logger(jaeger.StdLogger),
    )
    if err != nil {
        return nil, nil, err
    }
    
    opentracing.SetGlobalTracer(tracer)
    return tracer, closer, nil
}

// HTTP middleware for tracing
func TracingMiddleware(tracer opentracing.Tracer) fiber.Handler {
    return func(c *fiber.Ctx) error {
        span := tracer.StartSpan(c.Path())
        defer span.Finish()
        
        ctx := opentracing.ContextWithSpan(c.Context(), span)
        c.SetUserContext(ctx)
        
        span.SetTag("http.method", c.Method())
        span.SetTag("http.url", c.OriginalURL())
        
        err := c.Next()
        
        span.SetTag("http.status_code", c.Response().StatusCode())
        if err != nil {
            span.SetTag("error", true)
            span.LogKV("error", err.Error())
        }
        
        return err
    }
}
```

### Prometheus Metrics

```go
// internal/metrics/prometheus.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"service", "method", "endpoint", "status"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
        },
        []string{"service", "method", "endpoint"},
    )
    
    GRPCRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "grpc_requests_total",
            Help: "Total number of gRPC requests",
        },
        []string{"service", "method", "status"},
    )
    
    DBConnectionsActive = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "db_connections_active",
            Help: "Number of active database connections",
        },
    )
    
    OrdersCreatedTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "orders_created_total",
            Help: "Total number of orders created",
        },
    )
)

// Middleware for HTTP metrics
func MetricsMiddleware(serviceName string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Response().StatusCode())
        
        HTTPRequestsTotal.WithLabelValues(
            serviceName,
            c.Method(),
            c.Path(),
            status,
        ).Inc()
        
        HTTPRequestDuration.WithLabelValues(
            serviceName,
            c.Method(),
            c.Path(),
        ).Observe(duration)
        
        return err
    }
}
```

## Testing Strategy

### Unit Tests (Go)

```go
// internal/service/catalog_service_test.go
package service_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourusername/bookstore/catalog-service/internal/domain"
    "github.com/yourusername/bookstore/catalog-service/internal/service"
    "github.com/yourusername/bookstore/catalog-service/internal/repository/mocks"
)

func TestGetBook(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.MockCatalogRepository)
    svc := service.NewCatalogService(mockRepo)
    
    expectedBook := &domain.Book{
        ID:    "123",
        Title: "Test Book",
        Price: 29.99,
    }
    
    mockRepo.On("GetByID", mock.Anything, "123").
        Return(expectedBook, nil)
    
    // Act
    book, err := svc.GetBook(context.Background(), "123")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedBook, book)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests (Go)

```go
// tests/integration/catalog_test.go
package integration_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/postgres"
)

type CatalogTestSuite struct {
    suite.Suite
    pgContainer *postgres.PostgresContainer
    repo        repository.CatalogRepository
}

func (suite *CatalogTestSuite) SetupSuite() {
    ctx := context.Background()
    
    // Start PostgreSQL container
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("test_db"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    suite.Require().NoError(err)
    suite.pgContainer = pgContainer
    
    // Connect to database
    connStr, err := pgContainer.ConnectionString(ctx)
    suite.Require().NoError(err)
    
    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    suite.Require().NoError(err)
    
    // Run migrations
    err = db.AutoMigrate(&domain.Book{}, &domain.Author{})
    suite.Require().NoError(err)
    
    suite.repo = postgres.NewCatalogRepository(db)
}

func (suite *CatalogTestSuite) TearDownSuite() {
    suite.pgContainer.Terminate(context.Background())
}

func (suite *CatalogTestSuite) TestCreateAndGetBook() {
    ctx := context.Background()
    
    book := &domain.Book{
        Title: "Integration Test Book",
        ISBN:  "1234567890123",
        Price: 39.99,
    }
    
    // Create
    err := suite.repo.Create(ctx, book)
    suite.Require().NoError(err)
    suite.Require().NotEmpty(book.ID)
    
    // Get
    retrieved, err := suite.repo.GetByID(ctx, book.ID)
    suite.Require().NoError(err)
    suite.Equal(book.Title, retrieved.Title)
    suite.Equal(book.ISBN, retrieved.ISBN)
}

func TestCatalogTestSuite(t *testing.T) {
    suite.Run(t, new(CatalogTestSuite))
}
```

### E2E Tests (Frontend - Cypress)

```typescript
// cypress/e2e/checkout.cy.ts
describe('Checkout Flow', () => {
  beforeEach(() => {
    // Login
    cy.visit('/login');
    cy.get('[data-cy=email]').type('test@example.com');
    cy.get('[data-cy=password]').type('password123');
    cy.get('[data-cy=login-button]').click();
    cy.url().should('include', '/');
  });

  it('should complete order successfully', () => {
    // Search for book
    cy.get('[data-cy=search-input]').type('Distributed Systems');
    cy.get('[data-cy=search-button]').click();
    
    // Add to cart
    cy.get('[data-cy=book-card]').first().click();
    cy.get('[data-cy=add-to-cart-button]').click();
    cy.get('[data-cy=cart-badge]').should('contain', '1');
    
    // Go to cart
    cy.get('[data-cy=cart-icon]').click();
    cy.url().should('include', '/cart');
    cy.get('[data-cy=cart-item]').should('have.length', 1);
    
    // Proceed to checkout
    cy.get('[data-cy=checkout-button]').click();
    cy.url().should('include', '/checkout');
    
    // Fill shipping address
    cy.get('[data-cy=address-line1]').type('123 Main St');
    cy.get('[data-cy=city]').type('New York');
    cy.get('[data-cy=postal-code]').type('10001');
    
    // Enter payment
    cy.get('[data-cy=card-number]').type('4242424242424242');
    cy.get('[data-cy=card-expiry]').type('12/25');
    cy.get('[data-cy=card-cvc]').type('123');
    
    // Place order
    cy.get('[data-cy=place-order-button]').click();
    
    // Verify success
    cy.url().should('include', '/orders/');
    cy.get('[data-cy=order-confirmation]').should('be.visible');
    cy.get('[data-cy=order-status]').should('contain', 'Order Confirmed');
  });
});
```

## Security Best Practices

### JWT Implementation

```go
// pkg/jwt/jwt.go
package jwt

import (
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string   `json:"user_id"`
    Email  string   `json:"email"`
    Roles  []string `json:"roles"`
    jwt.RegisteredClaims
}

func GenerateToken(userID, email string, roles []string, secret string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Roles:  roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "bookstore-api",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrSignatureInvalid
}
```

### Input Validation

```go
// pkg/validator/validator.go
package validator

import (
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

type CreateBookRequest struct {
    ISBN        string  `json:"isbn" validate:"required,len=13"`
    Title       string  `json:"title" validate:"required,min=1,max=500"`
    Description string  `json:"description" validate:"max=5000"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"gte=0"`
}

func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}
```

### Rate Limiting

```go
// internal/middleware/ratelimit.go
package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": "Rate limit exceeded",
            })
        },
    })
}
```

## Performance Optimization

### Database Query Optimization

```go
// Use indexes
CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_price ON books(price);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

// Use EXPLAIN ANALYZE
EXPLAIN ANALYZE SELECT * FROM books WHERE price > 20 ORDER BY title;

// Preload relationships to avoid N+1
books := []domain.Book{}
db.Preload("Authors").Preload("Categories").Find(&books)
```

### Caching Strategy

```go
// internal/cache/redis.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/redis/go-redis/v9"
)

type CatalogCache struct {
    client *redis.Client
}

func (c *CatalogCache) GetBook(ctx context.Context, bookID string) (*domain.Book, error) {
    key := fmt.Sprintf("book:%s", bookID)
    
    val, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Cache miss
    } else if err != nil {
        return nil, err
    }
    
    var book domain.Book
    if err := json.Unmarshal([]byte(val), &book); err != nil {
        return nil, err
    }
    
    return &book, nil
}

func (c *CatalogCache) SetBook(ctx context.Context, book *domain.Book) error {
    key := fmt.Sprintf("book:%s", book.ID)
    
    data, err := json.Marshal(book)
    if err != nil {
        return err
    }
    
    return c.client.Set(ctx, key, data, 15*time.Minute).Err()
}
```

## Troubleshooting

### Common Issues

**Service can't connect to database:**
```bash
# Check database is running
docker ps | grep postgres

# Check connection string
echo $DATABASE_URL

# Test connection
psql $DATABASE_URL

# Check network
docker network inspect bookstore-network
```

**JWT token invalid:**
```bash
# Verify JWT_SECRET matches across services
kubectl get secret jwt-secret -n production -o yaml

# Check token expiration
# Decode JWT at jwt.io
```

**RabbitMQ consumer not receiving messages:**
```bash
# Check RabbitMQ management UI
http://localhost:15672

# Check exchange and queue bindings
# Verify routing keys match

# Check consumer logs
kubectl logs -f deployment/notification-service -n production
```

**High latency:**
```bash
# Check Jaeger traces
http://localhost:16686

# Check Prometheus metrics
http://localhost:9090

# Check database slow queries
SELECT * FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;

# Check Redis cache hit rate
INFO stats
```

## Project Roadmap & Tasks

### Phase 1: Foundation (Week 1-2) ✅ IN PROGRESS

#### Infrastructure Setup
- [ ] Initialize Git repository with proper .gitignore
- [ ] Create project folder structure
- [ ] Setup Docker Compose for local development
- [ ] Configure PostgreSQL with init scripts
- [ ] Configure Redis
- [ ] Configure RabbitMQ with management UI
- [ ] Create Makefile with common commands
- [ ] Write setup documentation in README.md

#### Protobuf Definitions
- [ ] Define `common.proto` (shared types)
- [ ] Define `catalog.proto`
- [ ] Define `user.proto`
- [ ] Define `cart.proto`
- [ ] Define `order.proto`
- [ ] Define `payment.proto`
- [ ] Define `inventory.proto`
- [ ] Create script to generate Go code from protos
- [ ] Test protobuf compilation

#### User Service (Golang)
- [ ] Create project structure
- [ ] Implement domain models (User, Role, Address)
- [ ] Create PostgreSQL repository
- [ ] Write database migrations
- [ ] Implement JWT generation/validation
- [ ] Create HTTP handlers (register, login, profile)
- [ ] Create gRPC server (ValidateToken, GetUser)
- [ ] Add middleware (auth, logging, metrics)
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test locally

#### Catalog Service (Golang)
- [ ] Create project structure
- [ ] Implement domain models (Book, Author, Category)
- [ ] Create PostgreSQL repository
- [ ] Write database migrations
- [ ] Implement search functionality
- [ ] Create HTTP handlers (CRUD, search)
- [ ] Create gRPC server
- [ ] Add Redis caching
- [ ] Integrate Jaeger tracing
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test locally

### Phase 2: Core Services (Week 3-4)

#### API Gateway (Golang)
- [ ] Create project structure
- [ ] Implement request routing
- [ ] Add JWT validation middleware
- [ ] Implement rate limiting (Redis-backed)
- [ ] Add CORS handling
- [ ] Implement health check aggregation
- [ ] Add circuit breaker for downstream services
- [ ] Integrate Jaeger tracing
- [ ] Add Prometheus metrics
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test end-to-end flow

#### Cart Service (Golang)
- [ ] Create project structure
- [ ] Implement domain models
- [ ] Create Redis repository
- [ ] Implement cart operations (add, update, remove)
- [ ] Create HTTP handlers
- [ ] Create gRPC server
- [ ] Integrate with Catalog Service (gRPC)
- [ ] Add session management
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test locally

#### Order Service (Golang)
- [ ] Create project structure
- [ ] Implement domain models
- [ ] Create PostgreSQL repository
- [ ] Write database migrations
- [ ] Implement saga orchestrator
- [ ] Create RabbitMQ event publisher
- [ ] Create RabbitMQ event consumer
- [ ] Create HTTP handlers
- [ ] Create gRPC server
- [ ] Integrate with other services (gRPC)
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test order creation flow

### Phase 3: Payment & Inventory (Week 5)

#### Payment Service (Node.js/TypeScript)
- [ ] Initialize Node.js project with TypeScript
- [ ] Setup project structure
- [ ] Implement Stripe integration
- [ ] Create domain models
- [ ] Create PostgreSQL repository
- [ ] Implement payment processing logic
- [ ] Create HTTP handlers
- [ ] Create gRPC server
- [ ] Create RabbitMQ event consumer
- [ ] Create RabbitMQ event publisher
- [ ] Integrate Jaeger tracing
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test payment flow

#### Inventory Service (Golang)
- [ ] Create project structure
- [ ] Implement domain models
- [ ] Create PostgreSQL repository
- [ ] Write database migrations
- [ ] Implement stock reservation logic
- [ ] Create HTTP handlers
- [ ] Create gRPC server
- [ ] Create RabbitMQ event consumer
- [ ] Create RabbitMQ event publisher
- [ ] Implement low stock alerts
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test inventory flow

### Phase 4: Notifications & Reviews (Week 6)

#### Notification Service (Node.js/TypeScript)
- [ ] Initialize Node.js project
- [ ] Setup project structure
- [ ] Integrate SendGrid/SES for email
- [ ] Create email templates (Handlebars)
- [ ] Create PostgreSQL repository
- [ ] Implement notification history
- [ ] Create RabbitMQ consumers for all events
- [ ] Implement template rendering
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test email sending

#### Review Service (Python/FastAPI)
- [ ] Initialize Python project
- [ ] Setup FastAPI project structure
- [ ] Implement domain models (SQLAlchemy)
- [ ] Create PostgreSQL repository
- [ ] Write database migrations (Alembic)
- [ ] Implement sentiment analysis (NLTK/TextBlob)
- [ ] Create REST endpoints
- [ ] Create gRPC server
- [ ] Integrate Jaeger tracing
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test review submission

### Phase 5: ML & Admin (Week 7)

#### Recommendation Service (Python/FastAPI)
- [ ] Initialize Python project
- [ ] Setup project structure
- [ ] Implement collaborative filtering (scikit-learn)
- [ ] Implement content-based filtering
- [ ] Create PostgreSQL repository with pgvector
- [ ] Implement recommendation caching
- [ ] Create REST endpoints
- [ ] Create gRPC server
- [ ] Create RabbitMQ event consumer
- [ ] Train initial ML models
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test recommendations

#### Admin Service (Golang)
- [ ] Create project structure
- [ ] Implement analytics aggregation
- [ ] Create HTTP handlers
- [ ] Integrate with all services (gRPC)
- [ ] Implement dashboard metrics
- [ ] Implement bulk operations
- [ ] Write unit tests
- [ ] Create Dockerfile
- [ ] Add to docker-compose.yml
- [ ] Test admin functions

### Phase 6: Frontend (Week 8-9)

#### Customer App (React + shadcn/ui)
- [ ] Initialize Vite + React + TypeScript project
- [ ] Install and configure shadcn/ui
- [ ] Setup Tailwind CSS
- [ ] Create project structure
- [ ] Implement API client with Axios
- [ ] Setup TanStack Query
- [ ] Setup Zustand stores
- [ ] Create UI components
  - [ ] Layout components
  - [ ] Book components
  - [ ] Cart components
  - [ ] Order components
  - [ ] Auth components
- [ ] Create pages
  - [ ] Home page
  - [ ] Book details page
  - [ ] Search page
  - [ ] Cart page
  - [ ] Checkout page
  - [ ] Order history page
  - [ ] Profile page
- [ ] Implement authentication flow
- [ ] Implement shopping flow
- [ ] Add loading states and error handling
- [ ] Write component tests
- [ ] Write E2E tests (Cypress)
- [ ] Create production Dockerfile (Nginx)
- [ ] Add to docker-compose.yml
- [ ] Test locally

### Phase 7: Kubernetes (Week 10-11)

#### K8s Base Configuration
- [ ] Create namespace manifests
- [ ] Create ConfigMaps for each service
- [ ] Create Secrets for sensitive data
- [ ] Create NetworkPolicies
- [ ] Setup PostgreSQL StatefulSet
- [ ] Setup Redis Deployment
- [ ] Setup RabbitMQ StatefulSet

#### Service Deployments
- [ ] Create Deployment for API Gateway
- [ ] Create Deployment for Catalog Service
- [ ] Create Deployment for User Service
- [ ] Create Deployment for Cart Service
- [ ] Create Deployment for Order Service
- [ ] Create Deployment for Payment Service
- [ ] Create Deployment for Inventory Service
- [ ] Create Deployment for Notification Service
- [ ] Create Deployment for Review Service
- [ ] Create Deployment for Recommendation Service
- [ ] Create Deployment for Admin Service
- [ ] Create Deployment for Frontend

#### K8s Services & Ingress
- [ ] Create Service manifests for all microservices
- [ ] Install Nginx Ingress Controller
- [ ] Create Ingress manifest
- [ ] Configure SSL/TLS certificates

#### HPA & Autoscaling
- [ ] Create HPA for API Gateway
- [ ] Create HPA for Catalog Service
- [ ] Create HPA for Order Service
- [ ] Create HPA for other services
- [ ] Test autoscaling

### Phase 8: Observability (Week 12)

#### Distributed Tracing
- [ ] Deploy Jaeger to K8s
- [ ] Configure all services with Jaeger
- [ ] Test trace propagation
- [ ] Create Jaeger dashboards

#### Metrics & Monitoring
- [ ] Deploy Prometheus to K8s
- [ ] Deploy Grafana to K8s
- [ ] Configure ServiceMonitors
- [ ] Create Grafana dashboards
  - [ ] Service health dashboard
  - [ ] Business metrics dashboard
  - [ ] Infrastructure dashboard

#### Logging
- [ ] Deploy Elasticsearch to K8s
- [ ] Deploy Fluentd DaemonSet
- [ ] Deploy Kibana
- [ ] Configure log collection
- [ ] Create Kibana dashboards
- [ ] Setup log retention policies

### Phase 9: AWS Deployment (Week 13-14)

#### AWS Infrastructure
- [ ] Create AWS account and configure CLI
- [ ] Setup VPC with public/private subnets
- [ ] Create EKS cluster
- [ ] Configure node groups
- [ ] Setup RDS PostgreSQL instances
- [ ] Setup ElastiCache Redis
- [ ] Setup Amazon MQ (RabbitMQ)
- [ ] Create S3 buckets
- [ ] Setup ECR repositories
- [ ] Configure CloudFront distribution
- [ ] Setup Route53 hosted zone
- [ ] Request ACM certificates

#### CI/CD Pipeline
- [ ] Create GitHub Actions workflow
- [ ] Configure build jobs
- [ ] Configure test jobs
- [ ] Configure Docker build and push to ECR
- [ ] Configure K8s deployment
- [ ] Setup staging environment
- [ ] Setup production environment
- [ ] Configure deployment approval

#### Security & Compliance
- [ ] Configure AWS Secrets Manager
- [ ] Setup IAM roles and policies
- [ ] Configure Security Groups
- [ ] Setup VPC peering if needed
- [ ] Enable CloudTrail
- [ ] Configure AWS WAF
- [ ] Setup backup policies

### Phase 10: Testing & Documentation (Week 15)

#### Testing
- [ ] Run full test suite
- [ ] Perform load testing (k6)
- [ ] Perform security testing
- [ ] Test disaster recovery
- [ ] Test auto-scaling under load
- [ ] Test failover scenarios

#### Documentation
- [ ] Complete API documentation (Swagger/OpenAPI)
- [ ] Write deployment guide
- [ ] Write troubleshooting guide
- [ ] Create architecture diagrams
- [ ] Write developer onboarding guide
- [ ] Document monitoring and alerting
- [ ] Create runbooks for common operations

#### Final Deliverables
- [ ] Record demo video
- [ ] Create presentation
- [ ] Write final project report
- [ ] Submit to Interactiva Virtual

---

## Current Status

**Last Updated**: [DATE]

**Completed**: 0/200+ tasks

**Current Phase**: Phase 1 - Foundation

**Next Steps**:
1. Initialize repository structure
2. Setup Docker Compose environment
3. Define protobuf contracts
4. Start User Service development

---

## Notes for Claude Code

- **Language Versions**: Go 1.21+, Node 18+ LTS, Python 3.11+
- **Package Management**: Go modules, npm, pip
- **Environment Variables**: Never commit `.env` files, use `.env.example` templates
- **Database Migrations**: Use `golang-migrate` for Go services, Alembic for Python
- **Code Generation**: Use `protoc` for gRPC stubs
- **Naming Conventions**:
  - Go: camelCase (private), PascalCase (public)
  - TypeScript: camelCase (variables/functions), PascalCase (types/components)
  - Python: snake_case (variables/functions), PascalCase (classes)
  - Database: snake_case (tables and columns)
  - URLs: kebab-case (routes)
  - K8s resources: kebab-case
- **Git Workflow**: Feature branches, PR reviews, squash merges
- **Commit Convention**: Conventional Commits (feat:, fix:, docs:, etc.)

## Quick Reference

**Service Ports**:
- API Gateway: 8080
- Catalog: 8081 (HTTP), 50051 (gRPC)
- User: 8082 (HTTP), 50052 (gRPC)
- Cart: 8083 (HTTP), 50053 (gRPC)
- Order: 8084 (HTTP), 50054 (gRPC)
- Payment: 8085 (HTTP), 50055 (gRPC)
- Inventory: 8086 (HTTP), 50056 (gRPC)
- Notification: 8087
- Review: 8088 (HTTP), 50058 (gRPC)
- Recommendation: 8089 (HTTP), 50059 (gRPC)
- Admin: 8090 (HTTP), 50060 (gRPC)
- Frontend: 3000
- PostgreSQL: 5432
- Redis: 6379
- RabbitMQ: 5672 (AMQP), 15672 (Management)
- Jaeger: 16686 (UI), 6831 (Agent)
- Prometheus: 9090
- Grafana: 3001

**Essential Commands**:
```bash
# Local development
make services-start         # Start all services
make services-stop          # Stop all services
make proto-gen             # Generate protobuf code
make test                  # Run tests
make lint                  # Run linters

# Database
make migrate-up            # Run migrations
make migrate-down          # Rollback migrations

# Kubernetes
make k8s-deploy           # Deploy to K8s
kubectl get pods -n production
kubectl logs -f <pod-name> -n production
kubectl port-forward svc/<service> 8080:8080
```

---

**Remember**: This is a learning project focused on distributed systems principles. Every architectural decision should be documented with ADRs (Architecture Decision Records). Focus on understanding WHY each pattern is used, not just HOW to implement it.
