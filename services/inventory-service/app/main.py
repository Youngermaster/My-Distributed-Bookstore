"""
Inventory Service - FastAPI Application

Manages stock tracking, reservations, and inventory movements.
"""
import logging
from contextlib import asynccontextmanager
from datetime import datetime

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.core.config import settings
from app.db.base import init_db, close_db
from app.api.v1.endpoints import inventory
from app.tasks.reservation_expiry import start_background_tasks

logging.basicConfig(
    level=getattr(logging, settings.LOG_LEVEL),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for startup and shutdown events."""
    # Startup
    logger.info(f"Starting {settings.APP_NAME} v{settings.APP_VERSION}")

    # Initialize database tables (development only - use Alembic in production)
    if settings.DEBUG:
        logger.info("Initializing database tables (DEBUG mode)...")
        await init_db()

    # Start background tasks
    logger.info("Starting background tasks...")
    start_background_tasks()

    logger.info(f"{settings.APP_NAME} started successfully on {settings.HOST}:{settings.PORT}")

    yield

    # Shutdown
    logger.info("Shutting down...")
    await close_db()
    logger.info("Shutdown complete")


app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    description="""
    Inventory Service for the Distributed Bookstore.
    
    ## Features
    * Real-time stock tracking
    * Stock reservations for orders
    * Stock movement history
    * Low stock alerts
    * Automatic reservation expiry
    """,
    lifespan=lifespan,
    debug=settings.DEBUG,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=settings.CORS_ALLOW_CREDENTIALS,
    allow_methods=settings.CORS_ALLOW_METHODS,
    allow_headers=settings.CORS_ALLOW_HEADERS,
)


@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": settings.APP_NAME,
        "version": settings.APP_VERSION,
        "timestamp": datetime.utcnow()
    }


@app.get("/")
async def root():
    return {
        "service": settings.APP_NAME,
        "version": settings.APP_VERSION,
        "docs": "/docs",
        "endpoints": {
            "health": "/health",
            "api_docs": "/docs",
            "inventory": "/api/v1/inventory"
        }
    }


# Include API routers
app.include_router(inventory.router, prefix="/api/v1")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level=settings.LOG_LEVEL.lower()
    )
