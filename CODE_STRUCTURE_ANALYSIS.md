# Code Structure Analysis - Distributed Bookstore

**Date:** October 27, 2025  
**Status:** Current Implementation Review

## 📊 Implementation Status Overview

### ✅ **Fully Implemented Services (3)**

1. **User Service** - Go (Port 8082)
2. **Catalog Service** - Go (Port 8081) ✨ **NEW**
3. **Review Service** - Python/FastAPI (Port 8088)
4. **Inventory Service** - Python/FastAPI (Port 8086)

### 🚧 **Services Needing Implementation (2)**

1. **Cart Service** - Go (Port 8083) - Only boilerplate
2. **Order Service** - Go (Port 8084) - Only boilerplate

---

## 🏗️ Standard Go Service Architecture Pattern

Based on **user-service** and **catalog-service**, all Go services follow this structure:

```
service-name/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point, wires everything together
├── internal/
│   ├── config/                        # Configuration management
│   │   ├── config.go                  # Load env vars, settings
│   │   └── database.go                # DB initialization, migrations, seeding
│   ├── domain/                        # Domain entities (GORM models)
│   │   └── models.go                  # All database models
│   ├── repository/                    # Data access layer
│   │   └── *_repository.go            # One file per entity
│   ├── service/                       # Business logic layer
│   │   ├── *_service.go               # Service implementation
│   │   └── dto.go                     # Request/Response DTOs
│   ├── handler/                       # HTTP handlers
│   │   └── http/
│   │       ├── routes.go              # Route definitions
│   │       ├── *_handler.go           # Handler per entity
│   │       ├── responses.go           # Common response structures
│   │       └── helpers.go             # Helper functions
│   └── middleware/                    # Custom middleware (auth, logging, etc.)
│       └── auth.go
├── pkg/                               # Shared utilities (reusable across services)
│   ├── jwt/                           # JWT utilities
│   └── password/                      # Password hashing
├── migrations/                        # SQL migration files (optional)
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 📋 Detailed Service Comparison

### 1. **Catalog Service** (Recently Implemented by Partner)

#### Structure:
```
catalog-service/
├── cmd/server/main.go                 ✅ Complete
├── internal/
│   ├── config/
│   │   ├── config.go                  ✅ Complete
│   │   └── database.go                ✅ Complete (with auto-migrate & seeding)
│   ├── domain/
│   │   └── models.go                  ✅ Complete (Book, Author, Category, Publisher)
│   ├── repository/
│   │   ├── book_repository.go         ✅ Complete (CRUD + Search + Filters)
│   │   ├── author_repository.go       ✅ Complete
│   │   ├── category_repository.go     ✅ Complete
│   │   └── publisher_repository.go    ✅ Complete
│   ├── service/
│   │   ├── catalog_service.go         ✅ Complete (All business logic)
│   │   └── dto.go                     ✅ Complete (Request/Response objects)
│   └── handler/http/
│       ├── routes.go                  ✅ Complete (All REST routes)
│       ├── book_handler.go            ✅ Complete
│       ├── responses.go               ✅ Complete
│       └── helpers.go                 ✅ Complete
└── migrations/                        ✅ Present
```

#### Key Features:
- ✅ Full CRUD for Books, Authors, Categories, Publishers
- ✅ Advanced filtering (by category, author, publisher, price range)
- ✅ Search functionality
- ✅ Stock management
- ✅ Pagination support
- ✅ Many-to-many relationships (books ↔ authors, books ↔ categories)
- ✅ Database seeding for development
- ✅ Clean separation of concerns

#### API Endpoints:
```
Books:
  GET    /api/v1/books              # List with filters
  GET    /api/v1/books/search       # Search books
  GET    /api/v1/books/:id          # Get by ID
  POST   /api/v1/books              # Create
  PUT    /api/v1/books/:id          # Update
  DELETE /api/v1/books/:id          # Delete
  PATCH  /api/v1/books/:id/stock    # Update stock

