"""
Review API endpoints.
"""
from typing import Annotated
from uuid import UUID
import logging

from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.ext.asyncio import AsyncSession

from app.db.base import get_db
from app.services.review_service import ReviewService
from app.schemas.review import (
    ReviewCreateRequest,
    ReviewUpdateRequest,
    ReviewResponse,
    ReviewListResponse,
    ReviewStatsResponse,
    ReviewVoteRequest,
    MessageResponse
)
from app.core.config import settings

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/reviews", tags=["reviews"])


@router.post(
    "",
    response_model=ReviewResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new review",
    description="Create a new review for a book. Automatically analyzes sentiment."
)
async def create_review(
    review_data: ReviewCreateRequest,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Create a new review with automatic sentiment analysis."""
    try:
        review = await ReviewService.create_review(db, review_data)
        return review
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    except Exception as e:
        logger.error(f"Error creating review: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create review"
        )


@router.get(
    "/{review_id}",
    response_model=ReviewResponse,
    summary="Get a review by ID",
    description="Retrieve a specific review by its ID."
)
async def get_review(
    review_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Get a specific review."""
    review = await ReviewService.get_review(db, review_id)
    if not review:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Review not found"
        )
    return review


@router.get(
    "/book/{book_id}",
    response_model=ReviewListResponse,
    summary="Get all reviews for a book",
    description="Retrieve all reviews for a specific book with pagination."
)
async def get_book_reviews(
    book_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)],
    page: int = Query(1, ge=1, description="Page number"),
    page_size: int = Query(
        default=settings.DEFAULT_PAGE_SIZE,
        ge=1,
        le=settings.MAX_PAGE_SIZE,
        description="Number of items per page"
    )
):
    """Get all reviews for a book with pagination."""
    skip = (page - 1) * page_size

    reviews, total = await ReviewService.get_reviews_by_book(
        db, book_id, skip=skip, limit=page_size
    )

    total_pages = (total + page_size - 1) // page_size

    return ReviewListResponse(
        reviews=reviews,
        total=total,
        page=page,
        page_size=page_size,
        total_pages=total_pages
    )


@router.get(
    "/book/{book_id}/stats",
    response_model=ReviewStatsResponse,
    summary="Get review statistics for a book",
    description="Get aggregated statistics including average rating and sentiment distribution."
)
async def get_book_stats(
    book_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Get review statistics for a book."""
    stats = await ReviewService.get_book_stats(db, book_id)
    return stats


@router.put(
    "/{review_id}",
    response_model=ReviewResponse,
    summary="Update a review",
    description="Update an existing review. Sentiment is re-analyzed if content changes."
)
async def update_review(
    review_id: UUID,
    review_data: ReviewUpdateRequest,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Update a review."""
    review = await ReviewService.update_review(db, review_id, review_data)
    if not review:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Review not found"
        )
    return review


@router.delete(
    "/{review_id}",
    response_model=MessageResponse,
    summary="Delete a review",
    description="Delete a review by its ID."
)
async def delete_review(
    review_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Delete a review."""
    success = await ReviewService.delete_review(db, review_id)
    if not success:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Review not found"
        )
    return MessageResponse(message="Review deleted successfully")


@router.post(
    "/{review_id}/vote",
    response_model=MessageResponse,
    summary="Vote on review helpfulness",
    description="Mark a review as helpful or not helpful."
)
async def vote_on_review(
    review_id: UUID,
    vote_data: ReviewVoteRequest,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Vote on review helpfulness."""
    success = await ReviewService.vote_on_review(
        db, review_id, vote_data.user_id, vote_data.is_helpful
    )

    if not success:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Review not found"
        )

    return MessageResponse(
        message="Vote recorded successfully",
        detail=f"Marked as {'helpful' if vote_data.is_helpful else 'not helpful'}"
    )
