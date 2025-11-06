# Distributed Bookstore System# Distributed Bookstore System



A production-ready, cloud-native distributed bookstore platform deployed on AWS EKS with 12 microservices, implementing distributed systems principles from "Distributed Systems" by Tanenbaum & van Steen.A production-ready distributed bookstore built with microservices architecture on AWS EKS, implementing industry best practices and distributed systems principles from "Distributed Systems" by Tanenbaum & van Steen.



## 🎯 Overview## 🎯 Overview



Complete e-commerce platform for selling books, built with microservices architecture:This project demonstrates a real-world implementation of distributed systems concepts with **11 microservices** using a **polyglot architecture**:



- **12 Microservices** - 6 Go, 4 Python, 2 TypeScript services- **6 Go Services**: API Gateway, Catalog, User, Cart, Order, Inventory, Admin

- **Polyglot Architecture** - Right tool for each job- **2 TypeScript Services**: Payment (Stripe), Notification (SendGrid)

- **Cloud Native** - Containerized with Docker, orchestrated on AWS EKS- **2 Python Services**: Review (ML/Sentiment Analysis), Recommendation (ML)

- **Modern Frontend** - React 19 with TypeScript, TanStack Router/Query- **Event-Driven Architecture** with RabbitMQ

- **API Gateway Pattern** - Unified entry point with rate limiting- **Comprehensive Observability** with Jaeger, Prometheus, Grafana

- **ML-Powered** - Sentiment analysis and personalized recommendations- **React Frontend** with shadcn/ui

- **Production Deployed** - Running on AWS EKS with 3-node cluster

Built following principles from "Distributed Systems" by Tanenbaum & van Steen.

## 🚀 Production Deployment

## ✨ Key Features

### Live System

- **AWS EKS Cluster:** my-bookstore (us-east-1)### Core Services

