from sqlalchemy import Connection, text
from sqlalchemy.ext.asyncio import AsyncConnection
from app.configuration import settings

import asyncio
import pytest_asyncio

import pytest
from alembic.config import Config
from alembic.migration import MigrationContext
from alembic.operations import Operations
from alembic.script import ScriptDirectory
from app.database import Base, get_db_session, sessionmanager


@pytest.fixture(scope="session")
def event_loop(request):
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()

def run_migrations(connection: Connection):
    config = Config("./alembic.ini")
    config.set_main_option("sqlalchemy.url", settings.database_url)
    script = ScriptDirectory.from_config(config)

    def upgrade(rev, context):
        return script._upgrade_revs("head", rev)

    context = MigrationContext.configure(connection, opts={"target_metadata": Base.metadata, "fn": upgrade})

    with context.begin_transaction():
        with Operations.context(context):
            context.run_migrations()

# @pytest.fixture(scope="session", autouse=True)
@pytest_asyncio.fixture(scope="function", autouse=False)
async def setup_database():
    # Run alembic migrations on test DB
    async with sessionmanager.connect() as connection:
        await connection.run_sync(run_migrations)

    yield

    # Teardown
    await sessionmanager.close()

async def truncate_all_tables(connection: AsyncConnection):
    await connection.execute(text("COMMIT"))  # Ensure any open txn is committed
    await connection.execute(text("SET session_replication_role = replica;"))

    for table in reversed(Base.metadata.sorted_tables):
        await connection.execute(text(f'TRUNCATE TABLE "{table.name}" RESTART IDENTITY CASCADE;'))

    await connection.execute(text("SET session_replication_role = origin;"))

@pytest_asyncio.fixture(scope="function")
async def transactional_session():
    async with sessionmanager.session() as session:
        try:
            await session.begin()
            yield session
        finally:
            await session.rollback()
            connection = await session.connection()
            await truncate_all_tables(connection)

@pytest_asyncio.fixture(scope="function")
async def db_session(transactional_session):
    yield transactional_session

@pytest_asyncio.fixture(scope="function", autouse=False)
async def session_override(app, db_session):
    async def get_db_session_override():
        yield db_session

    app.dependency_overrides[get_db_session] = get_db_session_override
