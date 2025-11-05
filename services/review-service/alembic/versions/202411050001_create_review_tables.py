"""create reviews and review_votes tables

Revision ID: 202411050001
Revises: 
Create Date: 2024-11-05 00:01:00
"""

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql


# revision identifiers, used by Alembic.
revision = "202411050001"
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.create_table(
        "reviews",
        sa.Column("id", postgresql.UUID(as_uuid=True), primary_key=True, nullable=False),
        sa.Column("book_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("rating", sa.Integer(), nullable=False),
        sa.Column("title", sa.String(length=255), nullable=False),
        sa.Column("content", sa.Text(), nullable=False),
        sa.Column("sentiment_score", sa.Float(), nullable=True),
        sa.Column("sentiment_label", sa.String(length=20), nullable=True),
        sa.Column("verified_purchase", sa.Boolean(), nullable=False, server_default=sa.text("false")),
        sa.Column("helpful_votes", sa.Integer(), nullable=False, server_default=sa.text("0")),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False, server_default=sa.text("now()")),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False, server_default=sa.text("now()")),
        sa.CheckConstraint("rating >= 1 AND rating <= 5", name="rating_range"),
        sa.CheckConstraint("helpful_votes >= 0", name="helpful_votes_non_negative"),
        sa.UniqueConstraint("book_id", "user_id", name="unique_review_per_user_per_book"),
        comment="Book reviews and ratings table",
    )

    op.create_index("ix_reviews_book_id", "reviews", ["book_id"])
    op.create_index("ix_reviews_user_id", "reviews", ["user_id"])
    op.create_index("ix_reviews_created_at", "reviews", ["created_at"])

    op.create_table(
        "review_votes",
        sa.Column("review_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("user_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("is_helpful", sa.Boolean(), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False, server_default=sa.text("now()")),
        sa.PrimaryKeyConstraint("review_id", "user_id", name="pk_review_votes"),
        sa.ForeignKeyConstraint(["review_id"], ["reviews.id"], ondelete="CASCADE"),
        comment="Votes on review helpfulness",
    )

    op.create_index("ix_review_votes_review_id", "review_votes", ["review_id"])
    op.create_index("ix_review_votes_user_id", "review_votes", ["user_id"])


def downgrade() -> None:
    op.drop_index("ix_review_votes_user_id", table_name="review_votes")
    op.drop_index("ix_review_votes_review_id", table_name="review_votes")
    op.drop_table("review_votes")

    op.drop_index("ix_reviews_created_at", table_name="reviews")
    op.drop_index("ix_reviews_user_id", table_name="reviews")
    op.drop_index("ix_reviews_book_id", table_name="reviews")
    op.drop_table("reviews")
