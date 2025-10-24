# Recommendation Service

## Overview

Provides personalized book recommendations using collaborative filtering and content-based filtering algorithms.

## Technology Stack

- **Language**: Python 3.11+
- **Framework**: FastAPI
- **ORM**: SQLAlchemy
- **Database**: PostgreSQL 15 + pgvector
- **ML Libraries**: scikit-learn, numpy, pandas
- **Messaging**: RabbitMQ
- **Ports**: HTTP: 8089, gRPC: 50059

## Responsibilities

- Personalized book recommendations
- Collaborative filtering (user-based)
- Content-based filtering (book similarities)
- "Customers who bought this also bought" recommendations
- Trending books calculation
- User interaction tracking

## Database Schema

**user_interactions**: id, user_id, book_id, interaction_type, score, created_at
**recommendations_cache**: id, user_id, book_ids (array), algorithm, score, expires_at, created_at
**book_embeddings**: book_id, embedding (vector(128)), updated_at

## Interaction Types & Weights

- View: 1
- Add to cart: 3
- Purchase: 5
- Review: 4

## REST API Endpoints

```
GET /api/v1/recommendations              # Get personalized recommendations
GET /api/v1/books/{bookId}/similar       # Get similar books
GET /api/v1/recommendations/trending     # Get trending books
```

## gRPC Methods

```protobuf
rpc GetRecommendations(GetRecommendationsRequest) returns (BookList);
rpc GetSimilarBooks(GetSimilarBooksRequest) returns (BookList);
```

## ML Algorithms

### Collaborative Filtering (KNN)
```python
from sklearn.neighbors import NearestNeighbors

def collaborative_filtering(user_id, n_recommendations=10):
    # Create user-item matrix
    user_item_matrix = create_user_item_matrix()
    
    # Train KNN model
    model = NearestNeighbors(metric='cosine', algorithm='brute')
    model.fit(user_item_matrix)
    
    # Find similar users
    distances, indices = model.kneighbors(
        user_item_matrix[user_id], 
        n_neighbors=20
    )
    
    return aggregate_recommendations(indices)
```

### Content-Based Filtering
Uses book embeddings (genres, authors, descriptions) with cosine similarity.

## Events Consumed

- `order.completed` - Update user interaction history
- `review.submitted` - Update interaction scores
- `user.registered` - Initialize recommendation cache

## Environment Variables

```bash
PORT=8089
GRPC_PORT=50059
DATABASE_URL=postgresql://bookstore:password@postgres:5432/recommendations_db
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/
CATALOG_SERVICE_URL=catalog-service:50051
MODEL_RETRAIN_INTERVAL=3600  # 1 hour
```

## Getting Started

```bash
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8089
pytest
```

## Next Steps

- [ ] Implement interaction tracking
- [ ] Build collaborative filtering
- [ ] Build content-based filtering
- [ ] Add recommendation caching
- [ ] Implement gRPC server
- [ ] Add model training pipeline
- [ ] Write tests
