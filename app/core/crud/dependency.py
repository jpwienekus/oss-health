import logging
import datetime
from sqlalchemy import select
from core.models import Dependency as DependencyDBModel
from core.dependency_repository_resolvers.base import  package_repository_resolvers

from sqlalchemy.ext.asyncio import AsyncSession

logger = logging.getLogger()

async def resolve_pending_dependencies(db_session: AsyncSession, batch_size: int, offset: int):
    pending = (
        await db_session.scalars(
            select(DependencyDBModel)
            .where(DependencyDBModel.github_url_resolved == False)
            .offset(offset)
            .limit(batch_size)
        )
    ).all()

    logger.info(f"Number of pending: {len(pending)}")

    for dependency in pending:
        ecosystem = dependency.ecosystem.lower()
        resolver = package_repository_resolvers[ecosystem]
        logger.info('$' * 100)
        logger.info(ecosystem)
        
        if resolver:
            url = await resolver(dependency.name)
            dependency.github_url = url
            dependency.github_url_resolved = True
            dependency.github_url_checked_at = datetime.datetime.now(datetime.UTC)

    await db_session.commit()

