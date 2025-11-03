"""
Pydantic schemas for user interactions.
"""

from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime
from uuid import UUID


class InteractionBase(BaseModel):
    """Base schema for user interactions."""

    book_id: UUID = Field(..., description="ID of the book")
    interaction_type: str = Field(..., description="Type of interaction: view, add_to_cart, purchase, review, wishlist")
    metadata: Optional[str] = Field(None, description="Optional metadata as JSON string")


class InteractionCreate(InteractionBase):
    """Schema for creating a new interaction."""

    pass


class InteractionResponse(InteractionBase):
    """Schema for interaction responses."""

    id: UUID
    user_id: UUID
    score: float
    created_at: datetime

    class Config:
        from_attributes = True


class InteractionStats(BaseModel):
    """Statistics about user interactions."""

    total_interactions: int
    views: int
    purchases: int
    reviews: int
    wishlists: int
    cart_additions: int
