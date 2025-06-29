import pytest
from sqlalchemy.ext.asyncio import AsyncSession

from core.crud.repository import add_repository_ids, get_repositories
from core.models.user import User as UserDBModel


@pytest.mark.asyncio
async def test_add_repository_ids(db_session: AsyncSession, test_user: UserDBModel):
    repos_to_add = [
        {"id": 1111, "url": "https://github.com/test/repo1"},
        {"id": 2222, "url": "https://github.com/test/repo2"},
    ]

    await add_repository_ids(db_session, test_user.id, repos_to_add)
    repos = await get_repositories(db_session, test_user.id)

    assert len(repos) == 2
    assert repos[0].github_id == 1111
    assert repos[0].url == "https://github.com/test/repo1"
    assert repos[1].github_id == 2222
    assert repos[1].url == "https://github.com/test/repo2"
