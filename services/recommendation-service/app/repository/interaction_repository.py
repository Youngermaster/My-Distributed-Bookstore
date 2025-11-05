"""
Repository for user interaction data access.
"""

from sqlalchemy.orm import Session
from sqlalchemy import func, and_
from typing import List, Optional, Dict
from uuid import UUID
from datetime import datetime, timedelta

from app.models import UserInteraction
from app.schemas import InteractionCreate


class InteractionRepository:
    """Handles database operations for user interactions."""

    def __init__(self, db: Session):
        self.db = db

    def create(self, user_id: UUID, interaction: InteractionCreate) -> UserInteraction:
        """Create a new user interaction."""
        # Define score mapping
        score_map = {
            "view": 1.0,
            "wishlist": 2.0,
            "add_to_cart": 3.0,
            "review": 4.0,
            "purchase": 5.0,
        }

        db_interaction = UserInteraction(
            user_id=user_id,
            book_id=interaction.book_id,
            interaction_type=interaction.interaction_type,
            score=score_map.get(interaction.interaction_type, 1.0),
            metadata=interaction.metadata,
        )

        self.db.add(db_interaction)
        self.db.commit()
        self.db.refresh(db_interaction)

        return db_interaction

    def get_by_user(
        self,
        user_id: UUID,
        limit: Optional[int] = None,
        interaction_type: Optional[str] = None,
    ) -> List[UserInteraction]:
        """Get all interactions for a user."""
        query = self.db.query(UserInteraction).filter(UserInteraction.user_id == user_id)

        if interaction_type:
            query = query.filter(UserInteraction.interaction_type == interaction_type)

        query = query.order_by(UserInteraction.created_at.desc())

        if limit:
            query = query.limit(limit)

        return query.all()

    def get_user_book_ids(self, user_id: UUID) -> List[UUID]:
        """Get all book IDs that a user has interacted with."""
        results = (
            self.db.query(UserInteraction.book_id)
            .filter(UserInteraction.user_id == user_id)
            .distinct()
            .all()
        )

        return [row.book_id for row in results]

    def get_interaction_stats(self, user_id: UUID) -> Dict[str, int]:
        """Get statistics about user's interactions."""
        stats = (
            self.db.query(
                UserInteraction.interaction_type,
                func.count(UserInteraction.id).label("count"),
            )
            .filter(UserInteraction.user_id == user_id)
            .group_by(UserInteraction.interaction_type)
            .all()
        )

        result = {
            "total_interactions": 0,
            "views": 0,
            "purchases": 0,
            "reviews": 0,
            "wishlists": 0,
            "cart_additions": 0,
        }

        for interaction_type, count in stats:
            result["total_interactions"] += count
            if interaction_type == "view":
                result["views"] = count
            elif interaction_type == "purchase":
                result["purchases"] = count
            elif interaction_type == "review":
                result["reviews"] = count
            elif interaction_type == "wishlist":
                result["wishlists"] = count
            elif interaction_type == "add_to_cart":
                result["cart_additions"] = count

        return result

    def get_weighted_book_scores(self, user_id: UUID) -> Dict[UUID, float]:
        """
        Get weighted scores for all books the user has interacted with.

        Returns a dictionary mapping book_id to total weighted score.
        """
        results = (
            self.db.query(
                UserInteraction.book_id,
                func.sum(UserInteraction.score).label("total_score"),
            )
            .filter(UserInteraction.user_id == user_id)
            .group_by(UserInteraction.book_id)
            .all()
        )

        return {row.book_id: float(row.total_score) for row in results}

    def get_users_who_interacted_with_book(self, book_id: UUID) -> List[UUID]:
        """Get all users who have interacted with a specific book."""
        results = (
            self.db.query(UserInteraction.user_id)
            .filter(UserInteraction.book_id == book_id)
            .distinct()
            .all()
        )

        return [row.user_id for row in results]

    def get_recent_interactions(self, days: int = 30) -> List[UserInteraction]:
        """Get all interactions from the last N days."""
        cutoff_date = datetime.utcnow() - timedelta(days=days)

        return (
            self.db.query(UserInteraction)
            .filter(UserInteraction.created_at >= cutoff_date)
            .order_by(UserInteraction.created_at.desc())
            .all()
        )

    def get_popular_books(self, limit: int = 10, days: Optional[int] = None) -> List[tuple]:
        """
        Get most popular books based on weighted interactions.

        Returns list of tuples: (book_id, total_score, interaction_count)
        """
        query = self.db.query(
            UserInteraction.book_id,
            func.sum(UserInteraction.score).label("total_score"),
            func.count(UserInteraction.id).label("interaction_count"),
        )

        if days:
            cutoff_date = datetime.utcnow() - timedelta(days=days)
            query = query.filter(UserInteraction.created_at >= cutoff_date)

        results = (
            query.group_by(UserInteraction.book_id)
            .order_by(func.sum(UserInteraction.score).desc())
            .limit(limit)
            .all()
        )

        return [(row.book_id, float(row.total_score), row.interaction_count) for row in results]
