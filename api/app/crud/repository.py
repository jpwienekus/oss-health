from datetime import datetime
from typing import List, Sequence
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.models import Repository as RepositoryDBModel


async def get_repository(
    db_session: AsyncSession, user_id: int
) -> Sequence[RepositoryDBModel]:
    return (
        await db_session.scalars(
            select(RepositoryDBModel)
            .where(RepositoryDBModel.user_id == user_id)
            .order_by(RepositoryDBModel.updated_at.desc())
        )
    ).all()


async def upsert_user_repositories(
    db_session: AsyncSession, user_id: int, repos: List[dict]
):
    for repo in repos:
        github_id = repo.get("id")
        existing = (
            await db_session.scalars(
                select(RepositoryDBModel).where(
                    RepositoryDBModel.github_id == github_id
                )
            )
        ).first()

        if existing:
            updated_at = datetime.strptime(
                repo.get("updated_at", existing.updated_at), "%Y-%m-%dT%H:%M:%SZ"
            )
            existing.name = repo.get("name", existing.name)
            existing.description = repo.get("description", existing.description)
            existing.updated_at = updated_at
            existing.url = repo.get("url", existing.url)
            existing.open_issues = repo.get("open_issues", existing.open_issues)
            existing.score = 0
        else:
            date = repo.get("updated_at")
            updated_at = datetime.strptime(date, "%Y-%m-%dT%H:%M:%SZ") if date else None
            new_repo = RepositoryDBModel(
                github_id=repo.get("id"),
                name=repo.get("name"),
                description=repo.get("description"),
                updated_at=updated_at,
                user_id=user_id,
                url=repo.get("url"),
                open_issues=repo.get("open_issues"),
                score=0,
            )
            db_session.add(new_repo)

    await db_session.commit()
