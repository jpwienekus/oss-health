import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.crud.repository_dependency_version import (
    replace_repository_dependency_versions,
)
from app.models.dependency import Dependency as DependencyDBModel
from app.models.relationships import (
    RepositoryDependencyVersion as RepositoryDependencyVersionDBModel,
)
from app.models.repository import Repository as RepositoryDBModel
from app.models.version import Version as VersionDBModel


@pytest.mark.asyncio
async def test_replace_repository_dependency_versions_inserts_new(
    db_session: AsyncSession, test_repository: RepositoryDBModel
):
    dependency_version_pairs = [("requests", "2.31.0", "pypi")]

    result = await replace_repository_dependency_versions(
        db_session=db_session,
        repository_id=test_repository.id,
        dependency_version_pairs=dependency_version_pairs,
    )

    # 1 dependency, 1 version, 1 relationship
    assert len(result) == 1
    _, name, version_str, ecosystem = result[0]
    assert name == "requests"
    assert version_str == "2.31.0"
    assert ecosystem == "pypi"

    # Check dependency exists
    dependency_result = await db_session.execute(
        select(DependencyDBModel).where(
            DependencyDBModel.name == "requests", DependencyDBModel.ecosystem == "pypi"
        )
    )
    dependency = dependency_result.scalar_one_or_none()
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
    removed_links_result = await db_session.execute(
        select(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == test_repository.id
        )
    )
    links = removed_links_result.scalars().all()
    assert len(links) == 1
    assert links[0].dependency_id == dependency.id
    assert links[0].version_id == version.id


@pytest.mark.asyncio
async def test_replace_repository_dependency_versions_replaces_old_links(
    db_session: AsyncSession, test_repository: RepositoryDBModel
):
    # Setup first pair
    first_pair = [("flask", "2.0.0", "pypi")]
    await replace_repository_dependency_versions(
        db_session, test_repository.id, first_pair
    )

    # Now replace with new pair
    second_pair = [("fastapi", "0.95.0", "pypi")]
    await replace_repository_dependency_versions(
        db_session, test_repository.id, second_pair
    )

    # Should only have one new link
    removed_links_result = await db_session.execute(
        select(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == test_repository.id
        )
    )
    links = removed_links_result.scalars().all()
    assert len(links) == 1

    # Check it's fastapi not flask
    dependency_result = await db_session.execute(
        select(DependencyDBModel).where(DependencyDBModel.id == links[0].dependency_id)
    )
    dependency = dependency_result.scalar_one()
    assert dependency.name == "fastapi"


@pytest.mark.asyncio
# ruff: noqa: E501
async def test_replace_repository_dependency_versions_skips_existing_dependency_and_version(
    db_session: AsyncSession, test_repository: RepositoryDBModel
):
    # Create dep/version manually
    dependency = DependencyDBModel(name="existing", ecosystem="pypi")
    db_session.add(dependency)
    await db_session.flush()

    version = VersionDBModel(version="1.0.0", dependency_id=dependency.id)
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
    dependency_count = await db_session.execute(
        select(DependencyDBModel).where(DependencyDBModel.name == "existing")
    )
    assert len(dependency_count.scalars().all()) == 1

    version_count = await db_session.execute(
        select(VersionDBModel).where(
            VersionDBModel.version == "1.0.0",
            VersionDBModel.dependency_id == dependency.id,
        )
    )
    assert len(version_count.scalars().all()) == 1


@pytest.mark.asyncio
async def test_replace_repository_dependency_versions_handles_empty_list(
    db_session: AsyncSession, test_repository: RepositoryDBModel
):
    # Pre-insert link
    await replace_repository_dependency_versions(
        db_session, test_repository.id, [("something", "0.1", "pypi")]
    )

    # Now replace with empty
    result = await replace_repository_dependency_versions(
        db_session, test_repository.id, []
    )

    # Should remove previous links
    removed_links_result = await db_session.execute(
        select(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == test_repository.id
        )
    )
    links = removed_links_result.scalars().all()
    assert len(links) == 0
    assert result == []
