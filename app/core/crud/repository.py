import random
from datetime import datetime, timezone
from typing import List, Sequence

from sqlalchemy import select
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
                url=repository.get("clone_url"),
                scan_day=random.randint(1, 7),
                scan_hour=random.randint(1, 24),
                scan_status="pending"
            )
            for repository in tracked_repositories
        ]
    )

    await db_session.commit()


async def get_repository(db_session: AsyncSession, repository_id: int, user_id: int):
    return (
        await db_session.scalars(
            select(RepositoryDBModel).where(
                RepositoryDBModel.id == repository_id,
                RepositoryDBModel.user_id == user_id,
            )
        )
    ).first()


async def update_scanned_date(
    db_session: AsyncSession, repository_id: int, user_id: int
):
    repository = await get_repository(db_session, repository_id, user_id)

    if repository:
        repository.last_scanned_at = datetime.now(timezone.utc)

    await db_session.commit()