- **Frontend URL:** [Live Application](http://ab1a1c3c5b1ca49768c26f26e92ca780-844836377.us-east-1.elb.amazonaws.com)- **API Gateway**: Single entry point, JWT validation, rate limiting, circuit breaker

- **Status:** 12/12 services deployed and healthy- **Catalog Service**: Book management, search, categories, authors, publishers

- **Infrastructure:** 3 × t3.medium nodes, PostgreSQL StatefulSet, Redis cache- **User Service**: Authentication, RBAC, profile management, addresses

- **Cart Service**: Session & persistent carts, real-time pricing

## ✨ Key Features- **Order Service**: Saga orchestration, order lifecycle, distributed transactions

- **Payment Service**: Stripe integration, refunds, webhook handling

### Implemented Services (10/12 Production-Ready)- **Inventory Service**: Stock tracking, reservations, low-stock alerts

- **Notification Service**: Email/SMS via SendGrid/Twilio, template rendering

#### Core E-Commerce- **Review Service**: Reviews, ratings, ML sentiment analysis (NLTK)

- ✅ **Catalog Service** (Go, 8081) - Books, authors, categories, publishers with search- **Recommendation Service**: Collaborative & content-based filtering

- ✅ **User Service** (Go, 8082) - JWT authentication, RBAC, wishlist management- **Admin Service**: Analytics, reporting, system monitoring

- ✅ **Cart Service** (Go, 8083) - Shopping cart with Redis caching

- ✅ **Order Service** (Go, 8084) - Order processing and management### Architecture Patterns

- ✅ **Admin Service** (Go, 8090) - Analytics dashboard, sales metrics- Database per Service

- Event-Driven (Saga choreography)

#### Advanced Features- CQRS for complex queries

- ✅ **Inventory Service** (Python, 8086) - Stock tracking with auto-expiry reservations- API Gateway pattern

- ✅ **Review Service** (Python, 8088) - ML sentiment analysis (NLTK/TextBlob)- Circuit Breaker

- ✅ **Recommendation Service** (Python, 8089) - Personalized recommendations- Service Discovery (K8s DNS)

- ⚠️ **Payment Service** (TypeScript, 8085) - Stripe integration (scaffold)

- ⚠️ **Notification Service** (TypeScript, 8087) - Email/SMS (scaffold)## 🛠 Tech Stack



#### Infrastructure### Backend Microservices (Polyglot)

- ✅ **API Gateway** (Go, 8080) - Routing, rate limiting, health checks**Go Services** (1.21+):

- ✅ **Frontend** (React 19) - SPA with Nginx reverse proxy- Fiber v2 (web framework), GORM (ORM)

- gRPC + Protocol Buffers

### Architecture Patterns- JWT authentication

- ✅ **Database per Service** - 6 PostgreSQL databases for data isolation

- ✅ **API Gateway Pattern** - Centralized routing and authentication**TypeScript Services** (Node.js 18+):

- ✅ **Clean Architecture** - Domain, Repository, Service, Handler layers- Express.js, Prisma ORM

- ✅ **Microservices** - Independent deployment and scaling- Stripe SDK, SendGrid

- ✅ **Cloud-Native** - Kubernetes orchestration on AWS EKS

- ✅ **ML Integration** - Sentiment analysis and recommendations**Python Services** (3.11+):

- FastAPI, SQLAlchemy

## 🛠 Tech Stack- scikit-learn, NLTK, pandas



### Backend Services### Frontend

- React 18+ with TypeScript

#### Go Services (6) - Fiber Framework- shadcn/ui components

| Service | Port | Features |- TanStack Query v5

|---------|------|----------|- Zustand (state)

| API Gateway | 8080 | Request routing, rate limiting, CORS |- Tailwind CSS

| Catalog | 8081 | Book catalog with PostgreSQL |

| User | 8082 | JWT auth, RBAC, wishlist |### Infrastructure & Data

| Cart | 8083 | Redis-backed shopping cart |- **Databases**: PostgreSQL 15 (per-service)

| Order | 8084 | Order processing |- **Caching**: Redis 7

| Admin | 8090 | Analytics and dashboard |- **Messaging**: RabbitMQ 3.12

- **Containers**: Docker & Docker Compose

**Tech:** Go 1.21+, Fiber, GORM, PostgreSQL, Redis, JWT, bcrypt- **Orchestration**: Kubernetes (AWS EKS)

- **Storage**: AWS S3, CloudFront CDN

#### Python Services (4) - FastAPI

| Service | Port | Features |### Observability

|---------|------|----------|- **Tracing**: Jaeger

| Inventory | 8086 | Stock management, background tasks |- **Metrics**: Prometheus + Grafana

| Review | 8088 | ML sentiment analysis (NLTK/TextBlob) |- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)

| Recommendation | 8089 | Collaborative & content-based filtering |

## 🚀 Quick Start

**Tech:** Python 3.11+, FastAPI, SQLAlchemy 2.0 (async), PostgreSQL, NLTK, TextBlob

### Prerequisites

#### TypeScript Services (2) - Express- Docker & Docker Compose

| Service | Port | Features |- Go 1.21+ (optional, for local dev)

|---------|------|----------|- Node.js 18+ (optional, for local dev)

| Payment | 8085 | Stripe integration (scaffold) |- Python 3.11+ (optional, for local dev)

| Notification | 8087 | SendGrid/Twilio (scaffold) |- kubectl (for K8s deployment)



**Tech:** Node.js 18+, Express, Prisma, TypeScript 5+### Start All Services



### Frontend```bash

- **Framework:** React 19 + TypeScript 5# Clone the repository

- **Build Tool:** Vitegit clone <repo-url>

- **Router:** TanStack Router (file-based routing)cd My-Distributed-Bookstore

- **State Management:** TanStack Query (server), Zustand (client)

- **UI:** ShadcnUI + Tailwind CSS# Start all services with Docker Compose

