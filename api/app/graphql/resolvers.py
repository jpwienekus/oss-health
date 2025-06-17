from typing import List

import httpx
import strawberry
from sqlalchemy.ext.asyncio import AsyncSession
from strawberry.types import Info

from app.auth.jwt_utils import decode_token
from app.crud.repository import add_repository_ids, get_repositories, get_repository, update_scanned_date
from app.crud.repository_dependency_version import (
    replace_repository_dependency_versions,
)
from app.crud.user import get_access_token, get_user
from app.crud.vulnerability import replace_version_vulnerabilities
from app.graphql.types import Dependency, GitHubRepository
from app.services.osv_api import get_dependency_version_vulnerability
from app.services.scanner import get_repository_dependencies


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

async def get_repositories_for_user(user_id: int, db_session: AsyncSession) -> List[GitHubRepository]:
    repositories = await get_repository_information_from_github(db_session, user_id)

    if not repositories:
        return []

    repositories_by_id = {repo.get("id"): repo for repo in repositories}
    tracked_repositories = await get_repositories(db_session, user_id)
    results: List[GitHubRepository] = []

    for repository in tracked_repositories:
        total_vulnerabilities = sum(len(d.version.vulnerabilities) for d in repository.dependency_versions)

        if repository.github_id in repositories_by_id:
            results.append(
                GitHubRepository.from_model(
                    repositories_by_id[repository.github_id],
                    id=repository.id,
                    score=repository.score,
                    number_of_dependencies=len(repository.dependency_versions),
                    number_of_vulnerabilities=total_vulnerabilities,
                    scanned_date=repository.scanned_date
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
    async def manual_scan_debug(self, info: Info, repository_id: int) -> List[GitHubRepository]:
        user_id = get_user_id(info)
        db = info.context["db"]

        repository = await get_repository(db, repository_id, user_id)

        if not repository:
            return []

        dependencies = get_repository_dependencies(
            repository.clone_url
        )
        dependency_versions_to_check = await replace_repository_dependency_versions(
            db, repository_id, dependencies
        )
        dependency_version_vulnerabilities = await get_dependency_version_vulnerability(dependency_versions_to_check)
        await replace_version_vulnerabilities(db, dependency_version_vulnerabilities)
        await update_scanned_date(db, repository_id, user_id)

        # print(test)

        return await get_repositories_for_user(user_id, db)

    @strawberry.field
    async def repositories(self, info: Info) -> List[GitHubRepository]:
        user_id = get_user_id(info)
        db = info.context["db"]

        return await get_repositories_for_user(user_id, db)



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
