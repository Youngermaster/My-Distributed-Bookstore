# Inventory Service (FastAPI) ✅

## Overview

The Inventory Service manages real-time stock tracking, reservations, and inventory movements for the distributed bookstore system.

## ✨ Features

- ✅ **Real-time stock tracking** (available & reserved quantities)
- ✅ **Stock reservations** for pending orders
- ✅ **Automatic reservation expiry** (15-minute timeout with background task)
- ✅ **Stock movement history** (complete audit trail)
- ✅ **Low stock alerts**
- ✅ **Stock adjustment operations** (add, subtract, set)
- 🔜 **Multi-warehouse support** (future enhancement)
- ✅ **FastAPI 0.115+ with standard dependencies**
- ✅ **Async database operations** with SQLAlchemy 2.0
- ✅ **Complete REST API** with 9 endpoints
- ✅ **Background tasks** for reservation expiry
- ✅ **Docker ready** with docker-compose
- ✅ **gRPC folder structure** (ready for future implementation)

## 🛠 Tech Stack

- **Python**: 3.11+
- **Framework**: FastAPI 0.115+ with standard dependencies
- **Database**: PostgreSQL 15 with async support
- **ORM**: SQLAlchemy 2.0 (async)
- **Server**: Uvicorn (ASGI server)
- **Validation**: Pydantic v2

## 📁 Project Structure

```
inventory-service/
├── app/
│   ├── api/v1/endpoints/
│   │   └── inventory.py       # ✅ 9 REST endpoints
│   ├── core/
│   │   └── config.py          # ✅ Settings configured
│   ├── db/
│   │   └── base.py            # ✅ Async SQLAlchemy setup
│   ├── models/
│   │   └── inventory.py       # ✅ 3 models (Inventory, Reservation, StockMovement)
│   ├── schemas/
│   │   └── inventory.py       # ✅ Request/Response schemas
│   ├── services/
│   │   └── inventory_service.py # ✅ Business logic
│   ├── tasks/
│   │   └── reservation_expiry.py # ✅ Background task
│   ├── grpc/                   # 📁 Folder ready (future)
│   └── main.py                 # ✅ Complete FastAPI app
├── alembic/                    # 🔜 Migrations (TODO)
├── tests/                      # 🔜 Tests (TODO)
├── requirements.txt            # ✅ Dependencies
├── Dockerfile                  # ✅ Container definition
├── docker-compose.yml          # ✅ Local development
├── .env.example                # ✅ Environment template
├── QUICKSTART.md               # ✅ Implementation guide
└── README.md                   # This file (✅ Updated)
```

## 🚀 Getting Started

### Option 1: Python venv (Recommended for Development)

```bash
# Create virtual environment
python3.11 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Copy environment file
cp .env.example .env

# Start PostgreSQL
docker-compose up postgres -d

# Run the service
uvicorn app.main:app --reload --port 8086
```

### Option 2: Docker Compose (Full Stack)

```bash
# Start everything
docker-compose up

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## 📡 API Endpoints (9 total) ✅

### Inventory Management

```http
POST   /api/v1/inventory                        # Create inventory record
GET    /api/v1/inventory/{book_id}              # Get stock level
POST   /api/v1/inventory/{book_id}/adjust       # Adjust stock (admin)
GET    /api/v1/inventory/low-stock              # Get low stock items
```

### Reservations (for Orders)

```http
POST   /api/v1/inventory/reserve                # Reserve stock for order
POST   /api/v1/inventory/release/{order_id}     # Release reservation
POST   /api/v1/inventory/commit/{order_id}      # Commit reservation (after payment)
```

### Stock Movements

```http
GET    /api/v1/inventory/{book_id}/movements    # Get movement history (paginated)
```

### System

```http
GET    /health                                   # Health check
GET    /                                         # Service info
```

**OpenAPI Docs**: `http://localhost:8086/docs`

## 🗄️ Planned Database Schema

