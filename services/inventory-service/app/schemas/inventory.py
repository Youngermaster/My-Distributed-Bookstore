"""
Pydantic schemas for Inventory Service API contracts.

This module defines request and response models for all inventory operations.
"""

from datetime import datetime
from typing import Optional, List
from uuid import UUID
from pydantic import BaseModel, Field, ConfigDict, field_validator


# ============================================================================
# Inventory Schemas
# ============================================================================

class InventoryBase(BaseModel):
    """Base schema for inventory data."""
    book_id: UUID = Field(..., description="ID of the book")
    available_quantity: int = Field(..., ge=0, description="Available stock quantity")
    reserved_quantity: int = Field(default=0, ge=0, description="Reserved stock quantity")
    reorder_level: int = Field(default=10, ge=0, description="Low stock threshold")


class InventoryCreate(BaseModel):
    """Schema for creating new inventory record."""
    book_id: UUID = Field(..., description="ID of the book")
    initial_quantity: int = Field(..., ge=0, description="Initial stock quantity")
    reorder_level: int = Field(default=10, ge=0, description="Low stock threshold")

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "book_id": "123e4567-e89b-12d3-a456-426614174000",
                "initial_quantity": 100,
                "reorder_level": 15
            }
        }
    )


class InventoryAdjustRequest(BaseModel):
    """Schema for adjusting inventory levels (admin operation)."""
    adjustment_type: str = Field(..., description="add, subtract, or set")
    quantity: int = Field(..., description="Quantity to adjust")
    reason: Optional[str] = Field(None, max_length=500, description="Reason for adjustment")

    @field_validator('adjustment_type')
    @classmethod
    def validate_adjustment_type(cls, v: str) -> str:
        if v not in ['add', 'subtract', 'set']:
            raise ValueError("adjustment_type must be 'add', 'subtract', or 'set'")
        return v

    @field_validator('quantity')
    @classmethod
    def validate_quantity(cls, v: int, info) -> int:
        adjustment_type = info.data.get('adjustment_type')
        if adjustment_type in ['add', 'set'] and v < 0:
            raise ValueError(f"quantity must be non-negative for '{adjustment_type}'")
        if adjustment_type == 'subtract' and v <= 0:
            raise ValueError("quantity must be positive for 'subtract'")
        return v

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "adjustment_type": "add",
                "quantity": 50,
                "reason": "Received new shipment"
            }
        }
    )


class InventoryResponse(BaseModel):
    """Schema for inventory response."""
    id: UUID
    book_id: UUID
    available_quantity: int
    reserved_quantity: int
    total_quantity: int
    reorder_level: int
    is_low_stock: bool
    is_in_stock: bool
    last_restocked_at: Optional[datetime] = None
    updated_at: datetime

    model_config = ConfigDict(from_attributes=True)


class LowStockItem(BaseModel):
    """Schema for low stock alert."""
    book_id: UUID
    available_quantity: int
    reorder_level: int
    deficit: int = Field(..., description="How many units below reorder level")

    model_config = ConfigDict(from_attributes=True)


class LowStockListResponse(BaseModel):
    """Schema for list of low stock items."""
    items: List[LowStockItem]
    total: int


# ============================================================================
# Reservation Schemas
# ============================================================================

class ReservationItemRequest(BaseModel):
    """Schema for a single item in a reservation request."""
    book_id: UUID = Field(..., description="ID of the book to reserve")
    quantity: int = Field(..., gt=0, description="Quantity to reserve")


class ReserveStockRequest(BaseModel):
    """Schema for reserving stock for an order."""
    order_id: UUID = Field(..., description="ID of the order")
    items: List[ReservationItemRequest] = Field(..., min_length=1, description="Items to reserve")

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "order_id": "123e4567-e89b-12d3-a456-426614174000",
                "items": [
                    {
                        "book_id": "223e4567-e89b-12d3-a456-426614174000",
                        "quantity": 2
                    },
                    {
                        "book_id": "323e4567-e89b-12d3-a456-426614174000",
                        "quantity": 1
                    }
                ]
            }
        }
    )


class ReservationResponse(BaseModel):
    """Schema for reservation response."""
    id: UUID
    book_id: UUID
    order_id: UUID
    quantity: int
    status: str
    expires_at: datetime
    created_at: datetime

    model_config = ConfigDict(from_attributes=True)


class ReserveStockResponse(BaseModel):
    """Schema for stock reservation result."""
    success: bool
    order_id: UUID
    reservations: List[ReservationResponse]
    expires_at: datetime
    message: str = Field(default="Stock reserved successfully")


class ReleaseReservationRequest(BaseModel):
    """Schema for releasing a reservation."""
    order_id: UUID = Field(..., description="ID of the order whose reservation to release")


class ReleaseReservationResponse(BaseModel):
    """Schema for reservation release result."""
    success: bool
    order_id: UUID
    released_quantity: int
    message: str


class CommitReservationRequest(BaseModel):
    """Schema for committing a reservation (after payment)."""
    order_id: UUID = Field(..., description="ID of the order whose reservation to commit")


class CommitReservationResponse(BaseModel):
    """Schema for reservation commit result."""
    success: bool
    order_id: UUID
    committed_quantity: int
    message: str


# ============================================================================
# Stock Movement Schemas
# ============================================================================

class StockMovementResponse(BaseModel):
    """Schema for stock movement response."""
    id: UUID
    book_id: UUID
    movement_type: str
    quantity: int
    reference_type: Optional[str] = None
    reference_id: Optional[UUID] = None
    notes: Optional[str] = None
    created_at: datetime

    model_config = ConfigDict(from_attributes=True)


class StockMovementListResponse(BaseModel):
    """Schema for list of stock movements."""
    movements: List[StockMovementResponse]
    total: int
    page: int
    page_size: int
    total_pages: int


# ============================================================================
# Health Check Schema
# ============================================================================

class HealthCheckResponse(BaseModel):
    """Schema for health check response."""
    status: str = Field(..., description="Service health status")
    service: str = Field(..., description="Service name")
    timestamp: datetime = Field(..., description="Current timestamp")
    version: str = Field(default="1.0.0", description="Service version")