- **HTTP Client:** Axiosmake services-start



### Infrastructure# View logs

- **Container Registry:** AWS ECRmake logs

- **Orchestration:** AWS EKS (Kubernetes 1.28+)

- **Database:** PostgreSQL 15 StatefulSet (6 databases)# Check available commands

- **Cache:** Redis 7 Deploymentmake help

- **Load Balancer:** AWS ELB (for frontend)```

- **Storage:** AWS EBS (gp2, 10Gi per volume)

### Service URLs

## 📋 Quick StartOnce started, services are available at:

- **API Gateway**: http://localhost:8080

### Prerequisites- **Frontend**: http://localhost:3000

- Docker & Docker Compose- **Jaeger UI**: http://localhost:16686

- Node.js 18+ (for frontend)- **RabbitMQ Management**: http://localhost:15672

- Go 1.21+ (for backend development)- **Prometheus**: http://localhost:9090

- Python 3.11+ (for Python services)- **Grafana**: http://localhost:3001

- kubectl (for Kubernetes operations)

- AWS CLI (for EKS deployment)### Test the APIs



### Local Development**Register a User:**

```bash

#### Option 1: Docker Compose (Recommended)curl -X POST http://localhost:8082/api/v1/auth/register \

  -H "Content-Type: application/json" \

```bash  -d '{

# Start all backend services    "email": "john@example.com",

cd services/api-gateway    "password": "SecurePass123!",

docker compose up -d    "full_name": "John Doe"

  }'

# Wait 30 seconds for database seeding```

sleep 30

**Login:**

# Start frontend```bash

cd ../../frontend/customer-appcurl -X POST http://localhost:8082/api/v1/auth/login \

npm install  -H "Content-Type: application/json" \

npm run dev  -d '{

# Open http://localhost:5173    "email": "john@example.com",

```    "password": "SecurePass123!"

  }'

#### Option 2: Manual Start```



```bash**Create a Book:**

# Terminal 1 - PostgreSQL```bash

docker run --name bookstore-postgres \curl -X POST http://localhost:8081/api/v1/books \

  -e POSTGRES_USER=bookstore \  -H "Content-Type: application/json" \

  -e POSTGRES_PASSWORD=dev_password \  -d '{

  -e POSTGRES_DB=catalog_db \    "isbn": "9780134190440",

  -p 5432:5432 -d postgres:15-alpine    "title": "The Go Programming Language",

    "description": "The authoritative resource to writing clear and idiomatic Go",

# Terminal 2 - Catalog Service    "price": 44.99,

cd services/catalog-service    "stock_quantity": 50,

export DATABASE_URL=postgresql://bookstore:dev_password@localhost:5432/catalog_db?sslmode=disable    "language": "en",

go run cmd/server/main.go    "pages": 400,

    "format": "paperback"

# Terminal 3 - API Gateway  }'

cd services/api-gateway```

export CATALOG_SERVICE_URL=http://localhost:8081

go run cmd/server/main.go**List Books:**

```bash

# Terminal 4 - Frontendcurl "http://localhost:8081/api/v1/books?limit=10&offset=0"

cd frontend/customer-app```

npm install && npm run dev

```## 📁 Project Structure



### Test API Endpoints```

My-Distributed-Bookstore/

