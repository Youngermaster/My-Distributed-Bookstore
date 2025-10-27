# Quick Start Guide - Full Stack Bookstore

This guide will help you quickly start and test the complete bookstore application stack:
- **Frontend**: React + TypeScript + ShadcnUI (Port 5173)
- **API Gateway**: Go Fiber (Port 8080)
- **Catalog Service**: Go Fiber + PostgreSQL (Port 8081)
- **Database**: PostgreSQL (Port 5432)

## Prerequisites

- Docker and Docker Compose installed
- Node.js 18+ and npm installed
- Go 1.21+ installed (for local development)

## Option 1: Quick Start with Docker (Recommended)

### Start Backend Services

```bash
# Navigate to api-gateway directory
cd services/api-gateway

# Start all backend services (PostgreSQL, Catalog Service, API Gateway)
docker compose up -d

# Wait 30 seconds for services to initialize and seed data
sleep 30

# Check logs to verify everything is running
docker compose logs -f
```

You should see:
```
✅ catalog-service started
✅ API Gateway started
✅ Database seeded with sample data
```

### Start Frontend

```bash
# Open a new terminal
cd frontend/customer-app

# Install dependencies (first time only)
npm install

# Start the development server
npm run dev
```

Frontend will be available at: **http://localhost:5173**

## Option 2: Manual Start (All Services)

### Terminal 1 - PostgreSQL
```bash
docker run --name bookstore-postgres \
  -e POSTGRES_USER=bookstore \
  -e POSTGRES_PASSWORD=dev_password \
  -e POSTGRES_DB=catalog_db \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### Terminal 2 - Catalog Service
```bash
cd services/catalog-service

export DATABASE_URL=postgresql://bookstore:dev_password@localhost:5432/catalog_db?sslmode=disable
export HTTP_PORT=8081
export ENV=development

go run cmd/server/main.go
```

### Terminal 3 - API Gateway
```bash
cd services/api-gateway

export PORT=8080
export CATALOG_SERVICE_URL=http://localhost:8081
export ENV=development

go run cmd/server/main.go
```

### Terminal 4 - Frontend
```bash
cd frontend/customer-app

npm install
npm run dev
```

## Testing the Application

### 1. Verify Backend is Running

Open your browser or use curl to test the API Gateway:

```bash
# Health check
curl http://localhost:8080/health | jq

# List books
curl http://localhost:8080/api/v1/catalog/books | jq

# List categories
curl http://localhost:8080/api/v1/catalog/categories | jq

# Search books
curl "http://localhost:8080/api/v1/catalog/books/search?q=distributed" | jq
```

### 2. Test Frontend

Open your browser and navigate to: **http://localhost:5173**

#### Home Page (`/`)
- ✅ See hero section with search bar
- ✅ See "Featured Books" section with 3 books
- ✅ See "Browse by Genre" section with 5 categories

#### Test Search
- Type "distributed" in the search bar
- Should show search results with the "Distributed Systems" book
- Search is debounced (waits 300ms after you stop typing)

#### Test Genre Navigation
- Click on any genre card (e.g., "Programming")
- Should navigate to `/genres/programming`
- Should show all books in that category
- Breadcrumb shows: Home > Genres > Programming

#### Test Book Details
- Click on any book card
- Should navigate to `/books/{id}`
- Should show book details, author info, etc.

#### Test Author Page
- From a book detail page, click on an author name
- Should navigate to `/authors/{id}`
- Should show author biography and their books

## Sample Data

After starting the services, the database is automatically seeded with:

### Books (3 total)
1. **Building Microservices** by Martin Fowler - $49.99
2. **Clean Code** by Robert C. Martin - $44.99
3. **Distributed Systems** by Andrew S. Tanenbaum & Maarten van Steen - $89.99

### Categories (5 total)
- Programming 💻
- Distributed Systems 🌐
- Software Architecture 🏗️
- Databases 🗄️
- Cloud Computing ☁️

### Authors (5 total)
- Martin Fowler
- Robert C. Martin
- Eric Evans
- Andrew S. Tanenbaum
- Maarten van Steen

### Publishers (3 total)
- O'Reilly Media
- Manning Publications
- Addison-Wesley

## Architecture Flow

```
Browser (http://localhost:5173)
    ↓
React App (Vite Dev Server)
    ↓ HTTP/REST
API Gateway (http://localhost:8080)
    ↓ HTTP Proxy
Catalog Service (http://localhost:8081)
    ↓ GORM
PostgreSQL (localhost:5432)
```

## Component Structure

### Frontend Components Created
- **BookCard** - Displays individual books with cover, title, price, and actions
- **GenreCard** - Displays category/genre with icon and navigation
- **SearchBar** - Debounced search input
- **BookGrid** - Responsive grid layout for books with loading/empty states

### Frontend Pages Created
- **Home** (`/`) - Hero section, featured books, and genre grid
- **Genres** (`/genres`) - All categories in a grid
- **GenreDetail** (`/genres/:slug`) - Books in specific genre
- **AuthorDetail** (`/authors/:id`) - Author info and their books
- **BookList** (`/books`) - All books (existing)
- **BookDetail** (`/books/:id`) - Book details (existing)

## Stopping the Services

### Docker Compose
```bash
cd services/api-gateway
docker compose down

# To also remove volumes (database data)
docker compose down -v
```

### Manual Services
- Press `Ctrl+C` in each terminal running a service
- Stop PostgreSQL: `docker stop bookstore-postgres`

## Troubleshooting

### Backend not responding
```bash
# Check if services are running
docker compose ps

# Check logs
docker compose logs catalog-service
docker compose logs api-gateway

# Restart services
docker compose restart
```

### Frontend can't connect to API
```bash
# Verify API Gateway is running
curl http://localhost:8080/health

# Check frontend .env file
cat frontend/customer-app/.env
# Should have: VITE_API_URL=http://localhost:8080
```

### Database connection errors
```bash
# Reset database
docker compose down -v
docker compose up -d

# Wait 30 seconds for seed data to be inserted
```

### Port already in use
```bash
# Check what's using the port
lsof -i :5173  # Frontend
lsof -i :8080  # API Gateway
lsof -i :8081  # Catalog Service
lsof -i :5432  # PostgreSQL

# Kill the process or change the port in environment variables
```

## Next Steps

Now that everything is running, you can:

1. **Browse the catalog** - Explore books and genres
2. **Test the API** - Use curl or Postman to test API endpoints
3. **Develop new features** - Add more services (Cart, Order, User services)
4. **Customize the frontend** - Modify components and add new pages
5. **Add more data** - Insert more books, authors, and categories

## Useful Commands

```bash
# View API Gateway logs
docker compose logs -f api-gateway

# View Catalog Service logs
docker compose logs -f catalog-service

# View database logs
docker compose logs -f postgres

# Execute SQL in database
docker exec -it bookstore-postgres psql -U bookstore -d catalog_db

# Rebuild services after code changes
docker compose up -d --build

# Frontend build for production
cd frontend/customer-app
npm run build
```

## Documentation

- **TESTING_GUIDE.md** - Comprehensive testing guide
- **services/catalog-service/README.md** - Catalog Service API documentation
- **services/api-gateway/README.md** - API Gateway documentation
- **CLAUDE.md** - Full project architecture and guidelines

Happy coding! 🚀
