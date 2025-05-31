from models import User as UserDBModel
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession


async def get_user(db_session: AsyncSession, github_id: int):
    return (await db_session.scalars(select(UserDBModel).where(UserDBModel.github_id == github_id))).first()


async def add_user(db_session: AsyncSession, github_id: int, github_username: str):
    db_session.add(UserDBModel(github_id=github_id, github_username=github_username))
    await db_session.commit()
