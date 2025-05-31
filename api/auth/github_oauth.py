from fastapi import APIRouter
from fastapi.responses import RedirectResponse
import httpx
from configuration import settings
from dependencies.core import DBSessionDep
from crud.user import get_user, add_user

router = APIRouter()


@router.get("auth/github/login")
async def login():
    return RedirectResponse(
        f"https://github.com/login/oauth/authorize?client_id={settings.github_client_id}&scope=read:user"
    )


@router.get("/auth/github/callback")
async def callback(code: str, db_session: DBSessionDep):
    async with httpx.AsyncClient() as client:
        token_response = await client.post(
            "https://github.com/login/oauth/access_token",
            headers={"Accept": "application/json"},
            data={
                "client_id": settings.github_client_id,
                "client_secret": settings.github_client_secret,
                "code": code,
            },
        )

        token_json = token_response.json()
        access_token = token_json["access_token"]

        user_response = await client.get(
            "https://api.github.com/user",
            headers={"Authorization": f"token {access_token}"},
        )
        github_user = user_response.json()

    user = await get_user(db_session, github_user["login"])

    if not user:
        await add_user(db_session, github_user["id"], github_user["login"])
