"""Alembic environment configuration for the inventory service."""

from logging.config import fileConfig

from alembic import context
from sqlalchemy import create_engine, pool

from app.core.config import settings
from app.db.base import Base

# Import models so metadata is registered
from app.models import inventory  # noqa: F401


config = context.config

if config.config_file_name is not None:
    fileConfig(config.config_file_name)


def _get_sync_database_url() -> str:
    url = str(settings.DATABASE_URL)
    if "+asyncpg" in url:
        url = url.replace("+asyncpg", "")
    return url


target_metadata = Base.metadata


def run_migrations_offline() -> None:
    context.configure(
        url=_get_sync_database_url(),
        target_metadata=target_metadata,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
    )

    with context.begin_transaction():
        context.run_migrations()


def run_migrations_online() -> None:
    connectable = create_engine(
        _get_sync_database_url(),
        poolclass=pool.NullPool,
    )

    with connectable.connect() as connection:
        context.configure(connection=connection, target_metadata=target_metadata)

        with context.begin_transaction():
            context.run_migrations()

    connectable.dispose()


if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online()
