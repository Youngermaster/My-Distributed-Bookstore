"""
Pydantic schemas for recommendations.
"""

from pydantic import BaseModel, Field
from typing import List, Optional
from uuid import UUID


class RecommendationItem(BaseModel):
    """A single recommendation item."""

    book_id: UUID
    score: float = Field(..., description="Relevance score (0.0 to 1.0)")
    reason: Optional[str] = Field(None, description="Explanation for why this was recommended")


class RecommendationResponse(BaseModel):
    """Response containing recommended books."""

    user_id: UUID
    recommendations: List[RecommendationItem]
    algorithm: str = Field(..., description="Algorithm used: tag_based, collaborative, popular, hybrid")
    total: int = Field(..., description="Total number of recommendations")


class SimilarBooksRequest(BaseModel):
    """Request for finding similar books."""

    book_id: UUID
    limit: int = Field(10, ge=1, le=50, description="Maximum number of similar books to return")


class SimilarBooksResponse(BaseModel):
    """Response containing similar books."""

    book_id: UUID
    similar_books: List[RecommendationItem]
    total: int


class UserPreferenceCreate(BaseModel):
    """Schema for creating/updating user preferences."""

    preferred_genres: Optional[List[str]] = None
    preferred_authors: Optional[List[UUID]] = None
    min_price: Optional[float] = Field(None, ge=0)
    max_price: Optional[float] = Field(None, ge=0)
    preferred_languages: Optional[List[str]] = None
    excluded_genres: Optional[List[str]] = None


class UserPreferenceResponse(BaseModel):
    """Schema for user preference responses."""

    id: UUID
    user_id: UUID
    preferred_genres: Optional[List[str]]
    preferred_authors: Optional[List[UUID]]
    min_price: Optional[float]
    max_price: Optional[float]
    preferred_languages: Optional[List[str]]
    excluded_genres: Optional[List[str]]

    class Config:
        from_attributes = True
