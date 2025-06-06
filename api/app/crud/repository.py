from typing import Sequence
from sqlalchemy import delete, select
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import Repository as RepositoryDBModel


async def get_repositories(db_session: AsyncSession) -> Sequence[RepositoryDBModel]:
    return (await db_session.scalars(select(RepositoryDBModel))).all()


async def sync_repository_ids(db_session: AsyncSession, user_id, target_ids: list[int]):
    id_set = set(target_ids)

    existing_ids = set(
        (await db_session.scalars(select(RepositoryDBModel.github_id))).all()
    )

    ids_to_add = id_set - existing_ids
    ids_to_delete = existing_ids - id_set

    db_session.add_all(
        [RepositoryDBModel(github_id=id, user_id=user_id, score=0) for id in ids_to_add]
    )

    if ids_to_delete:
        await db_session.execute(
            delete(RepositoryDBModel).where(
                RepositoryDBModel.github_id.in_(ids_to_delete)
            )
        )

    await db_session.commit()
