"""
User interaction model for tracking user behavior with books.
"""

from sqlalchemy import Column, String, Float, DateTime, Index
from sqlalchemy.dialects.postgresql import UUID
import uuid
from datetime import datetime

from app.core.database import Base


class UserInteraction(Base):
    """
    Tracks user interactions with books for recommendation purposes.

    Interaction types:
    - view: User viewed book details
    - add_to_cart: User added book to cart
    - purchase: User purchased book
    - review: User reviewed book
    - wishlist: User added book to wishlist
    """

    __tablename__ = "user_interactions"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    book_id = Column(UUID(as_uuid=True), nullable=False, index=True)

    # Interaction type: view, add_to_cart, purchase, review, wishlist
    interaction_type = Column(String(50), nullable=False)

    # Weight/score for the interaction (purchase=5, review=4, add_to_cart=3, wishlist=2, view=1)
    score = Column(Float, nullable=False, default=1.0)

    # Metadata (optional JSON field for additional data)
    metadata = Column(String, nullable=True)

    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)

    # Indexes
    __table_args__ = (
        Index("idx_user_book", "user_id", "book_id"),
        Index("idx_interaction_type", "interaction_type"),
        Index("idx_created_at", "created_at"),
    )

    def __repr__(self):
        return f"<UserInteraction(user={self.user_id}, book={self.book_id}, type={self.interaction_type})>"
