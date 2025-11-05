"""
Main recommendation service with hybrid approach.
"""

from typing import List, Optional
from uuid import UUID
from sqlalchemy.orm import Session
from collections import defaultdict

from app.repository import InteractionRepository, RecommendationRepository
from app.services.strategies import (
    TagBasedStrategy,
    CollaborativeStrategy,
    PopularStrategy,
)
from app.schemas import (
    RecommendationItem,
    RecommendationResponse,
    InteractionCreate,
    InteractionResponse,
    InteractionStats,
    UserPreferenceCreate,
    UserPreferenceResponse,
)
from app.core.config import get_settings

settings = get_settings()


class RecommendationService:
    """
    Main service for generating personalized recommendations.

    Uses a hybrid approach combining multiple strategies with configurable weights.
    """

    def __init__(self, db: Session):
        self.db = db
        self.interaction_repo = InteractionRepository(db)
        self.recommendation_repo = RecommendationRepository(db)

        # Initialize strategies
        self.tag_strategy = TagBasedStrategy(
            self.interaction_repo, self.recommendation_repo
        )
        self.collaborative_strategy = CollaborativeStrategy(
            self.interaction_repo,
            min_common_books=settings.min_interactions_for_collaborative,
        )
        self.popular_strategy = PopularStrategy(self.interaction_repo)

    def get_recommendations(
        self, user_id: UUID, limit: int = None
    ) -> RecommendationResponse:
        """
        Get personalized recommendations for a user using hybrid approach.

        First checks cache, then generates new recommendations if needed.
        """
        if limit is None:
            limit = settings.default_recommendations_count

        # Check cache
        cached = self.recommendation_repo.get_cached_recommendations(user_id)
        if cached:
            # Convert cached recommendations to response
            items = [
                RecommendationItem(
                    book_id=book_id,
                    score=cached.score,
                    reason=f"Cached recommendations using {cached.algorithm}",
                )
                for book_id in cached.book_ids[:limit]
            ]

            return RecommendationResponse(
                user_id=user_id,
                recommendations=items,
                algorithm=cached.algorithm,
                total=len(items),
            )

        # Generate new recommendations
        recommendations = self._generate_hybrid_recommendations(user_id, limit)

        if recommendations:
            # Cache the results
            book_ids = [rec["book_id"] for rec in recommendations]
            avg_score = sum(rec["score"] for rec in recommendations) / len(
                recommendations
            )

            self.recommendation_repo.save_recommendations(
                user_id=user_id,
                book_ids=book_ids,
                algorithm="hybrid",
                score=avg_score,
                ttl_hours=24,
            )

        # Convert to response format
        items = [
            RecommendationItem(
                book_id=rec["book_id"],
                score=rec["score"],
                reason=rec.get("reason"),
            )
            for rec in recommendations
        ]

        return RecommendationResponse(
            user_id=user_id,
            recommendations=items,
            algorithm="hybrid",
            total=len(items),
        )

    def _generate_hybrid_recommendations(
        self, user_id: UUID, limit: int
    ) -> List[dict]:
        """
        Generate recommendations using a weighted combination of strategies.
        """
        # Get recommendations from each strategy
        tag_recs = self.tag_strategy.recommend(user_id, limit=limit * 2)
        collab_recs = self.collaborative_strategy.recommend(user_id, limit=limit * 2)
        popular_recs = self.popular_strategy.recommend(user_id, limit=limit * 2)

        # If we have no personalized recommendations, fall back to popular
        if not tag_recs and not collab_recs:
            return popular_recs[:limit]

        # Combine recommendations with weights
        combined_scores = defaultdict(lambda: {"score": 0.0, "reasons": []})

        # Add tag-based recommendations
        for rec in tag_recs:
            book_id = rec["book_id"]
            combined_scores[book_id]["score"] += rec["score"] * settings.tag_weight
            combined_scores[book_id]["reasons"].append(rec.get("reason", ""))

        # Add collaborative recommendations
        for rec in collab_recs:
            book_id = rec["book_id"]
            combined_scores[book_id]["score"] += (
                rec["score"] * settings.collaborative_weight
            )
            combined_scores[book_id]["reasons"].append(rec.get("reason", ""))

        # Add popular recommendations (with lower weight)
        for rec in popular_recs:
            book_id = rec["book_id"]
            combined_scores[book_id]["score"] += rec["score"] * settings.popular_weight
            if not combined_scores[book_id]["reasons"]:
                combined_scores[book_id]["reasons"].append(rec.get("reason", ""))

        # Convert to list and sort by score
        recommendations = []
        for book_id, data in combined_scores.items():
            # Combine reasons (take first non-empty reason)
            reason = next(
                (r for r in data["reasons"] if r), "Recommended for you"
            )

            recommendations.append({
                "book_id": book_id,
                "score": data["score"],
                "reason": reason,
            })

        recommendations.sort(key=lambda x: x["score"], reverse=True)

        return recommendations[:limit]

    def get_similar_books(self, book_id: UUID, limit: int = 10) -> List[RecommendationItem]:
        """Get books similar to a given book based on tags."""
        # Get tags for the source book
        tags = self.recommendation_repo.get_book_tags(book_id)

        if not tags:
            return []

        # Find books with similar tags
        similar_books_scores = defaultdict(float)

        for tag in tags:
            similar_books = self.recommendation_repo.get_books_by_tag(
                tag_type=tag.tag_type,
                tag_value=tag.tag_value,
                limit=50,
            )

            for similar_book_id in similar_books:
                if similar_book_id != book_id:
                    similar_books_scores[similar_book_id] += tag.weight

        # Convert to list and sort
        recommendations = [
            {"book_id": book_id, "score": score}
            for book_id, score in similar_books_scores.items()
        ]

        # Normalize scores
        if recommendations:
            max_score = max(rec["score"] for rec in recommendations)
            for rec in recommendations:
                rec["score"] = rec["score"] / max_score if max_score > 0 else 0

        recommendations.sort(key=lambda x: x["score"], reverse=True)

        # Convert to response format
        return [
            RecommendationItem(
                book_id=rec["book_id"],
                score=rec["score"],
                reason="Similar content and tags",
            )
            for rec in recommendations[:limit]
        ]

    def get_trending_books(self, limit: int = 10, days: int = 7) -> List[RecommendationItem]:
        """Get trending books from the last N days."""
        trending = self.popular_strategy.get_trending(limit=limit, days=days)

        return [
            RecommendationItem(
                book_id=rec["book_id"],
                score=rec["score"],
                reason=rec.get("reason", "Trending"),
            )
            for rec in trending
        ]

    # Interaction Methods

    def track_interaction(
        self, user_id: UUID, interaction: InteractionCreate
    ) -> InteractionResponse:
        """Track a user interaction with a book."""
        db_interaction = self.interaction_repo.create(user_id, interaction)

        # Invalidate recommendation cache since user behavior changed
        self.recommendation_repo.invalidate_cache(user_id)

        return InteractionResponse.from_orm(db_interaction)

    def get_user_interactions(self, user_id: UUID, limit: Optional[int] = None) -> List[InteractionResponse]:
        """Get all interactions for a user."""
        interactions = self.interaction_repo.get_by_user(user_id, limit=limit)

        return [InteractionResponse.from_orm(interaction) for interaction in interactions]

    def get_interaction_stats(self, user_id: UUID) -> InteractionStats:
        """Get statistics about user's interactions."""
        stats = self.interaction_repo.get_interaction_stats(user_id)

        return InteractionStats(**stats)

    # User Preference Methods

    def get_user_preferences(self, user_id: UUID) -> Optional[UserPreferenceResponse]:
        """Get user preferences."""
        preferences = self.recommendation_repo.get_user_preferences(user_id)

        return UserPreferenceResponse.from_orm(preferences) if preferences else None

    def update_user_preferences(
        self, user_id: UUID, preferences: UserPreferenceCreate
    ) -> UserPreferenceResponse:
        """Create or update user preferences."""
        db_preferences = self.recommendation_repo.create_or_update_preferences(
            user_id, preferences.dict(exclude_unset=True)
        )

        # Invalidate cache since preferences changed
        self.recommendation_repo.invalidate_cache(user_id)

        return UserPreferenceResponse.from_orm(db_preferences)

    def delete_user_preferences(self, user_id: UUID) -> None:
        """Delete user preferences."""
        self.recommendation_repo.delete_user_preferences(user_id)
        self.recommendation_repo.invalidate_cache(user_id)
