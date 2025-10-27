# Communication Architecture Clarifications

Great questions! Let me clarify the exact communication flow and architecture decisions.

## 1. API Gateway - Custom vs AWS Service

**Answer: We're building a CUSTOM API Gateway in Golang** ✅

**Why NOT AWS API Gateway or EC2:**
- **AWS API Gateway** is a managed service that would add unnecessary cost and complexity for a learning project
- We want to **learn** how to build an API Gateway ourselves
- Our custom Go API Gateway will run **inside the EKS cluster** as a Kubernetes Deployment
- This gives us full control and understanding of request routing, rate limiting, JWT validation, etc.

**Architecture:**
```
Internet
    ↓
CloudFront (CDN for static assets)
    ↓
Application Load Balancer (AWS ALB)
    ↓
Nginx Ingress Controller (in EKS)
    ↓
API Gateway Service (Custom Go app, Deployment in K8s)
    ↓
Microservices (all in same EKS cluster)
```

**API Gateway Responsibilities:**
- Single entry point for all client requests
- JWT token validation
- Rate limiting (Redis-backed)
- Request routing to appropriate microservices
- CORS handling
- Request/response logging with correlation IDs
- Circuit breaker for downstream services
- Health check aggregation

---

## 2. Complete Communication Matrix

Here's the **complete breakdown** of how every service communicates:

### 📊 Communication Types

| Source | Target | Protocol | Sync/Async | Purpose |
|--------|--------|----------|------------|---------|
| **Frontend** | **API Gateway** | REST/HTTP | Sync | All user requests |
| **API Gateway** | **Catalog Service** | gRPC | Sync | Get books, search |
| **API Gateway** | **User Service** | gRPC | Sync | Auth, profile |
| **API Gateway** | **Cart Service** | gRPC | Sync | Cart operations |
| **API Gateway** | **Order Service** | gRPC | Sync | Create order, get orders |
| **API Gateway** | **Review Service** | gRPC | Sync | Submit/get reviews |
| **API Gateway** | **Recommendation Service** | gRPC | Sync | Get recommendations |
| **Cart Service** | **Catalog Service** | gRPC | Sync | Validate book exists, get price |
| **Order Service** | **Catalog Service** | gRPC | Sync | Get book details |
| **Order Service** | **User Service** | gRPC | Sync | Get user/address info |
| **Order Service** | **RabbitMQ** | AMQP | Async | Publish `order.created` event |
| **Payment Service** | **RabbitMQ** | AMQP | Async | Consume `order.created`, Publish `payment.completed` |
| **Inventory Service** | **RabbitMQ** | AMQP | Async | Consume `order.created`, Publish `inventory.reserved` |
| **Notification Service** | **RabbitMQ** | AMQP | Async | Consume ALL events, send emails |
| **Order Service** | **RabbitMQ** | AMQP | Async | Consume `payment.completed`, `inventory.reserved` |
| **Recommendation Service** | **RabbitMQ** | AMQP | Async | Consume `order.completed`, `review.submitted` |
| **Admin Service** | **All Services** | gRPC | Sync | Get data for dashboards |
| **Review Service** | **RabbitMQ** | AMQP | Async | Publish `review.submitted` |
| **Catalog Service** | **RabbitMQ** | AMQP | Async | Publish `catalog.price_updated` |
| **Inventory Service** | **RabbitMQ** | AMQP | Async | Publish `inventory.low_stock` |

---

## 3. Detailed Communication Flows

### 🔄 Flow 1: User Browses Books (Synchronous - gRPC)

```
┌──────────┐
│ Frontend │
│ (React)  │
└────┬─────┘
     │ HTTP GET /api/v1/books?page=1
     ▼
┌──────────────┐
│ API Gateway  │
│   (Golang)   │
└────┬─────────┘
     │ gRPC: SearchBooks(query, page)
     ▼
┌──────────────┐
│  Catalog     │
│  Service     │
│  (Golang)    │
└──────────────┘
     │
     ▼ Query PostgreSQL
┌──────────────┐
│ catalog_db   │
└──────────────┘
```

**Why gRPC?**
- Fast, efficient binary protocol
- Strong typing with Protocol Buffers
- Lower latency than REST
- Built-in load balancing

