# Deployment Guide - Distributed Bookstore

This guide covers how to run all microservices and the frontend application.

## Prerequisites

- **Docker** and **Docker Compose** installed
- **Node.js** 18+ and **npm** (for frontend)
- **Go** 1.21+ (if running services without Docker)
- **Python** 3.11+ (if running Python services without Docker)
- **Git** (for version control)

---

## Quick Start (All Services with Docker Compose)

### 1. Start All Backend Services

Run from the **root directory** of the project:

```bash
# Start all services in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

This will start:
- **PostgreSQL databases** for catalog, user, and order services
- **Redis** for cart service
- **RabbitMQ** for message queuing (notification service)
- **All microservices** on their respective ports

### 2. Start the Frontend

```bash
# Navigate to frontend directory
cd frontend/customer-app

# Install dependencies (first time only)
npm install

# Start development server
npm run dev
```

Frontend will be available at: **http://localhost:5173**

---

## Service Ports & Endpoints

| Service | Port | Type | Database |
|---------|------|------|----------|
| **Catalog Service** | 8080 | Go/Fiber | PostgreSQL |
| **User Service** | 8081 | Go/Fiber | PostgreSQL |
| **Admin Service** | 8082 | Go/Fiber | - |
| **Cart Service** | 8083 | Go/Fiber | Redis |
| **Order Service** | 8084 | Go/Fiber | PostgreSQL |
| **Inventory Service** | 8085 | Python/FastAPI | PostgreSQL |
| **Review Service** | 8086 | Python/FastAPI | PostgreSQL |
| **Recommendation Service** | 8087 | Python/FastAPI | - |
| **Payment Service** | 8088 | Node/Express | - |
| **Notification Service** | 8089 | Node/Express | RabbitMQ |
| **API Gateway** | 8090 | Go/Fiber | - |
| **Frontend** | 5173 | React/Vite | - |

---

## Running Services Individually

### Catalog Service

```bash
cd services/catalog-service

# With Docker Compose (includes PostgreSQL)
docker-compose up -d

# Or manually with Go (requires PostgreSQL running)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=catalog
export DB_PASSWORD=catalog123
export DB_NAME=catalog
export PORT=8080

go run cmd/server/main.go
```

**Health Check:** `curl http://localhost:8080/health`

### User Service

```bash
cd services/user-service

# With Docker Compose
docker-compose up -d

# Or manually
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=userservice
export DB_PASSWORD=userpass123
export DB_NAME=userdb
export PORT=8081

go run cmd/server/main.go
```

**Test Endpoint:** `curl http://localhost:8081/api/users`

### Cart Service

```bash
cd services/cart-service

# With Docker Compose (includes Redis)
docker-compose up -d

# Or manually (requires Redis running)
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=
export REDIS_DB=0
export PORT=8083

go run cmd/server/main.go
```

**Test Endpoint:** `curl http://localhost:8083/api/cart/123e4567-e89b-12d3-a456-426614174000`

### Order Service

```bash
cd services/order-service

# With Docker Compose
docker-compose up -d

# Or manually
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=orderuser
export DB_PASSWORD=orderpass123
export DB_NAME=orderdb
export PORT=8084

go run cmd/server/main.go
```

**Health Check:** `curl http://localhost:8084/health`

### Inventory Service (Python/FastAPI)

```bash
cd services/inventory-service

# With Docker Compose
docker-compose up -d

# Or manually
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt

export DATABASE_URL=postgresql://inventoryuser:inventorypass@localhost:5432/inventorydb
export PORT=8085

uvicorn app.main:app --host 0.0.0.0 --port 8085 --reload
```

**Docs:** `http://localhost:8085/docs`

### Review Service (Python/FastAPI)

```bash
cd services/review-service

# With Docker Compose
docker-compose up -d

# Or manually
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt

export DATABASE_URL=postgresql://reviewuser:reviewpass@localhost:5432/reviewdb
export PORT=8086

uvicorn app.main:app --host 0.0.0.0 --port 8086 --reload
```

**Docs:** `http://localhost:8086/docs`

