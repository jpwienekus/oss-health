from typing import Sequence
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import Repository as RepositoryDBModel


async def get_repositories(db_session: AsyncSession) -> Sequence[RepositoryDBModel]:
    return (await db_session.scalars(select(RepositoryDBModel))).all()


async def sync_repository_ids(db_session: AsyncSession, user_id, target_ids: list[int]):
    db_session.add_all(
        [RepositoryDBModel(github_id=id, user_id=user_id, score=0) for id in target_ids]
    )

    await db_session.commit()