Authors:
  GET    /api/v1/authors
  GET    /api/v1/authors/:id
  POST   /api/v1/authors
  PUT    /api/v1/authors/:id
  DELETE /api/v1/authors/:id

Categories:
  GET    /api/v1/categories
  GET    /api/v1/categories/:id
  POST   /api/v1/categories
  PUT    /api/v1/categories/:id
  DELETE /api/v1/categories/:id

Publishers:
  GET    /api/v1/publishers
  GET    /api/v1/publishers/:id
  POST   /api/v1/publishers
  PUT    /api/v1/publishers/:id
  DELETE /api/v1/publishers/:id
```

---

### 2. **User Service** (Reference Implementation)

#### Structure:
```
user-service/
├── cmd/server/main.go                 ✅ Complete (387 lines)
├── internal/
│   ├── config/                        ✅ Complete
│   ├── database/                      ✅ Complete (separate from config)
│   ├── domain/                        ✅ Complete (User, Role, Address, Session)
│   ├── dto/                           ✅ Complete (separate from service)
│   ├── repository/
│   │   └── user_repository.go         ✅ Complete (single file, all repos)
│   ├── service/                       ✅ Complete
│   └── middleware/
│       └── auth.go                    ✅ Complete (JWT middleware)
└── pkg/
    ├── jwt/                           ✅ Complete
    └── password/                      ✅ Complete
```

#### Key Differences from Catalog:
- **Separate `database/` folder** instead of putting it in `config/`
- **Separate `dto/` folder** instead of `dto.go` in `service/`
- **All repositories in one file** instead of separate files
- **Has `pkg/` utilities** (JWT, password hashing)
- **Has authentication middleware**

---

### 3. **Cart Service** (Needs Implementation)

#### Current State:
```
cart-service/
├── cmd/server/main.go                 ⚠️ Boilerplate only (128 lines of TODOs)
├── go.mod                             ✅ Complete (all dependencies listed)
└── docker-compose.yml                 ✅ Complete
```

#### Dependencies Already in go.mod:
- ✅ Fiber v2
- ✅ GORM + PostgreSQL driver
- ✅ JWT
- ✅ Redis client (go-redis/v9)
- ✅ gRPC + Protobuf
- ✅ RabbitMQ (streadway/amqp)
- ✅ Jaeger client
- ✅ Prometheus client

#### What Needs to be Built:
```
cart-service/
├── internal/
│   ├── config/
│   │   ├── config.go                  # Load env vars
│   │   └── redis.go                   # Redis initialization
│   ├── domain/
│   │   └── models.go                  # Cart, CartItem
│   ├── repository/
│   │   ├── cart_repository.go         # Redis operations
│   │   └── cart_backup_repository.go  # PostgreSQL backup (optional)
│   ├── service/
│   │   ├── cart_service.go            # Business logic
│   │   └── dto.go                     # Request/Response
│   ├── handler/http/
│   │   ├── routes.go
│   │   ├── cart_handler.go
│   │   └── responses.go
│   └── grpc_client/
│       └── catalog_client.go          # gRPC client to Catalog Service
```

---

### 4. **Order Service** (Needs Implementation)

#### Current State:
```
order-service/
├── cmd/server/main.go                 ⚠️ Boilerplate only (128 lines of TODOs)
├── go.mod                             ✅ Complete (all dependencies)
└── docker-compose.yml                 ✅ Complete
```

#### What Needs to be Built:
```
order-service/
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── database.go                # PostgreSQL + migrations
│   ├── domain/
│   │   └── models.go                  # Order, OrderItem, OrderStatusHistory
│   ├── repository/
│   │   └── order_repository.go        # CRUD + status management
│   ├── service/
│   │   ├── order_service.go           # Order creation logic
│   │   ├── saga_orchestrator.go       # Saga pattern implementation ⭐
│   │   └── dto.go
│   ├── handler/http/
│   │   ├── routes.go
│   │   ├── order_handler.go
│   │   └── responses.go
│   ├── events/                        # Event handling
│   │   ├── publisher.go               # RabbitMQ publisher
│   │   └── consumer.go                # RabbitMQ consumer
│   └── grpc_client/
│       ├── payment_client.go          # gRPC to Payment Service
│       └── inventory_client.go        # gRPC to Inventory Service
```

---

## 🎯 Recommended Architecture Pattern for Cart & Order Services

### Pattern Choice: **Follow Catalog Service Structure**

**Why?**
1. ✅ More recent implementation (your partner's work)
2. ✅ Better separation: one repository file per entity
3. ✅ Clear handler structure with separate files
4. ✅ Well-organized with `responses.go` and `helpers.go`
5. ✅ Consistent with modern Go project layout

### Pattern to Follow:

```go
// 1. Domain Layer (internal/domain/models.go)
type Cart struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Items     []CartItem
    Total     float64
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 2. Repository Layer (internal/repository/cart_repository.go)
type CartRepository interface {
    Create(ctx context.Context, cart *domain.Cart) error
    GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
    // ... more methods
}

