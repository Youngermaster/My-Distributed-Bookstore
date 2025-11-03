"""
Collaborative filtering recommendation strategy.
"""

from typing import List, Dict, Set
from uuid import UUID
from collections import defaultdict

from app.services.strategies.base_strategy import RecommendationStrategy
from app.repository import InteractionRepository


class CollaborativeStrategy(RecommendationStrategy):
    """
    Recommends books based on similar users' preferences.

    This implements a simple user-user collaborative filtering:
    1. Find users who have interacted with the same books as the target user
    2. Calculate similarity between target user and other users
    3. Recommend books that similar users liked but target user hasn't seen
    """

    def __init__(
        self,
        interaction_repo: InteractionRepository,
        min_common_books: int = 2,
    ):
        self.interaction_repo = interaction_repo
        self.min_common_books = min_common_books

    @property
    def name(self) -> str:
        return "collaborative"

    def recommend(
        self,
        user_id: UUID,
        limit: int = 10,
        exclude_books: List[UUID] = None,
    ) -> List[Dict[str, any]]:
        """Generate collaborative filtering recommendations."""
        # Get user's interacted books
        user_books = set(self.interaction_repo.get_user_book_ids(user_id))

        if len(user_books) < self.min_common_books:
            # Not enough data for collaborative filtering
            return []

        # Find similar users
        similar_users = self._find_similar_users(user_id, user_books)

        if not similar_users:
            return []

        # Get recommendations from similar users
        recommendations = self._get_recommendations_from_similar_users(
            user_books, similar_users
        )

        # Normalize scores
        recommendations = self._normalize_scores(recommendations)

        # Filter and limit
        exclude_list = list(user_books)
        if exclude_books:
            exclude_list.extend(exclude_books)

        return self._filter_books(recommendations, exclude_list, limit)

    def _find_similar_users(
        self, user_id: UUID, user_books: Set[UUID]
    ) -> Dict[UUID, float]:
        """
        Find users similar to the target user.

        Uses Jaccard similarity: |A ∩ B| / |A ∪ B|

        Args:
            user_id: Target user ID
            user_books: Set of books the target user has interacted with

        Returns:
            Dict mapping similar user_id to similarity score
        """
        similar_users = {}

        # For each book the user has interacted with, find other users
        # who also interacted with it
        other_users: Dict[UUID, Set[UUID]] = defaultdict(set)

        for book_id in user_books:
            users = self.interaction_repo.get_users_who_interacted_with_book(book_id)
            for other_user_id in users:
                if other_user_id != user_id:
                    other_users[other_user_id].add(book_id)

        # Calculate Jaccard similarity for each candidate user
        for other_user_id, common_books in other_users.items():
            if len(common_books) >= self.min_common_books:
                # Get all books the other user has interacted with
                other_user_books = set(
                    self.interaction_repo.get_user_book_ids(other_user_id)
                )

                # Jaccard similarity
                intersection = len(user_books & other_user_books)
                union = len(user_books | other_user_books)

                if union > 0:
                    similarity = intersection / union
                    similar_users[other_user_id] = similarity

        return similar_users

    def _get_recommendations_from_similar_users(
        self,
        user_books: Set[UUID],
        similar_users: Dict[UUID, float],
    ) -> List[Dict]:
        """
        Get book recommendations from similar users.

        Args:
            user_books: Books the target user has already interacted with
            similar_users: Dict mapping user_id to similarity score

        Returns:
            List of recommendation dicts
        """
        # Aggregate scores for each book recommended by similar users
        book_scores: Dict[UUID, float] = defaultdict(float)
        book_recommenders: Dict[UUID, int] = defaultdict(int)

        for similar_user_id, similarity in similar_users.items():
            # Get books with weighted scores for this similar user
            weighted_books = self.interaction_repo.get_weighted_book_scores(
                similar_user_id
            )

            for book_id, interaction_score in weighted_books.items():
                if book_id not in user_books:
                    # Weight by similarity and interaction score
                    book_scores[book_id] += similarity * interaction_score
                    book_recommenders[book_id] += 1

        # Create recommendations
        recommendations = []
        for book_id, score in book_scores.items():
            num_recommenders = book_recommenders[book_id]
            reason = f"Recommended by {num_recommenders} similar user{'s' if num_recommenders > 1 else ''}"

            recommendations.append({
                "book_id": book_id,
                "score": score,
                "reason": reason,
            })

        return recommendations
