"""
Review database models using SQLAlchemy 2.0.
"""
from datetime import datetime
from typing import Optional
from uuid import uuid4

from sqlalchemy import String, Integer, Float, Boolean, Text, CheckConstraint, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.dialects.postgresql import UUID

from app.db.base import Base


class Review(Base):
    """Review model for storing book reviews and ratings."""

    __tablename__ = "reviews"

    # Primary key
    id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
        index=True
    )

    # Foreign keys
    book_id: Mapped[UUID] = mapped_column(UUID(as_uuid=True), nullable=False, index=True)
    user_id: Mapped[UUID] = mapped_column(UUID(as_uuid=True), nullable=False, index=True)

    # Review content
    rating: Mapped[int] = mapped_column(Integer, nullable=False)
    title: Mapped[str] = mapped_column(String(255), nullable=False)
    content: Mapped[str] = mapped_column(Text, nullable=False)

    # Sentiment analysis (ML generated)
    sentiment_score: Mapped[Optional[float]] = mapped_column(Float, nullable=True)
    sentiment_label: Mapped[Optional[str]] = mapped_column(String(20), nullable=True)

    # Metadata
    verified_purchase: Mapped[bool] = mapped_column(Boolean, default=False)
    helpful_votes: Mapped[int] = mapped_column(Integer, default=0)

    # Timestamps
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        nullable=False
    )

    # Constraints
    __table_args__ = (
        CheckConstraint('rating >= 1 AND rating <= 5', name='rating_range'),
        CheckConstraint('helpful_votes >= 0', name='helpful_votes_non_negative'),
        UniqueConstraint('book_id', 'user_id', name='unique_review_per_user_per_book'),
        {'comment': 'Book reviews and ratings table'}
    )

    def __repr__(self) -> str:
        return f"<Review(id={self.id}, book_id={self.book_id}, rating={self.rating})>"


class ReviewVote(Base):
    """Model for tracking helpful/not helpful votes on reviews."""

    __tablename__ = "review_votes"

    # Composite primary key
    review_id: Mapped[UUID] = mapped_column(UUID(as_uuid=True), primary_key=True)
    user_id: Mapped[UUID] = mapped_column(UUID(as_uuid=True), primary_key=True)

    # Vote
    is_helpful: Mapped[bool] = mapped_column(Boolean, nullable=False)

    # Timestamp
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow, nullable=False)

    __table_args__ = (
        {'comment': 'Votes on review helpfulness'}
    )

    def __repr__(self) -> str:
        return f"<ReviewVote(review_id={self.review_id}, user_id={self.user_id}, helpful={self.is_helpful})>"
