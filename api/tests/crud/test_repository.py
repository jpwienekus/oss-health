import pytest
from app.crud.user import add_user
from app.models import Repository as RepositoryDBModel
from app.crud.repository import get_repositories, add_repository_ids  # your module/functions

@pytest.mark.asyncio
async def test_add_repository_ids_and_get(db_session):
    user_id = 1
    tracked_repositories = [
        {"id": 101, "clone_url": "https://github.com/user/repo1.git"},
        {"id": 102, "clone_url": "https://github.com/user/repo2.git"},
    ]

    await add_user(db_session=db_session, github_id=1, github_username="test_github_username", access_token="test_access_token")
    # Add repos
    await add_repository_ids(db_session=db_session, user_id=user_id, tracked_repositories=tracked_repositories)

    # Fetch repos back
    repos = await get_repositories(db_session, user_id)

    assert len(repos) == 2
    assert repos[0].github_id == 101
    assert repos[0].clone_url == "https://github.com/user/repo1.git"
    assert repos[0].user_id == user_id
    assert repos[1].github_id == 102
    assert repos[1].clone_url == "https://github.com/user/repo2.git"
    assert repos[1].user_id == user_id