### inventory table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| book_id | UUID | Book identifier (unique) |
| available_quantity | INTEGER | Available stock |
| reserved_quantity | INTEGER | Reserved for orders |
| reorder_level | INTEGER | Low stock threshold |
| last_restocked_at | TIMESTAMP | Last restock time |
| updated_at | TIMESTAMP | Last update time |

### reservations table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| book_id | UUID | Book identifier |
| order_id | UUID | Order identifier |
| quantity | INTEGER | Reserved quantity |
| status | VARCHAR | pending/committed/released |
| expires_at | TIMESTAMP | Auto-release time |
| created_at | TIMESTAMP | Creation time |

### stock_movements table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| book_id | UUID | Book identifier |
| movement_type | VARCHAR | restock/sale/adjustment/reservation |
| quantity | INTEGER | Quantity changed |
| reference_type | VARCHAR | order/purchase/adjustment |
| reference_id | UUID | Related record ID |
| notes | TEXT | Additional info |
| created_at | TIMESTAMP | Movement time |

## ⚙️ Configuration

Environment variables:

```bash
PORT=8086
DATABASE_URL=postgresql+asyncpg://bookstore:dev_password@localhost:5432/inventory_db
LOG_LEVEL=INFO
LOW_STOCK_THRESHOLD=10              # Alert when stock below this
RESERVATION_EXPIRY_MINUTES=15       # Auto-release reservations after
```

## 🔄 Reservation Flow

1. **Reserve** - When order is created:
   - Check if stock available
   - Create reservation with expiry time
   - Decrease available_quantity
   - Increase reserved_quantity

2. **Commit** - When payment succeeds:
   - Mark reservation as committed
   - Decrease reserved_quantity
   - Stock movement recorded

3. **Release** - When order is cancelled OR reservation expires:
   - Delete/mark reservation as released
   - Increase available_quantity
   - Decrease reserved_quantity

## 📊 Implementation Status

| Component | Status |
|-----------|--------|
| Project Structure | ✅ Complete |
| requirements.txt | ✅ Complete |
| Configuration | ✅ Complete |
| Main App | ✅ Complete |
| Database Models | ✅ Complete (3 models) |
| Database Setup | ✅ Complete (async SQLAlchemy 2.0) |
| Business Logic | ✅ Complete (full service layer) |
| API Endpoints | ✅ Complete (9 endpoints) |
| Background Tasks | ✅ Complete (reservation expiry) |
| Tests | 🔜 TODO |
| Migrations | 🔜 TODO (use Alembic in production) |
| Docker | ✅ Complete |

## 📚 Next Steps

✅ **Core Implementation Complete!**

The Inventory Service is now fully functional with all core features implemented.

**Optional Enhancements**:
1. Add comprehensive tests (pytest suite)
2. Setup Alembic database migrations
3. Implement gRPC server for inter-service communication
4. Add RabbitMQ event publishers/consumers
5. Implement Redis caching for hot inventory data
6. Add Prometheus metrics
7. Add Jaeger distributed tracing

**Reference**: This service was built using the same patterns as the Review Service for consistency across all FastAPI microservices.

## 🔗 Integration with Other Services

### Events to Publish (RabbitMQ)
- `inventory.updated` - Stock level changed
- `inventory.low_stock` - Stock below threshold
- `inventory.reserved` - Stock reserved for order
- `inventory.reservation_failed` - Insufficient stock

### Events to Consume
- `order.created` - Reserve stock
- `order.cancelled` - Release reservation
- `payment.completed` - Commit reservation
- `catalog.book_created` - Initialize inventory

### gRPC Methods (Future)
```protobuf
rpc CheckStock(CheckStockRequest) returns (CheckStockResponse);
rpc ReserveStock(ReserveStockRequest) returns (ReservationResult);
rpc ReleaseReservation(ReleaseReservationRequest) returns (ReleaseResponse);
rpc CommitReservation(CommitReservationRequest) returns (CommitResponse);
```

## 📄 License

Apache License 2.0
