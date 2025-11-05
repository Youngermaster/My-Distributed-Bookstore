# Service Interactions and Distributed Systems Best Practices

## Overview

This document describes the architecture, inter-service communication patterns, and distributed systems best practices implemented in the Distributed Bookstore system, based on principles from Andrew S. Tanenbaum's "Distributed Systems" textbook.

---

## System Architecture

### Service Inventory

#### **Golang Services (Go + Fiber)**

1. **API Gateway** (Port: 8080)
2. **Catalog Service** (Port: 8081)
3. **User Service** (Port: 8082)
4. **Cart Service** (Port: 8083)
5. **Order Service** (Port: 8084)
6. **Notification Service** (Port: 8088)

#### **FastAPI Services (Python)**

7. **Inventory Service** (Port: 8086)
8. **Review Service** (Port: 8087)
9. **Recommendation Service** (Port: 8085)

---

## Communication Patterns

### 1. Synchronous Communication (HTTP/REST)

#### **API Gateway as Entry Point**

The API Gateway serves as the **single entry point** for all client requests, implementing the **API Gateway pattern**:

```
Client → API Gateway (8080) → Backend Services
```

**Gateway Routing:**

```go
/api/v1/catalog/*        → Catalog Service (8081)
/api/v1/users/*          → User Service (8082)
/api/v1/cart/*           → Cart Service (8083)
/api/v1/orders/*         → Order Service (8084)
/api/v1/recommendations/* → Recommendation Service (8085)
/api/v1/inventory/*      → Inventory Service (8086)
/api/v1/reviews/*        → Review Service (8087)
/api/v1/notifications/*  → Notification Service (8088)
```

**Implementation:**

- **Proxy Pattern**: The gateway uses proxy clients for each service
- **Request Forwarding**: Headers, body, and query parameters are preserved
- **Timeout Management**: 30-second timeout per request
- **Rate Limiting**: Configurable rate limits (default: disabled in dev)

**Example:**

```go
// API Gateway forwards requests via proxy clients
catalogClient := proxy.NewCatalogClient(cfg.CatalogServiceURL, cfg.RequestTimeout)
responseBody, statusCode, err := catalogClient.ProxyRequest(method, path, body, headers)
```

---

### 2. Asynchronous Communication (Event-Driven with RabbitMQ)

The system implements **publish-subscribe pattern** for asynchronous communication using RabbitMQ with a **topic exchange**.

#### **Event Publishers (Golang Services)**

**User Service:**

- **Events Published:**
  - `user.registered` - New user registration
  - `user.logged_in` - User login
  - `wishlist.item_added` - Item added to wishlist
  - `wishlist.item_removed` - Item removed from wishlist

**Cart Service:**

- **Events Published:**
  - `cart.item_added` - Item added to cart
  - `cart.item_removed` - Item removed from cart
  - `cart.cleared` - Cart cleared

**Event Structure:**

```go
type NotificationEvent struct {
    EventID    string                 `json:"event_id"`
    EventType  string                 `json:"event_type"`     // Routing key
    Source     string                 `json:"source"`         // Service name
    OccurredAt time.Time              `json:"occurred_at"`
    UserID     string                 `json:"user_id,omitempty"`
    Email      string                 `json:"email,omitempty"`
    Payload    map[string]interface{} `json:"payload,omitempty"`
}
```

**Publishing Mechanism:**

```go
// RabbitMQ Connection with Topic Exchange
exchange := "bookstore.events"
conn, _ := amqp091.Dial(rabbitMQURL)
channel, _ := conn.Channel()

// Declare topic exchange
channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)

// Publish with routing key
channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
    ContentType: "application/json",
    Body:        eventJSON,
    MessageId:   event.EventID,
    Timestamp:   event.OccurredAt,
})
```

#### **Event Consumer (Notification Service)**

The Notification Service subscribes to specific routing keys:

```go
routingKeys := []string{
    "user.registered",
    "user.logged_in",
    "cart.item_added",
    "cart.item_removed",
    "wishlist.item_added",
}

// Queue binding with topic patterns
for _, key := range routingKeys {
    channel.QueueBind(queueName, key, exchange, false, nil)
}
```

**Consumer Implementation:**

- **Prefetch Count**: 10 messages (QoS)
- **Manual Acknowledgment**: Messages acknowledged after successful processing
- **Error Handling**: Failed messages are nack'd and requeued
- **Graceful Shutdown**: Context-based cancellation

**Notification Dispatcher:**

```go
// Processes events and sends notifications
dispatcher.HandleEvent(ctx, event)
  → Log-based notification (development)
  → Email, SMS, Push (production-ready structure)
```

---

## Data Storage Patterns

### **Database per Service Pattern**