type cartRepository struct {
    redis *redis.Client
    db    *gorm.DB  // Optional backup
}

// 3. Service Layer (internal/service/cart_service.go)
type CartService interface {
    AddItem(ctx context.Context, req AddItemRequest) (*domain.Cart, error)
    // ... more methods
}

type cartService struct {
    cartRepo    repository.CartRepository
    catalogGRPC grpc_client.CatalogClient  // NEW: gRPC client
}

// 4. Handler Layer (internal/handler/http/cart_handler.go)
type CartHandler struct {
    service service.CartService
}

func (h *CartHandler) AddItem(c *fiber.Ctx) error {
    // Parse request, call service, return response
}
```

---

## 📦 Key Components Comparison

### Config Loading

**Catalog Service Pattern** (Recommended):
```go
// internal/config/config.go
type Config struct {
    DBHost        string
    DBPort        string
    DBUser        string
    DBPassword    string
    DBName        string
    HTTPPort      string
    GRPCPort      string
    Env           string
}

func Load() (*Config, error) {
    // Load from env vars
}
```

**User Service Pattern**:
```go
// internal/config/config.go
type Config struct {
    // Similar but includes JWT-specific settings
    JWTSecret         string
    JWTAccessDuration time.Duration
}
```

### Database Initialization

**Catalog Service** (in `config/database.go`):
```go
func InitDatabase(cfg *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(&domain.Book{}, &domain.Author{}, ...)
}

func SeedDatabase(db *gorm.DB) error {
    // Insert sample data for development
}
```

**User Service** (in separate `database/` folder):
```go
// internal/database/database.go
func InitDB(cfg *config.Config) (*gorm.DB, error) {
    // Similar implementation
}
```

### Repository Pattern

**Catalog Service** (One file per entity):
```go
// internal/repository/book_repository.go
type BookRepository interface {
    Create(ctx context.Context, book *domain.Book) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
    // ... more methods
}

type bookRepository struct {
    db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
    return &bookRepository{db: db}
}
```

**User Service** (All repos in one file):
```go
// internal/repository/user_repository.go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// ... methods ...