### Recommendation Service (Python/FastAPI)

```bash
cd services/recommendation-service

# With Docker Compose
docker-compose up -d

# Or manually
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt

export PORT=8087

uvicorn app.main:app --host 0.0.0.0 --port 8087 --reload
```

**Docs:** `http://localhost:8087/docs`

### Payment Service (Node.js/Express)

```bash
cd services/payment-service

# With Docker Compose
docker-compose up -d

# Or manually
npm install

export PORT=8088

npm run dev
```

**Test:** `curl http://localhost:8088/health`

### Notification Service (Node.js/Express)

```bash
cd services/notification-service

# With Docker Compose (includes RabbitMQ)
docker-compose up -d

# Or manually (requires RabbitMQ)
npm install

export PORT=8089
export RABBITMQ_URL=amqp://localhost

npm run dev
```

**Test:** `curl http://localhost:8089/health`

---

## Environment Configuration

### Root `docker-compose.yml` Environment Variables

The main `docker-compose.yml` in the root directory configures all services. Key environment variables:

**PostgreSQL Databases:**
- Catalog DB: `catalog:catalog123@catalog-db:5432/catalog`
- User DB: `userservice:userpass123@user-db:5432/userdb`
- Order DB: `orderuser:orderpass123@order-db:5432/orderdb`
- Inventory DB: `inventoryuser:inventorypass@inventory-db:5432/inventorydb`
- Review DB: `reviewuser:reviewpass@review-db:5432/reviewdb`

**Redis:**
- Address: `cart-redis:6379`
- No password by default

**RabbitMQ:**
- URL: `amqp://rabbitmq:5672`

### Frontend Configuration

Create `frontend/customer-app/.env` (if not exists):

```env
VITE_API_BASE_URL=http://localhost:8080
```

Adjust the base URL if using API Gateway:

```env
VITE_API_BASE_URL=http://localhost:8090
```

---

## Testing the Full Stack

### 1. Verify All Services Are Running

```bash
# Check Docker containers
docker-compose ps

# All services should show "Up" status
```

### 2. Test Cart Service

```bash
# Get cart (creates new if doesn't exist)
curl http://localhost:8083/api/cart/123e4567-e89b-12d3-a456-426614174000

# Add item to cart
curl -X POST http://localhost:8083/api/cart/123e4567-e89b-12d3-a456-426614174000/items \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": "550e8400-e29b-41d4-a716-446655440000",
    "quantity": 2
  }'
```

### 3. Test Order Service

```bash
# Create order
curl -X POST http://localhost:8084/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "items": [
      {
        "book_id": "550e8400-e29b-41d4-a716-446655440000",
        "quantity": 2,
        "price": 29.99
      }
    ],
    "shipping_address": "123 Main St, City, Country",
    "payment_method": "credit_card"
  }'

# List all orders
curl http://localhost:8084/api/orders
```

### 4. Test Frontend Integration

1. Open browser: **http://localhost:5173**
2. Register/Login
3. Browse books at `/books`
4. Click "Add to Cart" on any book
5. Navigate to Cart (🛒 icon in navbar)
6. Adjust quantities, remove items
7. Click "Proceed to Checkout" (creates order)
8. View orders at `/orders` (📦 icon in user menu)
9. Click on an order to see details
10. Test "Cancel Order" functionality

---

## Stopping Services

### Stop All Services

```bash
# Stop and remove containers
docker-compose down

# Stop and remove containers + volumes (clears databases)
docker-compose down -v
```

### Stop Individual Services

```bash
cd services/cart-service
docker-compose down

cd ../order-service
docker-compose down
```

### Stop Frontend

```bash
# Ctrl+C in the terminal running npm run dev
```

---

## Troubleshooting

### Port Already in Use

If you get "port already allocated" errors:

```bash
# Find process using the port (example: 8083)
# Windows PowerShell:
netstat -ano | findstr :8083

# Kill the process
taskkill /PID <PID> /F

# Or change port in docker-compose.yml or service config
```

### Database Connection Issues

