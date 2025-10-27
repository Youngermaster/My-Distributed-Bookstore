# Complete Testing Guide: Frontend → API Gateway → Catalog Service

## What We've Built

### ✅ Completed

1. **Catalog Service** (Port 8081)
   - Full CRUD for Books, Authors, Categories, Publishers
   - Search and filtering
   - PostgreSQL database with GORM
   - Auto-migrations and seed data
   - Docker support

2. **API Gateway** (Port 8080)
   - HTTP proxy to Catalog Service
   - Rate limiting (100 req/min)
   - CORS handling
   - Request logging
   - Health checks
   - Docker Compose with full stack

3. **Frontend Setup**
   - React 19 + TypeScript
   - Vite
   - TanStack Query
   - ShadcnUI components
   - API client updated to use API Gateway

### Architecture Flow

```
Browser (React App)
    ↓ HTTP
API Gateway (Port 8080)
    ↓ HTTP Proxy
Catalog Service (Port 8081)
    ↓ GORM
PostgreSQL (Port 5432)
```

## Quick Start - Full Stack Testing

### Option 1: Using Docker Compose (Recommended)

```bash
# Start everything (Postgres + Catalog + API Gateway)
cd services/api-gateway
docker compose up -d

# Wait 30 seconds for services to start
docker compose logs -f

# You should see:
# ✅ catalog-service started
# ✅ API Gateway started
```

**Services Running:**
- PostgreSQL: localhost:5432
- Catalog Service: localhost:8081
- API Gateway: localhost:8080

### Option 2: Manual Start

**Terminal 1 - PostgreSQL:**
```bash
docker run --name bookstore-postgres \
  -e POSTGRES_USER=bookstore \
  -e POSTGRES_PASSWORD=dev_password \
  -e POSTGRES_DB=catalog_db \
  -p 5432:5432 \
  -d postgres:15-alpine
```

**Terminal 2 - Catalog Service:**
```bash
cd services/catalog-service
export DATABASE_URL=postgresql://bookstore:dev_password@localhost:5432/catalog_db?sslmode=disable
export HTTP_PORT=8081
export ENV=development
go run cmd/server/main.go
```

**Terminal 3 - API Gateway:**
```bash
cd services/api-gateway
export PORT=8080
export CATALOG_SERVICE_URL=http://localhost:8081
export ENV=development
go run cmd/server/main.go
```

**Terminal 4 - Frontend:**
```bash
cd frontend/customer-app
pnpm install
pnpm run dev
```

## Manual API Testing

### 1. Test API Gateway Health

```bash
curl http://localhost:8080/health | jq
```

Expected:
```json
{
  "service": "api-gateway",
  "services": {
    "catalog": {
      "healthy": true
    }
  },
  "status": "ok"
}
```

### 2. Test Gateway Info

```bash
curl http://localhost:8080/ | jq
```

### 3. List Books Through Gateway

```bash
curl http://localhost:8080/api/v1/catalog/books | jq
```

You should see 3 sample books:
- Building Microservices
- Clean Code
- Distributed Systems: Principles and Paradigms

### 4. Search Books

```bash
curl "http://localhost:8080/api/v1/catalog/books/search?q=distributed" | jq
```

### 5. List Categories

```bash
curl http://localhost:8080/api/v1/catalog/categories | jq
```

You should see 5 categories:
- Programming
- Distributed Systems
- Software Architecture
- Databases
- Cloud Computing

### 6. List Authors

```bash
curl http://localhost:8080/api/v1/catalog/authors | jq
```

## Frontend Testing ✅ COMPLETED

The frontend has been fully implemented with all components and pages!

### Pages Implemented

1. **Home Page** (`/`) ✅
   - Displays featured books (first 10)
   - Shows categories grid (first 8)
   - Integrated search bar with real-time search
   - Hero section with gradient background
   - Responsive layout

2. **Genres Page** (`/genres`) ✅
   - Lists all categories in a grid
   - Shows category icons and names
   - Click to navigate to genre detail page
   - Loading states and error handling

3. **Genre Detail Page** (`/genres/:slug`) ✅
   - Displays all books in a specific category
   - Breadcrumb navigation (Home > Genres > [Genre Name])
   - Gradient header with category info
   - Book count display
   - Empty state handling

