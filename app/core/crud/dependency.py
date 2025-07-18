import math
from typing import Sequence, Tuple

from sqlalchemy import asc, desc, func, select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from api.graphql.inputs import (
    DependencyFilter,
    DependencySortInput,
    PaginationInput,
    SortDirection,
)
from core.models import Dependency as DependencyDBModel


async def get_dependencies_paginated(
    db_session: AsyncSession,
    filter: DependencyFilter,
    sort: DependencySortInput,
    pagination: PaginationInput,
) -> Tuple[int, int, int, int, Sequence[DependencyDBModel]]:
    filters = []

    if filter.name:
        safe_search = filter.name.replace("%", r"\%").replace("_", r"\_")
        filters.append(DependencyDBModel.name.ilike(f"%{safe_search}%"))

    if filter.statuses:
        filters.append(DependencyDBModel.scan_status.in_(filter.statuses))

    offset = (pagination.page - 1) * pagination.page_size
    sort_direction = asc if sort.direction == SortDirection.ASC else desc

    query = (
        select(DependencyDBModel)
        .options(selectinload(DependencyDBModel.dependency_repository))
        .order_by(
            sort_direction(getattr(DependencyDBModel, sort.field.value)),
            sort_direction(DependencyDBModel.name),
        )
    )

    if filters:
        query = query.where(*filters)

    query = query.offset(offset).limit(pagination.page_size)

    status_count_query = select(
        DependencyDBModel.scan_status, func.count().label("count")
    ).group_by(DependencyDBModel.scan_status)

    if filters:
        status_count_query = status_count_query.where(*filters)

    results = await db_session.execute(query)

    status_counts_raw = await db_session.execute(status_count_query)
    status_counts = {status: count for status, count in status_counts_raw.all()}

    int(status_counts.get("completed", 0))
    int(status_counts.get("pending", 0))
    int(status_counts.get("failed", 0))
    calc_total = sum(status_counts.values())

    return (
        math.ceil(calc_total / pagination.page_size),
        int(status_counts.get("completed", 0)),
        int(status_counts.get("pending", 0)),
        int(status_counts.get("failed", 0)),
        results.scalars().all(),
    )