```bash├── services/

# Health check│   ├── api-gateway/             # Go - API Gateway

curl http://localhost:8080/health | jq│   ├── catalog-service/         # Go - Book catalog

│   ├── user-service/            # Go - Authentication & users

# List books│   ├── cart-service/            # Go - Shopping cart

curl http://localhost:8080/api/v1/catalog/books | jq│   ├── order-service/           # Go - Order processing & saga

│   ├── payment-service/         # TypeScript - Stripe payments

# Search books│   ├── inventory-service/       # Go - Stock management

curl "http://localhost:8080/api/v1/catalog/books/search?q=distributed" | jq│   ├── notification-service/    # TypeScript - Email/SMS

│   ├── review-service/          # Python - Reviews & ML

# Register user│   ├── recommendation-service/  # Python - ML recommendations

curl -X POST http://localhost:8080/api/v1/users/auth/register \│   └── admin-service/           # Go - Admin & analytics

  -H "Content-Type: application/json" \├── frontend/

  -d '{"email":"test@example.com","password":"Test123!","full_name":"Test User"}' | jq│   └── customer-app/            # React + TypeScript

├── proto/                       # Protobuf definitions

# Login├── infrastructure/

curl -X POST http://localhost:8080/api/v1/users/auth/login \│   └── k8s/                     # Kubernetes manifests

  -H "Content-Type: application/json" \├── scripts/                     # Build & deployment scripts

  -d '{"email":"test@example.com","password":"Test123!"}' | jq├── docs/                        # Documentation

```├── docker-compose.yml           # Local development

├── Makefile                     # Development commands

## 🎯 Key API Endpoints├── CLAUDE.md                    # Full project documentation

└── README.md

### Catalog Service```

```

GET    /api/v1/catalog/books              # List books (with filters)Each service follows clean architecture with:

GET    /api/v1/catalog/books/search       # Search books- `cmd/` - Entry points

GET    /api/v1/catalog/books/:id          # Get book details- `internal/` - Business logic (domain, repository, service, handler)

POST   /api/v1/catalog/books              # Create book (admin)- `pkg/` - Shared utilities

GET    /api/v1/catalog/categories         # List categories- `proto/` - gRPC definitions

GET    /api/v1/catalog/authors            # List authors- `migrations/` - Database migrations (if applicable)

```

## Documentation

### User Service

```- [DEVELOPMENT.md](docs/DEVELOPMENT.md) - Detailed development guide

POST   /api/v1/users/auth/register        # Register new user- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Complete project overview

POST   /api/v1/users/auth/login           # Login (returns JWT)- [CLAUDE.md](CLAUDE.md) - AI assistant guidance

GET    /api/v1/users/me                   # Get current user (protected)

GET    /api/v1/users/me/wishlist          # Get wishlist## Available Commands

POST   /api/v1/users/me/wishlist          # Add to wishlist

``````bash

make help          # Show all available commands

### Cart Servicemake up-build      # Build and start all services

```make down          # Stop all services

GET    /api/v1/cart/:cartId               # Get cartmake logs          # View all service logs

POST   /api/v1/cart/:cartId/items         # Add item to cartmake health        # Check service health

PATCH  /api/v1/cart/:cartId/items/:bookId # Update quantitymake clean         # Clean up everything

DELETE /api/v1/cart/:cartId/items/:bookId # Remove item```

```

## Architecture Highlights

### Order Service

```### Microservices Pattern

POST   /api/v1/orders                     # Create order- Each service has its own database (database-per-service)

GET    /api/v1/orders/:orderId            # Get order details- Services communicate via REST APIs (gRPC ready)

POST   /api/v1/orders/:orderId/cancel     # Cancel order- Stateless design for horizontal scaling

```

### Clean Architecture

### Review ServiceEach service follows:

```1. **Domain Layer** - Business entities

POST   /api/v1/reviews                    # Create review (auto-sentiment)2. **Repository Layer** - Data access abstraction

GET    /api/v1/reviews/book/:bookId       # Get book reviews3. **Service Layer** - Business logic

GET    /api/v1/reviews/book/:bookId/stats # Get review statistics4. **Handler Layer** - HTTP endpoints

POST   /api/v1/reviews/:reviewId/vote     # Vote on review5. **Middleware Layer** - Auth, logging, CORS

```

### Security

### Recommendation Service- JWT authentication with expiration

```- bcrypt password hashing

GET    /api/v1/recommendations/me                # Personalized recommendations- Role-based access control

