"""
Business logic service for Inventory management.

This service handles:
- Inventory CRUD operations
- Stock reservations for orders
- Reservation commits and releases
- Stock adjustments and movements
- Low stock detection
"""

import logging
from datetime import datetime, timedelta
from typing import List, Tuple, Optional
from uuid import UUID
from sqlalchemy import select, func, and_
from sqlalchemy.ext.asyncio import AsyncSession
from app.models.inventory import Inventory, Reservation, StockMovement
from app.schemas.inventory import (
    InventoryCreate, InventoryAdjustRequest, ReserveStockRequest,
    LowStockItem
)
from app.core.config import settings

logger = logging.getLogger(__name__)


class InventoryService:
    """Service class for inventory-related operations."""

    # ========================================================================
    # Inventory Management
    # ========================================================================

    @staticmethod
    async def create_inventory(
        db: AsyncSession,
        inventory_data: InventoryCreate
    ) -> Inventory:
        """
        Create new inventory record for a book.

        Args:
            db: Database session
            inventory_data: Inventory creation data

        Returns:
            Created inventory record

        Raises:
            ValueError: If inventory already exists for this book
        """
        # Check if inventory already exists
        existing_query = select(Inventory).where(
            Inventory.book_id == inventory_data.book_id
        )
        result = await db.execute(existing_query)
        existing = result.scalar_one_or_none()

        if existing:
            raise ValueError(f"Inventory already exists for book {inventory_data.book_id}")

        # Create inventory
        inventory = Inventory(
            book_id=inventory_data.book_id,
            title=inventory_data.title,
            short_description=inventory_data.short_description,
            available_quantity=inventory_data.initial_quantity,
            reserved_quantity=0,
            reorder_level=inventory_data.reorder_level,
            last_restocked_at=datetime.utcnow() if inventory_data.initial_quantity > 0 else None
        )

        db.add(inventory)
        await db.flush()

        # Record initial stock movement
        if inventory_data.initial_quantity > 0:
            movement = StockMovement(
                book_id=inventory_data.book_id,
                movement_type="restock",
                quantity=inventory_data.initial_quantity,
                reference_type="initial_stock",
                notes="Initial inventory creation"
            )
            db.add(movement)

        await db.commit()
        await db.refresh(inventory)

        logger.info(f"Created inventory for book {inventory_data.book_id} with {inventory_data.initial_quantity} units")
        return inventory

    @staticmethod
    async def get_inventory(db: AsyncSession, book_id: UUID) -> Optional[Inventory]:
        """Get inventory by book ID."""
        query = select(Inventory).where(Inventory.book_id == book_id)
        result = await db.execute(query)
        return result.scalar_one_or_none()

    @staticmethod
    async def adjust_stock(
        db: AsyncSession,
        book_id: UUID,
        adjustment: InventoryAdjustRequest
    ) -> Inventory:
        """
        Adjust inventory stock level (admin operation).

        Args:
            db: Database session
            book_id: Book to adjust
            adjustment: Adjustment details

        Returns:
            Updated inventory

        Raises:
            ValueError: If inventory doesn't exist or invalid adjustment
        """
        # Get inventory
        inventory = await InventoryService.get_inventory(db, book_id)
        if not inventory:
            raise ValueError(f"Inventory not found for book {book_id}")

        old_quantity = inventory.available_quantity
        movement_quantity = 0

        # Apply adjustment
        if adjustment.adjustment_type == "add":
            inventory.available_quantity += adjustment.quantity
            movement_quantity = adjustment.quantity
            inventory.last_restocked_at = datetime.utcnow()

        elif adjustment.adjustment_type == "subtract":
            if inventory.available_quantity < adjustment.quantity:
                raise ValueError(
                    f"Insufficient stock. Available: {inventory.available_quantity}, "
                    f"Requested: {adjustment.quantity}"
                )
            inventory.available_quantity -= adjustment.quantity
            movement_quantity = -adjustment.quantity

        elif adjustment.adjustment_type == "set":
            movement_quantity = adjustment.quantity - inventory.available_quantity
            inventory.available_quantity = adjustment.quantity
            if adjustment.quantity > old_quantity:
                inventory.last_restocked_at = datetime.utcnow()

        inventory.updated_at = datetime.utcnow()

        # Record movement
        movement = StockMovement(
            book_id=book_id,
            movement_type="adjustment",
            quantity=movement_quantity,
            reference_type="manual_adjustment",
            notes=adjustment.reason or f"Manual {adjustment.adjustment_type} adjustment"
        )
        db.add(movement)

        await db.commit()
        await db.refresh(inventory)

        logger.info(
            f"Adjusted inventory for book {book_id}: {old_quantity} -> {inventory.available_quantity}"
        )
        return inventory

    @staticmethod
    async def get_low_stock_items(
        db: AsyncSession,
        threshold: Optional[int] = None
    ) -> List[LowStockItem]:
        """
        Get all items with stock below their reorder level.

        Args:
            db: Database session
            threshold: Optional override for reorder level threshold

        Returns:
            List of low stock items
        """
        if threshold is not None:
            query = select(Inventory).where(
                Inventory.available_quantity < threshold
            )
        else:
            query = select(Inventory).where(
                Inventory.available_quantity < Inventory.reorder_level
            )

        result = await db.execute(query)
        inventories = result.scalars().all()

        low_stock_items = [
            LowStockItem(
                book_id=inv.book_id,
                title=inv.title,
                short_description=inv.short_description,
                available_quantity=inv.available_quantity,
                reorder_level=inv.reorder_level,
                deficit=inv.reorder_level - inv.available_quantity
            )
            for inv in inventories
        ]

        return low_stock_items

    # ========================================================================
    # Reservation Management
    # ========================================================================

    @staticmethod
    async def reserve_stock(
        db: AsyncSession,
        reserve_request: ReserveStockRequest
    ) -> Tuple[List[Reservation], datetime]:
        """
        Reserve stock for an order.

        Args:
            db: Database session
            reserve_request: Reservation request with order ID and items

        Returns:
            Tuple of (list of created reservations, expiry time)

        Raises:
            ValueError: If insufficient stock or reservation already exists
        """
        # Check if reservation already exists for this order
        existing_query = select(Reservation).where(
            Reservation.order_id == reserve_request.order_id
        )
        result = await db.execute(existing_query)
        if result.scalar_one_or_none():
            raise ValueError(f"Reservation already exists for order {reserve_request.order_id}")

        # Validate stock availability for all items
        for item in reserve_request.items:
            inventory = await InventoryService.get_inventory(db, item.book_id)
            if not inventory:
                raise ValueError(f"Inventory not found for book {item.book_id}")

            if inventory.available_quantity < item.quantity:
                raise ValueError(
                    f"Insufficient stock for book {item.book_id}. "
                    f"Available: {inventory.available_quantity}, Requested: {item.quantity}"
                )

        # Create reservations and update inventory
        reservations = []
        expiry_time = datetime.utcnow() + timedelta(minutes=settings.RESERVATION_EXPIRY_MINUTES)

        for item in reserve_request.items:
            # Update inventory
            inventory = await InventoryService.get_inventory(db, item.book_id)
            inventory.available_quantity -= item.quantity
            inventory.reserved_quantity += item.quantity
            inventory.updated_at = datetime.utcnow()

            # Create reservation
            reservation = Reservation(
                book_id=item.book_id,
                order_id=reserve_request.order_id,
                quantity=item.quantity,
                status="pending",
                expires_at=expiry_time
            )
            db.add(reservation)
            reservations.append(reservation)

            # Record movement
            movement = StockMovement(
                book_id=item.book_id,
                movement_type="reservation",
                quantity=-item.quantity,
                reference_type="order",
                reference_id=reserve_request.order_id,
                notes=f"Stock reserved for order {reserve_request.order_id}"
            )
            db.add(movement)

        await db.commit()
        for res in reservations:
            await db.refresh(res)

        logger.info(f"Reserved stock for order {reserve_request.order_id}, expires at {expiry_time}")
        return reservations, expiry_time

    @staticmethod
    async def release_reservation(
        db: AsyncSession,
        order_id: UUID,
        mark_expired: bool = False,
    ) -> int:
        """
        Release a reservation (order cancelled or expired).

        Args:
            db: Database session
            order_id: Order ID whose reservation to release

        Returns:
            Total quantity released

        Raises:
            ValueError: If reservation not found
        """
        # Get all reservations for order
        query = select(Reservation).where(
            and_(
                Reservation.order_id == order_id,
                Reservation.status == "pending"
            )
        )
        result = await db.execute(query)
        reservations = result.scalars().all()

        if not reservations:
            raise ValueError(f"No pending reservation found for order {order_id}")

        total_released = 0

        target_status = "expired" if mark_expired else "released"

        for reservation in reservations:
            # Update inventory
            inventory = await InventoryService.get_inventory(db, reservation.book_id)
            if inventory:
                inventory.reserved_quantity -= reservation.quantity
                inventory.available_quantity += reservation.quantity
                inventory.updated_at = datetime.utcnow()

            # Update reservation status
            reservation.status = target_status
            total_released += reservation.quantity

            # Record movement
            movement = StockMovement(
                book_id=reservation.book_id,
                movement_type="reservation_release",
                quantity=reservation.quantity,
                reference_type="order",
                reference_id=order_id,
                notes=(
                    f"Reservation expired for order {order_id}"
                    if mark_expired
                    else f"Reservation released for order {order_id}"
                ),
            )
            db.add(movement)

        await db.commit()

        action = "expired" if mark_expired else "released"
        logger.info(f"{action.capitalize()} {total_released} units for order {order_id}")
        return total_released

    @staticmethod
    async def commit_reservation(
        db: AsyncSession,
        order_id: UUID
    ) -> int:
        """
        Commit a reservation (payment successful).

        Args:
            db: Database session
            order_id: Order ID whose reservation to commit

        Returns:
            Total quantity committed

        Raises:
            ValueError: If reservation not found
        """
        # Get all reservations for order
        query = select(Reservation).where(
            and_(
                Reservation.order_id == order_id,
                Reservation.status == "pending"
            )
        )
        result = await db.execute(query)
        reservations = result.scalars().all()

        if not reservations:
            raise ValueError(f"No pending reservation found for order {order_id}")

        total_committed = 0

        for reservation in reservations:
            # Update inventory (decrease reserved quantity)
            inventory = await InventoryService.get_inventory(db, reservation.book_id)
            if inventory:
                inventory.reserved_quantity -= reservation.quantity
                inventory.updated_at = datetime.utcnow()

            # Update reservation status
            reservation.status = "committed"
            total_committed += reservation.quantity

            # Record movement
            movement = StockMovement(
                book_id=reservation.book_id,
                movement_type="sale",
                quantity=-reservation.quantity,
                reference_type="order",
                reference_id=order_id,
                notes=f"Sale completed for order {order_id}"
            )
            db.add(movement)

        await db.commit()

        logger.info(f"Committed {total_committed} units for order {order_id}")
        return total_committed

    @staticmethod
    async def expire_old_reservations(db: AsyncSession) -> int:
        """
        Find and expire old reservations (background task).

        Returns:
            Number of reservations expired
        """
        query = (
            select(Reservation.order_id)
            .where(
                and_(
                    Reservation.status == "pending",
                    Reservation.expires_at <= datetime.utcnow(),
                )
            )
            .distinct()
        )
        result = await db.execute(query)
        order_ids = [row.order_id for row in result.all()]

        expired_count = 0

        for order_id in order_ids:
            try:
                released = await InventoryService.release_reservation(
                    db,
                    order_id,
                    mark_expired=True,
                )
                if released > 0:
                    expired_count += 1
            except ValueError:
                continue
            except Exception as exc:
                logger.error(f"Failed to expire reservation for order {order_id}: {exc}")

        if expired_count > 0:
            logger.info(f"Expired reservations for {expired_count} order(s)")

        return expired_count

    # ========================================================================
    # Stock Movements
    # ========================================================================

    @staticmethod
    async def get_stock_movements(
        db: AsyncSession,
        book_id: UUID,
        skip: int = 0,
        limit: int = 50
    ) -> Tuple[List[StockMovement], int]:
        """
        Get stock movement history for a book.

        Args:
            db: Database session
            book_id: Book ID
            skip: Pagination offset
            limit: Page size

        Returns:
            Tuple of (movements list, total count)
        """
        # Count query
        count_query = select(func.count()).select_from(StockMovement).where(
            StockMovement.book_id == book_id
        )
        count_result = await db.execute(count_query)
        total = count_result.scalar_one()

        # Data query
        query = select(StockMovement).where(
            StockMovement.book_id == book_id
        ).order_by(
            StockMovement.created_at.desc()
        ).offset(skip).limit(limit)

        result = await db.execute(query)
        movements = result.scalars().all()

        return list(movements), total