Each service has its own database/storage:

| Service        | Storage            | Purpose                                       |
| -------------- | ------------------ | --------------------------------------------- |
| Catalog        | PostgreSQL         | Books, authors, categories, publishers        |
| User           | PostgreSQL         | Users, roles, sessions, addresses, wishlist   |
| Cart           | Redis              | Shopping cart items (TTL: configurable hours) |
| Order          | PostgreSQL         | Orders, order items, order status             |
| Inventory      | PostgreSQL         | Stock levels, reservations, movements         |
| Review         | PostgreSQL         | Reviews, ratings, sentiment analysis          |
| Recommendation | Redis + PostgreSQL | Cache + user preferences, book tags           |

### **Storage Characteristics**

**PostgreSQL Services:**

- **ACID Transactions**: Order, User, Catalog, Inventory
- **Foreign Keys**: Referential integrity
- **Migrations**: Gorm AutoMigrate (dev), Alembic (Python services)
- **Connection Pooling**: Configurable max connections

**Redis Services:**

- **In-Memory**: Cart and Recommendation caching
- **TTL**: Automatic expiration (cart: 24-72 hours)
- **Pub/Sub**: Future real-time notifications
- **Serialization**: JSON storage for complex objects

---

## Distributed Systems Best Practices (Tanenbaum Principles)

### 1. **Transparency**

#### **Location Transparency**

- Clients interact only with API Gateway, unaware of backend service locations
- Kubernetes Services provide DNS-based discovery (e.g., `catalog-service.bookstore-dev.svc.cluster.local:8081`)
- Environment variables externalize service URLs

#### **Access Transparency**

- Uniform REST API interface across all services
- Consistent response formats (JSON)
- Standard HTTP methods (GET, POST, PUT, DELETE, PATCH)

#### **Failure Transparency** (Partial Implementation)

- Health checks (`/health`, `/ready`) for liveness/readiness probes
- Graceful shutdown with signal handling (SIGINT, SIGTERM)
- Error propagation with proper HTTP status codes

**Example:**

```go
// Graceful shutdown pattern
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
<-quit
app.Shutdown() // Closes connections gracefully
```

### 2. **Scalability**

#### **Horizontal Scaling**

```yaml
# Kubernetes HPA (Horizontal Pod Autoscaler)
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
spec:
  minReplicas: 1
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

**Scaled Services:**

- Catalog Service: 2+ replicas
- User Service: 2+ replicas
- Cart Service: 2+ replicas
- Order Service: 2+ replicas

#### **Stateless Design**

- All Golang services are stateless
- Session data stored in database (User Service)
- Cart data in Redis (shared across instances)
- No in-memory state prevents scaling issues

#### **Load Balancing**

- Kubernetes Services provide round-robin load balancing
- API Gateway distributes load across service replicas

### 3. **Fault Tolerance**

#### **Replication**

```yaml
spec:
  replicas: 2 # Multiple instances
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0 # Zero-downtime deployments
```

#### **Health Monitoring**

```go
// Liveness probe: Is service alive?
livenessProbe:
  httpGet:
    path: /health
    port: 8081
  initialDelaySeconds: 30
  periodSeconds: 10

// Readiness probe: Can service handle traffic?
readinessProbe:
  httpGet:
    path: /ready
    port: 8081
  initialDelaySeconds: 10
  periodSeconds: 5
```

#### **Message Durability**

- RabbitMQ queues are **durable** (survive broker restart)
- Messages persisted to disk
- Manual acknowledgment prevents message loss
- Dead-letter queues for failed messages (can be added)

```go
// Durable queue declaration
channel.QueueDeclare(
    queueName,
    true,   // durable
    false,  // delete when unused
    false,  // exclusive
    false,  // no-wait
    nil,
)
```

### 4. **Consistency**

#### **Eventual Consistency**

- Event-driven updates are **eventually consistent**
- Example: Cart clearing triggers async notification
- Acceptable for non-critical notifications

#### **Strong Consistency**

- Database transactions for critical operations
- Order creation, payment processing, inventory reservation
- ACID properties maintained at service level

**Example:**

```go
// Transaction in Order Service
tx := db.Begin()
order, err := tx.CreateOrder(...)
if err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

### 5. **Security**

#### **Authentication & Authorization**

```go
// JWT-based authentication
app.Use(middleware.AuthMiddleware(jwtService))

// Token validation
claims, err := jwtService.ValidateToken(token)
if err != nil {
    return ErrUnauthorized
}
```

**Security Features:**

- **JWT Tokens**: Access (15 min) + Refresh (7 days)
- **Password Hashing**: bcrypt with configurable cost (default: 10)
- **Role-Based Access Control**: User roles (customer, admin, etc.)
- **HTTPS**: Kubernetes Ingress with TLS (production)
- **Secrets Management**: Kubernetes Secrets for sensitive data

