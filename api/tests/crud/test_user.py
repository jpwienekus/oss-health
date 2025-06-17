import pytest
from sqlalchemy.ext.asyncio import AsyncSession

from app.crud.user import (
    add_user,
    get_access_token,
    get_user,
    get_user_by_github_id,
    update_access_token,
)


@pytest.mark.asyncio
async def test_add_and_get_user(db_session: AsyncSession):
    user = await add_user(
        db_session=db_session,
        github_id=42,
        github_username="testuser42",
        access_token="secret_token_42",
    )

    # Fetch by user ID
    fetched_user = await get_user(db_session, user.id)

    assert fetched_user is not None
    assert fetched_user.id == user.id
    assert fetched_user.github_id == 42
    assert fetched_user.github_username == "testuser42"


@pytest.mark.asyncio
async def test_get_user_by_github_id(db_session: AsyncSession):
    await add_user(
        db_session=db_session,
        github_id=12345,
        github_username="octocat",
        access_token="token123",
    )

    fetched_user = await get_user_by_github_id(db_session, 12345)

    assert fetched_user is not None
    assert fetched_user.github_username == "octocat"
    assert fetched_user.github_id == 12345


@pytest.mark.asyncio
async def test_get_access_token(db_session: AsyncSession):
    user = await add_user(
        db_session=db_session,
        github_id=99,
        github_username="secrettest",
        access_token="access-token-99",
    )

    token = await get_access_token(db_session, user.id)

    assert token == "access-token-99"


@pytest.mark.asyncio
async def test_get_access_token_user_not_found(db_session: AsyncSession):
    token = await get_access_token(db_session, user_id=9999)
    assert token is None


@pytest.mark.asyncio
async def test_update_access_token(db_session: AsyncSession):
    await add_user(
        db_session=db_session,
        github_id=555,
        github_username="update-test",
        access_token="old-token",
    )

    await update_access_token(db_session, github_id=555, access_token="new-token")

    updated_user = await get_user_by_github_id(db_session, 555)
    assert updated_user
    assert updated_user.access_token == "new-token"


@pytest.mark.asyncio
async def test_update_access_token_user_not_found(db_session: AsyncSession):
    await update_access_token(db_session, github_id=9999, access_token="noop-token")

    user = await get_user_by_github_id(db_session, 9999)
    assert user is None
