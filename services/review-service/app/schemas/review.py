"""
Pydantic schemas for request/response validation.
"""
from datetime import datetime
from typing import Optional
from uuid import UUID

from pydantic import BaseModel, Field, ConfigDict


# ============================================================================
# Request Schemas
# ============================================================================

class ReviewCreateRequest(BaseModel):
    """Schema for creating a new review."""

    book_id: UUID = Field(..., description="ID of the book being reviewed")
    user_id: UUID = Field(..., description="ID of the user creating the review")
    rating: int = Field(..., ge=1, le=5, description="Rating from 1 to 5 stars")
    title: str = Field(..., min_length=1, max_length=255, description="Review title")
    content: str = Field(..., min_length=10, description="Review content")
    verified_purchase: bool = Field(default=False, description="Whether this is a verified purchase")

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "book_id": "123e4567-e89b-12d3-a456-426614174000",
                "user_id": "123e4567-e89b-12d3-a456-426614174001",
                "rating": 5,
                "title": "Excellent book!",
                "content": "This book was incredibly insightful and well-written. Highly recommended for anyone interested in distributed systems.",
                "verified_purchase": True
            }
        }
    )


class ReviewUpdateRequest(BaseModel):
    """Schema for updating an existing review."""

    rating: Optional[int] = Field(None, ge=1, le=5, description="Rating from 1 to 5 stars")
    title: Optional[str] = Field(None, min_length=1, max_length=255, description="Review title")
    content: Optional[str] = Field(None, min_length=10, description="Review content")

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "rating": 4,
                "title": "Updated: Very good book",
                "content": "After further reflection, I still think this is a great book, though not perfect."
            }
        }
    )


class ReviewVoteRequest(BaseModel):
    """Schema for voting on review helpfulness."""

    user_id: UUID = Field(..., description="ID of the user voting")
    is_helpful: bool = Field(..., description="True if helpful, False if not helpful")

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "user_id": "123e4567-e89b-12d3-a456-426614174002",
                "is_helpful": True
            }
        }
    )


# ============================================================================
# Response Schemas
# ============================================================================

class ReviewResponse(BaseModel):
    """Schema for review response."""

    id: UUID
    book_id: UUID
    user_id: UUID
    rating: int
    title: str
    content: str
    sentiment_score: Optional[float] = None
    sentiment_label: Optional[str] = None
    verified_purchase: bool
    helpful_votes: int
    created_at: datetime
    updated_at: datetime

    model_config = ConfigDict(from_attributes=True)


class ReviewListResponse(BaseModel):
    """Schema for paginated list of reviews."""

    reviews: list[ReviewResponse]
    total: int
    page: int
    page_size: int
    total_pages: int


class ReviewStatsResponse(BaseModel):
    """Schema for review statistics."""

    book_id: UUID
    total_reviews: int
    average_rating: float
    rating_distribution: dict[int, int]  # {rating: count}
    sentiment_distribution: dict[str, int]  # {sentiment_label: count}


class HealthResponse(BaseModel):
    """Health check response."""

    status: str
    service: str
    version: str
    timestamp: datetime


class MessageResponse(BaseModel):
    """Generic message response."""

    message: str
    detail: Optional[str] = None