#### **Secret Management**

```yaml
# Kubernetes Secret
apiVersion: v1
kind: Secret
metadata:
  name: postgres-credentials
type: Opaque
data:
  DATABASE_URL: <base64-encoded>

# Environment variable injection
env:
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: postgres-credentials
        key: DATABASE_URL
```

### 6. **Communication Efficiency**

#### **Connection Pooling**

```go
// Database connection pooling
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

#### **Resource Limits**

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

#### **Compression & Optimization**

- JSON serialization for network efficiency
- Redis binary protocol for cache
- GZIP compression in API Gateway (can be enabled)

### 7. **Observability**

#### **Structured Logging**

```go
// Fiber logger middleware
app.Use(logger.New(logger.Config{
    Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
    TimeFormat: "2006-01-02 15:04:05",
}))
```

**Python Services:**

```python
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
```

#### **Health Checks**

All services expose:

- `/health` - Basic health status
- `/ready` - Readiness check (DB, Redis, RabbitMQ connectivity)

#### **Metrics (Ready for Prometheus)**

- Service endpoints tracked
- Response times logged
- Error rates visible in logs
- Infrastructure ready for Prometheus scraping

### 8. **Naming & Discovery**

#### **Service Discovery**

```yaml
# Kubernetes Service
apiVersion: v1
kind: Service
metadata:
  name: catalog-service
  namespace: bookstore-dev
spec:
  selector:
    app: catalog-service
  ports:
    - port: 8081
      targetPort: 8081
```

**DNS Resolution:**

- Internal: `catalog-service.bookstore-dev.svc.cluster.local:8081`
- Short form: `catalog-service:8081` (same namespace)

**Environment-based Configuration:**

```bash
CATALOG_SERVICE_URL=http://catalog-service:8081
USER_SERVICE_URL=http://user-service:8082
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

---

## Service-Specific Interactions

### **User Registration Flow**

```
1. Client → API Gateway → User Service
   POST /api/v1/users/auth/register

2. User Service:
   - Hash password (bcrypt)
   - Create user in PostgreSQL
   - Assign default role

3. User Service → RabbitMQ:
   Publish: user.registered {userID, email, name}

4. Notification Service ← RabbitMQ:
   Consume: user.registered
   Send welcome email (log in dev)

5. User Service → Client:
   201 Created {user, accessToken, refreshToken}
```

### **Add to Cart Flow**

```
1. Client → API Gateway → Cart Service
   POST /api/v1/cart
   Headers: Authorization: Bearer <JWT>
   Body: {bookId, quantity, title, price}

2. Cart Service:
   - Extract userId from JWT
   - Validate quantity limits
   - Store in Redis with TTL
   - Calculate total

3. Cart Service → RabbitMQ:
   Publish: cart.item_added {userId, bookId, quantity}

4. Notification Service ← RabbitMQ:
   Consume: cart.item_added
   Log event (or send push notification)

5. Cart Service → Client:
   200 OK {cart: {items, total, itemCount}}
```

### **Create Order Flow**

```
1. Client → API Gateway → Order Service
   POST /api/v1/orders
   Body: {userId, items[], shippingAddress, paymentMethod}

2. Order Service:
   - Validate items
   - Calculate totals
   - Create order in PostgreSQL (PENDING status)
   - Store order items

3. Order Service → Client:
   201 Created {order}

Note: Full order flow would involve:
- Inventory Service: Reserve stock
- Payment Service: Process payment
- Notification Service: Order confirmation
- Order Service: Update status → CONFIRMED/FAILED
```

### **Book Recommendations Flow**

```
1. Client → API Gateway → Recommendation Service
   GET /api/v1/recommendations/me
   Headers: Authorization: Bearer <JWT>

2. Recommendation Service (FastAPI):
   - Extract userId
   - Check Redis cache
   - If cached: return immediately
   - If not cached:
     * Query PostgreSQL for user preferences
     * Query book tags
     * Apply recommendation strategies:
       - Tag-based (similar genres)
       - Collaborative filtering (similar users)
       - Popularity-based
     * Store in Redis (TTL)

3. Recommendation Service → Client:
   200 OK {recommendations: [{bookId, score, reason}]}
```

### **Book Review with Sentiment Analysis Flow**

```
1. Client → API Gateway → Review Service
   POST /api/v1/reviews
   Body: {bookId, rating, comment}

2. Review Service (FastAPI):
   - Validate rating (1-5)
   - Store in PostgreSQL
   - Async ML processing:
     * NLTK/TextBlob sentiment analysis
     * Calculate polarity score
     * Classify: positive/neutral/negative
     * Update review with sentiment

3. Review Service → Client:
   201 Created {review, sentiment: "positive", polarity: 0.75}
```

