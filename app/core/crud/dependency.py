from typing import Sequence, Tuple
from sqlalchemy import  asc, desc, func, select
from sqlalchemy.ext.asyncio import AsyncSession
from api.graphql.inputs import DependencyFilter, DependencySortInput, PaginationInput, SortDirection
from core.models import Dependency as DependencyDBModel
import math

async def get_dependencies_paginated(
    db_session: AsyncSession,
    filter: DependencyFilter,
    sort: DependencySortInput,
    pagination: PaginationInput
) -> Tuple[int, Sequence[DependencyDBModel]]:
    filters = []

    if filter.name:
        safe_search = filter.name.replace('%', r'\%').replace('_', r'\_')
        filters.append(DependencyDBModel.name.ilike(f"%{safe_search}%"))

    offset = (pagination.page - 1) * pagination.page_size
    sort_direction = asc if sort.direction == SortDirection.ASC else desc


    base_query = select(DependencyDBModel).order_by(sort_direction(getattr(DependencyDBModel, sort.field.value)))

    if filters:
        base_query = base_query.where(*filters)

    count_query = select(func.count()).select_from(base_query.subquery())
    query = base_query.offset(offset).limit(pagination.page_size)

    total = (await db_session.execute(count_query)).scalar_one()
    results = await db_session.execute(query)

    return (math.ceil(total/pagination.page_size), results.scalars().all())
