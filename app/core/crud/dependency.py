from sqlalchemy import and_, asc, desc, select
from sqlalchemy.ext.asyncio import AsyncSession
from api.graphql.inputs import DependencyFilter, DependencySortInput, PaginationInput, SortDirection
from api.graphql.types import DependencyConnection, DependencyEdge, DependencyType, PageInfo
from core.crud.utils import paginate_query
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

    if filters:
        query = query.where(and_(*filters))

    return await paginate_query(
        session=db_session,
        base_query=query,
        model=None,
        sort_column=getattr(DependencyDBModel, sort.field.value),
        limit=pagination.limit,
        before=pagination.before,
        after=pagination.after,
        descending=sort.direction == SortDirection.DESC
    )

    # order_clause = asc(sort_column) if sort.direction == SortDirection.ASC else desc(sort_column)
    # query.order_by(order_clause)

    # result = await db_session.execute(query)
    # dependencies = result.scalars().all()

    # has_next_page = len(dependencies) > pagination.limit
    # items = dependencies[:pagination.limit]

    # if is_reverse:
    #     items = list(reversed(items))

    # edges = [
    #     DependencyEdge(
    #         node=DependencyType.from_model(dependency),
    #         cursor=getattr(dependency, sort.field.value)
    #     ) for dependency in items
    # ]

    # start_cursor = getattr(items[0], sort.field.value) if items else None
    # end_cursor = getattr(items[-1], sort.field.value) if items else None

    # has_previous_page = False
    #
    # if start_cursor is not None:
    #     has_previous_query = select(DependencyDBModel).where(
    #                 and_(
    #                     *filters, 
    #                     sort_column < (start_cursor if direction == SortDirection.ASC else sort_column > start_cursor)
    #                 )
    #             ).limit(1)
    #     print("eeeeeeeeeeeeeeee")
    #     print(has_previous_query)
    #     print(start_cursor)
    #     has_previous_result = await db_session.execute(has_previous_query)
    #     has_previous_page = bool(has_previous_result.scalars().first())
    #     print("33333333333")
    #     print(has_previous_page)


    # if pagination.after is not None and start_cursor is not None:
    #     has_previous_result = await db_session.execute(
    #         (
    #             select(DependencyDBModel).where(
    #                 and_(*filters, sort_column < start_cursor if sort.direction == "ASC" else sort_column > start_cursor)
    #             ).order_by(sort_order).limit(1)
    #         )
    #     )
    #     has_previous_page = bool(has_previous_result.scalars())


    return DependencyConnection(
        edges=edges,
        page_info=PageInfo(
            has_next_page=has_next_page,
            has_previous_page=has_previous_page,
            start_cursor=start_cursor,
            end_cursor=end_cursor
        )
    )

