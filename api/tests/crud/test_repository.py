import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.crud.user import add_user
from app.models import Repository as RepositoryDBModel
from app.crud.repository import get_repositories, add_repository_ids

@pytest.mark.asyncio
async def test_add_repository_ids(db_session: AsyncSession):
    repos_to_add = [
        {"id": 1111, "clone_url": "https://github.com/test/repo1"},
        {"id": 2222, "clone_url": "https://github.com/test/repo2"},
    ]

    user = await add_user(db_session=db_session, github_id=1, github_username="test_github_username", access_token="test_access_token")
    await add_repository_ids(db_session, user.id, repos_to_add)

    # Directly query to confirm they were added
    result = await db_session.execute(
        select(RepositoryDBModel).where(RepositoryDBModel.user_id == user.id)
    )
    repos = result.scalars().all()

    assert len(repos) == 2
    assert {repo.github_id for repo in repos} == {1111, 2222}


@pytest.mark.asyncio
async def test_get_repositories(db_session: AsyncSession):
    user = await add_user(db_session=db_session, github_id=1, github_username="test_github_username", access_token="test_access_token")
    # Add test data manually
    repo = RepositoryDBModel(
        github_id=3333,
        user_id=user.id,
        score=0,
        clone_url="https://github.com/test/repo3"
    )
    db_session.add(repo)
    await db_session.commit()

    # Now use the function to fetch it
    repos = await get_repositories(db_session, user.id)

    assert len(repos) == 1
    assert repos[0].github_id == 3333
    assert repos[0].clone_url == "https://github.com/test/repo3"
