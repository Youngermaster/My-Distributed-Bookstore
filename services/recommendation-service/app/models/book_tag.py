"""
Book tag model for content-based recommendations.
"""

from sqlalchemy import Column, String, Float, DateTime, Index
from sqlalchemy.dialects.postgresql import UUID
import uuid
from datetime import datetime

from app.core.database import Base


class BookTag(Base):
    """
    Tags associated with books for content-based filtering.

    Tags can include:
    - Genres (e.g., "science-fiction", "mystery")
    - Topics (e.g., "artificial-intelligence", "historical")
    - Authors (e.g., "author:isaac-asimov")
    - Publishers (e.g., "publisher:oreilly")
    - Attributes (e.g., "bestseller", "award-winner")
    """

    __tablename__ = "book_tags"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    book_id = Column(UUID(as_uuid=True), nullable=False, index=True)

    # Tag type: genre, topic, author, publisher, attribute
    tag_type = Column(String(50), nullable=False)

    # Tag value (e.g., "science-fiction", "author:isaac-asimov")
    tag_value = Column(String(255), nullable=False)

    # Relevance weight (0.0 to 1.0)
    weight = Column(Float, nullable=False, default=1.0)

    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)

    # Indexes
    __table_args__ = (
        Index("idx_book_tag", "book_id", "tag_type", "tag_value", unique=True),
        Index("idx_tag_value", "tag_value"),
        Index("idx_tag_type", "tag_type"),
    )

    def __repr__(self):
        return f"<BookTag(book={self.book_id}, type={self.tag_type}, value={self.tag_value})>"
