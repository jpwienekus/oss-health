import asyncio
from contextlib import ExitStack

import pytest
import pytest_asyncio
from sqlalchemy import Connection, text
from sqlalchemy.ext.asyncio import AsyncConnection, AsyncSession

from alembic.config import Config
from alembic.migration import MigrationContext
from alembic.operations import Operations
from alembic.script import ScriptDirectory
from config.settings import settings
from api.main import app as actual_app
from core.database import Base, DatabaseSessionManager, get_db_session
from core.models.dependency import Dependency as DependencyDBModel
from core.models.repository import Repository as RepositoryDBModel
from core.models.user import User as UserDBModel
from core.models.version import Version as VersionDBModel


@pytest.fixture(autouse=False)
def app():
    with ExitStack():
        yield actual_app


@pytest.fixture(scope="session")
def event_loop():
    loop = asyncio.new_event_loop()
    yield loop
    loop.close()


def run_migrations(connection: Connection):
    config = Config("./alembic.ini")
    config.set_main_option("sqlalchemy.url", settings.database_url)
    script = ScriptDirectory.from_config(config)

    def upgrade(rev, _):
        return script._upgrade_revs("head", rev)

    context = MigrationContext.configure(
        connection, opts={"target_metadata": Base.metadata, "fn": upgrade}
    )

    with context.begin_transaction():
        with Operations.context(context):
            context.run_migrations()


async def truncate_all_tables(connection: AsyncConnection):
    await connection.execute(text("COMMIT"))
    await connection.execute(text("SET session_replication_role = replica;"))

    for table in reversed(Base.metadata.sorted_tables):
        await connection.execute(
            text(f'TRUNCATE TABLE "{table.name}" RESTART IDENTITY CASCADE;')
        )

    await connection.execute(text("SET session_replication_role = origin;"))


@pytest_asyncio.fixture(scope="function")
async def db_session(app):
    test_sessionmanager = DatabaseSessionManager(
        settings.database_url, {"echo": settings.echo_sql}
    )

    async with test_sessionmanager.connect() as connection:
        await connection.run_sync(run_migrations)

    async def get_test_db_session():
        async with test_sessionmanager.session() as session:
            yield session

    app.dependency_overrides[get_db_session] = get_test_db_session

    async with test_sessionmanager.session() as session:
        try:
            await session.begin()
            yield session
        finally:
            await session.rollback()
            conn = await session.connection()
            await truncate_all_tables(conn)
            await test_sessionmanager.close()


@pytest_asyncio.fixture
async def test_user(db_session: AsyncSession) -> UserDBModel:
    user = UserDBModel(
        github_id=1,
        github_username="test_github_username",
        access_token="test_access_token",
    )
    db_session.add(user)
    await db_session.flush()
    return user


@pytest_asyncio.fixture
async def test_repository(
    db_session: AsyncSession, test_user: UserDBModel
) -> RepositoryDBModel:
    repo = RepositoryDBModel(
        github_id=12345,
        user_id=test_user.id,
        clone_url="https://github.com/example/repo",
    )
    db_session.add(repo)
    await db_session.flush()
    return repo


@pytest_asyncio.fixture
async def test_version(
    db_session: AsyncSession, version_str: str = "1.0.0"
) -> VersionDBModel:
    dependency = DependencyDBModel(name="example-package", ecosystem="pypi")
    db_session.add(dependency)
    await db_session.flush()

    version = VersionDBModel(version=version_str, dependency_id=dependency.id)
    db_session.add(version)
    await db_session.flush()

    return version
