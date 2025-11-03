"""
Database models for the Recommendation Service.

Each model is in its own file for better modularity and maintainability.
"""

from app.models.user_interaction import UserInteraction
from app.models.book_tag import BookTag
from app.models.recommendation_cache import RecommendationCache
from app.models.user_preference import UserPreference

__all__ = [
    "UserInteraction",
    "BookTag",
    "RecommendationCache",
    "UserPreference",
]