GET    /api/v1/recommendations/similar/:bookId   # Similar books- Input validation

GET    /api/v1/recommendations/trending          # Trending books- SQL injection prevention

POST   /api/v1/recommendations/interactions      # Track user interaction- CORS configuration

```

### Observability

### Inventory Service- Centralized logging service

```- Structured logging with zerolog

GET    /api/v1/inventory/:bookId          # Get stock level- Health check endpoints

POST   /api/v1/inventory/reserve          # Reserve stock for order- Distributed tracing support

POST   /api/v1/inventory/commit/:orderId  # Commit reservation

GET    /api/v1/inventory/low-stock        # Get low stock items## 🔌 Service Ports

```

| Service | HTTP | gRPC | Database |

## 📊 Database Design|---------|------|------|----------|

| API Gateway | 8080 | - | - |

### Database-per-Service Pattern| Catalog | 8081 | 50051 | catalog_db |

| User | 8082 | 50052 | users_db |

| Database | Service | Tables | Records || Cart | 8083 | 50053 | Redis |

|----------|---------|--------|---------|| Order | 8084 | 50054 | orders_db |

| `catalog_db` | Catalog | books, authors, categories, publishers, book_authors, book_categories | 3 books, 5 authors || Payment | 8085 | 50055 | payments_db |

| `userdb` | User | users, roles, addresses, sessions, wishlist | Dynamic || Inventory | 8086 | 50056 | inventory_db |

| `orderdb` | Order | orders, order_items, order_status_history | Dynamic || Notification | 8087 | - | notifications_db |

| `inventory_db` | Inventory | inventory, reservations, stock_movements | Dynamic || Review | 8088 | 50058 | reviews_db |

| `reviews_db` | Review | reviews, review_votes | Dynamic || Recommendation | 8089 | 50059 | recommendations_db |

| `bookstore` | Admin | (queries other databases for analytics) | - || Admin | 8090 | 50060 | admin_db |



### Sample Data (Auto-Seeded)**Infrastructure:**

- PostgreSQL: 5432

**Books:**- Redis: 6379

1. Building Microservices by Martin Fowler - $49.99- RabbitMQ: 5672 (AMQP), 15672 (Management)

2. Clean Code by Robert C. Martin - $44.99- Jaeger: 16686 (UI), 6831 (Agent)

3. Distributed Systems by Tanenbaum & van Steen - $89.99- Prometheus: 9090

- Grafana: 3001

**Categories:**

- Programming 💻## Environment Variables

- Distributed Systems 🌐

- Software Architecture 🏗️Each service can be configured via environment variables. See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for details.

- Databases 🗄️

- Cloud Computing ☁️## Distributed Systems Principles



## 🏗️ ArchitectureThis project implements key concepts from Tanenbaum & van Steen:



```- **Transparency**: Location-independent service access

┌─────────────┐- **Scalability**: Stateless services, database per service

│   Browser   │- **Fault Tolerance**: Health checks, graceful shutdown

└──────┬──────┘- **Consistency**: Strong consistency for critical operations

       │- **Security**: Authentication, authorization, encryption

┌──────▼──────────┐- **Communication**: REST APIs, structured messaging

│  Frontend       │

│  React 19       │## 🎯 Development Roadmap

│  AWS LoadBalancer

└──────┬──────────┘### Phase 1: Foundation ✅

       │- [x] Project scaffolding

┌──────▼──────────┐- [x] 11 microservices structure

│  API Gateway    │- [x] Proto definitions

│  Port 8080      │- [x] Docker configurations

│  (Go/Fiber)     │- [x] Frontend scaffold

└───┬─────────┬───┘

    │         │### Phase 2: Core Implementation (In Progress)

    ├─────────┼─────────────┐- [ ] Implement Go service logic

    │         │             │- [ ] Implement TypeScript services

┌───▼───┐ ┌──▼───┐ ┌───▼────┐ ┌────▼─────┐- [ ] Implement Python ML services

