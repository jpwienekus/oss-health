from typing import List

import httpx
import strawberry
from sqlalchemy.ext.asyncio import AsyncSession
from strawberry.types import Info

from api.auth.jwt_utils import decode_token
from api.graphql.inputs import DependencyFilter, DependencySortInput, PaginationInput
from api.graphql.types import DependencyConnection, GitHubRepository
from core.crud.dependency import get_dependencies_paginated
from core.crud.repository import (
    add_repository_ids,
    get_repositories,
)
from core.crud.user import get_access_token, get_user


def get_user_id(info: Info) -> int:
    token = (
        info.context["request"].headers.get("Authorization", "").replace("Bearer ", "")
    )
    return decode_token(token)


async def get_repository_information_from_github(
    db: AsyncSession, user_id: int
) -> List[dict]:
    async with httpx.AsyncClient() as client:
        access_token = await get_access_token(db, user_id)
        gh_response = await client.get(
            "https://api.github.com/user/repos?per_page=100&type=public&sort=updated",
            headers={"authorization": f"token {access_token}"},
        )
        repo_data = gh_response.json()

    return repo_data


async def get_repositories_for_user(
    user_id: int, db_session: AsyncSession
) -> List[GitHubRepository]:
    repositories = await get_repository_information_from_github(db_session, user_id)

    if not repositories:
        return []

    repositories_by_id = {repo.get("id"): repo for repo in repositories}
    tracked_repositories = await get_repositories(db_session, user_id)
    results: List[GitHubRepository] = []

    for repository in tracked_repositories:
        total_vulnerabilities = sum(
            len(d.version.vulnerabilities) for d in repository.dependency_versions
        )

        if repository.github_id in repositories_by_id:
            results.append(
                GitHubRepository.from_model(
                    repositories_by_id[repository.github_id],
                    id=repository.id,
                    score=repository.score,
                    number_of_dependencies=len(repository.dependency_versions),
                    number_of_vulnerabilities=total_vulnerabilities,
                    last_scanned_at=repository.last_scanned_at,
                )
            )

    return results


@strawberry.type
class Query:
    @strawberry.field
    async def username(self, info: Info) -> str:
        user_id = get_user_id(info)
        db = info.context["db"]
        user = await get_user(db, user_id)
        return user.github_username if user is not None else ""

    @strawberry.field
    async def github_repositories(self, info: Info) -> List[GitHubRepository]:
        user_id = get_user_id(info)
        db = info.context["db"]

        repositories = await get_repository_information_from_github(db, user_id)

        return [GitHubRepository.from_model(repo) for repo in repositories]

    @strawberry.field
    async def repositories(self, info: Info) -> List[GitHubRepository]:
        user_id = get_user_id(info)
        db = info.context["db"]

        return await get_repositories_for_user(user_id, db)

    @strawberry.field
    async def dependencies(self, info: Info, filter: DependencyFilter, sort: DependencySortInput, pagination: PaginationInput) -> DependencyConnection | None:
    # async def dependencies(self, info: Info, filter: DependencyFilter) -> DependencyConnection | None:
        return await get_dependencies_paginated(info.context["db"], filter, sort, pagination)
        # return None


@strawberry.type
class Mutation:
    @strawberry.mutation
    async def save_selected_repositories(
        self, info: Info, selected_github_repository_ids: List[int]
    ) -> List[GitHubRepository]:
        user_id = get_user_id(info)
        db = info.context["db"]

        repositories = await get_repository_information_from_github(db, user_id)
        tracked_repositories = [
            r for r in repositories if r.get("id") in selected_github_repository_ids
        ]

        await add_repository_ids(db, user_id, tracked_repositories)

        return [GitHubRepository.from_model(repo) for repo in tracked_repositories]
