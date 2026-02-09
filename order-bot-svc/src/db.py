from collections.abc import AsyncGenerator

from sqlalchemy.ext.asyncio import AsyncAttrs, AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.orm import DeclarativeBase
from src.config import settings


class Base(AsyncAttrs, DeclarativeBase):
    """
    Use AsyncAttrs.awaitable_attrs to lazily load related entities
    Example:
        async def get_cart_items(cart: Cart) -> CartItem:
            return cart.awaitable_attrs.items
    """
    pass


engine = create_async_engine(settings.database_url, pool_pre_ping=True)
SessionLocal = async_sessionmaker(bind=engine, autoflush=False, autocommit=False)


async def get_db_session() -> AsyncGenerator[AsyncSession, None]:
    async with SessionLocal() as db:
        yield db
