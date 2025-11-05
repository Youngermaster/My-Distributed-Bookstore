"""
Tag-based recommendation strategy (content-based filtering).
"""

from typing import List, Dict
from uuid import UUID
from collections import defaultdict

from app.services.strategies.base_strategy import RecommendationStrategy
from app.repository import InteractionRepository, RecommendationRepository


class TagBasedStrategy(RecommendationStrategy):
    """
    Recommends books based on tags/categories of books the user has interacted with.

    This is a simple content-based filtering approach that:
    1. Finds all books the user has interacted with
    2. Extracts tags from those books
    3. Finds other books with similar tags
    4. Scores them based on tag overlap and interaction weight
    """

    def __init__(
        self,
        interaction_repo: InteractionRepository,
        recommendation_repo: RecommendationRepository,
    ):
        self.interaction_repo = interaction_repo
        self.recommendation_repo = recommendation_repo

    @property
    def name(self) -> str:
        return "tag_based"

    def recommend(
        self,
        user_id: UUID,
        limit: int = 10,
        exclude_books: List[UUID] = None,
    ) -> List[Dict[str, any]]:
        """Generate tag-based recommendations."""
        # Get user's interacted books with their scores
        weighted_books = self.interaction_repo.get_weighted_book_scores(user_id)

        if not weighted_books:
            # No interactions yet, return empty
            return []

        # Get all book IDs the user has interacted with (to exclude)
        user_book_ids = list(weighted_books.keys())

        # Build user's tag profile
        tag_scores = self._build_tag_profile(weighted_books)

        if not tag_scores:
            return []

        # Find candidate books based on tags
        candidate_books = self._find_candidate_books(tag_scores)

        # Score candidates
        recommendations = self._score_candidates(candidate_books, tag_scores)

        # Normalize scores
        recommendations = self._normalize_scores(recommendations)

        # Filter out books user already interacted with and apply limit
        exclude_list = user_book_ids.copy()
        if exclude_books:
            exclude_list.extend(exclude_books)

        return self._filter_books(recommendations, exclude_list, limit)

    def _build_tag_profile(self, weighted_books: Dict[UUID, float]) -> Dict[str, float]:
        """
        Build a user's tag profile based on their interactions.

        Args:
            weighted_books: Dict mapping book_id to interaction weight

        Returns:
            Dict mapping tag_value to aggregated score
        """
        tag_scores = defaultdict(float)

        for book_id, interaction_weight in weighted_books.items():
            # Get tags for this book
            tags = self.recommendation_repo.get_book_tags(book_id)

            for tag in tags:
                # Aggregate score: interaction_weight * tag_weight
                tag_scores[tag.tag_value] += interaction_weight * tag.weight

        return dict(tag_scores)

    def _find_candidate_books(self, tag_scores: Dict[str, float]) -> Dict[UUID, List[str]]:
        """
        Find candidate books that have tags in the user's profile.

        Args:
            tag_scores: User's tag profile

        Returns:
            Dict mapping book_id to list of matching tags
        """
        candidate_books = defaultdict(list)

        # Get top tags (sorted by score)
        top_tags = sorted(tag_scores.items(), key=lambda x: x[1], reverse=True)[:20]

        for tag_value, _ in top_tags:
            # Find books with this tag
            # Note: We assume tags have both type and value, but for simplicity
            # we're just using tag_value. In production, you'd want to be more specific.
            books = self.recommendation_repo.get_books_by_tag(
                tag_type="genre", tag_value=tag_value, limit=50
            )

            for book_id in books:
                candidate_books[book_id].append(tag_value)

        return dict(candidate_books)

    def _score_candidates(
        self,
        candidate_books: Dict[UUID, List[str]],
        tag_scores: Dict[str, float],
    ) -> List[Dict]:
        """
        Score candidate books based on tag overlap.

        Args:
            candidate_books: Dict mapping book_id to matching tags
            tag_scores: User's tag profile with scores

        Returns:
            List of recommendation dicts
        """
        recommendations = []

        for book_id, matching_tags in candidate_books.items():
            # Calculate score as sum of matching tag scores
            score = sum(tag_scores.get(tag, 0) for tag in matching_tags)

            # Create reason
            top_tags = sorted(matching_tags, key=lambda t: tag_scores.get(t, 0), reverse=True)[:3]
            reason = f"Similar to books with tags: {', '.join(top_tags)}"

            recommendations.append({
                "book_id": book_id,
                "score": score,
                "reason": reason,
            })

        return recommendations
