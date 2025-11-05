"""
Health check endpoints.
"""

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from typing import Dict

from app.core.database import get_db
from app.core.config import get_settings

settings = get_settings()
router = APIRouter()


@router.get("/health", tags=["health"])
def health_check() -> Dict[str, str]:
    """Health check endpoint."""
    return {
        "status": "healthy",
        "service": settings.app_name,
        "version": settings.app_version,
    }


@router.get("/ready", tags=["health"])
def readiness_check(db: Session = Depends(get_db)) -> Dict[str, str]:
    """
    Readiness check endpoint.

    Verifies database connectivity.
    """
    try:
        # Try to execute a simple query
        db.execute("SELECT 1")
        return {
            "status": "ready",
            "service": settings.app_name,
            "database": "connected",
        }
    except Exception as e:
        return {
            "status": "not ready",
            "service": settings.app_name,
            "database": "disconnected",
            "error": str(e),
        }
