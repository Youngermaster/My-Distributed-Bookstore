# Review Service (FastAPI)

## Overview

The Review Service handles book reviews, ratings, and ML-powered sentiment analysis using FastAPI, SQLAlchemy 2.0, and async PostgreSQL.

## ✨ Features

- ✅ **Create, Read, Update, Delete reviews**
- ✅ **Automatic sentiment analysis** using TextBlob and NLTK
- ✅ **Review voting system** (helpful/not helpful)
- ✅ **Review statistics and aggregations**
- ✅ **Rating distribution analytics**
- ✅ **Async database operations** with SQLAlchemy 2.0
- ✅ **Pydantic v2 validation**
- ✅ **FastAPI with all standard dependencies** (uvicorn, performance extras)
- ✅ **Complete REST API** with OpenAPI documentation
- 🔜 **gRPC support** (folder structure ready)

## 🛠 Tech Stack

- **Python**: 3.11+
- **Framework**: FastAPI 0.115+ with standard dependencies
- **Database**: PostgreSQL 15 with async support (asyncpg)
- **ORM**: SQLAlchemy 2.0 (async)
- **Migrations**: Alembic
- **Validation**: Pydantic v2
- **ML**: NLTK, TextBlob, scikit-learn
- **Server**: Uvicorn (ASGI server with performance extras)

## 📁 Project Structure

```
review-service/
├── app/
│   ├── api/
│   │   └── v1/
│   │       └── endpoints/
│   │           └── reviews.py      # Review endpoints
│   ├── core/
│   │   └── config.py              # Settings with Pydantic
│   ├── db/
│   │   └── base.py                # Database session & base
│   ├── models/
│   │   └── review.py              # SQLAlchemy models
│   ├── schemas/
│   │   └── review.py              # Pydantic schemas
│   ├── services/
│   │   └── review_service.py     # Business logic
│   ├── ml/
│   │   └── sentiment.py           # Sentiment analysis
│   ├── grpc/                      # gRPC (future)
│   └── main.py                    # FastAPI app
├── alembic/                       # Database migrations
├── tests/                         # Tests
├── requirements.txt               # Dependencies
├── Dockerfile                     # Container definition
├── docker-compose.yml             # Local development
├── .env.example                   # Environment template
└── README.md
```

## 🚀 Getting Started

### Option 1: Using Python venv (Recommended for Development)

```bash
# Create virtual environment
python3.11 -m venv venv

# Activate virtual environment
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Copy environment file
cp .env.example .env

# Edit .env with your settings
nano .env

# Start PostgreSQL (if not running)
# You can use docker-compose for just the database:
docker-compose up postgres -d

# Run database migrations (TODO: setup Alembic)
# alembic upgrade head

# Start the service
uvicorn app.main:app --reload --port 8088

# Or using FastAPI CLI (comes with fastapi[standard])
fastapi dev app/main.py

# Or run directly
python -m app.main
```

### Option 2: Using Docker Compose (Full Stack)

```bash
# Start service + PostgreSQL
docker-compose up

# View logs
docker-compose logs -f review-service

# Stop services
docker-compose down

# Rebuild after code changes
docker-compose up --build
```

## 📡 API Endpoints

Once running, visit:
- **OpenAPI Docs**: http://localhost:8088/docs
- **ReDoc**: http://localhost:8088/redoc
- **Health Check**: http://localhost:8088/health

### Available Endpoints

#### Reviews

```http
POST   /api/v1/reviews                    # Create review
GET    /api/v1/reviews/{review_id}        # Get review
PUT    /api/v1/reviews/{review_id}        # Update review
DELETE /api/v1/reviews/{review_id}        # Delete review

GET    /api/v1/reviews/book/{book_id}            # Get book reviews (paginated)
GET    /api/v1/reviews/book/{book_id}/stats      # Get book statistics

POST   /api/v1/reviews/{review_id}/vote          # Vote on review
```

## 📝 Example Usage

### Create a Review

```bash
curl -X POST http://localhost:8088/api/v1/reviews \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "123e4567-e89b-12d3-a456-426614174001",
    "rating": 5,
    "title": "Excellent book!",
    "content": "This book was incredibly insightful and well-written. Highly recommended for anyone interested in distributed systems.",
    "verified_purchase": true
  }'
```