---

### 🔄 Flow 2: Add to Cart (Synchronous - gRPC)

```
┌──────────┐
│ Frontend │
└────┬─────┘
     │ HTTP POST /api/v1/cart/items
     │ { bookId: "123", quantity: 2 }
     ▼
┌──────────────┐
│ API Gateway  │ ← Validates JWT token
└────┬─────────┘
     │ gRPC: AddToCart(userId, bookId, quantity)
     ▼
┌──────────────┐
│ Cart Service │
│  (Golang)    │
└────┬─────────┘
     │ gRPC: GetBook(bookId)
     │ (to validate book exists & get current price)
     ▼
┌──────────────┐
│  Catalog     │
│  Service     │
└──────────────┘
     │
     ▼ Store cart in Redis
┌──────────────┐
│    Redis     │
└──────────────┘
```

**Why gRPC between Cart → Catalog?**
- Cart needs to validate book exists
- Cart needs current price
- Synchronous call makes sense (user is waiting)

---

### 🔄 Flow 3: Place Order (Saga Pattern - Mixed gRPC + RabbitMQ)

This is the **most complex flow** - it uses **BOTH** gRPC (synchronous) and RabbitMQ (asynchronous):

```
┌──────────┐
│ Frontend │
└────┬─────┘
     │ HTTP POST /api/v1/orders
     ▼
┌──────────────┐
│ API Gateway  │
└────┬─────────┘
     │ gRPC: CreateOrder(userId, cartId)
     ▼
┌──────────────────────────────────────────────────────────┐
│                    Order Service (Golang)                 │
│                                                            │
│  Step 1: Create order in database (status: PENDING)      │
│                                                            │
│  Step 2: Synchronous gRPC calls to validate:             │
│    ┌─────────────────────────────────────────┐           │
│    │ gRPC → User Service: GetAddress(addressId)│          │
│    │ gRPC → Catalog Service: GetBooks(bookIds) │          │
│    │ gRPC → Cart Service: GetCart(userId)      │          │
│    └─────────────────────────────────────────┘           │
│                                                            │
│  Step 3: If all valid, save order to PostgreSQL          │
│                                                            │
└────┬───────────────────────────────────────────────────────┘
     │
     │ Step 4: Publish async event
     ▼
┌─────────────────────────────────────────────────────────┐
│                    RabbitMQ                              │
│                                                          │
│  Exchange: orders (fanout)                              │
│                                                          │
│  Event: order.created                                   │
│  {                                                       │
│    order_id: "uuid",                                    │
│    user_id: "uuid",                                     │
│    items: [...],                                        │
│    total: 99.99                                         │
│  }                                                       │
└────┬────────────────┬─────────────────┬─────────────────┘
     │                │                 │
     ▼                ▼                 ▼
┌──────────┐  ┌──────────────┐  ┌──────────────────┐
│ Payment  │  │  Inventory   │  │  Notification    │
│ Service  │  │   Service    │  │    Service       │
│(Node.js) │  │  (Golang)    │  │   (Node.js)      │
└────┬─────┘  └────┬─────────┘  └────┬─────────────┘
     │             │                  │
     │ Process     │ Reserve          │ Send email
     │ payment     │ stock            │
     │ via Stripe  │                  │
     │             │                  │
     ▼             ▼                  ▼
┌──────────┐  ┌──────────────┐  ┌──────────────────┐
│ Publish: │  │  Publish:    │  │  (No publish,    │
│ payment. │  │  inventory.  │  │   just consume)  │
│completed │  │  reserved    │  │                  │
└────┬─────┘  └────┬─────────┘  └──────────────────┘
     │             │
     └─────────────┴──────────────────┐
                                      │
                                      ▼
                           ┌─────────────────────┐
                           │     RabbitMQ        │
                           │  (events collected) │
                           └──────────┬──────────┘
                                      │
                                      ▼
                           ┌─────────────────────┐
                           │   Order Service     │
                           │   (Listening for    │
                           │   confirmation)     │
                           │                     │
                           │ If both successful: │
                           │ Update status to    │
                           │ CONFIRMED           │
                           │                     │
                           │ If any failed:      │
                           │ Compensate          │
                           │ (rollback)          │
                           └─────────────────────┘
```