│Catalog│ │ User │ │  Cart  │ │  Order   │- [ ] Database migrations

│ :8081 │ │:8082 │ │ :8083  │ │  :8084   │- [ ] Event bus (RabbitMQ)

│  Go   │ │  Go  │ │   Go   │ │   Go     │

└───┬───┘ └──┬───┘ └───┬────┘ └────┬─────┘### Phase 3: Frontend

    │        │         │           │- [ ] Initialize React app

┌───▼───┐ ┌──▼───┐ ┌──▼──┐   ┌────▼─────┐- [ ] Install shadcn/ui

│catalog│ │userdb│ │Redis│   │ orderdb  │- [ ] Implement pages & components

│  _db  │ │      │ │:6379│   │          │- [ ] API integration

└───────┘ └──────┘ └─────┘   └──────────┘

### Phase 4: Observability

┌───────────────────────────────────────────┐- [ ] Jaeger distributed tracing

│      Python Services (FastAPI)            │- [ ] Prometheus metrics

├──────────┬──────────┬──────────────────┬──┤- [ ] Grafana dashboards

│Inventory │  Review  │  Recommendation  │  │- [ ] ELK Stack logging

│  :8086   │  :8088   │      :8089       │  │

└──────┬───┴────┬─────┴──────────────────┘  │### Phase 5: Deployment

       │        │                            │- [ ] Kubernetes manifests

┌──────▼──┐ ┌──▼──────┐                     │- [ ] AWS EKS deployment

│inventory│ │reviews  │                     │- [ ] CI/CD pipeline (GitHub Actions)

│   _db   │ │   _db   │                     │- [ ] Production monitoring

└─────────┘ └─────────┘                     │

```## 🤝 Getting Started with Development



### Request FlowEach microservice has its own README with:

1. **Browser** → Frontend (React SPA)- Service overview and responsibilities

2. **Frontend** → API Gateway (via LoadBalancer)- Technology stack

3. **API Gateway** → Backend Service (internal routing)- Database schema

4. **Backend Service** → Database (PostgreSQL/Redis)- API endpoints

5. **Response** ← Reverse path- gRPC methods

- Events published/consumed

## 📚 Documentation- Environment variables

- Next steps for implementation

- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Comprehensive architecture guide with:

  - Service implementation details**Start with any service:**

  - Go & Python architecture patterns```bash

  - API Gateway integrationcd services/catalog-service

  - Frontend architecturecat README.md

  - Database design```

  - Testing guide

**Run individual service:**

- **[EKS_DEPLOYMENT_GUIDE.md](EKS_DEPLOYMENT_GUIDE.md)** - Complete AWS EKS deployment guide:```bash

  - Prerequisites and setupcd services/catalog-service

  - Cluster creationdocker-compose up

  - Building and pushing Docker images```

  - Deploying all 12 services

  - Database configuration## License

  - Troubleshooting

Apache License 2.0

- **Service-Specific READMEs:**

  - `services/catalog-service/README.md`## Contributing

  - `services/user-service/README.md`

  - `services/inventory-service/README.md`Contributions are welcome! Please read the development guide first.

  - `services/review-service/README.md`

  - And more...## Author



## 🧪 TestingBuilt with best practices in distributed systems architecture.


### Frontend Testing
- Browse catalog: http://localhost:5173
- Search books (debounced)
- Genre navigation
- Book details with recommendations
- User registration and login
- Shopping cart
- Wishlist management

### Backend Testing
```bash
# Test catalog search
curl "http://localhost:8080/api/v1/catalog/books/search?q=distributed" | jq

# Test user registration
curl -X POST http://localhost:8080/api/v1/users/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com","password":"Test123!","full_name":"New User"}' | jq