### Get Reviews for a Book

```bash
curl http://localhost:8088/api/v1/reviews/book/123e4567-e89b-12d3-a456-426614174000?page=1&page_size=20
```

### Get Book Statistics

```bash
curl http://localhost:8088/api/v1/reviews/book/123e4567-e89b-12d3-a456-426614174000/stats
```

## 🤖 ML Sentiment Analysis

Reviews are automatically analyzed for sentiment:

- **Positive**: sentiment_score > 0.3
- **Neutral**: sentiment_score between -0.3 and 0.3
- **Negative**: sentiment_score < -0.3

The sentiment is calculated using TextBlob's polarity analysis and stored with each review.

## 🗄️ Database Schema

### reviews table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| book_id | UUID | Book identifier |
| user_id | UUID | User identifier |
| rating | INTEGER | 1-5 stars |
| title | VARCHAR(255) | Review title |
| content | TEXT | Review content |
| sentiment_score | FLOAT | ML sentiment score (-1 to 1) |
| sentiment_label | VARCHAR(20) | positive/neutral/negative |
| verified_purchase | BOOLEAN | Is verified purchase |
| helpful_votes | INTEGER | Number of helpful votes |
| created_at | TIMESTAMP | Creation time |
| updated_at | TIMESTAMP | Last update time |

**Constraints:**
- Unique constraint on (book_id, user_id)
- Check constraint: rating between 1 and 5
- Check constraint: helpful_votes >= 0

### review_votes table

| Column | Type | Description |
|--------|------|-------------|
| review_id | UUID | Review identifier (PK) |
| user_id | UUID | User identifier (PK) |
| is_helpful | BOOLEAN | Helpful vote |
| created_at | TIMESTAMP | Vote time |

## ⚙️ Configuration

Edit `.env` file or set environment variables:

```bash
# Application
DEBUG=True
PORT=8088

# Database
DATABASE_URL=postgresql+asyncpg://bookstore:dev_password@localhost:5432/reviews_db

# CORS
CORS_ORIGINS=["http://localhost:3000","http://localhost:8080"]

# Logging
LOG_LEVEL=INFO
```

## 🧪 Testing

```bash
# Run tests
pytest

# Run with coverage
pytest --cov=app tests/

# Run specific test file
pytest tests/api/test_reviews.py -v
```

## 📊 Development

### Code Quality

```bash
# Format code with Black
black app/

# Lint with Ruff
ruff check app/

# Type checking with mypy
mypy app/
```

### Database Migrations (Alembic)

```bash
# Create a migration
alembic revision --autogenerate -m "Add reviews table"

# Apply migrations
alembic upgrade head

# Rollback migration
alembic downgrade -1
```

## 🔜 Future Enhancements

- [ ] Setup Alembic migrations
- [ ] Add gRPC server for inter-service communication
- [ ] Implement RabbitMQ event consumers
- [ ] Add caching layer (Redis)
- [ ] Implement rate limiting
- [ ] Add comprehensive tests
- [ ] Add authentication middleware
- [ ] Implement review moderation
- [ ] Add image upload support for reviews
- [ ] Implement review replies/comments

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check connection
psql -h localhost -U bookstore -d reviews_db

# View logs
docker-compose logs postgres
```

### NLTK Data Issues

```bash
# Download NLTK data manually
python -c "import nltk; nltk.download('brown'); nltk.download('punkt')"
```

## 📚 API Documentation

When the service is running, comprehensive interactive documentation is available:

- **Swagger UI**: http://localhost:8088/docs
- **ReDoc**: http://localhost:8088/redoc
- **OpenAPI JSON**: http://localhost:8088/openapi.json

## 🏗️ Architecture Notes

- **Async/Await**: All database operations use async/await for better performance
- **Dependency Injection**: FastAPI's dependency injection for database sessions
- **Separation of Concerns**: Clear separation between models, schemas, services, and endpoints
- **Type Safety**: Full type hints throughout the codebase
- **Pydantic Validation**: Automatic request/response validation
- **SQL Alchemy 2.0**: Modern async SQLAlchemy with type-safe queries

## 📄 License

Apache License 2.0
