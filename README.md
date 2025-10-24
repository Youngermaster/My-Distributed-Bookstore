# Distributed Bookstore System

A production-ready distributed bookstore built with microservices architecture on AWS EKS, implementing industry best practices and distributed systems principles from "Distributed Systems" by Tanenbaum & van Steen.

## 🎯 Overview

This project demonstrates a real-world implementation of distributed systems concepts with **11 microservices** using a **polyglot architecture**:

- **6 Go Services**: API Gateway, Catalog, User, Cart, Order, Inventory, Admin
- **2 TypeScript Services**: Payment (Stripe), Notification (SendGrid)
- **2 Python Services**: Review (ML/Sentiment Analysis), Recommendation (ML)
- **Event-Driven Architecture** with RabbitMQ
- **Comprehensive Observability** with Jaeger, Prometheus, Grafana
- **React Frontend** with shadcn/ui

Built following principles from "Distributed Systems" by Tanenbaum & van Steen.

## ✨ Key Features

### Core Services
- **API Gateway**: Single entry point, JWT validation, rate limiting, circuit breaker
- **Catalog Service**: Book management, search, categories, authors, publishers
- **User Service**: Authentication, RBAC, profile management, addresses
- **Cart Service**: Session & persistent carts, real-time pricing
- **Order Service**: Saga orchestration, order lifecycle, distributed transactions
- **Payment Service**: Stripe integration, refunds, webhook handling
- **Inventory Service**: Stock tracking, reservations, low-stock alerts
- **Notification Service**: Email/SMS via SendGrid/Twilio, template rendering
- **Review Service**: Reviews, ratings, ML sentiment analysis (NLTK)
- **Recommendation Service**: Collaborative & content-based filtering
- **Admin Service**: Analytics, reporting, system monitoring

### Architecture Patterns
- Database per Service
- Event-Driven (Saga choreography)
- CQRS for complex queries
- API Gateway pattern
- Circuit Breaker
- Service Discovery (K8s DNS)

## 🛠 Tech Stack

### Backend Microservices (Polyglot)
**Go Services** (1.21+):
- Fiber v2 (web framework), GORM (ORM)
- gRPC + Protocol Buffers
- JWT authentication

**TypeScript Services** (Node.js 18+):
- Express.js, Prisma ORM
- Stripe SDK, SendGrid

**Python Services** (3.11+):
- FastAPI, SQLAlchemy
- scikit-learn, NLTK, pandas

### Frontend
- React 18+ with TypeScript
- shadcn/ui components
- TanStack Query v5
- Zustand (state)
- Tailwind CSS

### Infrastructure & Data
- **Databases**: PostgreSQL 15 (per-service)
- **Caching**: Redis 7
- **Messaging**: RabbitMQ 3.12
- **Containers**: Docker & Docker Compose
- **Orchestration**: Kubernetes (AWS EKS)
- **Storage**: AWS S3, CloudFront CDN

### Observability
- **Tracing**: Jaeger
- **Metrics**: Prometheus + Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (optional, for local dev)
- Node.js 18+ (optional, for local dev)
- Python 3.11+ (optional, for local dev)
- kubectl (for K8s deployment)

### Start All Services

```bash
# Clone the repository
git clone <repo-url>
cd My-Distributed-Bookstore

# Start all services with Docker Compose
make services-start

# View logs
make logs

# Check available commands
make help
```

### Service URLs
Once started, services are available at:
- **API Gateway**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Jaeger UI**: http://localhost:16686
- **RabbitMQ Management**: http://localhost:15672
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001

### Test the APIs

**Register a User:**
```bash
curl -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Create a Book:**
```bash
curl -X POST http://localhost:8081/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "isbn": "9780134190440",
    "title": "The Go Programming Language",
    "description": "The authoritative resource to writing clear and idiomatic Go",
    "price": 44.99,
    "stock_quantity": 50,
    "language": "en",
    "pages": 400,
    "format": "paperback"
  }'
