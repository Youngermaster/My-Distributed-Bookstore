"""
Background task for automatically expiring old reservations.

This module provides a periodic task that runs every minute to find and
expire reservations that have passed their expiry time.
"""

import logging
import asyncio
from datetime import datetime
from app.db.base import AsyncSessionLocal
from app.services.inventory_service import InventoryService

logger = logging.getLogger(__name__)


async def expire_reservations_task():
    """
    Background task to expire old reservations.

    This task runs periodically to:
    1. Find all reservations past their expiry time
    2. Release the reserved stock back to available
    3. Mark reservations as expired

    Runs every 60 seconds by default.
    """
    logger.info("Starting reservation expiry background task")

    while True:
        try:
            async with AsyncSessionLocal() as db:
                expired_count = await InventoryService.expire_old_reservations(db)

                if expired_count > 0:
                    logger.info(f"Expired {expired_count} reservations at {datetime.utcnow()}")

        except Exception as e:
            logger.error(f"Error in reservation expiry task: {e}")

        # Wait 60 seconds before next check
        await asyncio.sleep(60)


def start_background_tasks():
    """
    Start all background tasks.

    This function should be called during application startup.
    """
    logger.info("Initializing background tasks...")

    # Create task for reservation expiry
    asyncio.create_task(expire_reservations_task())

    logger.info("Background tasks started successfully")
