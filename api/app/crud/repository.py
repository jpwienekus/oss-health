from typing import List, Sequence
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import Repository as RepositoryDBModel


async def get_repositories(
    db_session: AsyncSession, user_id: int
) -> Sequence[RepositoryDBModel]:
    return (
        await db_session.scalars(
            select(RepositoryDBModel).where(RepositoryDBModel.user_id == user_id)
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
                clone_url=repository.get("clone_url"),
            )
            for repository in tracked_repositories
        ]
    )

    await db_session.commit()
