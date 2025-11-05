"""
API v1 router configuration.
"""

from fastapi import APIRouter

from app.api.v1.endpoints import health, recommendations

api_router = APIRouter()

# Include all endpoint routers
api_router.include_router(health.router, tags=["health"])
api_router.include_router(
    recommendations.router,
    prefix="/recommendations",
    tags=["recommendations"],
)
