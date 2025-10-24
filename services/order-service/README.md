# Order Service

## Overview

Manages order creation, saga orchestration for distributed transactions, and order lifecycle management.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Messaging**: RabbitMQ
- **Ports**: HTTP: 8084, gRPC: 50054

## Responsibilities

- Order creation and management
- Saga orchestration (choreography pattern)
- Order status tracking
- Order history
- Coordination with Payment and Inventory services
- Compensating transactions

## Database Schema

**orders**: id, user_id, status, total_amount, shipping_address_id, payment_method, tracking_number, created_at, updated_at
**order_items**: id, order_id, book_id, quantity, unit_price, subtotal
**order_status_history**: id, order_id, previous_status, new_status, changed_by, notes, changed_at

## Order Status Flow

```
PENDING → PAYMENT_PROCESSING → PAID → PREPARING → SHIPPED → DELIVERED
    ↓              ↓              ↓         ↓
CANCELLED      CANCELLED      CANCELLED  CANCELLED
```

## REST API Endpoints

```
POST   /api/v1/orders             # Create order from cart
GET    /api/v1/orders             # List user orders
GET    /api/v1/orders/:id         # Get order details
POST   /api/v1/orders/:id/cancel  # Cancel order
GET    /api/v1/orders/:id/tracking # Get tracking info
```

## Saga Orchestration

1. Order Service creates order (status: PENDING)
2. Publishes: `order.created`
3. Payment Service processes payment
4. Inventory Service reserves stock
5. Order Service updates status based on events
6. Notification Service sends confirmation

## Events Published

- `order.created`
- `order.confirmed`
- `order.cancelled`
- `order.shipped`
- `order.delivered`

## Events Consumed

- `payment.completed`
- `payment.failed`
- `inventory.reserved`
- `inventory.reservation_failed`

## Next Steps

- [ ] Implement order models
- [ ] Create saga orchestrator
- [ ] Implement event publishers/consumers
- [ ] Add compensating transactions
- [ ] Write comprehensive tests
