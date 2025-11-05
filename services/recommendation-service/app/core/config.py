"""
Configuration settings for the Recommendation Service.
"""

from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    # Application
    app_name: str = "Recommendation Service"
    app_version: str = "1.0.0"
    debug: bool = False

    # Server
    host: str = "0.0.0.0"
    port: int = 8089
    grpc_port: int = 50059

    # Database
    database_url: str = "postgresql://bookstore:dev_password@localhost:5432/recommendations_db"

    # External Services (gRPC)
    catalog_service_url: str = "localhost:50051"
    user_service_url: str = "localhost:50052"
    order_service_url: str = "localhost:50054"

    # Redis (for caching)
    redis_url: str = "redis://localhost:6379/0"
    cache_ttl: int = 3600  # 1 hour

    # Recommendation Settings
    default_recommendations_count: int = 10
    min_interactions_for_collaborative: int = 3
    similarity_threshold: float = 0.3

    # Strategy Weights (for hybrid approach)
    tag_weight: float = 0.4
    collaborative_weight: float = 0.4
    popular_weight: float = 0.2

    # CORS
    cors_origins: list = ["*"]

    # Logging
    log_level: str = "INFO"

    # Observability
    jaeger_agent_host: str = "localhost"
    jaeger_agent_port: int = 6831

    class Config:
        env_file = ".env"
        case_sensitive = False


@lru_cache()
def get_settings() -> Settings:
    """Get cached settings instance."""
    return Settings()