```bash
# Check if PostgreSQL container is running
docker ps | grep postgres

# View logs
docker logs catalog-db

# Restart database
docker-compose restart catalog-db
```

### Redis Connection Issues

```bash
# Check Redis container
docker ps | grep redis

# Test Redis connection
docker exec -it cart-redis redis-cli ping
# Should return: PONG

# View Redis data
docker exec -it cart-redis redis-cli
> KEYS *
> GET cart:<uuid>
```

### Frontend API Errors

1. **Check CORS:** Ensure backend services allow `http://localhost:5173`
2. **Check API Base URL:** Verify `.env` file has correct `VITE_API_BASE_URL`
3. **Check service health:** `curl http://localhost:8083/health`

### Clean Restart

```bash
# Stop everything
docker-compose down -v

# Remove old images
docker-compose rm -f

# Rebuild and start
docker-compose up --build -d

# Check logs
docker-compose logs -f
```

---

## Database Migrations

### Order Service (Auto-Migration)

The order service automatically migrates on startup via GORM:

```go
// In cmd/server/main.go
db.AutoMigrate(&domain.Order{}, &domain.OrderItem{})
```

To manually reset:

```bash
# Connect to order database
docker exec -it order-db psql -U orderuser -d orderdb

# Drop tables
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

# Restart service (will recreate tables)
docker-compose restart order-service
```

### Catalog/User Services

Check their respective READMEs for migration scripts.

---

## Production Considerations

### 1. Environment Variables

- **Never commit** `.env` files with production credentials
- Use **secrets management** (Docker secrets, Kubernetes secrets, AWS Secrets Manager)
- Set strong passwords for databases

### 2. Security

- Enable **authentication/authorization** on all endpoints
- Use **HTTPS/TLS** for all communication
- Implement **rate limiting**
- Add **input validation** and **sanitization**

### 3. Monitoring

- Add **health check endpoints** to all services
- Implement **logging** (ELK stack, Grafana Loki)
- Use **metrics** (Prometheus + Grafana)
- Set up **alerts** for service failures

### 4. Scaling

- Use **Kubernetes** for orchestration
- Implement **load balancing** (NGINX, Traefik)
- Add **caching** (Redis for read-heavy operations)
- Use **message queues** for async operations (already have RabbitMQ)

### 5. Database Optimization

- Add **indexes** on frequently queried fields
- Implement **connection pooling**
- Use **read replicas** for scalability
- Set up **automated backups**

---

## Development Workflow

### 1. Start Backend Services

```bash
docker-compose up -d catalog-service user-service cart-service order-service
```

### 2. Start Frontend in Dev Mode

```bash
cd frontend/customer-app
npm run dev
```

### 3. Make Changes

- Edit code in your IDE
- Frontend hot-reloads automatically (Vite)
- Backend requires restart: `docker-compose restart <service-name>`

### 4. View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f cart-service

# Frontend logs are in the terminal running npm run dev
```

### 5. Test Changes

- Use **Postman** or **curl** for API testing
- Check browser console for frontend errors
- Verify data in databases:

```bash
# PostgreSQL
docker exec -it order-db psql -U orderuser -d orderdb
SELECT * FROM orders;

# Redis
docker exec -it cart-redis redis-cli
KEYS *
```

---

## Additional Resources

- **Catalog Service README:** `services/catalog-service/README.md`
- **Cart Service README:** `services/cart-service/README.md`
- **Order Service README:** `services/order-service/README.md`
- **Inventory Quickstart:** `services/inventory-service/QUICKSTART.md`
- **Cart/Order Integration Guide:** `CART_ORDER_INTEGRATION.md`
- **Project Summary:** `PROJECT_SUMMARY.md`
- **Implementation Status:** `IMPLEMENTATION_STATUS.md`

---

## Support

For issues or questions:
1. Check service-specific READMEs
2. Review logs: `docker-compose logs -f <service-name>`
3. Verify environment variables in `docker-compose.yml`
4. Ensure all prerequisites are installed

---

**Happy Coding! 🚀📚**
