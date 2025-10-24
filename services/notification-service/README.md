# Notification Service

## Overview

Handles all notifications (email, SMS, push) via RabbitMQ event consumption. Manages templates and delivery tracking.

## Technology Stack

- **Language**: TypeScript
- **Runtime**: Node.js 18+ LTS
- **Framework**: Express.js
- **Email**: SendGrid / AWS SES
- **SMS**: Twilio (optional)
- **Database**: PostgreSQL 15
- **ORM**: Prisma
- **Messaging**: RabbitMQ
- **Templates**: Handlebars
- **Port**: HTTP: 8087

## Responsibilities

- Email notifications (SendGrid/SES)
- SMS notifications (Twilio - optional)
- Push notifications (future)
- Notification templates management
- Delivery tracking and history
- RabbitMQ consumer for all notification events
- Template rendering with Handlebars

## Database Schema

**notifications**: id, user_id, type, channel, recipient, subject, body, status, sent_at, error_message, metadata (JSONB), created_at
**notification_templates**: id, name, type, subject_template, body_template, variables (JSONB)

## Notification Types

- Order confirmation
- Order shipped
- Order delivered
- Payment receipt
- Password reset
- Welcome email
- Low stock alert (admin)

## Events Consumed (RabbitMQ)

- `user.registered` → Send welcome email
- `order.created` → Send order confirmation
- `order.shipped` → Send shipping notification
- `order.delivered` → Send delivery confirmation
- `payment.completed` → Send payment receipt
- `inventory.low_stock` → Alert admin

## Environment Variables

```bash
PORT=8087
NODE_ENV=development
DATABASE_URL=postgresql://bookstore:password@postgres:5432/notifications_db
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/
SENDGRID_API_KEY=SG.xxx
FROM_EMAIL=noreply@bookstore.com
TWILIO_ACCOUNT_SID=xxx
TWILIO_AUTH_TOKEN=xxx
TWILIO_PHONE_NUMBER=+1234567890
```

## Getting Started

```bash
npm install
npm run dev
npm test
```

## Next Steps

- [ ] Implement SendGrid integration
- [ ] Create email templates
- [ ] Set up RabbitMQ consumers
- [ ] Add template rendering
- [ ] Implement delivery tracking
- [ ] Add SMS support (optional)
- [ ] Write tests
