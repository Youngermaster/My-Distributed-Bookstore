# Inventory Service (FastAPI)

## Overview

The Inventory Service manages real-time stock tracking, reservations, and inventory movements for the distributed bookstore system.

## ✨ Features

- 🔜 **Real-time stock tracking** (available & reserved quantities)
- 🔜 **Stock reservations** for pending orders
- 🔜 **Automatic reservation expiry** (15-minute timeout)
- 🔜 **Stock movement history** (audit trail)
- 🔜 **Low stock alerts**
- 🔜 **Multi-warehouse support** (future)
- ✅ **FastAPI 0.115+ with standard dependencies**
- ✅ **Async database operations** with SQLAlchemy 2.0
- ✅ **Docker ready**
- ✅ **gRPC folder structure** (for future inter-service communication)

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
│   ├── api/v1/endpoints/       # API endpoints (TODO)
│   ├── core/
│   │   └── config.py          # ✅ Settings configured
│   ├── db/                     # Database setup (TODO)
│   ├── models/                 # SQLAlchemy models (TODO)
│   ├── schemas/                # Pydantic schemas (TODO)
│   ├── services/               # Business logic (TODO)
│   ├── grpc/                   # gRPC (future)
│   └── main.py                 # ✅ FastAPI app
├── alembic/                    # Migrations (TODO)
├── tests/                      # Tests (TODO)
├── requirements.txt            # ✅ Dependencies
├── Dockerfile                  # ✅ Container definition
├── docker-compose.yml          # ✅ Local development
├── .env.example                # ✅ Environment template
├── QUICKSTART.md               # ✅ Implementation guide
└── README.md                   # This file
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

## 📡 Planned API Endpoints

### Inventory Management

```http
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
GET    /api/v1/inventory/{book_id}/movements    # Get movement history
```

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
| Main App | ✅ Basic |
| Database Models | 🔜 TODO |
| Database Setup | 🔜 TODO |
| Business Logic | 🔜 TODO |
| API Endpoints | 🔜 TODO |
| Tests | 🔜 TODO |
| Migrations | 🔜 TODO |
| Docker | ✅ Complete |

## 📚 Next Steps

See `QUICKSTART.md` for detailed implementation guide.

The service follows the same patterns as the Review Service, so you can reference that for:
- Database setup (`app/db/base.py`)
- Model patterns (`app/models/`)
- Service patterns (`app/services/`)
- Endpoint patterns (`app/api/v1/endpoints/`)

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
