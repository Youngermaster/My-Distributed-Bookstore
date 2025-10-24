# Payment Service

## Overview

Handles payment processing via Stripe API, payment method management, refunds, and PCI compliance.

## Technology Stack

- **Language**: TypeScript
- **Runtime**: Node.js 18+ LTS
- **Framework**: Express.js
- **Payment Gateway**: Stripe
- **Database**: PostgreSQL 15
- **ORM**: Prisma
- **Messaging**: RabbitMQ
- **Ports**: HTTP: 8085, gRPC: 50055

## Responsibilities

- Payment processing via Stripe API
- Payment method management
- Refund processing
- Payment status tracking
- PCI compliance handling
- Stripe webhook handling
- Payment intent creation and confirmation

## Database Schema

**payments**: id, order_id, user_id, amount, currency, payment_method, status, stripe_payment_intent_id, stripe_charge_id, failure_reason, metadata (JSONB), created_at, updated_at

**refunds**: id, payment_id, amount, reason, status, stripe_refund_id, created_at

## REST API Endpoints

```
POST   /api/v1/payments                # Create payment intent
POST   /api/v1/payments/:id/confirm    # Confirm payment
GET    /api/v1/payments/:id            # Get payment status
POST   /api/v1/payments/:id/refund     # Issue refund
POST   /webhook/stripe                 # Stripe webhook handler
```

## gRPC Methods

```protobuf
rpc ProcessPayment(ProcessPaymentRequest) returns (PaymentResult);
rpc GetPaymentStatus(GetPaymentStatusRequest) returns (PaymentStatus);
rpc RefundPayment(RefundPaymentRequest) returns (RefundResult);
```

## Events Published

- `payment.processing` - Payment initiated
- `payment.completed` - Payment successful
- `payment.failed` - Payment failed
- `payment.refunded` - Refund issued

## Events Consumed

- `order.created` - Process payment for order
- `order.cancelled` - Refund payment if applicable

## Environment Variables

```bash
# Server
PORT=8085
GRPC_PORT=50055
NODE_ENV=development

# Database
DATABASE_URL=postgresql://bookstore:password@postgres:5432/payments_db

# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# RabbitMQ
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/

# Observability
JAEGER_AGENT_HOST=jaeger
JAEGER_AGENT_PORT=6831
LOG_LEVEL=info
```

## Getting Started

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build
npm run build

# Run production
npm start

# Run tests
npm test
```

## Stripe Integration

### Payment Flow

1. Client requests payment intent creation
2. Service creates Stripe PaymentIntent
3. Client confirms payment with Stripe.js
4. Stripe webhook notifies service of status
5. Service updates database and publishes event

### Webhook Handling

The service handles these Stripe webhook events:
- `payment_intent.succeeded`
- `payment_intent.payment_failed`
- `charge.refunded`

## Next Steps

- [ ] Implement Stripe integration
- [ ] Create payment controllers
- [ ] Set up Prisma ORM
- [ ] Implement webhook handlers
- [ ] Add gRPC server
- [ ] Set up RabbitMQ consumers/publishers
- [ ] Add comprehensive error handling
- [ ] Implement idempotency
- [ ] Write tests
- [ ] Add request validation