4. **Author Detail Page** (`/authors/:id`) ✅
   - Author information card with avatar
   - Biography section (if available)
   - Birth date display
   - List of all books by the author
   - Breadcrumb navigation

### Components Implemented

All components created in `frontend/customer-app/src/components/`:

**BookCard.tsx** ✅
- Displays book cover image (or placeholder)
- Shows title, authors, and price
- Add to cart button
- Add to wishlist button
- Stock quantity warnings
- Out of stock overlay
- Hover effects

**GenreCard.tsx** ✅
- Category icon based on slug
- Category name and book count
- Hover animations and scale effect
- Chevron icon for navigation
- Click handler

**SearchBar.tsx** ✅
- Debounced search input (300ms)
- Search icon
- Clear button when text is entered
- Responsive sizing
- Real-time search integration

**BookGrid.tsx** ✅
- Responsive grid layout (1-5 columns)
- Loading spinner
- Empty state with helpful message
- Click handlers for books
- Add to cart and wishlist actions

### Starting the Frontend

```bash
# Navigate to frontend directory
cd frontend/customer-app

# Install dependencies (if not already done)
pnpm install

# Start the development server
pnpm run dev

# Frontend will be available at http://localhost:5173
```

### Frontend Routes

- **`/`** - Home page with featured books and genres
- **`/books`** - All books list
- **`/books/:id`** - Book details page
- **`/genres`** - All genres/categories
- **`/genres/:slug`** - Books in specific genre (e.g., `/genres/programming`)
- **`/authors/:id`** - Author profile and their books
- **`/login`** - Login page
- **`/register`** - Registration page
- **`/wishlist`** - User's wishlist (protected)
- **`/admin/books`** - Admin book management (protected, admin only)

### Testing the Complete Flow

1. **Start all services:**
   ```bash
   # Terminal 1 - Start backend services
   cd services/api-gateway
   docker compose up -d

   # Wait 30 seconds for services to initialize
   docker compose logs -f
   ```

2. **Start frontend:**
   ```bash
   # Terminal 2 - Start React app
   cd frontend/customer-app
   pnpm run dev
   ```

3. **Test the flow:**
   - Open browser: http://localhost:5173
   - You should see the Home page with:
     - Hero section with search bar
     - Featured books section (3 sample books)
     - Browse by Genre section (5 categories)

4. **Test Search:**
   - Type "distributed" in the search bar
   - Should see "Distributed Systems" book in results
   - Search is debounced (waits 300ms after you stop typing)

5. **Test Genre Navigation:**
   - Click on "Programming" genre card
   - Should navigate to `/genres/programming`
   - Should show all programming books
   - Breadcrumb shows: Home > Genres > Programming

6. **Test Author Navigation:**
   - Click on any book card
   - In book details, click on an author name
   - Should navigate to `/authors/{id}`
   - Should show author bio and their books

### Sample Data Available

The backend is seeded with:

**3 Books:**
- Building Microservices by Martin Fowler ($49.99)
- Clean Code by Robert C. Martin ($44.99)
- Distributed Systems by Andrew S. Tanenbaum & Maarten van Steen ($89.99)

**5 Categories:**
- Programming
- Distributed Systems
- Software Architecture
- Databases
- Cloud Computing

**5 Authors:**
- Martin Fowler
- Robert C. Martin
- Eric Evans
- Andrew S. Tanenbaum
- Maarten van Steen

### API Integration

The API client in `src/lib/api.ts` uses:

```typescript
// All endpoints go through API Gateway
baseURL: "http://localhost:8080"

// Books API
booksAPI.list({ page: 1, page_size: 10 })           // GET /api/v1/catalog/books
booksAPI.search("query")                             // GET /api/v1/catalog/books/search?q=query
booksAPI.get(id)                                     // GET /api/v1/catalog/books/:id

// Categories API
categoriesAPI.list()                                 // GET /api/v1/catalog/categories
categoriesAPI.get(id)                                // GET /api/v1/catalog/categories/:id

// Authors API
authorsAPI.list()                                    // GET /api/v1/catalog/authors
authorsAPI.get(id)                                   // GET /api/v1/catalog/authors/:id
```

## Testing Checklist

### Backend Testing

