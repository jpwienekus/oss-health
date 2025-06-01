from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import User as UserDBModel


async def get_user(db_session: AsyncSession, github_id: int):
    return (await db_session.scalars(select(UserDBModel).where(UserDBModel.github_id == github_id))).first()


async def add_user(db_session: AsyncSession, github_id: int, github_username: str):
    user = UserDBModel(github_id=github_id, github_username=github_username)
    db_session.add(user)
    await db_session.commit()

    return user
