# Inventory Service

## Overview

Manages real-time stock tracking, reservations, and inventory movements.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Messaging**: RabbitMQ
- **Ports**: HTTP: 8086, gRPC: 50056

## Responsibilities

- Real-time stock tracking
- Stock reservation for orders
- Stock release on cancellation
- Low stock alerts
- Stock movement history
- Multi-warehouse support (future)

## Database Schema

**inventory**: id, book_id, available_quantity, reserved_quantity, reorder_level, last_restocked_at, updated_at
**stock_movements**: id, book_id, movement_type, quantity, reference_type, reference_id, notes, created_at
**reservations**: id, book_id, order_id, quantity, status, expires_at, created_at

## REST API Endpoints

```
GET    /api/v1/inventory/:bookId        # Get stock level
POST   /api/v1/inventory/:bookId/adjust # Adjust stock (admin)
GET    /api/v1/inventory/low-stock      # Get low stock items
```

## gRPC Methods

```protobuf
rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
rpc ReserveStock(ReserveStockRequest) returns (ReservationResult);
rpc ReleaseReservation(ReleaseReservationRequest) returns (ReleaseResponse);
rpc CommitReservation(CommitReservationRequest) returns (CommitResponse);
```

## Events Published

- `inventory.updated`
- `inventory.low_stock`
- `inventory.reserved`
- `inventory.reservation_failed`

## Events Consumed

- `order.created` - Reserve stock
- `order.cancelled` - Release reservation
- `payment.completed` - Commit reservation
- `catalog.book_created` - Initialize inventory

## Next Steps

- [ ] Implement inventory models
- [ ] Create reservation logic
- [ ] Add stock movement tracking
- [ ] Implement low stock alerts
- [ ] Add event consumers
- [ ] Write tests
