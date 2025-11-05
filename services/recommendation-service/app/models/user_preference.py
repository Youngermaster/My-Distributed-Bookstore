"""
User preference model for explicit user preferences.
"""

from sqlalchemy import Column, String, Float, DateTime, Text
from sqlalchemy.dialects.postgresql import UUID, ARRAY
import uuid
from datetime import datetime

from app.core.database import Base


class UserPreference(Base):
    """
    Explicit user preferences for personalized recommendations.

    This includes user-selected preferences like favorite genres,
    preferred authors, price range, etc.
    """

    __tablename__ = "user_preferences"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), nullable=False, unique=True, index=True)

    # Preferred genres (e.g., ["science-fiction", "mystery", "thriller"])
    preferred_genres = Column(ARRAY(String), nullable=True)

    # Preferred authors (UUIDs)
    preferred_authors = Column(ARRAY(UUID(as_uuid=True)), nullable=True)

    # Price range preferences
    min_price = Column(Float, nullable=True)
    max_price = Column(Float, nullable=True)

    # Preferred languages
    preferred_languages = Column(ARRAY(String), nullable=True, default=["English"])

    # Exclude categories (e.g., ["horror", "romance"])
    excluded_genres = Column(ARRAY(String), nullable=True)

    # Additional preferences as JSON string
    extra_data = Column(Text, nullable=True)

    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)

    def __repr__(self):
        return f"<UserPreference(user={self.user_id}, genres={self.preferred_genres})>"
