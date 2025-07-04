from operator import or_
from sqlalchemy import and_, asc, desc, select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.sql import operators
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

    # if filter.name:
    #     filters.append(DependencyDBModel.name.ilike(f"%{filter.name}"))
    #
    # if filter.ecosystem:
    #     filters.append(DependencyDBModel.ecosystem.ilike(f"%{filter.ecosystem}"))
    #
    # if filter.github_url_resolve_failed is not None:
    #     filters.append(DependencyDBModel.github_url_resolve_failed == filter.github_url_resolve_failed)
    #
    # if filters:
    #     query = query.where(and_(*filters))
    # op_str = '>'
    # after: int | None = 761
    after: int | None = None
    # before: int | None = None
    before: int | None = 10


    id = 761
    op = operators.gt
    sortOrder = asc

    if after is not None:
        id = after
        op = operators.gt
        sortOrder = asc
    elif before is not None:
        id = before
        op = operators.lt
        sortOrder = desc

    ecosystem = 'npm'
    # op = getattr(operators, op_str)
    limit = 10

    query = (
        select(DependencyDBModel)
        .where(
            or_(
                op(DependencyDBModel.ecosystem, ecosystem),
                and_(
                    DependencyDBModel.ecosystem == ecosystem,
                    op(DependencyDBModel.id, id)
                )
            )
        )
        .order_by(sortOrder(DependencyDBModel.ecosystem), sortOrder(DependencyDBModel.id))
        .limit(limit + 2)
    )

    result = (await db_session.execute(query)).scalars().all()
    print('=======')

    items = result[:limit]


    has_next_page = False
    has_previous_page = False

    if before is not None: # Handling previous page means there is a next
        has_next_page = True
    elif after is not None and len(result) > limit:
        has_next_page = True

    if after is not None: # Handling next page means there is a previous
        has_previous_page = True
    elif before is not None and len(result) > limit:
        has_previous_page = True

    # first page 
    if before is None and after is None:
        has_next_page = len(result > limit)
        has_previous_page = False


    if before:
        items = list(reversed(items))


    for item in items:
        print(item.id, item.name)

    print("#########")
    print(has_next_page, has_previous_page)
        


    return {}
    # return await paginate_query(
    #     session=db_session,
    #     base_query=query,
    #     model=None,
    #     sort_column=getattr(DependencyDBModel, sort.field.value),
    #     limit=pagination.limit,
    #     before=pagination.before,
    #     after=pagination.after,
    #     descending=sort.direction == SortDirection.DESC
    # )

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

