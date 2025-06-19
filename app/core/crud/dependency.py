import datetime
import logging

from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession

from core.dependency_repository_resolvers.base import package_repository_resolvers
from core.models import Dependency as DependencyDBModel
from core.utils.limiter_config import LIMITERS

logger = logging.getLogger(__name__)


async def resolve_pending_dependencies(
    db_session: AsyncSession, batch_size: int, offset: int, ecosystem: str
):
    pending = (
        await db_session.scalars(
            select(DependencyDBModel)
            .where(
                DependencyDBModel.github_url_resolved == False,
                func.lower(DependencyDBModel.ecosystem) == ecosystem.lower()
            )
            .offset(offset)
            .limit(batch_size)
        )
    ).all()

    logger.info(f"Number of pending: {len(pending)}")

    for dependency in pending:
        ecosystem = dependency.ecosystem.lower()
        resolver = package_repository_resolvers[ecosystem]

        LIMITERS[ecosystem].wail_until_allowed()

        if resolver:
            url = await resolver(dependency.name)
            dependency.github_url = url
            dependency.github_url_resolved = True
            dependency.github_url_checked_at = datetime.datetime.now(datetime.UTC)

    await db_session.commit()
