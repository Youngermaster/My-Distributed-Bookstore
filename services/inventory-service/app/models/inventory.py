"""
Database models for Inventory Service.

This module defines SQLAlchemy models for:
- Inventory: Real-time stock tracking per book
- Reservation: Temporary stock holds for pending orders
- StockMovement: Audit trail of all stock changes
"""

from datetime import datetime, timedelta
from typing import Optional
from uuid import uuid4
from sqlalchemy import String, Integer, Text, CheckConstraint, Index
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.dialects.postgresql import UUID, TIMESTAMP
from app.db.base import Base


class Inventory(Base):
    """
    Inventory model for tracking book stock levels.

    Tracks both available (can be sold) and reserved (held for pending orders) quantities.
    """
    __tablename__ = "inventory"

    id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
        index=True
    )

    book_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        nullable=False,
        unique=True,  # One inventory record per book
        index=True
    )

    available_quantity: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        default=0,
        comment="Quantity available for sale"
    )

    reserved_quantity: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        default=0,
        comment="Quantity reserved for pending orders"
    )

    reorder_level: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        default=10,
        comment="Threshold for low stock alerts"
    )

    last_restocked_at: Mapped[Optional[datetime]] = mapped_column(
        TIMESTAMP(timezone=True),
        nullable=True,
        comment="Last time stock was added"
    )

    updated_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        nullable=False
    )

    __table_args__ = (
        CheckConstraint(
            'available_quantity >= 0',
            name='available_quantity_non_negative'
        ),
        CheckConstraint(
            'reserved_quantity >= 0',
            name='reserved_quantity_non_negative'
        ),
        CheckConstraint(
            'reorder_level >= 0',
            name='reorder_level_non_negative'
        ),
        Index('idx_inventory_book_id', 'book_id'),
        Index('idx_inventory_low_stock', 'available_quantity', 'reorder_level'),
    )

    @property
    def total_quantity(self) -> int:
        """Total quantity (available + reserved)."""
        return self.available_quantity + self.reserved_quantity

    @property
    def is_low_stock(self) -> bool:
        """Check if stock is below reorder level."""
        return self.available_quantity < self.reorder_level

    @property
    def is_in_stock(self) -> bool:
        """Check if any stock is available."""
        return self.available_quantity > 0


class Reservation(Base):
    """
    Reservation model for temporary stock holds during order processing.

    When an order is created, stock is reserved (moved from available to reserved).
    Reservations automatically expire after a configured timeout.
    """
    __tablename__ = "reservations"

    id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
        index=True
    )

    book_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        nullable=False,
        index=True
    )

    order_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        nullable=False,
        unique=True,  # One reservation per order
        index=True
    )

    quantity: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        comment="Quantity reserved"
    )

    status: Mapped[str] = mapped_column(
        String(20),
        nullable=False,
        default="pending",
        comment="pending, committed, released, expired"
    )

    expires_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        nullable=False,
        comment="Auto-release time"
    )

    created_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        default=datetime.utcnow,
        nullable=False
    )

    __table_args__ = (
        CheckConstraint(
            'quantity > 0',
            name='reservation_quantity_positive'
        ),
        CheckConstraint(
            "status IN ('pending', 'committed', 'released', 'expired')",
            name='valid_reservation_status'
        ),
        Index('idx_reservations_order_id', 'order_id'),
        Index('idx_reservations_book_id', 'book_id'),
        Index('idx_reservations_expires_at', 'expires_at'),
        Index('idx_reservations_status', 'status'),
    )

    @property
    def is_expired(self) -> bool:
        """Check if reservation has expired."""
        return datetime.utcnow() > self.expires_at

    @classmethod
    def create_with_expiry(cls, book_id: UUID, order_id: UUID, quantity: int, expiry_minutes: int = 15):
        """Factory method to create reservation with automatic expiry time."""
        return cls(
            book_id=book_id,
            order_id=order_id,
            quantity=quantity,
            expires_at=datetime.utcnow() + timedelta(minutes=expiry_minutes)
        )


class StockMovement(Base):
    """
    Stock movement model for audit trail of all inventory changes.

    Records every change to inventory levels with context about why the change occurred.
    """
    __tablename__ = "stock_movements"

    id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
        index=True
    )

    book_id: Mapped[UUID] = mapped_column(
        UUID(as_uuid=True),
        nullable=False,
        index=True
    )

    movement_type: Mapped[str] = mapped_column(
        String(50),
        nullable=False,
        comment="restock, sale, adjustment, reservation, reservation_release"
    )

    quantity: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        comment="Quantity changed (positive or negative)"
    )

    reference_type: Mapped[Optional[str]] = mapped_column(
        String(50),
        nullable=True,
        comment="order, purchase_order, manual_adjustment"
    )

    reference_id: Mapped[Optional[UUID]] = mapped_column(
        UUID(as_uuid=True),
        nullable=True,
        comment="ID of related entity (order_id, etc.)"
    )

    notes: Mapped[Optional[str]] = mapped_column(
        Text,
        nullable=True,
        comment="Additional context about the movement"
    )

    created_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        default=datetime.utcnow,
        nullable=False,
        index=True
    )

    __table_args__ = (
        CheckConstraint(
            "movement_type IN ('restock', 'sale', 'adjustment', 'reservation', 'reservation_release', 'reservation_commit')",
            name='valid_movement_type'
        ),
        Index('idx_stock_movements_book_id', 'book_id'),
        Index('idx_stock_movements_created_at', 'created_at'),
        Index('idx_stock_movements_reference', 'reference_type', 'reference_id'),
    )
