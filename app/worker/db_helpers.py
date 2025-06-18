from functools import wraps
from typing import Awaitable, Callable, TypeVar

from sqlalchemy.ext.asyncio import AsyncSession

from core.database import sessionmanager

F = TypeVar("F", bound=Callable[..., Awaitable[None]])


def with_db_session(func: F) -> F:
    @wraps(func)
    async def wrapper(*args, **kwargs):
        if not sessionmanager._sessionmaker:
            return

        session: AsyncSession = sessionmanager._sessionmaker()

        try:
            return await func(*args, db_session=session, **kwargs)
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()

    return wrapper  # type: ignore
