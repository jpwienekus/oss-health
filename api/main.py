import logging
import sys
import strawberry
from configuration import settings
from contextlib import asynccontextmanager
from fastapi import FastAPI

# from database import sessionmanager
from strawberry.fastapi import GraphQLRouter

logging.basicConfig(
    stream=sys.stdout,
    level=logging.DEBUG if settings.log_level == "DEBUG" else logging.INFO,
)


# https://github.com/ThomasAitken/demo-fastapi-async-sqlalchemy/blob/main/backend/requirements.txt
@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Function that handles startup and shutdown events.
    To understand more, read https://fastapi.tiangolo.com/advanced/events/
    """
    yield
    # if sessionmanager._engine is not None:
    #     # Close the DB connection
    #     await sessionmanager.close()


@strawberry.type
class Query:
    @strawberry.field
    def hello(self) -> str:
        return "Hello World"


app = FastAPI(lifespan=lifespan, title=settings.project_name, docs_url="/api/docs")


@app.get("/")
async def root():
    return {"message": "Hello World"}


schema = strawberry.Schema(Query)
graphql_app = GraphQLRouter(schema)

app.include_router(graphql_app, prefix="/graphql")
