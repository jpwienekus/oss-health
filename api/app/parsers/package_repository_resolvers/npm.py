import httpx
from typing import Optional
from app.parsers.package_repository_resolvers.base import register_package_repository_resolver


@register_package_repository_resolver("npm")
async def get_npm_repo_url(name: str) -> Optional[str]:
    async with httpx.AsyncClient() as client:
        response = await client.get(f"https://registry.npmjs.org/{name}")

        if response.status_code == 200:
            repository = response.json().get("repository", {})
            url = repository.get("url") or ""

            return url.replace("git+", "").replace(".git", "").replace("git://", "https://")

    return None
