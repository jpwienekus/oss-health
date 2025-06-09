from typing import List
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.crud.repository import get_repository_with_dependencies_loaded
from app.models import Dependency as DependencyDBModel
from app.models.dependency import Dependency


async def add_dependencies_to_repository(
    db_session: AsyncSession,
    user_id: int,
    repository_id: int,
    dependencies: List[Dependency],
):
    repository = await get_repository_with_dependencies_loaded(db_session, user_id, repository_id)

    if not repository:
        raise ValueError(f"Repository with id {repository_id} not found.")

    repository.dependencies.clear()
    await db_session.flush()

    attached_dependencies = []

    for dependency in dependencies:
        result = await db_session.execute(
            select(DependencyDBModel).where(
                DependencyDBModel.name == dependency.name,
                DependencyDBModel.version == dependency.version,
                DependencyDBModel.ecosystem == dependency.ecosystem,
            )
        )
        existing_dependency = result.scalar_one_or_none()

        if existing_dependency:
            attached_dependencies.append(existing_dependency)
        else:
            db_session.add(dependency)
            await db_session.flush()
            attached_dependencies.append(dependency)

    repository.dependencies.extend(attached_dependencies)

    await db_session.commit()
