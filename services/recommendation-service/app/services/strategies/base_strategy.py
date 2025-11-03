"""
Base recommendation strategy interface.
"""

from abc import ABC, abstractmethod
from typing import List, Dict
from uuid import UUID


class RecommendationStrategy(ABC):
    """
    Abstract base class for recommendation strategies.

    All recommendation strategies should inherit from this class
    and implement the recommend method.
    """

    @abstractmethod
    def recommend(
        self,
        user_id: UUID,
        limit: int = 10,
        exclude_books: List[UUID] = None,
    ) -> List[Dict[str, any]]:
        """
        Generate recommendations for a user.

        Args:
            user_id: The user to generate recommendations for
            limit: Maximum number of recommendations to return
            exclude_books: List of book IDs to exclude from recommendations

        Returns:
            List of dictionaries with keys:
            - book_id: UUID of the recommended book
            - score: Float score (0.0 to 1.0) indicating relevance
            - reason: Optional string explaining why this was recommended
        """
        pass

    @property
    @abstractmethod
    def name(self) -> str:
        """Return the name of this strategy."""
        pass

    def _normalize_scores(self, recommendations: List[Dict]) -> List[Dict]:
        """
        Normalize scores to be between 0.0 and 1.0.

        Args:
            recommendations: List of recommendation dicts with 'score' key

        Returns:
            Same list with normalized scores
        """
        if not recommendations:
            return []

        scores = [rec["score"] for rec in recommendations]
        max_score = max(scores) if scores else 1.0
        min_score = min(scores) if scores else 0.0

        # Avoid division by zero
        if max_score == min_score:
            for rec in recommendations:
                rec["score"] = 1.0 if max_score > 0 else 0.0
            return recommendations

        # Normalize to 0-1 range
        for rec in recommendations:
            rec["score"] = (rec["score"] - min_score) / (max_score - min_score)

        return recommendations

    def _filter_books(
        self,
        recommendations: List[Dict],
        exclude_books: List[UUID] = None,
        limit: int = 10,
    ) -> List[Dict]:
        """
        Filter and limit recommendations.

        Args:
            recommendations: List of recommendation dicts
            exclude_books: Books to exclude
            limit: Maximum number to return

        Returns:
            Filtered and limited list
        """
        if exclude_books:
            exclude_set = set(exclude_books)
            recommendations = [
                rec for rec in recommendations if rec["book_id"] not in exclude_set
            ]

        # Sort by score descending
        recommendations.sort(key=lambda x: x["score"], reverse=True)

        return recommendations[:limit]