- [ ] PostgreSQL is running (port 5432)
- [ ] Catalog Service is running (port 8081)
- [ ] API Gateway is running (port 8080)
- [ ] Health check returns OK
- [ ] Can list books through gateway
- [ ] Can search books
- [ ] Can list categories
- [ ] Can list authors
- [ ] Rate limiting works (try 110 requests)

### Frontend Testing ✅ Ready to Test

- [ ] Home page loads at http://localhost:5173
- [ ] Featured books display correctly (should show 3 books)
- [ ] Categories grid displays (should show 5 genres)
- [ ] Search bar works with debouncing
- [ ] Searching for "distributed" shows correct results
- [ ] Genres page shows all categories (/genres)
- [ ] Clicking a genre navigates to genre detail page
- [ ] Genre detail page shows books in that category
- [ ] Clicking on book card navigates to book details
- [ ] Author page shows author info and books (/authors/:id)
- [ ] Navigation bar works (Books, Genres links)
- [ ] Loading states show (spinners)
- [ ] Empty states show when no data
- [ ] Error handling works (try navigating to /genres/invalid-slug)
- [ ] Breadcrumb navigation works on detail pages
- [ ] Responsive design works (try mobile view)

## Sample Data Available

After starting the services, you have:

**Publishers:**
- O'Reilly Media
- Manning Publications
- Addison-Wesley

**Authors:**
- Martin Fowler
- Robert C. Martin
- Eric Evans
- Andrew S. Tanenbaum
- Maarten van Steen

**Categories:**
- Programming
- Distributed Systems
- Software Architecture
- Databases
- Cloud Computing

**Books:**
1. Building Microservices (Martin Fowler, $49.99)
2. Clean Code (Robert C. Martin, $44.99)
3. Distributed Systems (Tanenbaum & van Steen, $89.99)

## Environment Variables

**Frontend (.env):**
```
VITE_API_URL=http://localhost:8080
```

**API Gateway:**
```
PORT=8080
CATALOG_SERVICE_URL=http://localhost:8081
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m
```

**Catalog Service:**
```
HTTP_PORT=8081
DATABASE_URL=postgresql://bookstore:dev_password@localhost:5432/catalog_db?sslmode=disable
ENV=development
```

## Next Steps

1. **Create Frontend Components:**
   ```bash
   cd frontend/customer-app/src/components
   # Create: BookCard.tsx, GenreCard.tsx, SearchBar.tsx, BookGrid.tsx
   ```

2. **Create Pages:**
   ```bash
   cd frontend/customer-app/src/pages
   # Create: Home.tsx, Genres.tsx, GenreDetail.tsx, AuthorDetail.tsx
   ```

3. **Update Routing:**
   ```typescript
   // In App.tsx or router config
   <Route path="/" element={<Home />} />
   <Route path="/genres" element={<Genres />} />
   <Route path="/genres/:slug" element={<GenreDetail />} />
   <Route path="/authors/:id" element={<AuthorDetail />} />
   ```

4. **Run and Test:**
   ```bash
   pnpm run dev
   # Open http://localhost:5173
   ```

## Troubleshooting

**API Gateway can't reach Catalog Service:**
```bash
# Check if catalog is running
curl http://localhost:8081/health

# Check Docker network (if using Docker)
docker network inspect bookstore-network
```

**Frontend can't reach API Gateway:**
```bash
# Check CORS (should allow *)
# Check API Gateway is running on 8080
curl http://localhost:8080/health

# Check frontend .env file
cat frontend/customer-app/.env
```

**Database connection errors:**
```bash
# Reset database
docker compose down -v
docker compose up -d

# Wait 30 seconds for seed data
```

## Project URLs

- Frontend: http://localhost:5173 (Vite dev server)
- API Gateway: http://localhost:8080
- Catalog Service: http://localhost:8081 (direct access, not recommended)
- PostgreSQL: localhost:5432

## Summary

You now have a fully functional backend stack:
- ✅ Catalog Service with full CRUD
- ✅ API Gateway with proxying and rate limiting
- ✅ Database with seed data
- ✅ Docker Compose setup

**What's left:** Build the React frontend pages and components to consume the APIs.

The API client is ready, ShadcnUI is configured, and you just need to create the UI components and pages as outlined above.
