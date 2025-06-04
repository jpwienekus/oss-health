import httpx
from datetime import datetime
from typing import List
import strawberry
from strawberry.types import Info
from app.auth.jwt_utils import decode_token
from app.crud.repository import get_repository, upsert_user_repositories
from app.crud.user import get_access_token, get_sync_time, get_user, update_sync_time
from app.models import Repository as RepositoryDBModel


def get_user_id(info: Info) -> int:
    token = (
        info.context["request"].headers.get("Authorization", "").replace("Bearer ", "")
    )
    return decode_token(token)


@strawberry.type
class RepositoryType:
    name: str
    description: str | None
    updated_at: datetime
    url: str | None
    open_issues: int | None
    score: int | None

    @classmethod
    def from_model(cls, model: RepositoryDBModel) -> "RepositoryType":
        return cls(
            name=model.name,
            description=model.description,
            updated_at=model.updated_at,
            url=model.url,
            open_issues=model.open_issues,
            score=model.score,
        )


@strawberry.type
class RepositoriesResponse:
    repositories: List[RepositoryType]
    sync_date: datetime | None


@strawberry.type
class Query:
    @strawberry.field
    async def username(self, info: Info) -> str:
        user_id = get_user_id(info)
        db = info.context["db"]
        user = await get_user(db, user_id)
        return user.github_username if user is not None else ''


    @strawberry.field
    async def repositories(self, info: Info) -> RepositoriesResponse:
        user_id = get_user_id(info)
        db = info.context["db"]
        result = await get_repository(db, user_id)
        repositories = [RepositoryType.from_model(repo) for repo in result]
        sync_date = await get_sync_time(db, user_id)

        return RepositoriesResponse(repositories=repositories, sync_date=sync_date)


@strawberry.type
class Mutation:
    @strawberry.mutation
    async def sync_repositories(self, info: Info) -> RepositoriesResponse:
        user_id = get_user_id(info)
        db = info.context["db"]

        async with httpx.AsyncClient() as client:
            access_token = await get_access_token(db, user_id)
            print('^' * 100)
            print(access_token)

            gh_response = await client.get(
                "https://api.github.com/user/repos?per_page=100&type=public&sort=updated",
                headers={"authorization": f"token {access_token}"},
            )
            repo_data = gh_response.json()

        print(repo_data)
        await upsert_user_repositories(db, user_id, repo_data)
        result = await get_repository(db, user_id)
        repositories = [RepositoryType.from_model(repo) for repo in result]
        sync_date = await update_sync_time(db, user_id)

        return RepositoriesResponse(repositories=repositories, sync_date=sync_date)
