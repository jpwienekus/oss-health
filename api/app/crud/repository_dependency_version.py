import logging
from typing import List

from sqlalchemy import delete, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.dependency import Dependency as DependencyDBModel
from app.models.relationships import (
    RepositoryDependencyVersion as RepositoryDependencyVersionDBModel,
)
from app.models.version import Version as VersionDBModel

logger = logging.getLogger()

async def replace_repository_dependency_versions(
    db_session: AsyncSession,
    repository_id: int,
    dependency_version_pairs: List[tuple[str, str, str]],  # (name, version_str, ecosystem)
):
    await db_session.execute(
        delete(RepositoryDependencyVersionDBModel).where(
            RepositoryDependencyVersionDBModel.repository_id == repository_id
        )
    )
    await db_session.flush()
    dependency_versions: List[tuple[int, str, str, str]] = []

    inserted_dependencies = 0
    existing_dependencies = 0
    inserted_versions = 0
    existing_versions = 0

    for name, version_str, ecosystem in dependency_version_pairs:
        dependency_result = await db_session.execute(
            select(DependencyDBModel).where(
                DependencyDBModel.name == name,
                DependencyDBModel.ecosystem == ecosystem,
            )
        )
        dependency = dependency_result.scalar_one_or_none()

        if not dependency:
            dependency = DependencyDBModel(name=name, ecosystem=ecosystem)
            db_session.add(dependency)
            await db_session.flush()
            inserted_dependencies += 1
        else:
            existing_dependencies += 1

        version_result = await db_session.execute(
            select(VersionDBModel).where(
                VersionDBModel.version == version_str,
                VersionDBModel.dependency_id == dependency.id
            )
        )
        version = version_result.scalar_one_or_none()

        if not version:
            version = VersionDBModel(version=version_str, dependency_id=dependency.id)
            db_session.add(version)
            await db_session.flush()
            inserted_versions += 1
        else:
            existing_versions += 1

        repository_dependency_version_link = RepositoryDependencyVersionDBModel(
            repository_id=repository_id,
            dependency_id=dependency.id,
            version_id=version.id,
        )

        db_session.add(repository_dependency_version_link)
        dependency_versions.append((version.id, name, version_str, ecosystem))

    await db_session.commit()
    logger.info(
        f"Dependencies: {inserted_dependencies} inserted, {existing_dependencies} existing; "
        f"Versions: {inserted_versions} inserted, {existing_versions} existing; "
    )

    return dependency_versions