**Why Mixed Communication?**
- **gRPC (Sync)**: For immediate validation (user/address/books exist)
- **RabbitMQ (Async)**: For long-running operations (payment, inventory) that don't block the user

---

### 🔄 Flow 4: Send Notification (Asynchronous - RabbitMQ Only)

```
┌──────────────────────────────────────────┐
│          ANY Service Publishes Event     │
│                                          │
│  Examples:                               │
│  • order.created                         │
│  • payment.completed                     │
│  • order.shipped                         │
│  • user.registered                       │
│  • inventory.low_stock                   │
└────────────────┬─────────────────────────┘
                 │
                 ▼
      ┌─────────────────────┐
      │      RabbitMQ       │
      │                     │
      │  Exchange:          │
      │  notifications      │
      │  (topic)            │
      └──────────┬──────────┘
                 │
                 ▼
      ┌─────────────────────┐
      │  Notification Svc   │
      │    (Node.js)        │
      │                     │
      │  Consumers listen   │
      │  to all exchanges   │
      └──────────┬──────────┘
                 │
                 ▼
         ┌──────────────┐
         │   SendGrid   │
         │   /AWS SES   │
         │              │
         │  Send email  │
         └──────────────┘
```

**Why RabbitMQ Only?**
- Notifications are **fire-and-forget**
- Don't need immediate response
- Decouples services completely
- If email service is down, messages queue up

---

## 4. Which Services Use RabbitMQ?

### ✅ Services WITH RabbitMQ Integration

| Service | Role | Events Published | Events Consumed |
|---------|------|------------------|-----------------|
| **Order Service** | Producer & Consumer | `order.created`, `order.confirmed`, `order.cancelled`, `order.shipped` | `payment.completed`, `payment.failed`, `inventory.reserved`, `inventory.reservation_failed` |
| **Payment Service** | Producer & Consumer | `payment.processing`, `payment.completed`, `payment.failed`, `payment.refunded` | `order.created`, `order.cancelled` |
| **Inventory Service** | Producer & Consumer | `inventory.reserved`, `inventory.reservation_failed`, `inventory.updated`, `inventory.low_stock` | `order.created`, `order.cancelled`, `payment.completed` |
| **Notification Service** | Consumer Only | None | **ALL EVENTS** (listens to everything) |
| **Catalog Service** | Producer | `catalog.book_created`, `catalog.price_updated`, `catalog.stock_updated` | None |
| **Review Service** | Producer | `review.submitted`, `review.updated` | `order.delivered` (to enable reviews) |
| **Recommendation Service** | Consumer | None | `order.completed`, `review.submitted`, `user.registered` |
| **Admin Service** | Consumer | None | `inventory.low_stock` (for alerts) |

### ❌ Services WITHOUT RabbitMQ Integration

| Service | Why No RabbitMQ? | Communication Type |
|---------|------------------|-------------------|
| **API Gateway** | Only routes requests | REST in, gRPC out |
| **User Service** | No async operations needed | gRPC only |
| **Cart Service** | Real-time operations only | gRPC + Redis |

---

## 5. RabbitMQ Exchange & Queue Configuration

Here's the **exact RabbitMQ setup**:

