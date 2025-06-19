from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from api.auth.github_oauth import router as auth_router
from config.settings import settings
from api.graphql.schema import graphql_app
from core.database import sessionmanager
from core.utils.loggin import configure_logging

configure_logging(settings.log_level)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Function that handles startup and shutdown events.
    To understand more, read https://fastapi.tiangolo.com/advanced/events/
    """
    yield
    if sessionmanager._engine is not None:
        # Close the DB connection
        await sessionmanager.close()


origins = ["http://localhost:5173"]

app = FastAPI(lifespan=lifespan, title=settings.project_name)
app.include_router(auth_router)
app.include_router(graphql_app, prefix="/graphql")
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
