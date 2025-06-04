from datetime import datetime
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import User as UserDBModel


async def get_user_by_github_id(db_session: AsyncSession, github_id: int):
    return (
        await db_session.scalars(
            select(UserDBModel).where(UserDBModel.github_id == github_id)
        )
    ).first()


async def get_user(db_session: AsyncSession, user_id: int):
    return (
        await db_session.scalars(select(UserDBModel).where(UserDBModel.id == user_id))
    ).first()


async def add_user(
    db_session: AsyncSession, github_id: int, github_username: str, access_token: str
):
    user = UserDBModel(
        github_id=github_id, github_username=github_username, access_token=access_token
    )
    db_session.add(user)
    await db_session.commit()

    return user


async def get_sync_time(db_session: AsyncSession, user_id: int):
    user = (
        await db_session.scalars(select(UserDBModel).where(UserDBModel.id == user_id))
    ).first()

    if user:
        return user.synced_at

    return None


async def get_access_token(db_session: AsyncSession, user_id: int):
    user = (
        await db_session.scalars(select(UserDBModel).where(UserDBModel.id == user_id))
    ).first()

    if user:
        return user.access_token

    return None


async def update_sync_time(db_session: AsyncSession, user_id: int):
    user = (
        await db_session.scalars(select(UserDBModel).where(UserDBModel.id == user_id))
    ).first()

    if user:
        now = datetime.now()
        user.synced_at = now
        await db_session.commit()
        return now


async def update_access_token(db_session: AsyncSession, github_id, access_token: str):
    user = (
        await db_session.scalars(select(UserDBModel).where(UserDBModel.github_id == github_id))
    ).first()

    if user:
        user.access_token = access_token
        await db_session.commit()
