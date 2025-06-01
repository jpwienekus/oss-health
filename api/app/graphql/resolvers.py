import httpx
from datetime import datetime
from typing import List
import strawberry
from strawberry.types import Info
from app.auth.jwt_utils import decode_token
from app.crud.repository import get_repository, upsert_user_repositories
from app.crud.user import get_user
from app.models import Repository as RepositoryDBModel
from app.configuration import settings


def get_user_id(info: Info) -> int:
    token = (
        info.context["request"].headers.get("Authorization", "").replace("Bearer ", "")
    )
    return int(decode_token(token)["sub"])


def get_access_token(info: Info) -> str:
    token = (
        info.context["request"].headers.get("Authorization", "").replace("Bearer ", "")
    )
    return decode_token(token)["access_token"]


@strawberry.type
class RepositoryType:
    name: str
    description: str | None
    updated_at: datetime

    @classmethod
    def from_model(cls, model: RepositoryDBModel) -> "RepositoryType":
        return cls(
            name=model.name, description=model.description, updated_at=model.updated_at
        )


@strawberry.type
class Query:
    @strawberry.field
    async def repositories(self, info: Info) -> List[RepositoryType]:
        user_id = get_user_id(info)
        db = info.context["db"]
        result = await get_repository(db, user_id)

        return [RepositoryType.from_model(repo) for repo in result]


@strawberry.type
class Mutation:
    @strawberry.mutation
    async def sync_repositories(self, info: Info) -> int:
        user_id = get_user_id(info)
        db = info.context["db"]

        async with httpx.AsyncClient() as client:
            access_token = get_access_token(info)

            gh_response = await client.get(
                "https://api.github.com/user/repos?per_page=100&type=public",
                headers={"authorization": f"token {access_token}"},
            )
            repo_data = gh_response.json()

        await upsert_user_repositories(db, user_id, repo_data)

        return len(repo_data)
