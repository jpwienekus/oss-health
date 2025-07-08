import random
from typing import List, Sequence

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from core.models import Repository as RepositoryDBModel
from core.models import (
    RepositoryDependencyVersion as RepositoryDependencyVersionDBModel,
)
from core.models import Version as VersionDBModel


async def get_repositories(
    db_session: AsyncSession, user_id: int
) -> Sequence[RepositoryDBModel]:
    return (
        await db_session.scalars(
            select(RepositoryDBModel)
            .options(
                selectinload(RepositoryDBModel.dependency_versions)
                .selectinload(RepositoryDependencyVersionDBModel.version)
                .selectinload(VersionDBModel.vulnerabilities)
            )
            .where(RepositoryDBModel.user_id == user_id)
        )
    ).all()


async def add_repository_ids(
    db_session: AsyncSession, user_id, tracked_repositories: List[dict]
):
    db_session.add_all(
        [
            RepositoryDBModel(
                github_id=repository.get("id"),
                user_id=user_id,
                score=0,
                url=repository.get("url"),
                scan_day=random.randint(0, 6),
                scan_hour=random.randint(0, 23),
                scan_status="pending",
            )
            for repository in tracked_repositories
        ]
    )

    await db_session.commit()


async def get_cron_info(db_session: AsyncSession):
    return (
        await db_session.execute(
            select(
                RepositoryDBModel.scan_day,
                RepositoryDBModel.scan_hour,
                func.count().label("total"),
            ).group_by(RepositoryDBModel.scan_day, RepositoryDBModel.scan_hour)
        )
    ).all()
