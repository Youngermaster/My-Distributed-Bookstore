"""
Recommendation cache model for storing pre-computed recommendations.
"""

from sqlalchemy import Column, String, Float, DateTime, Text
from sqlalchemy.dialects.postgresql import UUID, ARRAY
import uuid
from datetime import datetime

from app.core.database import Base


class RecommendationCache(Base):
    """
    Cached recommendations for users to improve performance.

    This allows pre-computing recommendations during off-peak hours
    or when user interactions occur, reducing real-time computation.
    """

    __tablename__ = "recommendation_cache"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), nullable=False, unique=True, index=True)

    # Array of recommended book IDs (ordered by relevance)
    book_ids = Column(ARRAY(UUID(as_uuid=True)), nullable=False)

    # Algorithm used: tag_based, collaborative, popular, hybrid
    algorithm = Column(String(50), nullable=False)

    # Overall recommendation score/confidence
    score = Column(Float, nullable=False, default=0.0)

    # Optional metadata about the recommendations
    extra_data = Column(Text, nullable=True)

    # Cache expiration
    expires_at = Column(DateTime, nullable=False)

    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)

    def __repr__(self):
        return f"<RecommendationCache(user={self.user_id}, algorithm={self.algorithm}, count={len(self.book_ids) if self.book_ids else 0})>"
