"""
Repository layer for data access.
"""

from app.repository.interaction_repository import InteractionRepository
from app.repository.recommendation_repository import RecommendationRepository

__all__ = [
    "InteractionRepository",
    "RecommendationRepository",
]