type AddressRepository struct {
    db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) *AddressRepository {
    return &AddressRepository{db: db}
}
```

---

## 🚀 Implementation Recommendations

### For **Cart Service**:

1. **Follow Catalog Service structure** with these additions:
   - Add `internal/config/redis.go` for Redis initialization
   - Add `internal/grpc_client/catalog_client.go` for gRPC calls
   - Use Redis as primary storage, PostgreSQL as backup (optional)

2. **New Components Needed**:
   ```go
   // Redis client initialization
   func InitRedis(cfg *Config) (*redis.Client, error)
   
   // gRPC client to Catalog Service
   type CatalogClient interface {
       GetBookPrice(ctx context.Context, bookID uuid.UUID) (float64, error)
       CheckStock(ctx context.Context, bookID uuid.UUID) (int, error)
   }
   ```

3. **Storage Strategy**:
   - **Anonymous users**: Redis only (session-based, TTL 7 days)
   - **Authenticated users**: Redis + PostgreSQL backup
   - On login: Merge anonymous cart → user cart

### For **Order Service**:

1. **Follow Catalog Service structure** with these additions:
   - Add `internal/events/` for RabbitMQ pub/sub
   - Add `internal/service/saga_orchestrator.go` for distributed transactions
   - Add `internal/grpc_client/` for Payment & Inventory services

2. **New Components Needed**:
   ```go
   // Saga orchestrator
   type SagaOrchestrator interface {
       ExecuteOrderSaga(ctx context.Context, order *domain.Order) error
       CompensateOrder(ctx context.Context, orderID uuid.UUID) error
   }
   
   // Event publisher
   type EventPublisher interface {
       PublishOrderCreated(order *domain.Order) error
       PublishOrderConfirmed(order *domain.Order) error
   }
   
   // Event consumer
   type EventConsumer interface {
       ConsumePaymentEvents() error
       ConsumeInventoryEvents() error
   }
   ```

3. **Saga Implementation**:
   ```
   1. Create Order (PENDING)
   2. Publish: order.created
   3. Reserve Inventory → wait for: inventory.reserved
   4. Process Payment → wait for: payment.completed
   5. If success: Update to CONFIRMED, publish: order.confirmed
   6. If failure: Compensate (release inventory, refund payment)
   ```

---

## 🔄 Differences to Reconcile

### Option 1: **Follow Catalog Pattern** (Recommended)
- ✅ More modular (separate repository files)
- ✅ Better for larger teams
- ✅ Easier to navigate
- ✅ More recent implementation

### Option 2: **Follow User Pattern**
- ✅ Has `pkg/` utilities (JWT, password)
- ✅ Separate `database/` and `dto/` folders
- ✅ More compact for small services

### **Recommended Hybrid Approach**:
```
service-name/
├── cmd/server/main.go              # Like both
├── internal/
│   ├── config/
│   │   ├── config.go               # Like both
│   │   └── database.go             # Like Catalog (or redis.go for Cart)
│   ├── domain/
│   │   └── models.go               # Like both
│   ├── repository/
│   │   └── *_repository.go         # Like Catalog (separate files)
│   ├── service/
│   │   ├── *_service.go            # Like Catalog
│   │   └── dto.go                  # Like Catalog
│   ├── handler/http/
│   │   ├── routes.go               # Like Catalog
│   │   ├── *_handler.go            # Like Catalog
│   │   ├── responses.go            # Like Catalog
│   │   └── helpers.go              # Like Catalog
│   ├── middleware/                 # Like User (if auth needed)
│   │   └── auth.go
│   └── grpc_client/                # NEW for Cart & Order
│       └── *_client.go
└── pkg/                            # Like User (if utilities needed)
    └── utils/
```

---

## 📝 Summary

### ✅ What You Have:
1. **Catalog Service** - Fully implemented, good reference
2. **User Service** - Fully implemented, slightly different pattern
3. **Proto definitions** - Complete for all services
4. **go.mod files** - All dependencies listed

### 🚧 What You Need to Build:
1. **Cart Service** - Full implementation following Catalog pattern
2. **Order Service** - Full implementation + Saga pattern

### 🎯 Next Steps:
1. **Review this analysis** with your partner
2. **Agree on pattern**: Use Catalog Service structure
3. **Start with Cart Service** (simpler, no saga)
4. **Then Order Service** (more complex, needs saga)
5. **Test integration** between all services

---

## 🤝 Collaboration Strategy

Since your partner implemented Catalog:
1. ✅ **Use their pattern** as the standard
2. ✅ **Ask them** about design decisions (why certain patterns)
3. ✅ **Maintain consistency** across all Go services
4. ✅ **Code review** each other's work

**Recommendation**: Have a quick 15-minute discussion to align on:
- Repository pattern (separate files vs single file)
- Config structure (database in config/ or separate?)
- DTO placement (dto.go in service/ or separate dto/?)
- Error handling approach
- Testing strategy

---

**Ready to implement? Let me know which service you want to start with!** 🚀
