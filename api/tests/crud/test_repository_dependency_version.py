import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.user import User as UserDBModel
from app.models.repository import Repository as RepositoryDBModel
from app.models.dependency import Dependency as DependencyDBModel
from app.models.version import Version as VersionDBModel
from app.models.relationships import RepositoryDependencyVersion as RepositoryDependencyVersionDBModel
from app.crud.repository_dependency_version import replace_repository_dependency_versions

async def insert_test_user(db_session: AsyncSession) -> UserDBModel: 
    user = UserDBModel(github_id=1, github_username="test_github_username", access_token="test_access_token")
    db_session.add(user)
    await db_session.flush()
    return user

async def insert_test_repository(db_session: AsyncSession, user_id: int) -> RepositoryDBModel:
    repo = RepositoryDBModel(github_id=12345, user_id=user_id, clone_url="https://github.com/example/repo")
    db_session.add(repo)
    await db_session.flush()
    return repo


@pytest.mark.asyncio

async def test_replace_repository_dependency_versions_inserts_new(db_session: AsyncSession):
    test_user= await insert_test_user(db_session)
    test_repository = await insert_test_repository(db_session, test_user.id)
    dep_version_pairs = [("requests", "2.31.0", "pypi")]

    result = await replace_repository_dependency_versions(
        db_session=db_session,
        repository_id=test_repository.id,
        dep_version_pairs=dep_version_pairs,
    )

    # 1 dependency, 1 version, 1 relationship
    assert len(result) == 1
    version_id, name, version_str, ecosystem = result[0]
    assert name == "requests"
    assert version_str == "2.31.0"
    assert ecosystem == "pypi"

    # Check dependency exists
    dep_result = await db_session.execute(
        select(DependencyDBModel).where(
            DependencyDBModel.name == "requests", DependencyDBModel.ecosystem == "pypi"
        )
    )
    dependency = dep_result.scalar_one_or_none()
    assert dependency is not None

    # Check version exists
    version_result = await db_session.execute(
        select(VersionDBModel).where(
            VersionDBModel.version == "2.31.0",
            VersionDBModel.dependency_id == dependency.id,
        )
    )
    version = version_result.scalar_one_or_none()
    assert version is not None

    # Check relationship exists
    rel_result = await db_session.execute(
        select(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == test_repository.id
        )
    )
    links = rel_result.scalars().all()
    assert len(links) == 1
    assert links[0].dependency_id == dependency.id
    assert links[0].version_id == version.id



@pytest.mark.asyncio
async def test_replace_repository_dependency_versions_replaces_old_links(db_session: AsyncSession):
    test_user= await insert_test_user(db_session)
    test_repository = await insert_test_repository(db_session, test_user.id)

    # Setup first pair
    first_pair = [("flask", "2.0.0", "pypi")]
    await replace_repository_dependency_versions(db_session, test_repository.id, first_pair)

    # Now replace with new pair
    second_pair = [("fastapi", "0.95.0", "pypi")]
    await replace_repository_dependency_versions(db_session, test_repository.id, second_pair)

    # Should only have one new link
    rel_result = await db_session.execute(
        select(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == test_repository.id
        )
    )
    links = rel_result.scalars().all()
    assert len(links) == 1

    # Check it's fastapi not flask
    dep_result = await db_session.execute(
        select(DependencyDBModel).where(DependencyDBModel.id == links[0].dependency_id)
    )
    dep = dep_result.scalar_one()
    assert dep.name == "fastapi"


@pytest.mark.asyncio
async def test_replace_repository_dependency_versions_skips_existing_dependency_and_version(db_session: AsyncSession):
    test_user= await insert_test_user(db_session)
    test_repository = await insert_test_repository(db_session, test_user.id)

    # Create dep/version manually
    dep = DependencyDBModel(name="existing", ecosystem="pypi")
    db_session.add(dep)
    await db_session.flush()

    version = VersionDBModel(version="1.0.0", dependency_id=dep.id)
    db_session.add(version)
    await db_session.flush()

    # Run function using existing dep/version
    result = await replace_repository_dependency_versions(
        db_session,
        test_repository.id,
        [("existing", "1.0.0", "pypi")],
    )

    # Still links correctly
    assert len(result) == 1

    # Confirm no duplicate dependencies or versions created
    dep_count = await db_session.execute(
        select(DependencyDBModel).where(DependencyDBModel.name == "existing")
    )
    assert len(dep_count.scalars().all()) == 1

    version_count = await db_session.execute(
        select(VersionDBModel).where(VersionDBModel.version == "1.0.0", VersionDBModel.dependency_id == dep.id)
    )
    assert len(version_count.scalars().all()) == 1


@pytest.mark.asyncio
async def test_replace_repository_dependency_versions_handles_empty_list(db_session: AsyncSession):
    test_user= await insert_test_user(db_session)
    test_repository = await insert_test_repository(db_session, test_user.id)
    # Pre-insert link
    await replace_repository_dependency_versions(
        db_session, test_repository.id, [("something", "0.1", "pypi")]
    )

    # Now replace with empty
    result = await replace_repository_dependency_versions(
        db_session, test_repository.id, []
    )

    # Should remove previous links
    rel_result = await db_session.execute(
        select(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == test_repository.id
        )
    )
    links = rel_result.scalars().all()
    assert len(links) == 0
    assert result == []
