import base64
from typing import Optional

from sqlalchemy import Select, asc, desc
from sqlalchemy.ext.asyncio import AsyncSession


def encode_cursor(value: str | int) -> str:
    return base64.urlsafe_b64encode(str(value).encode()).decode()


def decode_cursor(cursor: str) -> str:
    return base64.urlsafe_b64decode(cursor.encode()).decode()


def coerce_cursor_value(sort_column, value: str):
    python_type = sort_column.type.python_type
    return python_type(value)


async def paginate_query(
    session: AsyncSession,
    base_query: Select,
    model,  # ???
    sort_column,
    limit: int = 20,
    after: Optional[str] = None,
    before: Optional[str] = None,
    descending: bool = True,
) -> dict:
    after_value = (
        coerce_cursor_value(sort_column, decode_cursor(after)) if after else None
    )
    before_value = (
        coerce_cursor_value(sort_column, decode_cursor(before)) if before else None
    )
    actual_sort_desc = descending ^ bool(before)  # XOR: invert if before

    if after_value and not before_value:
        filter_clause = (
            sort_column < after_value if descending else sort_column > after_value
        )
    elif before_value and not after_value:
        filter_clause = (
            sort_column < before_value
            if actual_sort_desc
            else sort_column > before_value
        )
    elif after_value and before_value:
        raise ValueError("Cannot use both 'after' and 'before' at the same time.")
    else:
        filter_clause = None

    query = base_query.order_by(
        desc(sort_column) if actual_sort_desc else asc(sort_column)
    )

    if filter_clause is not None:
        query = query.filter(filter_clause)

    # Fetch one extra to determine if next/prev page exists
    result = await session.execute(query.limit(limit + 1))
    result = result.scalars().all()

    print(len(result))
    has_next_page = len(result) > limit
    # conservative estimate
    has_previous_page = bool(after or before)

    items = result[:limit]

    if before_value:
        items = list(reversed(items))

    paginated_items = items[:limit]
    start_cursor = (
        encode_cursor(getattr(paginated_items[0], sort_column.key))
        if paginated_items
        else None
    )
    end_cursor = (
        encode_cursor(getattr(paginated_items[-1], sort_column.key))
        if paginated_items
        else None
    )

    return {
        "edges": paginated_items,
        "page_info": {
            "has_next_page": has_next_page,
            "has_previous_page": has_previous_page,
            "start_cursor": start_cursor,
            "end_cursor": end_cursor,
        },
    }


async def has_records_before(
    session: AsyncSession,
    base_query: Select,
    sort_column,
    value,
    actual_sort_desc: bool,
) -> bool:
    filter_clause = sort_column > value if actual_sort_desc else sort_column < value

    test_query = base_query.filter(filter_clause).limit(1)
    result = await session.execute(test_query)
    return result.scalars().first() is not None
