import pytest
from sqlalchemy import insert
from sqlalchemy.ext.asyncio import AsyncSession

from core.crud.repository import add_repository_ids, get_cron_info, get_repositories
from core.models.repository import Repository as RepositoryDBModel
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


@pytest.mark.asyncio
async def test_get_cron_info(db_session: AsyncSession, test_user: UserDBModel):
    data = [
        {
            "user_id": test_user.id,
            "github_id": 1,
            "url": "https://github.com/a/repo1",
            "scan_day": 1,
            "scan_hour": 10,
        },  # Monday 10h
        {
            "user_id": test_user.id,
            "github_id": 2,
            "url": "https://github.com/a/repo2",
            "scan_day": 1,
            "scan_hour": 10,
        },  # Monday 10h
        {
            "user_id": test_user.id,
            "github_id": 3,
            "url": "https://github.com/a/repo3",
            "scan_day": 2,
            "scan_hour": 14,
        },  # Tuesday 14h
    ]

    await db_session.execute(insert(RepositoryDBModel), data)
    await db_session.commit()

    result = await get_cron_info(db_session)

    # Sort result for deterministic test validation
    result = sorted(result, key=lambda r: (r.scan_day, r.scan_hour))

    assert len(result) == 2

    assert result[0].scan_day == 1
    assert result[0].scan_hour == 10
    assert result[0].total == 2

    assert result[1].scan_day == 2
    assert result[1].scan_hour == 14
    assert result[1].total == 1
