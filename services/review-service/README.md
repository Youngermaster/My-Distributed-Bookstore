# Review Service

## Overview

Handles book reviews, ratings, sentiment analysis using NLTK/TextBlob, and review moderation.

## Technology Stack

- **Language**: Python 3.11+
- **Framework**: FastAPI
- **ORM**: SQLAlchemy
- **Database**: PostgreSQL 15
- **ML Libraries**: NLTK, TextBlob, scikit-learn
- **Messaging**: RabbitMQ (pika)
- **Ports**: HTTP: 8088, gRPC: 50058

## Responsibilities

- Book reviews and ratings management
- Sentiment analysis using ML
- Review moderation
- Rating aggregation
- Helpful votes tracking
- Verified purchase validation

## Database Schema

**reviews**: id, book_id, user_id, rating (1-5), title, content, sentiment_score, sentiment_label, verified_purchase, helpful_votes, created_at, updated_at
**review_votes**: review_id, user_id, is_helpful, created_at

## REST API Endpoints

```
GET    /api/v1/books/{bookId}/reviews     # List reviews
POST   /api/v1/books/{bookId}/reviews     # Submit review
PUT    /api/v1/reviews/{id}               # Update review
DELETE /api/v1/reviews/{id}               # Delete review
POST   /api/v1/reviews/{id}/vote          # Vote helpful/not helpful
GET    /api/v1/reviews/{id}               # Get review details
```

## gRPC Methods

```protobuf
rpc GetBookReviews(GetBookReviewsRequest) returns (ReviewList);
rpc GetAverageRating(GetAverageRatingRequest) returns (RatingInfo);
```

## ML Pipeline

```python
# Sentiment analysis using TextBlob
def analyze_sentiment(review_text):
    blob = TextBlob(review_text)
    polarity = blob.sentiment.polarity
    
    if polarity > 0.3:
        return "positive", polarity
    elif polarity < -0.3:
        return "negative", polarity
    else:
        return "neutral", polarity
```

## Events Published

- `review.submitted` - New review posted
- `review.updated` - Review edited
- `review.deleted` - Review removed

## Events Consumed

- `order.delivered` - Enable review submission for purchased books

## Environment Variables

```bash
PORT=8088
GRPC_PORT=50058
DATABASE_URL=postgresql://bookstore:password@postgres:5432/reviews_db
RABBITMQ_URL=amqp://bookstore:password@rabbitmq:5672/
LOG_LEVEL=INFO
NLTK_DATA_PATH=/app/nltk_data
```

## Getting Started

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Download NLTK data
python -c "import nltk; nltk.download('vader_lexicon')"

# Run development server
uvicorn app.main:app --reload --port 8088

# Run tests
pytest
```

## Next Steps

- [ ] Implement review models
- [ ] Create FastAPI routes
- [ ] Add sentiment analysis
- [ ] Implement gRPC server
- [ ] Add RabbitMQ consumers
- [ ] Write comprehensive tests
