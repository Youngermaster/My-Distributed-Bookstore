"""create inventory tables

Revision ID: 202411050002
Revises: 
Create Date: 2024-11-05 00:02:00
"""

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql


revision = "202411050002"
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.create_table(
        "inventory",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, primary_key=True),
        sa.Column("book_id", postgresql.UUID(as_uuid=True), nullable=False, unique=True),
        sa.Column("title", sa.String(length=200), nullable=False),
        sa.Column("short_description", sa.String(length=500), nullable=False),
        sa.Column("available_quantity", sa.Integer(), nullable=False, server_default=sa.text("0")),
        sa.Column("reserved_quantity", sa.Integer(), nullable=False, server_default=sa.text("0")),
        sa.Column("reorder_level", sa.Integer(), nullable=False, server_default=sa.text("10")),
        sa.Column("last_restocked_at", sa.DateTime(timezone=True), nullable=True),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=False, server_default=sa.text("now()")),
        sa.CheckConstraint("available_quantity >= 0", name="available_quantity_non_negative"),
        sa.CheckConstraint("reserved_quantity >= 0", name="reserved_quantity_non_negative"),
        sa.CheckConstraint("reorder_level >= 0", name="reorder_level_non_negative"),
        comment="Real-time stock tracking per book",
    )

    op.create_index("ix_inventory_book_id", "inventory", ["book_id"])
    op.create_index("ix_inventory_title", "inventory", ["title"])
    op.create_index("ix_inventory_low_stock", "inventory", ["available_quantity", "reorder_level"])

    op.create_table(
        "reservations",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("book_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("order_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("quantity", sa.Integer(), nullable=False),
        sa.Column("status", sa.String(length=20), nullable=False, server_default=sa.text("'pending'")),
        sa.Column("expires_at", sa.DateTime(timezone=True), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False, server_default=sa.text("now()")),
        sa.CheckConstraint("quantity > 0", name="reservation_quantity_positive"),
        sa.CheckConstraint("status IN ('pending', 'committed', 'released', 'expired')", name="valid_reservation_status"),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("order_id", "book_id", name="uq_reservations_order_book"),
        comment="Temporary stock holds for orders",
    )

    op.create_index("ix_reservations_order_id", "reservations", ["order_id"])
    op.create_index("ix_reservations_book_id", "reservations", ["book_id"])
    op.create_index("ix_reservations_status", "reservations", ["status"])
    op.create_index("ix_reservations_expires_at", "reservations", ["expires_at"])

    op.create_table(
        "stock_movements",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, primary_key=True),
        sa.Column("book_id", postgresql.UUID(as_uuid=True), nullable=False),
        sa.Column("movement_type", sa.String(length=50), nullable=False),
        sa.Column("quantity", sa.Integer(), nullable=False),
        sa.Column("reference_type", sa.String(length=50), nullable=True),
        sa.Column("reference_id", postgresql.UUID(as_uuid=True), nullable=True),
        sa.Column("notes", sa.Text(), nullable=True),
        sa.Column("created_at", sa.DateTime(timezone=True), nullable=False, server_default=sa.text("now()")),
        sa.CheckConstraint(
            "movement_type IN ('restock', 'sale', 'adjustment', 'reservation', 'reservation_release', 'reservation_commit')",
            name="valid_movement_type",
        ),
        comment="Audit trail of all inventory changes",
    )

    op.create_index("ix_stock_movements_book_id", "stock_movements", ["book_id"])
    op.create_index("ix_stock_movements_created_at", "stock_movements", ["created_at"])
    op.create_index("ix_stock_movements_reference", "stock_movements", ["reference_type", "reference_id"])


def downgrade() -> None:
    op.drop_index("ix_stock_movements_reference", table_name="stock_movements")
    op.drop_index("ix_stock_movements_created_at", table_name="stock_movements")
    op.drop_index("ix_stock_movements_book_id", table_name="stock_movements")
    op.drop_table("stock_movements")

    op.drop_index("ix_reservations_expires_at", table_name="reservations")
    op.drop_index("ix_reservations_status", table_name="reservations")
    op.drop_index("ix_reservations_book_id", table_name="reservations")
    op.drop_index("ix_reservations_order_id", table_name="reservations")
    op.drop_table("reservations")

    op.drop_index("ix_inventory_low_stock", table_name="inventory")
    op.drop_index("ix_inventory_title", table_name="inventory")
    op.drop_index("ix_inventory_book_id", table_name="inventory")
    op.drop_table("inventory")