```yaml
Exchanges:
  orders:
    type: fanout
    durable: true
    description: "Broadcast all order events to multiple consumers"
  
  payments:
    type: direct
    durable: true
    description: "Route payment events to specific consumers"
    routing_keys:
      - payment.completed
      - payment.failed
  
  inventory:
    type: topic
    durable: true
    description: "Route inventory events by topic"
    routing_keys:
      - inventory.reserved
      - inventory.updated
      - inventory.low_stock
  
  notifications:
    type: topic
    durable: true
    description: "All notification events"
    routing_keys:
      - notification.email.*
      - notification.sms.*
  
  catalog:
    type: topic
    durable: true
    description: "Catalog changes"
    routing_keys:
      - catalog.book_created
      - catalog.price_updated

Queues:
  # Order Service consumers
  order.payment_events:
    exchange: payments
    routing_key: payment.*
    consumer: Order Service
  
  order.inventory_events:
    exchange: inventory
    routing_key: inventory.*
    consumer: Order Service
  
  # Payment Service consumers
  payment.order_events:
    exchange: orders
    consumer: Payment Service
  
  # Inventory Service consumers
  inventory.order_events:
    exchange: orders
    consumer: Inventory Service
  
  inventory.payment_events:
    exchange: payments
    routing_key: payment.completed
    consumer: Inventory Service
  
  # Notification Service consumers (listens to EVERYTHING)
  notification.orders:
    exchange: orders
    consumer: Notification Service
  
  notification.payments:
    exchange: payments
    routing_key: payment.*
    consumer: Notification Service
  
  notification.inventory:
    exchange: inventory
    routing_key: inventory.low_stock
    consumer: Notification Service
  
  # Recommendation Service consumers
  recommendation.orders:
    exchange: orders
    routing_key: order.completed
    consumer: Recommendation Service
  
  recommendation.reviews:
    exchange: reviews
    consumer: Recommendation Service
  
  # Admin Service consumers
  admin.inventory_alerts:
    exchange: inventory
    routing_key: inventory.low_stock
    consumer: Admin Service
```

---

## 6. Decision Matrix: When to Use What?

| Scenario | Communication Type | Reason |
|----------|-------------------|--------|
| **Frontend → Backend** | REST/HTTP | Standard web protocol, easy to debug |
| **API Gateway → Services** | gRPC | Fast, efficient, type-safe |
| **Service needs immediate response** | gRPC | Synchronous, low latency |
| **Service can wait for response** | RabbitMQ | Asynchronous, decoupled |
| **One event → Multiple listeners** | RabbitMQ (fanout) | Pub/sub pattern |
| **Long-running operation** | RabbitMQ | Non-blocking, resilient |
| **Need to validate data before proceeding** | gRPC | Immediate validation |
| **Fire-and-forget notification** | RabbitMQ | No response needed |

---

## 7. Updated Architecture Diagram with Communication Protocols

```
┌────────────────────────────────────────────────────────────────┐
│                       AWS Cloud (EKS)                          │
│                                                                │
│  ┌──────────────┐                                             │
│  │   Frontend   │                                             │
│  │   (React)    │                                             │
│  └──────┬───────┘                                             │
│         │ REST/HTTP                                           │
│         ▼                                                      │
│  ┌──────────────────┐                                         │
│  │   API Gateway    │                                         │
│  │    (Golang)      │                                         │
│  └──────┬───────────┘                                         │
│         │ gRPC (all calls)                                    │
│         │                                                      │
│    ┌────┴────┬────────┬────────┬────────────┐                │
│    │         │        │        │            │                │
│    ▼         ▼        ▼        ▼            ▼                │
│ ┌─────┐  ┌─────┐  ┌─────┐  ┌────────┐  ┌────────┐          │
│ │User │  │Catal│  │Cart │  │ Order  │  │ Review │          │
│ │ Svc │  │og   │  │ Svc │  │  Svc   │  │  Svc   │          │
│ │gRPC │  │Svc  │  │gRPC │  │gRPC +  │  │ gRPC   │          │
│ │only │  │gRPC │  │only │  │RabbitMQ│  │  +MQ   │          │
│ └─────┘  └──┬──┘  └──┬──┘  └───┬────┘  └────────┘          │
│             │        │         │                             │
│             │gRPC    │gRPC     │RabbitMQ                     │
│             │        │         │Publish                      │
│             └────────┴─────────┤                             │
│                                │                             │
│                                ▼                             │
│                    ┌───────────────────────┐                │
│                    │      RabbitMQ         │                │
│                    │   (Message Broker)    │                │
│                    └───────┬───────────────┘                │
│                            │Consume                          │
│         ┌──────────────────┼────────────────────┐           │
│         │                  │                    │           │
│         ▼                  ▼                    ▼           │
│    ┌─────────┐      ┌──────────┐      ┌──────────────┐    │
│    │ Payment │      │Inventory │      │Notification  │    │
│    │ Service │      │ Service  │      │  Service     │    │
│    │(Node.js)│      │ (Golang) │      │  (Node.js)   │    │
│    │ RabbitMQ│      │ RabbitMQ │      │  RabbitMQ    │    │
│    └─────────┘      └──────────┘      └──────────────┘    │
│         │                  │                                │
│         └──────────────────┘                                │
│                  │RabbitMQ                                  │
│                  │Publish                                   │
│                  ▼                                          │
│            ┌──────────┐                                     │
│            │  Order   │                                     │
│            │ Service  │                                     │
│            │(consume) │                                     │
│            └──────────┘                                     │
│                                                              │
└──────────────────────────────────────────────────────────────┘

Legend:
━━━━  REST/HTTP
─ ─ ─  gRPC
┄┄┄┄┄  RabbitMQ (Async)
```

