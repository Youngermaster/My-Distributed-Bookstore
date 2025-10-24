# Inventory Service - Quick Implementation Guide

## What's Already Created

✅ Directory structure
✅ requirements.txt with FastAPI 0.115+
✅ app/core/config.py (Pydantic settings)

## Next Steps to Complete

The Inventory Service follows the same pattern as Review Service. Here's what you need to create:

### 1. Database Setup (`app/db/base.py`)
- Copy from review-service with same pattern
- Async SQLAlchemy 2.0 engine
- AsyncSession factory
- get_db() dependency

### 2. Models (`app/models/inventory.py`)
```python
# Tables needed:
# - inventory: book_id, available_quantity, reserved_quantity, reorder_level
# - stock_movements: track all inventory changes
# - reservations: order_id, book_id, quantity, status, expires_at
```

### 3. Schemas (`app/schemas/inventory.py`)
```python
# Request/Response schemas:
# - InventoryResponse
# - StockUpdateRequest
# - ReservationRequest
# - StockMovementResponse
```

### 4. Service (`app/services/inventory_service.py`)
```python
# Business logic:
# - check_stock(book_id)
# - update_stock(book_id, quantity, operation)
# - reserve_stock(order_id, items)
# - release_reservation(order_id)
# - commit_reservation(order_id)
# - get_low_stock_items()
```

### 5. Endpoints (`app/api/v1/endpoints/inventory.py`)
```python
# REST API:
# GET  /api/v1/inventory/{book_id}
# POST /api/v1/inventory/{book_id}/adjust
# POST /api/v1/inventory/reserve
# POST /api/v1/inventory/release/{order_id}
# POST /api/v1/inventory/commit/{order_id}
# GET  /api/v1/inventory/low-stock
```

### 6. Main App (`app/main.py`)
- FastAPI app with lifespan
- Include routers
- CORS middleware
- Health check endpoints

### 7. Docker Files
- Dockerfile (same pattern as review-service)
- docker-compose.yml (PostgreSQL + service)
- .env.example

## Quick Copy Commands

Since the pattern is identical to Review Service, you can:

```bash
# Copy and adapt the base files
cp ../review-service/app/db/base.py app/db/
cp ../review-service/Dockerfile .
cp ../review-service/.env.example .

# Then customize for inventory specifics
```

## Key Differences from Review Service

1. **No ML/Sentiment Analysis** - Inventory is pure data management
2. **Reservations System** - Temporary holds on stock for orders
3. **Stock Movements** - Audit trail of all inventory changes
4. **Low Stock Alerts** - Monitoring for reorder points
5. **Expiring Reservations** - Auto-release after timeout

## Run Instructions

```bash
# Setup venv
python3.11 -m venv venv
source venv/bin/activate

# Install
pip install -r requirements.txt

# Run
uvicorn app.main:app --reload --port 8086
```

