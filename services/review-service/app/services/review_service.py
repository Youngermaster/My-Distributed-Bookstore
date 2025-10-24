"""
Business logic for review operations.
"""
from typing import Optional
from uuid import UUID
import logging

from sqlalchemy import select, func, case
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.review import Review, ReviewVote
from app.schemas.review import (
    ReviewCreateRequest,
    ReviewUpdateRequest,
    ReviewStatsResponse
)
from app.ml.sentiment import analyze_sentiment

logger = logging.getLogger(__name__)


class ReviewService:
    """Service class for review-related operations."""

    @staticmethod
    async def create_review(
        db: AsyncSession,
        review_data: ReviewCreateRequest
    ) -> Review:
        """
        Create a new review with sentiment analysis.

        Args:
            db: Database session
            review_data: Review creation data

        Returns:
            Created review object

        Raises:
            ValueError: If review already exists for this user and book
        """
        # Check if review already exists
        existing_query = select(Review).where(
            Review.book_id == review_data.book_id,
            Review.user_id == review_data.user_id
        )
        result = await db.execute(existing_query)
        existing_review = result.scalar_one_or_none()

        if existing_review:
            raise ValueError("Review already exists for this user and book")

        # Analyze sentiment
        sentiment_label, sentiment_score = analyze_sentiment(review_data.content)

        # Create review
        review = Review(
            book_id=review_data.book_id,
            user_id=review_data.user_id,
            rating=review_data.rating,
            title=review_data.title,
            content=review_data.content,
            verified_purchase=review_data.verified_purchase,
            sentiment_score=sentiment_score,
            sentiment_label=sentiment_label,
        )

        db.add(review)
        await db.commit()
        await db.refresh(review)

        logger.info(f"Created review {review.id} for book {review.book_id}")
        return review

    @staticmethod
    async def get_review(db: AsyncSession, review_id: UUID) -> Optional[Review]:
        """Get a review by ID."""
        query = select(Review).where(Review.id == review_id)
        result = await db.execute(query)
        return result.scalar_one_or_none()

    @staticmethod
    async def get_reviews_by_book(
        db: AsyncSession,
        book_id: UUID,
        skip: int = 0,
        limit: int = 20
    ) -> tuple[list[Review], int]:
        """
        Get all reviews for a book with pagination.

        Returns:
            Tuple of (reviews list, total count)
        """
        # Get total count
        count_query = select(func.count()).select_from(Review).where(
            Review.book_id == book_id
        )
        total_result = await db.execute(count_query)
        total = total_result.scalar_one()

        # Get reviews
        query = (
            select(Review)
            .where(Review.book_id == book_id)
            .order_by(Review.created_at.desc())
            .offset(skip)
            .limit(limit)
        )
        result = await db.execute(query)
        reviews = list(result.scalars().all())

        return reviews, total

    @staticmethod
    async def update_review(
        db: AsyncSession,
        review_id: UUID,
        review_data: ReviewUpdateRequest
    ) -> Optional[Review]:
        """Update an existing review."""
        review = await ReviewService.get_review(db, review_id)
        if not review:
            return None

        # Update fields if provided
        if review_data.rating is not None:
            review.rating = review_data.rating

        if review_data.title is not None:
            review.title = review_data.title

        if review_data.content is not None:
            review.content = review_data.content
            # Reanalyze sentiment if content changed
            sentiment_label, sentiment_score = analyze_sentiment(review_data.content)
            review.sentiment_score = sentiment_score
            review.sentiment_label = sentiment_label

        await db.commit()
        await db.refresh(review)

        logger.info(f"Updated review {review_id}")
        return review

    @staticmethod
    async def delete_review(db: AsyncSession, review_id: UUID) -> bool:
        """Delete a review."""
        review = await ReviewService.get_review(db, review_id)
        if not review:
            return False

        await db.delete(review)
        await db.commit()

        logger.info(f"Deleted review {review_id}")
        return True

    @staticmethod
    async def vote_on_review(
        db: AsyncSession,
        review_id: UUID,
        user_id: UUID,
        is_helpful: bool
    ) -> bool:
        """
        Vote on review helpfulness.

        Returns:
            True if vote was recorded successfully
        """
        # Check if review exists
        review = await ReviewService.get_review(db, review_id)
        if not review:
            return False

        # Check if vote already exists
        vote_query = select(ReviewVote).where(
            ReviewVote.review_id == review_id,
            ReviewVote.user_id == user_id
        )
        result = await db.execute(vote_query)
        existing_vote = result.scalar_one_or_none()

        if existing_vote:
            # Update existing vote
            old_helpful = existing_vote.is_helpful
            existing_vote.is_helpful = is_helpful

            # Update helpful votes count
            if old_helpful != is_helpful:
                if is_helpful:
                    review.helpful_votes += 1
                else:
                    review.helpful_votes -= 1
        else:
            # Create new vote
            vote = ReviewVote(
                review_id=review_id,
                user_id=user_id,
                is_helpful=is_helpful
            )
            db.add(vote)

            # Update helpful votes count
            if is_helpful:
                review.helpful_votes += 1

        await db.commit()
        return True

    @staticmethod
    async def get_book_stats(
        db: AsyncSession,
        book_id: UUID
    ) -> ReviewStatsResponse:
        """Get statistics for all reviews of a book."""
        # Get rating distribution and average
        rating_query = select(
            Review.rating,
            func.count(Review.id).label('count')
        ).where(
            Review.book_id == book_id
        ).group_by(
            Review.rating
        )

        rating_result = await db.execute(rating_query)
        rating_rows = rating_result.all()

        rating_distribution = {row.rating: row.count for row in rating_rows}
        total_reviews = sum(rating_distribution.values())

        if total_reviews > 0:
            average_rating = sum(
                rating * count for rating, count in rating_distribution.items()
            ) / total_reviews
        else:
            average_rating = 0.0

        # Get sentiment distribution
        sentiment_query = select(
            Review.sentiment_label,
            func.count(Review.id).label('count')
        ).where(
            Review.book_id == book_id,
            Review.sentiment_label.isnot(None)
        ).group_by(
            Review.sentiment_label
        )

        sentiment_result = await db.execute(sentiment_query)
        sentiment_rows = sentiment_result.all()

        sentiment_distribution = {
            row.sentiment_label: row.count
            for row in sentiment_rows
            if row.sentiment_label
        }

        return ReviewStatsResponse(
            book_id=book_id,
            total_reviews=total_reviews,
            average_rating=round(average_rating, 2),
            rating_distribution=rating_distribution,
            sentiment_distribution=sentiment_distribution
        )