---

## Deployment Architecture (Kubernetes)

### **Namespace Isolation**

```
Namespace: bookstore-dev
- Services: 9 deployments
- Databases: PostgreSQL, Redis, RabbitMQ
- Configuration: ConfigMaps, Secrets
- Networking: Services, Ingress (optional)
```

### **Resource Organization**

```
infrastructure/k8s/
├── namespaces/          # Namespace definition
├── databases/           # PostgreSQL StatefulSet
├── messaging/           # RabbitMQ, Redis
├── services/            # Microservices deployments
│   ├── catalog-service/
│   ├── user-service/
│   ├── cart-service/
│   ├── order-service/
│   ├── inventory-service/
│   ├── review-service/
│   ├── recommendation-service/
│   ├── notification-service/
│   └── api-gateway/
├── frontend/            # React customer app
├── secrets/             # Credentials
├── configmaps/          # Shared configuration
└── ingress/             # External access (optional)
```

### **Deployment Strategy**

```bash
# Sequential deployment with readiness checks
1. PostgreSQL, Redis, RabbitMQ
2. Core services (Catalog, User)
3. Business services (Cart, Order, Inventory)
4. ML services (Review, Recommendation)
5. Integration layer (API Gateway)
6. Frontend (Customer app)

# Health validation at each stage
kubectl wait --for=condition=ready pod -l app=<service> --timeout=300s
```

---

## Best Practices Summary

### **Design Principles Applied**

1. **Microservices Architecture**

   - Single Responsibility: Each service owns one domain
   - Loose Coupling: Services interact via APIs and events
   - Independent Deployment: Services deployed separately

2. **API Gateway Pattern**

   - Single entry point for clients
   - Request routing and composition
   - Cross-cutting concerns (auth, logging, rate limiting)

3. **Event-Driven Architecture**

   - Asynchronous communication for non-critical operations
   - Decoupled services via message broker
   - Topic-based routing for flexibility

4. **Database per Service**

   - Each service owns its data
   - No shared databases
   - Technology diversity (PostgreSQL, Redis)

5. **Container Orchestration**

   - Kubernetes for deployment
   - Docker for containerization
   - Infrastructure as Code

6. **Configuration Management**

   - Environment variables for service URLs
   - Kubernetes ConfigMaps for settings
   - Kubernetes Secrets for credentials

7. **Observability**

   - Health checks for monitoring
   - Structured logging
   - Ready for metrics collection (Prometheus)

8. **Resilience**
   - Graceful shutdown
   - Rolling updates (zero downtime)
   - Multiple replicas for availability
   - Health probes for failure detection

---

## Trade-offs and Considerations

### **Advantages**

- ✅ **Scalability**: Individual services scale independently
- ✅ **Fault Isolation**: Service failure doesn't cascade
- ✅ **Technology Diversity**: Go for performance, Python for ML
- ✅ **Team Autonomy**: Teams own services end-to-end
- ✅ **Deployment Flexibility**: Deploy services independently

### **Challenges**

- ⚠️ **Distributed Complexity**: Multiple services to manage
- ⚠️ **Network Latency**: More network hops than monolith
- ⚠️ **Eventual Consistency**: Event-driven updates not immediate
- ⚠️ **Testing Complexity**: Integration testing requires all services
- ⚠️ **Monitoring Overhead**: Need to monitor multiple services

### **Future Improvements**

- 🔄 **Circuit Breaker**: Add Hystrix/resilience4go to API Gateway
- 🔄 **Distributed Tracing**: Implement Jaeger for request tracing
- 🔄 **Saga Pattern**: Implement distributed transactions for order flow
- 🔄 **Service Mesh**: Consider Istio for advanced networking
- 🔄 **API Versioning**: Implement versioning strategy
- 🔄 **Rate Limiting**: Per-user rate limits in API Gateway
- 🔄 **Caching**: Add Redis caching in API Gateway
- 🔄 **Dead Letter Queues**: Handle failed events systematically

---

## Conclusion

This distributed bookstore system demonstrates solid implementation of distributed systems principles from Tanenbaum's textbook:

- **Transparency**: Clients unaware of distributed nature
- **Scalability**: Horizontal scaling with Kubernetes
- **Fault Tolerance**: Replication and health monitoring
- **Consistency**: Trade-offs between strong and eventual consistency
- **Security**: Authentication, authorization, and secrets management
- **Communication**: Efficient synchronous and asynchronous patterns
- **Observability**: Logging, health checks, and monitoring readiness

The architecture balances practical implementation with theoretical best practices, creating a maintainable, scalable, and production-ready distributed system.
