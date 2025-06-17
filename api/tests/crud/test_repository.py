import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.crud.user import add_user
from app.models import Repository as RepositoryDBModel
from app.crud.repository import get_repositories, add_repository_ids
from app.models.user import User as UserDBModel

async def insert_test_user(db_session: AsyncSession) -> UserDBModel: 
    user = UserDBModel(github_id=1, github_username="test_github_username", access_token="test_access_token")
    db_session.add(user)
    await db_session.flush()
    return user


@pytest.mark.asyncio
async def test_add_repository_ids(db_session: AsyncSession):
    repos_to_add = [
        {"id": 1111, "clone_url": "https://github.com/test/repo1"},
        {"id": 2222, "clone_url": "https://github.com/test/repo2"},
    ]

    test_user = await insert_test_user(db_session)
    await add_repository_ids(db_session, test_user.id, repos_to_add)

    repos = await get_repositories(db_session, test_user.id)

    assert len(repos) == 2
    assert repos[0].github_id == 1111
    assert repos[0].clone_url == "https://github.com/test/repo1"
    assert repos[1].github_id == 2222
    assert repos[1].clone_url == "https://github.com/test/repo2"

