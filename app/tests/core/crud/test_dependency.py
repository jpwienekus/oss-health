import pytest
from sqlalchemy import insert
from sqlalchemy.ext.asyncio import AsyncSession

from core.models import Dependency as DependencyDBModel
from api.graphql.inputs import DependencyFilter, DependencySortField, DependencySortInput, PaginationInput, SortDirection

from core.crud.dependency import get_dependencies_paginated


@pytest.mark.asyncio
async def test_get_dependencies_paginated_basic(db_session: AsyncSession):
    data = [
        {"name": "dep-a", "scan_status": "completed"},
        {"name": "dep-b", "scan_status": "pending"},
        {"name": "dep-c", "scan_status": "failed"},
        {"name": "dep-d", "scan_status": "completed"},
    ]
    await db_session.execute(insert(DependencyDBModel), data)
    await db_session.commit()

    filter = DependencyFilter(name="", statuses=[])
    sort = DependencySortInput(field=DependencySortField.NAME, direction=SortDirection.ASC)
    pagination = PaginationInput(page=1, page_size=2)

    total_pages, completed, pending, failed, deps = await get_dependencies_paginated(
        db_session, filter, sort, pagination
    )

    assert total_pages == 2
    assert completed == 2
    assert pending == 1
    assert failed == 1
    assert len(deps) == 2
    assert deps[0].name == "dep-a"
    assert deps[1].name == "dep-b"


@pytest.mark.asyncio
async def test_get_dependencies_paginated_name_filter(db_session: AsyncSession):
    await db_session.execute(insert(DependencyDBModel), [
        {"name": "lib_sqlite", "scan_status": "completed"},
        {"name": "sqlite_parser", "scan_status": "pending"},
        {"name": "libcrypto", "scan_status": "failed"},
    ])
    await db_session.commit()

    filter = DependencyFilter(name="sqlite", statuses=None)
    sort = DependencySortInput(field=DependencySortField.NAME, direction=SortDirection.ASC)
    pagination = PaginationInput(page=1, page_size=10)

    _, _, _, _, deps = await get_dependencies_paginated(db_session, filter, sort, pagination)

    assert len(deps) == 2
    assert all("sqlite" in d.name for d in deps)


@pytest.mark.asyncio
async def test_get_dependencies_paginated_status_filter(db_session: AsyncSession):
    await db_session.execute(insert(DependencyDBModel), [
        {"name": "dep1", "scan_status": "completed"},
        {"name": "dep2", "scan_status": "pending"},
        {"name": "dep3", "scan_status": "failed"},
    ])
    await db_session.commit()

    filter = DependencyFilter(name=None, statuses=["pending", "failed"])
    sort = DependencySortInput(field=DependencySortField.NAME, direction=SortDirection.ASC)
    pagination = PaginationInput(page=1, page_size=10)

    _, completed, pending, failed, deps = await get_dependencies_paginated(
        db_session, filter, sort, pagination
    )

    assert completed == 0
    assert pending == 1
    assert failed == 1
    assert len(deps) == 2
    assert all(d.scan_status in ["pending", "failed"] for d in deps)


@pytest.mark.asyncio
async def test_get_dependencies_paginated_pagination_math(db_session: AsyncSession):
    await db_session.execute(insert(DependencyDBModel), [
        {"name": f"dep-{i}", "scan_status": "completed"} for i in range(5)
    ])
    await db_session.commit()

    filter = DependencyFilter(name=None, statuses=None)
    sort = DependencySortInput(field=DependencySortField.NAME, direction=SortDirection.ASC)
    pagination = PaginationInput(page=2, page_size=2)

    total_pages, completed, pending, failed, deps = await get_dependencies_paginated(
        db_session, filter, sort, pagination
    )

    assert total_pages == 3  # 5 items with 2 per page
    assert completed == 5
    assert pending == 0
    assert failed == 0
    assert len(deps) == 2