# Test review creation (with auto-sentiment)
curl -X POST http://localhost:8080/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{
    "book_id":"BOOK_UUID",
    "user_id":"USER_UUID",
    "rating":5,
    "title":"Excellent!",
    "content":"This book was amazing! Highly recommended."
  }' | jq

# Test recommendations
curl http://localhost:8080/api/v1/recommendations/me?limit=10 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" | jq
```

## 🚀 Deployment

### AWS EKS Production Deployment

**Full deployment guide:** See [EKS_DEPLOYMENT_GUIDE.md](EKS_DEPLOYMENT_GUIDE.md)

**Quick deployment:**
```bash
# 1. Configure AWS credentials
aws configure

# 2. Build and push all images to ECR
cd services/catalog-service
docker build -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest .
docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest
# Repeat for all services...

# 3. Deploy to Kubernetes
kubectl apply -f infrastructure/k8s/namespaces/
kubectl apply -f infrastructure/k8s/secrets/
kubectl apply -f infrastructure/k8s/databases/
kubectl apply -f infrastructure/k8s/services/
kubectl apply -f infrastructure/k8s/frontend/

# 4. Check deployment status
kubectl get pods -n bookstore-dev

# 5. Get frontend URL
kubectl get svc frontend -n bookstore-dev
```

**Deployment Script:** See `deploy-to-eks.ps1` for automated deployment

## 🎓 Distributed Systems Principles

This project implements key distributed systems concepts from Tanenbaum & van Steen:

### Communication
- ✅ RESTful HTTP APIs for client-server communication
- ✅ gRPC ports configured for inter-service communication
- ✅ API Gateway for request routing

### Processes & Architecture
- ✅ Microservices architecture (independent processes)
- ✅ Clean Architecture within each service
- ✅ Stateless services (horizontal scaling ready)
- ✅ Database-per-service pattern

### Naming
- ✅ Service discovery via Kubernetes DNS
- ✅ Environment-based configuration
- ✅ Consistent API versioning (/api/v1)

### Synchronization & Consistency
- ✅ Strong consistency for critical operations
- ✅ Transaction support via GORM
- ✅ Connection pooling

### Fault Tolerance
- ✅ Health check endpoints (/health, /ready)
- ✅ Graceful shutdown handling
- ✅ Database connection retry logic
- ✅ Multi-replica deployments

### Security
- ✅ JWT-based authentication
- ✅ Password hashing (bcrypt)
- ✅ Role-based access control (RBAC)
- ✅ SQL injection prevention
- ✅ CORS configuration

## 🔧 Development Workflow

### Adding a New Feature

**Backend (Go Service):**
```bash
# 1. Add domain model
vim services/catalog-service/internal/domain/models.go

# 2. Add repository method
vim services/catalog-service/internal/repository/book_repository.go

# 3. Add service method
vim services/catalog-service/internal/service/catalog_service.go

# 4. Add HTTP handler
vim services/catalog-service/internal/handler/http/book_handler.go

# 5. Test locally
go run cmd/server/main.go
```

**Frontend:**
```bash
# 1. Add API client method
vim src/lib/api.ts

# 2. Add component
vim src/components/BookRecommendations.tsx

# 3. Test locally
npm run dev
```

**Deploy to Production:**
```bash
# Build and push
docker build -t 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest .
docker push 905418472239.dkr.ecr.us-east-1.amazonaws.com/catalog-service:latest

# Restart deployment
kubectl rollout restart deployment/catalog-service -n bookstore-dev
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Based on distributed systems principles from "Distributed Systems" by Tanenbaum & van Steen
- Microservices patterns from "Building Microservices" by Sam Newman
- Clean Architecture concepts from Robert C. Martin

## 📞 Support

For questions or issues:
- Open an issue on GitHub
- Check the [ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed technical information
- Refer to [EKS_DEPLOYMENT_GUIDE.md](EKS_DEPLOYMENT_GUIDE.md) for deployment help

---

**Built with ❤️ for learning distributed systems in practice**
