from fastapi import APIRouter, HTTPException
from fastapi.responses import HTMLResponse, RedirectResponse
import httpx
from auth.jwt_utils import create_token_pair
from configuration import settings
from dependencies.core import DBSessionDep
from crud.user import get_user, add_user

router = APIRouter()

@router.get("/auth/github/login")
async def login(code_challenge: str):
    return RedirectResponse(
        f"https://github.com/login/oauth/authorize"
        f"?client_id={settings.github_client_id}"
        f"&scope=read:user"
        f"&response_type=code"
        f"&code_challenge={code_challenge}"
        f"&code_challenge_method=S256"
    )


@router.get("/auth/github/callback", response_class=HTMLResponse)
async def callback_html(code: str):
    frontend_url = "http://localhost:5173"
    return f"""
    <script>
        window.opener.postMessage({{"type": "github-oauth-code", "code": "{code}"}}, "{frontend_url}")
        window.close()
    </script>
    """

@router.post("/auth/github/token")
async def github_token_exchange(payload: dict, db_session: DBSessionDep):
    code = payload.get("code")
    code_verifier = payload.get("code_verifier")

    async with httpx.AsyncClient() as client:
        token_response = await client.post(
            "https://github.com/login/oauth/access_token",
            headers={"Accept": "application/json"},
            data={
                "client_id": settings.github_client_id,
                "client_secret": settings.github_client_secret,
                "code": code,
                "code_verifier": code_verifier,
            },
        )

        token_json = token_response.json()
        access_token = token_json.get("access_token")

        user_response = await client.get(
            "https://api.github.com/user",
            headers={"Authorization": f"token {access_token}"},
        )
        github_user = user_response.json()
        id = github_user.get("id")
        username = github_user.get("login")

    if not id:
        raise HTTPException(status_code=404, detail="Invalid user")

    user = await get_user(db_session, id)

    if not user:
        user = await add_user(db_session, id, username)

    jwt_token = create_token_pair(user.id)

    return { "access_token": jwt_token}
