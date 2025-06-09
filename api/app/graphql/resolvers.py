import httpx
from typing import List
from sqlalchemy.ext.asyncio import AsyncSession
import strawberry
from strawberry.types import Info
from app.auth.jwt_utils import decode_token
from app.crud.dependency import replace_repository_dependencies
from app.crud.repository import add_repository_ids, get_repositories
from app.crud.user import get_access_token, get_user
from app.graphql.types import Dependency, GitHubRepository
from app.services.osv_api import update_dependency_vulnerability
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


def chunk_list(data, chunk_size):
    """Yield successive chunks from data of size chunk_size."""
    for i in range(0, len(data), chunk_size):
        yield data[i : i + chunk_size]


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

        return [GitHubRepository.from_model(repo, 0) for repo in repositories]

    @strawberry.field
    async def debug_cloning(self, info: Info) -> List[Dependency]:
        user_id = get_user_id(info)
        db = info.context["db"]

        tracked_repositories = await get_repositories(db, user_id)
        debug_repo = tracked_repositories[1]
        dependencies = get_repository_dependencies(debug_repo.clone_url)
        attached_dependencies = await replace_repository_dependencies(
            db, user_id, debug_repo.id, dependencies
        )
        await update_dependency_vulnerability(db, attached_dependencies)

        return dependencies

    @strawberry.field
    async def repositories(self, info: Info) -> List[GitHubRepository]:
        user_id = get_user_id(info)
        db = info.context["db"]

        repositories = await get_repository_information_from_github(db, user_id)

        if len(repositories) == 0:
            return []

        tracked_repositories = await get_repositories(db, user_id)
        repository_score_map: dict[int, int] = {
            repository.github_id: repository.score
            for repository in tracked_repositories
        }
        tracked_repository_ids = repository_score_map.keys()

        return [
            GitHubRepository.from_model(
                repo, score=repository_score_map.get(repo.get("id", 0), 0)
            )
            for repo in repositories
            if repo.get("id") in tracked_repository_ids
        ]


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

        return [GitHubRepository.from_model(repo, 0) for repo in tracked_repositories]
