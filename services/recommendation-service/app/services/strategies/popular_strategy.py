"""
Popular books recommendation strategy.
"""

from typing import List, Dict
from uuid import UUID

from app.services.strategies.base_strategy import RecommendationStrategy
from app.repository import InteractionRepository


class PopularStrategy(RecommendationStrategy):
    """
    Recommends popular books based on overall interaction counts and scores.

    This is a simple but effective strategy that can be used:
    1. For cold start (new users with no interaction history)
    2. As a fallback when other strategies fail
    3. For trending books section
    """

    def __init__(
        self,
        interaction_repo: InteractionRepository,
        recency_days: int = 30,
    ):
        self.interaction_repo = interaction_repo
        self.recency_days = recency_days

    @property
    def name(self) -> str:
        return "popular"

    def recommend(
        self,
        user_id: UUID = None,
        limit: int = 10,
        exclude_books: List[UUID] = None,
    ) -> List[Dict[str, any]]:
        """
        Generate popular book recommendations.

        Note: user_id is optional for this strategy since it's not personalized.
        """
        # Get popular books from recent interactions
        popular_books = self.interaction_repo.get_popular_books(
            limit=limit * 2,  # Get more than needed in case we need to filter
            days=self.recency_days,
        )

        if not popular_books:
            return []

        # Convert to recommendation format
        recommendations = []
        for book_id, total_score, interaction_count in popular_books:
            recommendations.append({
                "book_id": book_id,
                "score": total_score,
                "reason": f"Popular with {interaction_count} interactions in the last {self.recency_days} days",
            })

        # Normalize scores
        recommendations = self._normalize_scores(recommendations)

        # Filter and limit
        return self._filter_books(recommendations, exclude_books, limit)

    def get_trending(
        self,
        limit: int = 10,
        days: int = 7,
    ) -> List[Dict]:
        """
        Get trending books from the last N days.

        This is a specialized version that focuses on very recent activity.
        """
        trending_books = self.interaction_repo.get_popular_books(
            limit=limit,
            days=days,
        )

        recommendations = []
        for book_id, total_score, interaction_count in trending_books:
            recommendations.append({
                "book_id": book_id,
                "score": total_score,
                "reason": f"Trending with {interaction_count} interactions in the last {days} days",
            })

        return self._normalize_scores(recommendations)
