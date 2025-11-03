# 🎯 Recommendation Service

> Intelligent book recommendations powered by tag-based filtering, collaborative filtering, and popularity-based strategies.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [API Endpoints](#api-endpoints)
- [Recommendation Strategies](#recommendation-strategies)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Database Models](#database-models)
- [Deployment](#deployment)

## 🌟 Overview

The Recommendation Service is a Python/FastAPI microservice that provides personalized book recommendations for users. It uses a hybrid approach combining multiple recommendation strategies to deliver accurate and relevant suggestions.

**Port**: 8089 (HTTP), 50059 (gRPC)  
**Database**: PostgreSQL (`recommendations_db`)  
**Cache**: Redis (for caching computed recommendations)  
**Framework**: FastAPI 0.104.1

## ✨ Features

- **Hybrid Recommendation Engine**: Combines tag-based, collaborative filtering, and popularity-based approaches
- **Personalized Recommendations**: Tailored suggestions based on user interaction history
- **Similar Books**: Content-based similarity using tags and metadata
- **Trending Books**: Real-time trending analysis based on recent interactions
- **Interaction Tracking**: Tracks user behavior (views, purchases, cart additions, etc.)
- **Recommendation Caching**: Improves performance with intelligent caching
- **User Preferences**: Explicit preference management (genres, authors, price range)
- **Modular Design**: One model per file for easy maintenance
- **Production Ready**: Kubernetes support, health checks, Docker containerization

## 🏗️ Architecture

```
recommendation-service/
├── app/
│   ├── api/v1/endpoints/      # API route handlers
│   │   ├── recommendations.py  # Recommendation endpoints
│   │   └── health.py          # Health checks
│   ├── core/                  # Core configuration
│   │   ├── config.py          # Settings management
│   │   └── database.py        # Database connection
│   ├── models/                # Database models (one per file)
│   │   ├── user_interaction.py
│   │   ├── book_tag.py
│   │   ├── recommendation_cache.py
│   │   └── user_preference.py
│   ├── schemas/               # Pydantic DTOs
│   │   ├── interaction.py
│   │   └── recommendation.py
│   ├── repository/            # Data access layer
│   │   ├── interaction_repository.py
│   │   └── recommendation_repository.py
│   ├── services/              # Business logic
│   │   ├── recommendation_service.py
│   │   └── strategies/
│   │       ├── base_strategy.py
│   │       ├── tag_based_strategy.py
│   │       ├── collaborative_strategy.py
│   │       └── popular_strategy.py
│   └── main.py               # FastAPI application
├── proto/                     # gRPC definitions
├── migrations/                # Database migrations
├── requirements.txt
├── Dockerfile
└── README.md
```

## 📡 API Endpoints

### Recommendations

#### GET `/api/v1/recommendations/me`
Get personalized recommendations for the current user.

**Headers:**
- `X-User-Id`: UUID of the user

**Query Parameters:**
- `limit` (optional): Number of recommendations (default: 10, max: 50)

**Response:**
```json
{
  "user_id": "uuid",
  "recommendations": [
    {
      "book_id": "uuid",
      "score": 0.95,
      "reason": "Similar to books with tags: science-fiction, thriller"
    }
  ],
  "algorithm": "hybrid",
  "total": 10
}
```

#### GET `/api/v1/recommendations/similar/{book_id}`
Get books similar to a specific book.

**Query Parameters:**
- `limit` (optional): Number of similar books (default: 10, max: 50)

**Response:**
```json
{
  "book_id": "uuid",
  "similar_books": [
    {
      "book_id": "uuid",
      "score": 0.88,
      "reason": "Similar content and tags"
    }
  ],
  "total": 10
}
```

#### GET `/api/v1/recommendations/trending`
Get trending books from recent interactions.

**Query Parameters:**
- `limit` (optional): Number of books (default: 10, max: 50)
- `days` (optional): Number of days to consider (default: 7, max: 30)

**Response:**
```json
[
  {
    "book_id": "uuid",
    "score": 0.92,
    "reason": "Trending with 245 interactions in the last 7 days"
  }
]
```

#### GET `/api/v1/recommendations/popular`
Get popular books based on overall interactions.

**Query Parameters:**
- `limit` (optional): Number of books (default: 10, max: 50)

### Interactions

#### POST `/api/v1/recommendations/interactions`
Track a user interaction with a book.

**Headers:**
- `X-User-Id`: UUID of the user

**Request Body:**
```json
{
  "book_id": "uuid",
  "interaction_type": "view",  // view, add_to_cart, purchase, review, wishlist
  "metadata": "{}"  // Optional JSON metadata
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "book_id": "uuid",
  "interaction_type": "view",
  "score": 1.0,
  "created_at": "2025-11-02T10:30:00Z"
}
```

#### GET `/api/v1/recommendations/interactions/me`
Get user's interaction history.

**Headers:**
- `X-User-Id`: UUID of the user

**Query Parameters:**
- `limit` (optional): Limit results (max: 100)

#### GET `/api/v1/recommendations/interactions/me/stats`
Get statistics about user's interactions.

**Response:**
```json
{
  "total_interactions": 150,
  "views": 80,
  "purchases": 15,
  "reviews": 10,
  "wishlists": 25,
  "cart_additions": 20
}
```

### User Preferences

#### GET `/api/v1/recommendations/preferences/me`
Get user's preferences.

#### PUT `/api/v1/recommendations/preferences/me`
Update user preferences.

**Request Body:**
```json
{
  "preferred_genres": ["science-fiction", "thriller"],
  "preferred_authors": ["uuid1", "uuid2"],
  "min_price": 10.00,
  "max_price": 50.00,
  "preferred_languages": ["English"],
  "excluded_genres": ["horror"]
}
```

#### DELETE `/api/v1/recommendations/preferences/me`
Delete user preferences.

### Health Checks

#### GET `/api/v1/health`
Service health status.

#### GET `/api/v1/ready`
Readiness check (includes database connectivity).

## 🧠 Recommendation Strategies

### 1. Tag-Based Strategy (Content-Based Filtering)

Recommends books with similar tags to what the user has interacted with.

**How it works:**
1. Builds user's tag profile based on interaction history
2. Weights tags by interaction type (purchase=5, review=4, cart=3, wishlist=2, view=1)
3. Finds books with matching tags
4. Scores candidates based on tag overlap

**Best for:**
- Users with clear genre preferences
- Cold start scenarios with some interaction data

### 2. Collaborative Filtering Strategy

Recommends books liked by similar users.

**How it works:**
1. Finds users with similar reading patterns (Jaccard similarity)
2. Identifies books that similar users liked
3. Weights recommendations by user similarity
4. Filters out books already seen

**Best for:**
- Users with substantial interaction history
- Discovering books outside user's usual preferences

### 3. Popular Strategy

Recommends trending or overall popular books.

**How it works:**
1. Aggregates weighted interactions over time window
2. Ranks books by popularity score
3. Considers recency for trending books

**Best for:**
- Cold start (new users)
- Fallback when other strategies fail
- Trending section

### 4. Hybrid Approach (Default)

Combines all three strategies with configurable weights.

**Default weights:**
- Tag-based: 40%
- Collaborative: 40%
- Popular: 20%

**Benefits:**
- More robust recommendations
- Balances personalization with discovery
- Handles various user scenarios

## 🚀 Getting Started

### Prerequisites

- Python 3.11+
- PostgreSQL 15
- Redis 7 (optional, for caching)

### Local Development

1. **Install dependencies:**
```bash
cd services/recommendation-service
pip install -r requirements.txt
```

2. **Set up environment:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Initialize database:**
```bash
# The service will auto-create tables on startup
```

4. **Run the service:**
```bash
python -m app.main
# or
uvicorn app.main:app --reload --port 8089
```

5. **Access the API:**
- API: http://localhost:8089
- Docs: http://localhost:8089/docs
- Redoc: http://localhost:8089/redoc

## ⚙️ Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `Recommendation Service` |
| `PORT` | HTTP port | `8089` |
| `GRPC_PORT` | gRPC port | `50059` |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://...` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `CACHE_TTL` | Cache expiration (seconds) | `3600` |
| `CATALOG_SERVICE_URL` | Catalog service gRPC URL | `localhost:50051` |
| `USER_SERVICE_URL` | User service gRPC URL | `localhost:50052` |
| `ORDER_SERVICE_URL` | Order service gRPC URL | `localhost:50054` |
| `DEFAULT_RECOMMENDATIONS_COUNT` | Default number of recommendations | `10` |
| `MIN_INTERACTIONS_FOR_COLLABORATIVE` | Minimum interactions for collaborative filtering | `3` |
| `TAG_WEIGHT` | Weight for tag-based strategy | `0.4` |
| `COLLABORATIVE_WEIGHT` | Weight for collaborative strategy | `0.4` |
| `POPULAR_WEIGHT` | Weight for popular strategy | `0.2` |
| `LOG_LEVEL` | Logging level | `INFO` |

## 📊 Database Models

### UserInteraction
Tracks user interactions with books for recommendation purposes.

**Columns:**
- `id`: UUID (PK)
- `user_id`: UUID (indexed)
- `book_id`: UUID (indexed)
- `interaction_type`: String (view, add_to_cart, purchase, review, wishlist)
- `score`: Float (1.0 to 5.0 based on type)
- `created_at`: Timestamp

**Indexes:**
- `(user_id, book_id)` - Composite index
- `interaction_type` - For filtering by type
- `created_at` - For time-based queries

### BookTag
Tags associated with books for content-based filtering.

**Columns:**
- `id`: UUID (PK)
- `book_id`: UUID
- `tag_type`: String (genre, topic, author, publisher, attribute)
- `tag_value`: String
- `weight`: Float (0.0 to 1.0)

**Indexes:**
- `(book_id, tag_type, tag_value)` - Unique composite
- `tag_value` - For tag lookups

### RecommendationCache
Cached recommendations for performance.

**Columns:**
- `id`: UUID (PK)
- `user_id`: UUID (unique)
- `book_ids`: UUID[] (array of recommended book IDs)
- `algorithm`: String
- `score`: Float
- `expires_at`: Timestamp

### UserPreference
Explicit user preferences.

**Columns:**
- `id`: UUID (PK)
- `user_id`: UUID (unique)
- `preferred_genres`: String[]
- `preferred_authors`: UUID[]
- `min_price`, `max_price`: Float
- `preferred_languages`: String[]
- `excluded_genres`: String[]

## 🐳 Docker

### Build Image
```bash
docker build -t recommendation-service:latest .
```

### Run Container
```bash
docker run -p 8089:8089 \
  -e DATABASE_URL=postgresql://... \
  -e REDIS_URL=redis://... \
  recommendation-service:latest
```

## ☸️ Kubernetes Deployment

### Apply Manifests
```bash
kubectl apply -f infrastructure/k8s/services/recommendation-service/
```

### Verify Deployment
```bash
kubectl get pods -n production -l app=recommendation-service
kubectl get svc -n production recommendation-service
kubectl logs -f deployment/recommendation-service -n production
```

### Scaling
```bash
# Manual scaling
kubectl scale deployment recommendation-service --replicas=5 -n production

# HPA handles automatic scaling (3-20 replicas)
kubectl get hpa recommendation-service-hpa -n production
```

## 🧪 Testing

### Run Tests
```bash
pytest
```

### Coverage
```bash
pytest --cov=app --cov-report=html
```

## 📚 API Documentation

Interactive API documentation is available at:
- **Swagger UI**: http://localhost:8089/docs
- **ReDoc**: http://localhost:8089/redoc
- **OpenAPI JSON**: http://localhost:8089/openapi.json

## 🔮 Future Enhancements

### Machine Learning Integration
- **Matrix Factorization**: SVD/ALS for collaborative filtering
- **Deep Learning**: Neural collaborative filtering (NCF)
- **NLP**: Book description embeddings with BERT/Sentence Transformers
- **Reinforcement Learning**: Multi-armed bandit for A/B testing strategies

### Advanced Features
- **Context-Aware Recommendations**: Time of day, device, location
- **Session-Based Recommendations**: RNN/GRU for sequential patterns
- **Cross-Domain Recommendations**: Leverage user behavior across categories
- **Explanation Generation**: Natural language explanations for recommendations

## 📝 Example Usage

### Track a Book View
```bash
curl -X POST http://localhost:8089/api/v1/recommendations/interactions \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 123e4567-e89b-12d3-a456-426614174000" \
  -d '{
    "book_id": "987fcdeb-51a2-43f1-9c7d-123456789abc",
    "interaction_type": "view"
  }'
```

### Get Personalized Recommendations
```bash
curl -X GET "http://localhost:8089/api/v1/recommendations/me?limit=10" \
  -H "X-User-Id: 123e4567-e89b-12d3-a456-426614174000"
```

### Get Trending Books
```bash
curl -X GET "http://localhost:8089/api/v1/recommendations/trending?limit=10&days=7"
```

## 🤝 Contributing

This service follows the modular architecture pattern established in the Distributed Bookstore system:
- **One model per file**: Easy to locate and understand
- **Repository pattern**: Clean separation of data access
- **Strategy pattern**: Pluggable recommendation algorithms
- **Dependency injection**: Testable and maintainable code

## 📄 License

Part of the Distributed Bookstore project.

---

**Built with** ❤️ **using FastAPI, SQLAlchemy, and Python 3.11**