---

## 8. Code Examples

### API Gateway → Service (gRPC Call)

```go
// API Gateway calling Catalog Service
package main

import (
    "context"
    pb "github.com/yourusername/bookstore/proto/catalog"
    "google.golang.org/grpc"
)

// In API Gateway
func (h *Handler) GetBooks(c *fiber.Ctx) error {
    // gRPC connection to Catalog Service
    conn, err := grpc.Dial(
        "catalog-service.production.svc.cluster.local:50051",
        grpc.WithInsecure(),
    )
    if err != nil {
        return err
    }
    defer conn.Close()
    
    client := pb.NewCatalogServiceClient(conn)
    
    // Make gRPC call
    resp, err := client.SearchBooks(context.Background(), &pb.SearchBooksRequest{
        Query: c.Query("q"),
        Page: int32(c.QueryInt("page", 1)),
        PageSize: 20,
    })
    if err != nil {
        return err
    }
    
    return c.JSON(resp)
}
```

### Order Service Publishing to RabbitMQ

```go
// Order Service publishing event
package events

import (
    "encoding/json"
    "github.com/streadway/amqp"
)

type OrderCreatedEvent struct {
    OrderID   string  `json:"order_id"`
    UserID    string  `json:"user_id"`
    Items     []Item  `json:"items"`
    Total     float64 `json:"total"`
    Timestamp int64   `json:"timestamp"`
}

func (p *Publisher) PublishOrderCreated(event OrderCreatedEvent) error {
    body, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    return p.channel.Publish(
        "orders",  // exchange
        "",        // routing key (fanout ignores this)
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
            DeliveryMode: amqp.Persistent,
        },
    )
}
```

### Payment Service Consuming from RabbitMQ

```typescript
// Payment Service (Node.js) consuming
import amqp from 'amqplib';

async function startConsumer() {
  const connection = await amqp.connect('amqp://rabbitmq:5672');
  const channel = await connection.createChannel();
  
  // Declare exchange
  await channel.assertExchange('orders', 'fanout', { durable: true });
  
  // Declare queue
  const queue = await channel.assertQueue('payment.order_events', {
    durable: true,
  });
  
  // Bind queue to exchange
  await channel.bindQueue(queue.queue, 'orders', '');
  
  console.log('Waiting for order events...');
  
  // Consume messages
  channel.consume(queue.queue, async (msg) => {
    if (msg) {
      const event = JSON.parse(msg.content.toString());
      console.log('Received order:', event.order_id);
      
      try {
        // Process payment
        await processPayment(event);
        
        // Acknowledge message
        channel.ack(msg);
      } catch (error) {
        console.error('Error processing payment:', error);
        // Reject and requeue
        channel.nack(msg, false, true);
      }
    }
  });
}
```

---

## Summary

### ✅ Communication Protocol Decisions:

1. **Frontend → API Gateway**: REST/HTTP (standard web)
2. **API Gateway → All Services**: gRPC (fast, efficient)
3. **Service ↔ Service (sync)**: gRPC (when response needed immediately)
4. **Service → Service (async)**: RabbitMQ (long operations, decoupling)

### ✅ RabbitMQ Services:
- Order Service (producer + consumer)
- Payment Service (producer + consumer)
- Inventory Service (producer + consumer)
- Notification Service (consumer only - listens to everything)
- Catalog Service (producer only)
- Review Service (producer + consumer)
- Recommendation Service (consumer only)
- Admin Service (consumer only for alerts)

### ✅ gRPC-Only Services:
- API Gateway
- User Service
- Cart Service

Does this clarify the communication architecture? Would you like me to expand on any specific flow or pattern?