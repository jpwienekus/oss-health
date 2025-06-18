import strawberry
from fastapi import Depends, Request
from sqlalchemy.ext.asyncio import AsyncSession
from strawberry.fastapi import GraphQLRouter

from api.graphql.resolvers import Mutation, Query
from core.database import get_db_session


async def get_context(
    request: Request, db_session: AsyncSession = Depends(get_db_session)
):
    return {"request": request, "db": db_session}


schema = strawberry.Schema(Query, Mutation)
graphql_app = GraphQLRouter(schema, context_getter=get_context)