```

**List Books:**
```bash
curl "http://localhost:8081/api/v1/books?limit=10&offset=0"
```

## 📁 Project Structure

```
My-Distributed-Bookstore/
├── services/
│   ├── api-gateway/             # Go - API Gateway
│   ├── catalog-service/         # Go - Book catalog
│   ├── user-service/            # Go - Authentication & users
│   ├── cart-service/            # Go - Shopping cart
│   ├── order-service/           # Go - Order processing & saga
│   ├── payment-service/         # TypeScript - Stripe payments
│   ├── inventory-service/       # Go - Stock management
│   ├── notification-service/    # TypeScript - Email/SMS
│   ├── review-service/          # Python - Reviews & ML
│   ├── recommendation-service/  # Python - ML recommendations
│   └── admin-service/           # Go - Admin & analytics
├── frontend/
│   └── customer-app/            # React + TypeScript
├── proto/                       # Protobuf definitions
├── infrastructure/
│   └── k8s/                     # Kubernetes manifests
├── scripts/                     # Build & deployment scripts
├── docs/                        # Documentation
├── docker-compose.yml           # Local development
├── Makefile                     # Development commands
├── CLAUDE.md                    # Full project documentation
└── README.md
```

Each service follows clean architecture with:
- `cmd/` - Entry points
- `internal/` - Business logic (domain, repository, service, handler)
- `pkg/` - Shared utilities
- `proto/` - gRPC definitions
- `migrations/` - Database migrations (if applicable)

## Documentation

- [DEVELOPMENT.md](docs/DEVELOPMENT.md) - Detailed development guide
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Complete project overview
- [CLAUDE.md](CLAUDE.md) - AI assistant guidance

## Available Commands

```bash
make help          # Show all available commands
make up-build      # Build and start all services
make down          # Stop all services
make logs          # View all service logs
make health        # Check service health
make clean         # Clean up everything
```

## Architecture Highlights

### Microservices Pattern
- Each service has its own database (database-per-service)
- Services communicate via REST APIs (gRPC ready)
- Stateless design for horizontal scaling

### Clean Architecture
Each service follows:
1. **Domain Layer** - Business entities
2. **Repository Layer** - Data access abstraction
3. **Service Layer** - Business logic
4. **Handler Layer** - HTTP endpoints
5. **Middleware Layer** - Auth, logging, CORS

### Security
- JWT authentication with expiration
- bcrypt password hashing
- Role-based access control
- Input validation
- SQL injection prevention
- CORS configuration

### Observability
- Centralized logging service
- Structured logging with zerolog
- Health check endpoints
- Distributed tracing support

## 🔌 Service Ports

| Service | HTTP | gRPC | Database |
|---------|------|------|----------|
| API Gateway | 8080 | - | - |
| Catalog | 8081 | 50051 | catalog_db |
| User | 8082 | 50052 | users_db |
| Cart | 8083 | 50053 | Redis |
| Order | 8084 | 50054 | orders_db |
| Payment | 8085 | 50055 | payments_db |
| Inventory | 8086 | 50056 | inventory_db |
| Notification | 8087 | - | notifications_db |
| Review | 8088 | 50058 | reviews_db |
| Recommendation | 8089 | 50059 | recommendations_db |
| Admin | 8090 | 50060 | admin_db |

**Infrastructure:**
- PostgreSQL: 5432
- Redis: 6379
- RabbitMQ: 5672 (AMQP), 15672 (Management)
- Jaeger: 16686 (UI), 6831 (Agent)
- Prometheus: 9090
- Grafana: 3001

## Environment Variables

Each service can be configured via environment variables. See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for details.

## Distributed Systems Principles

This project implements key concepts from Tanenbaum & van Steen:

- **Transparency**: Location-independent service access
- **Scalability**: Stateless services, database per service
- **Fault Tolerance**: Health checks, graceful shutdown
- **Consistency**: Strong consistency for critical operations
- **Security**: Authentication, authorization, encryption
- **Communication**: REST APIs, structured messaging

## 🎯 Development Roadmap

### Phase 1: Foundation ✅
- [x] Project scaffolding
- [x] 11 microservices structure
- [x] Proto definitions
- [x] Docker configurations
- [x] Frontend scaffold

### Phase 2: Core Implementation (In Progress)
- [ ] Implement Go service logic
- [ ] Implement TypeScript services
- [ ] Implement Python ML services
- [ ] Database migrations
- [ ] Event bus (RabbitMQ)

### Phase 3: Frontend
- [ ] Initialize React app
- [ ] Install shadcn/ui
- [ ] Implement pages & components
- [ ] API integration

### Phase 4: Observability
- [ ] Jaeger distributed tracing
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] ELK Stack logging

### Phase 5: Deployment
- [ ] Kubernetes manifests
- [ ] AWS EKS deployment
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Production monitoring

## 🤝 Getting Started with Development

Each microservice has its own README with:
- Service overview and responsibilities
- Technology stack
- Database schema
- API endpoints
- gRPC methods
- Events published/consumed
- Environment variables
- Next steps for implementation

**Start with any service:**
```bash
cd services/catalog-service
cat README.md
```

**Run individual service:**
```bash
cd services/catalog-service
docker-compose up
```

## License

Apache License 2.0

## Contributing

Contributions are welcome! Please read the development guide first.

## Author

Built with best practices in distributed systems architecture.
