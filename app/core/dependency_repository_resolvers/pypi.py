from typing import Optional

import httpx

from core.dependency_repository_resolvers.base import (
    register_package_repository_resolver,
)


@register_package_repository_resolver("PyPI")
async def get_pypi_repo_url(name: str) -> Optional[str]:
    async with httpx.AsyncClient() as client:
        response = await client.get(f"https://pypi.org/pypi/{name}/json")

        if response.status_code == 200:
            info = response.json()["info"]
            urls = info.get("project_urls") or {}

            return urls.get("Source") or urls.get("Homepage") or info.get("home_page")

    return None
