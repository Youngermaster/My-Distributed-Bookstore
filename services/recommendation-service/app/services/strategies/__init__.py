"""
Recommendation strategies for generating personalized book recommendations.
"""

from app.services.strategies.base_strategy import RecommendationStrategy
from app.services.strategies.tag_based_strategy import TagBasedStrategy
from app.services.strategies.collaborative_strategy import CollaborativeStrategy
from app.services.strategies.popular_strategy import PopularStrategy

__all__ = [
    "RecommendationStrategy",
    "TagBasedStrategy",
    "CollaborativeStrategy",
    "PopularStrategy",
]
