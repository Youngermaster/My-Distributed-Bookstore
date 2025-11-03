"""
Repository for recommendation cache and book tags.
"""

from sqlalchemy.orm import Session
from sqlalchemy import and_, func
from typing import List, Optional, Dict
from uuid import UUID
from datetime import datetime, timedelta

from app.models import RecommendationCache, BookTag, UserPreference


class RecommendationRepository:
    """Handles database operations for recommendations."""

    def __init__(self, db: Session):
        self.db = db

    # Recommendation Cache Methods

    def get_cached_recommendations(self, user_id: UUID) -> Optional[RecommendationCache]:
        """Get cached recommendations for a user if not expired."""
        cache = (
            self.db.query(RecommendationCache)
            .filter(
                and_(
                    RecommendationCache.user_id == user_id,
                    RecommendationCache.expires_at > datetime.utcnow(),
                )
            )
            .first()
        )

        return cache

    def save_recommendations(
        self,
        user_id: UUID,
        book_ids: List[UUID],
        algorithm: str,
        score: float,
        ttl_hours: int = 24,
    ) -> RecommendationCache:
        """Save or update cached recommendations for a user."""
        # Delete existing cache for this user
        self.db.query(RecommendationCache).filter(
            RecommendationCache.user_id == user_id
        ).delete()

        # Create new cache entry
        cache = RecommendationCache(
            user_id=user_id,
            book_ids=book_ids,
            algorithm=algorithm,
            score=score,
            expires_at=datetime.utcnow() + timedelta(hours=ttl_hours),
        )

        self.db.add(cache)
        self.db.commit()
        self.db.refresh(cache)

        return cache

    def invalidate_cache(self, user_id: UUID) -> None:
        """Invalidate cached recommendations for a user."""
        self.db.query(RecommendationCache).filter(
            RecommendationCache.user_id == user_id
        ).delete()
        self.db.commit()

    # Book Tag Methods

    def get_book_tags(self, book_id: UUID) -> List[BookTag]:
        """Get all tags for a book."""
        return self.db.query(BookTag).filter(BookTag.book_id == book_id).all()

    def get_books_by_tag(
        self, tag_type: str, tag_value: str, limit: int = 50
    ) -> List[UUID]:
        """Get books that have a specific tag."""
        results = (
            self.db.query(BookTag.book_id)
            .filter(
                and_(
                    BookTag.tag_type == tag_type,
                    BookTag.tag_value == tag_value,
                )
            )
            .order_by(BookTag.weight.desc())
            .limit(limit)
            .all()
        )

        return [row.book_id for row in results]

    def create_book_tag(
        self, book_id: UUID, tag_type: str, tag_value: str, weight: float = 1.0
    ) -> BookTag:
        """Create a new book tag."""
        tag = BookTag(
            book_id=book_id,
            tag_type=tag_type,
            tag_value=tag_value,
            weight=weight,
        )

        self.db.add(tag)
        self.db.commit()
        self.db.refresh(tag)

        return tag

    def get_tag_distribution(self, book_ids: List[UUID]) -> Dict[str, int]:
        """
        Get tag distribution for a list of books.

        Returns a dictionary mapping tag_value to count.
        """
        results = (
            self.db.query(
                BookTag.tag_value,
                func.count(BookTag.id).label("count"),
            )
            .filter(BookTag.book_id.in_(book_ids))
            .group_by(BookTag.tag_value)
            .all()
        )

        return {row.tag_value: row.count for row in results}

    # User Preference Methods

    def get_user_preferences(self, user_id: UUID) -> Optional[UserPreference]:
        """Get user preferences."""
        return (
            self.db.query(UserPreference)
            .filter(UserPreference.user_id == user_id)
            .first()
        )

    def create_or_update_preferences(
        self, user_id: UUID, preferences_data: dict
    ) -> UserPreference:
        """Create or update user preferences."""
        existing = self.get_user_preferences(user_id)

        if existing:
            # Update existing preferences
            for key, value in preferences_data.items():
                if hasattr(existing, key):
                    setattr(existing, key, value)
            existing.updated_at = datetime.utcnow()
            self.db.commit()
            self.db.refresh(existing)
            return existing
        else:
            # Create new preferences
            preferences = UserPreference(user_id=user_id, **preferences_data)
            self.db.add(preferences)
            self.db.commit()
            self.db.refresh(preferences)
            return preferences

    def delete_user_preferences(self, user_id: UUID) -> None:
        """Delete user preferences."""
        self.db.query(UserPreference).filter(
            UserPreference.user_id == user_id
        ).delete()
        self.db.commit()
