# Distributed Bookstore: Architecture & Best Practices

> A microservices-based bookstore implementing Tanenbaum's distributed systems principles with Go, Python, and Kubernetes.

## System Overview

**Architecture**: 9 microservices orchestrated by Kubernetes

- **6 Go Services** (Fiber): API Gateway, Catalog, User, Cart, Order, Notification
- **3 Python Services** (FastAPI): Inventory, Review, Recommendation

## Communication Architecture

### Synchronous: HTTP/REST via API Gateway

```
Client → API Gateway (8080) → Backend Services (8081-8088)
```

**Pattern**: Single entry point with proxy-based routing

- 30-second timeouts
- Header/body preservation
- Kubernetes DNS-based discovery

### Asynchronous: RabbitMQ Pub/Sub

```
Publishers (User, Cart) → Topic Exchange → Notification Consumer
```

**Events Published**:

- `user.registered`, `user.logged_in`
- `cart.item_added`, `cart.cleared`
- `wishlist.item_added/removed`

**Consumer**: Notification Service (email/SMS/push)

- Manual ACK, durable queues
- Prefetch: 10 messages

## Data Architecture: Database per Service

| Service                                 | Storage    | Tech Stack        |
| --------------------------------------- | ---------- | ----------------- |
| Catalog, User, Order, Inventory, Review | PostgreSQL | ACID transactions |
| Cart, Recommendation Cache              | Redis      | TTL-based expiry  |

## Tanenbaum's Principles Applied

### 1. **Transparency**

- **Location**: API Gateway hides service locations
- **Access**: Uniform REST APIs, consistent JSON responses
- **Failure**: Health checks (`/health`, `/ready`), graceful shutdown

### 2. **Scalability**

- **Horizontal**: HPA autoscaling (1-10 replicas, 70% CPU)
- **Stateless**: All services stateless, shared Redis/PostgreSQL
- **Load Balancing**: K8s round-robin across replicas

### 3. **Fault Tolerance**

- **Replication**: Multi-replica deployments (RollingUpdate, maxUnavailable=0)
- **Health Monitoring**: Liveness/readiness probes
- **Durability**: Persistent RabbitMQ queues, DB transactions

### 4. **Consistency**

- **Strong**: ACID for orders, payments, inventory
- **Eventual**: Async notifications acceptable

### 5. **Security**

- JWT auth (15min access + 7day refresh)
- bcrypt password hashing (cost: 10)
- RBAC, K8s Secrets, TLS-ready

### 6. **Observability**

- Structured logging (Fiber/Python logging)
- Health endpoints
- Prometheus-ready metrics

## Key Workflows

### User Registration

```
Client → Gateway → User Service → PostgreSQL (create user)
                                 → RabbitMQ (user.registered)
                                 → Notification Service (welcome email)
```

### Add to Cart

```
Client → Gateway → Cart Service → Redis (store cart + TTL)
                                 → RabbitMQ (cart.item_added)
```

### Order Creation

```
Client → Gateway → Order Service → PostgreSQL (PENDING order)
                                  → (Future: Inventory reserve, Payment process)
```

### ML Recommendations

```
Client → Gateway → Recommendation → Redis cache HIT? → Return
                                  ↓ MISS
                                  → PostgreSQL (preferences/tags)
                                  → Apply strategies (tag-based, collaborative, popular)
                                  → Cache in Redis
```

### Review with Sentiment

```
Client → Gateway → Review Service → PostgreSQL (store)
                                   → NLTK/TextBlob (analyze sentiment)
                                   → Update with polarity/classification
```

## Deployment (Kubernetes)

**Namespace**: `bookstore-dev`

**Sequential Deployment**:

1. Infrastructure: PostgreSQL, Redis, RabbitMQ
2. Core: Catalog, User services
3. Business: Cart, Order, Inventory
4. ML: Review, Recommendation
5. Gateway: API Gateway
6. Frontend: React app

**Resource Limits**: 128Mi-512Mi RAM, 100m-500m CPU per service

## Design Patterns

✅ **Microservices**: Single responsibility, loose coupling
✅ **API Gateway**: Single entry, cross-cutting concerns
✅ **Event-Driven**: Async decoupling via message broker
✅ **Database per Service**: Bounded contexts, no shared DB
✅ **CQRS-lite**: Read caching in Redis (Recommendations)

## Trade-offs

| Advantage                  | Challenge                      |
| -------------------------- | ------------------------------ |
| Independent scaling        | Network latency                |
| Fault isolation            | Distributed debugging          |
| Tech diversity (Go/Python) | Eventual consistency           |
| Team autonomy              | Integration testing complexity |

## Future Enhancements

🔄 Circuit breakers (resilience4go)
🔄 Distributed tracing (Jaeger)
🔄 Saga pattern for orders
🔄 Service mesh (Istio)
🔄 Dead letter queues

---

**Key Takeaway**: This architecture demonstrates production-ready distributed systems with proper separation of concerns, scalability, and fault tolerance while balancing complexity with practical implementation.
