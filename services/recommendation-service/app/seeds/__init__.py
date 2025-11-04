"""Database seed helpers for the recommendation service."""

import json
import logging
from datetime import datetime, timedelta
from typing import Iterable
from uuid import UUID, uuid5, NAMESPACE_URL

from sqlalchemy import func
from sqlalchemy.orm import Session

from app.models import BookTag, UserInteraction
from .data import BOOK_TAG_DATA, USER_INTERACTIONS_DATA

logger = logging.getLogger(__name__)


def book_uuid(isbn: str) -> UUID:
    """Generate the deterministic UUID used by the catalog service for a book."""

    return uuid5(NAMESPACE_URL, isbn)


def user_uuid(user_key: str) -> UUID:
    """Create a stable UUID for synthetic users."""

    return uuid5(NAMESPACE_URL, f"user:{user_key}")


def _ensure_tags(session: Session, data: Iterable[dict]) -> None:
    for entry in data:
        session.add(
            BookTag(
                book_id=book_uuid(entry["isbn"]),
                tag_type=entry["tag_type"],
                tag_value=entry["tag_value"],
                weight=float(entry.get("weight", 1.0)),
            )
        )


def _ensure_interactions(session: Session, data: Iterable[dict]) -> None:
    now = datetime.utcnow()

    for entry in data:
        extra = entry.get("extra_data")

        session.add(
            UserInteraction(
                user_id=user_uuid(entry["user"]),
                book_id=book_uuid(entry["isbn"]),
                interaction_type=entry["interaction_type"],
                score=float(entry.get("score", 1.0)),
                extra_data=json.dumps(extra) if extra else None,
                created_at=now - timedelta(days=int(entry.get("days_ago", 0))),
            )
        )


def seed_database(session: Session) -> None:
    """Populate the recommendation database with sample data if empty."""

    tag_count = session.query(func.count(BookTag.id)).scalar()
    interaction_count = session.query(func.count(UserInteraction.id)).scalar()

    if tag_count and tag_count > 0 and interaction_count and interaction_count > 0:
        logger.info("Recommendation seed data already present; skipping")
        return

    logger.info("Seeding recommendation service with book tags and user interactions")

    if not tag_count:
        _ensure_tags(session, BOOK_TAG_DATA)

    if not interaction_count:
        _ensure_interactions(session, USER_INTERACTIONS_DATA)

    session.commit()
    logger.info("Recommendation seed data inserted successfully")

