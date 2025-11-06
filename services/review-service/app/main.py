"""
Review Service - FastAPI Application Entry Point

This service handles book reviews, ratings, and ML-powered sentiment analysis.
"""
import logging
from contextlib import asynccontextmanager
from datetime import datetime

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.core.config import settings
from app.db.base import engine, Base
from app.api.v1.endpoints import reviews
from app.ml.sentiment import download_nltk_data
from app.schemas.review import HealthResponse

# Configure logging
logging.basicConfig(
    level=getattr(logging, settings.LOG_LEVEL),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Lifespan context manager for startup and shutdown events.
    Replaces deprecated @app.on_event decorators.
    """
    # Startup
    logger.info(f"Starting {settings.APP_NAME} v{settings.APP_VERSION}")

    # Ensure database tables exist (idempotent)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
        logger.info("Database schema ensured")

    # Download NLTK data for sentiment analysis
    download_nltk_data()

    logger.info(f"Service ready on {settings.HOST}:{settings.PORT}")

    yield  # Application runs here

    # Shutdown
    logger.info("Shutting down...")
    await engine.dispose()
    logger.info("Service stopped")


# Create FastAPI application
app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    description="""
    Review Service for the Distributed Bookstore.

    ## Features
    * Create, read, update, and delete book reviews
    * Automatic sentiment analysis using ML (NLTK & TextBlob)
    * Review voting system (helpful/not helpful)
    * Review statistics and aggregations
    * Rating distribution analytics

    ## ML Capabilities
    Reviews are automatically analyzed for sentiment and classified as:
    * **Positive**: Sentiment score > 0.3
    * **Neutral**: Sentiment score between -0.3 and 0.3
    * **Negative**: Sentiment score < -0.3
    """,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    lifespan=lifespan,
    debug=settings.DEBUG,
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=settings.CORS_ALLOW_CREDENTIALS,
    allow_methods=settings.CORS_ALLOW_METHODS,
    allow_headers=settings.CORS_ALLOW_HEADERS,
)


# ============================================================================
# Health Check Endpoints
# ============================================================================

@app.get(
    "/health",
    response_model=HealthResponse,
    tags=["health"],
    summary="Health check",
    description="Check if the service is running."
)
async def health_check():
    """Health check endpoint."""
    return HealthResponse(
        status="healthy",
        service=settings.APP_NAME,
        version=settings.APP_VERSION,
        timestamp=datetime.utcnow()
    )


@app.get(
    "/",
    tags=["root"],
    summary="Root endpoint",
    description="Get service information."
)
async def root():
    """Root endpoint."""
    return {
        "service": settings.APP_NAME,
        "version": settings.APP_VERSION,
        "docs": "/docs",
        "health": "/health"
    }


# ============================================================================
# API Routes
# ============================================================================

# Include review endpoints
app.include_router(reviews.router, prefix="/api/v1")


# ============================================================================
# Exception Handlers
# ============================================================================

@app.exception_handler(Exception)
async def global_exception_handler(request, exc):
    """Global exception handler."""
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={
            "detail": "Internal server error",
            "type": "internal_error"
        }
    )


if __name__ == "__main__":
    import uvicorn

    # Run with uvicorn
    # In production, use: uvicorn app.main:app --host 0.0.0.0 --port 8088 --workers 4
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level=settings.LOG_LEVEL.lower()
    )
