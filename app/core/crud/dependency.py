from sqlalchemy import and_, asc, desc, select
from sqlalchemy.ext.asyncio import AsyncSession
from api.graphql.inputs import DependencyFilter, DependencySortInput, PaginationInput
from api.graphql.types import DependencyConnection, DependencyEdge, DependencyType, PageInfo
from core.models import Dependency as DependencyDBModel

async def get_dependencies_paginated(
    db_session: AsyncSession,
    filter: DependencyFilter,
    sort: DependencySortInput,
    pagination: PaginationInput
):
    query = select(DependencyDBModel)
    filters = []

    if filter.name:
        filters.append(DependencyDBModel.name.ilike(f"%{filter.name}"))

    if filter.ecosystem:
        filters.append(DependencyDBModel.ecosystem.ilike(f"%{filter.ecosystem}"))

    if filter.github_url_resolve_failed is not None:
        filters.append(DependencyDBModel.github_url_resolve_failed == filter.github_url_resolve_failed)

    sort_column = getattr(DependencyDBModel, sort.field.value)
    sort_order = asc(sort_column) if sort.direction == "ASC" else desc(sort_column)

    if pagination.after is not None:
        op = ">" if sort.direction == "ASC" else "<"
        filters.append(getattr(DependencyDBModel, sort.field.value).op(op)(pagination.after))

    if filters:
        query = query.where(and_(*filters))

    query = query.order_by(sort_order).limit(pagination.limit + 1) # one extra to check if next row exists

    result = await db_session.execute(query)
    dependencies = result.scalars().all()

    has_next_page = len(dependencies) > pagination.limit
    items = dependencies[:pagination.limit]

    edges = [
        DependencyEdge(
            node=DependencyType.from_model(dependency),
            cursor=getattr(dependency, sort.field.value)
        ) for dependency in items
    ]

    start_cursor = getattr(items[0], sort.field.value) if items else None
    end_cursor = getattr(items[-1], sort.field.value) if items else None

    has_previous_page = False

    if pagination.after is not None and start_cursor is not None:
        has_previous_result = await db_session.execute(
            (
                select(DependencyDBModel).where(
                    and_(*filters, sort_column < start_cursor if sort.direction == "ASC" else sort_column > start_cursor)
                ).order_by(sort_order).limit(1)
            )
        )
        has_previous_page = bool(has_previous_result.scalars())


    return DependencyConnection(
        edges=edges,
        page_info=PageInfo(
            has_next_page=has_next_page,
            has_previous_page=has_previous_page,
            start_cursor=start_cursor,
            end_cursor=end_cursor
        )
    )

