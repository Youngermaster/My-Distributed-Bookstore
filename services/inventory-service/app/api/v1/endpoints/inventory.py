"""
REST API endpoints for Inventory Service.

This module provides HTTP endpoints for:
- Inventory management (get, create, adjust)
- Stock reservations (reserve, release, commit)
- Stock movements (history)
- Low stock alerts
"""

import logging
from typing import Annotated
from uuid import UUID
from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.ext.asyncio import AsyncSession
from app.db.base import get_db
from app.services.inventory_service import InventoryService
from app.schemas.inventory import (
    InventoryCreate, InventoryResponse, InventoryAdjustRequest,
    ReserveStockRequest, ReserveStockResponse, ReleaseReservationRequest,
    ReleaseReservationResponse, CommitReservationRequest, CommitReservationResponse,
    StockMovementListResponse, LowStockListResponse
)

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/inventory", tags=["inventory"])


# ============================================================================
# Inventory Management Endpoints
# ============================================================================

@router.get(
    "/low-stock",
    response_model=LowStockListResponse,
    summary="Get low stock items",
    description="Retrieve list of books with stock below reorder level"
)
async def get_low_stock(
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """Get all inventory items with stock below their reorder level."""
    try:
        low_stock_items = await InventoryService.get_low_stock_items(db)
        return LowStockListResponse(items=low_stock_items, total=len(low_stock_items))
    except Exception as e:
        logger.error(f"Error retrieving low stock items: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to retrieve low stock items"
        )

@router.post(
    "",
    response_model=InventoryResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Create inventory record",
    description="Initialize inventory tracking for a new book (admin operation)"
)
async def create_inventory(
    inventory_data: InventoryCreate,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """
    Create a new inventory record for a book.

    This endpoint is typically called when a new book is added to the catalog.
    """
    try:
        inventory = await InventoryService.create_inventory(db, inventory_data)
        return inventory
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    except Exception as e:
        logger.error(f"Error creating inventory: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to create inventory"
        )


@router.get(
    "/{book_id}",
    response_model=InventoryResponse,
    summary="Get inventory by book ID",
    description="Retrieve current stock levels for a specific book"
)
async def get_inventory(
    book_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """
    Get inventory information for a specific book.

    Returns available quantity, reserved quantity, and stock status.
    """
    inventory = await InventoryService.get_inventory(db, book_id)

    if not inventory:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Inventory not found for book {book_id}"
        )

    return inventory


@router.post(
    "/{book_id}/adjust",
    response_model=InventoryResponse,
    summary="Adjust stock level",
    description="Manually adjust inventory stock (admin operation)"
)
async def adjust_stock(
    book_id: UUID,
    adjustment: InventoryAdjustRequest,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """
    Adjust inventory stock level.

    Supports three adjustment types:
    - **add**: Increase stock by quantity (e.g., restocking)
    - **subtract**: Decrease stock by quantity (e.g., damaged goods)
    - **set**: Set stock to exact quantity (e.g., after physical count)
    """
    try:
        inventory = await InventoryService.adjust_stock(db, book_id, adjustment)
        return inventory
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    except Exception as e:
        logger.error(f"Error adjusting stock for book {book_id}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to adjust stock"
        )


@router.get(
    "/low-stock",
    response_model=LowStockListResponse,
    summary="Get low stock items",
    description="Retrieve all items with stock below reorder level"
)
async def get_low_stock_items(
    db: Annotated[AsyncSession, Depends(get_db)],
    threshold: int = Query(None, ge=0, description="Optional custom threshold")
):
    """
    Get all items with low stock.

    By default, uses each item's configured reorder_level.
    Optionally, specify a custom threshold to apply to all items.
    """
    try:
        items = await InventoryService.get_low_stock_items(db, threshold)
        return LowStockListResponse(
            items=items,
            total=len(items)
        )
    except Exception as e:
        logger.error(f"Error retrieving low stock items: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to retrieve low stock items"
        )


# ============================================================================
# Reservation Endpoints
# ============================================================================

@router.post(
    "/reserve",
    response_model=ReserveStockResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Reserve stock for order",
    description="Reserve inventory for a pending order (called by Order Service)"
)
async def reserve_stock(
    reserve_request: ReserveStockRequest,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """
    Reserve stock for an order.

    This endpoint is called when an order is created. Stock is moved from
    available to reserved and will automatically expire after the configured timeout.
    """
    try:
        reservations, expires_at = await InventoryService.reserve_stock(db, reserve_request)

        return ReserveStockResponse(
            success=True,
            order_id=reserve_request.order_id,
            reservations=reservations,
            expires_at=expires_at,
            message=f"Reserved stock for {len(reservations)} items"
        )
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(e)
        )
    except Exception as e:
        logger.error(f"Error reserving stock for order {reserve_request.order_id}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to reserve stock"
        )


@router.post(
    "/release/{order_id}",
    response_model=ReleaseReservationResponse,
    summary="Release reservation",
    description="Release reserved stock (order cancelled or expired)"
)
async def release_reservation(
    order_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """
    Release a stock reservation.

    This endpoint is called when:
    - Order is cancelled by user
    - Payment fails
    - Reservation expires (automatic)

    Stock is moved back from reserved to available.
    """
    try:
        released_quantity = await InventoryService.release_reservation(db, order_id)

        return ReleaseReservationResponse(
            success=True,
            order_id=order_id,
            released_quantity=released_quantity,
            message=f"Released {released_quantity} units"
        )
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e)
        )
    except Exception as e:
        logger.error(f"Error releasing reservation for order {order_id}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to release reservation"
        )


@router.post(
    "/commit/{order_id}",
    response_model=CommitReservationResponse,
    summary="Commit reservation",
    description="Commit reserved stock (payment successful)"
)
async def commit_reservation(
    order_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)]
):
    """
    Commit a stock reservation.

    This endpoint is called when payment is successful.
    Stock is removed from reserved quantity (completing the sale).
    """
    try:
        committed_quantity = await InventoryService.commit_reservation(db, order_id)

        return CommitReservationResponse(
            success=True,
            order_id=order_id,
            committed_quantity=committed_quantity,
            message=f"Committed {committed_quantity} units"
        )
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e)
        )
    except Exception as e:
        logger.error(f"Error committing reservation for order {order_id}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to commit reservation"
        )


# ============================================================================
# Stock Movement Endpoints
# ============================================================================

@router.get(
    "/{book_id}/movements",
    response_model=StockMovementListResponse,
    summary="Get stock movement history",
    description="Retrieve audit trail of all stock changes for a book"
)
async def get_stock_movements(
    book_id: UUID,
    db: Annotated[AsyncSession, Depends(get_db)],
    page: int = Query(1, ge=1, description="Page number"),
    page_size: int = Query(50, ge=1, le=100, description="Items per page")
):
    """
    Get stock movement history for a book.

    Returns paginated list of all stock changes with context about each change.
    Useful for auditing and troubleshooting inventory discrepancies.
    """
    try:
        skip = (page - 1) * page_size
        movements, total = await InventoryService.get_stock_movements(
            db, book_id, skip, page_size
        )

        total_pages = (total + page_size - 1) // page_size

        return StockMovementListResponse(
            movements=movements,
            total=total,
            page=page,
            page_size=page_size,
            total_pages=total_pages
        )
    except Exception as e:
        logger.error(f"Error retrieving movements for book {book_id}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to retrieve stock movements"
        )
