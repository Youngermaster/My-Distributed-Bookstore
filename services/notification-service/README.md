# Notification Service

## Overview

Golang microservice that listens to RabbitMQ events and dispatches transactional notifications (email/SMS/push). Currently ships with log-based senders for local development and persists an in-memory buffer of recent messages for quick inspection.

## Technology Stack

- **Language**: Go 1.21
- **Framework**: Fiber v2
- **Messaging**: RabbitMQ (topic exchange)
- **Logging**: Zerolog
- **Port**: HTTP 8089 (`/health`, `/ready`, `/api/v1/notifications/recent`)

## Responsibilities

- Consume RabbitMQ events (user, wishlist, cart, order)
- Build per-channel notification messages
- Dispatch via pluggable senders (log sender provided)
- Retain last N notifications in memory for diagnostics

## RabbitMQ Events

| Routing Key        | Purpose                          |
| ------------------ | -------------------------------- |
| `user.registered`  | Welcome message after sign-up    |
| `user.logged_in`   | Notify user about new logins     |
| `wishlist.added`   | Inform user about wishlist adds  |
| `wishlist.removed` | Notify on wishlist deletions     |
| `cart.item_added`  | Confirm cart updates             |
| `cart.item_updated`| Track cart quantity changes      |
| `cart.item_removed`| Confirm removal from cart        |
| `cart.cleared`     | Alert when cart is emptied       |

Additional routing keys can be bound via configuration.

## Environment Variables

```bash
HTTP_PORT=8089
ENV=development
SERVICE_NAME=notification-service
LOG_LEVEL=info
SHUTDOWN_TIMEOUT=10s
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/
RABBITMQ_EXCHANGE=bookstore.events
RABBITMQ_QUEUE=notification-service
RABBITMQ_ROUTING_KEYS=user.*,wishlist.*,cart.*,order.*
RABBITMQ_PREFETCH_COUNT=10
```

## Getting Started

```bash
go mod tidy
go run ./cmd/server
```

## Next Steps

- [ ] Wire real email/SMS providers (SendGrid, Twilio, etc.)
- [ ] Persist notification history (PostgreSQL)
- [ ] Template rendering for rich content
- [ ] Structured retries & dead-letter support
- [ ] Unit tests for dispatcher/senders
