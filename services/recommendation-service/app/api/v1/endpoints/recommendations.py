"""
Recommendation endpoints.
"""

from fastapi import APIRouter, Depends, HTTPException, Query, Header
from sqlalchemy.orm import Session
from typing import List, Optional
from uuid import UUID

from app.core.database import get_db
from app.services import RecommendationService
from app.schemas import (
    RecommendationResponse,
    RecommendationItem,
    SimilarBooksResponse,
    InteractionCreate,
    InteractionResponse,
    InteractionStats,
    UserPreferenceCreate,
    UserPreferenceResponse,
)

router = APIRouter()


def get_current_user_id(
    x_user_id: Optional[str] = Header(None, alias="X-User-Id")
) -> UUID:
    """
    Extract user ID from header.

    In production, this would validate JWT and extract user ID.
    For now, we expect the API Gateway to set the X-User-Id header.
    """
    if not x_user_id:
        raise HTTPException(status_code=401, detail="User ID not provided")

    try:
        return UUID(x_user_id)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid user ID format")


# Recommendation Endpoints


@router.get("/me", response_model=RecommendationResponse, tags=["recommendations"])
def get_my_recommendations(
    limit: int = Query(10, ge=1, le=50, description="Number of recommendations"),
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Get personalized recommendations for the current user.

    Uses a hybrid approach combining:
    - Tag-based (content-based filtering)
    - Collaborative filtering
    - Popular books
    """
    service = RecommendationService(db)
    return service.get_recommendations(user_id, limit=limit)


@router.get("/similar/{book_id}", response_model=SimilarBooksResponse, tags=["recommendations"])
def get_similar_books(
    book_id: UUID,
    limit: int = Query(10, ge=1, le=50, description="Number of similar books"),
    db: Session = Depends(get_db),
):
    """
    Get books similar to a given book based on tags and content.
    """
    service = RecommendationService(db)
    similar_books = service.get_similar_books(book_id, limit=limit)

    return SimilarBooksResponse(
        book_id=book_id,
        similar_books=similar_books,
        total=len(similar_books),
    )


@router.get("/trending", response_model=List[RecommendationItem], tags=["recommendations"])
def get_trending_books(
    limit: int = Query(10, ge=1, le=50, description="Number of trending books"),
    days: int = Query(7, ge=1, le=30, description="Number of days to consider"),
    db: Session = Depends(get_db),
):
    """
    Get trending books from the last N days.
    """
    service = RecommendationService(db)
    return service.get_trending_books(limit=limit, days=days)


@router.get("/popular", response_model=List[RecommendationItem], tags=["recommendations"])
def get_popular_books(
    limit: int = Query(10, ge=1, le=50, description="Number of popular books"),
    db: Session = Depends(get_db),
):
    """
    Get popular books based on overall interactions (last 30 days).
    """
    service = RecommendationService(db)
    return service.get_trending_books(limit=limit, days=30)


# Interaction Endpoints


@router.post("/interactions", response_model=InteractionResponse, status_code=201, tags=["interactions"])
def track_interaction(
    interaction: InteractionCreate,
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Track a user interaction with a book.

    Interaction types:
    - view: User viewed book details
    - add_to_cart: User added book to cart
    - purchase: User purchased book
    - review: User reviewed book
    - wishlist: User added book to wishlist
    """
    service = RecommendationService(db)
    return service.track_interaction(user_id, interaction)


@router.get("/interactions/me", response_model=List[InteractionResponse], tags=["interactions"])
def get_my_interactions(
    limit: Optional[int] = Query(None, ge=1, le=100, description="Limit results"),
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Get all interactions for the current user.
    """
    service = RecommendationService(db)
    return service.get_user_interactions(user_id, limit=limit)


@router.get("/interactions/me/stats", response_model=InteractionStats, tags=["interactions"])
def get_my_interaction_stats(
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Get statistics about the current user's interactions.
    """
    service = RecommendationService(db)
    return service.get_interaction_stats(user_id)


# User Preference Endpoints


@router.get("/preferences/me", response_model=Optional[UserPreferenceResponse], tags=["preferences"])
def get_my_preferences(
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Get the current user's preferences.
    """
    service = RecommendationService(db)
    return service.get_user_preferences(user_id)


@router.put("/preferences/me", response_model=UserPreferenceResponse, tags=["preferences"])
def update_my_preferences(
    preferences: UserPreferenceCreate,
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Create or update the current user's preferences.
    """
    service = RecommendationService(db)
    return service.update_user_preferences(user_id, preferences)


@router.delete("/preferences/me", status_code=204, tags=["preferences"])
def delete_my_preferences(
    user_id: UUID = Depends(get_current_user_id),
    db: Session = Depends(get_db),
):
    """
    Delete the current user's preferences.
    """
    service = RecommendationService(db)
    service.delete_user_preferences(user_id)
    return None
