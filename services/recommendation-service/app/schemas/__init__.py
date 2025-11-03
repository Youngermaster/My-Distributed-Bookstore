"""
Pydantic schemas for request/response validation.
"""

from app.schemas.interaction import (
    InteractionBase,
    InteractionCreate,
    InteractionResponse,
    InteractionStats,
)
from app.schemas.recommendation import (
    RecommendationItem,
    RecommendationResponse,
    SimilarBooksRequest,
    SimilarBooksResponse,
    UserPreferenceCreate,
    UserPreferenceResponse,
)

__all__ = [
    "InteractionBase",
    "InteractionCreate",
    "InteractionResponse",
    "InteractionStats",
    "RecommendationItem",
    "RecommendationResponse",
    "SimilarBooksRequest",
    "SimilarBooksResponse",
    "UserPreferenceCreate",
    "UserPreferenceResponse",
]
