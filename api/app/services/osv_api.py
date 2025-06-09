from typing import List

import httpx
from sqlalchemy.ext.asyncio import AsyncSession

from app.crud.vulnerability import update_dependency_vulnerabilities
from app.models.dependency import Dependency
from app.models.vulnerability import Vulnerability
from app.utils.chunking import chunk_list


async def update_dependency_vulnerability(
    db: AsyncSession, dependencies: List[Dependency]
):
    for chunk in chunk_list(dependencies, 500):
        queries = [
            {
                "package": {"name": dependency.name, "ecosystem": dependency.ecosystem},
                "version": dependency.version,
            }
            for dependency in chunk
        ]

        if not queries:
            continue

        async with httpx.AsyncClient() as client:
            response = await client.post(
                "https://api.osv.dev/v1/querybatch",
                json={"queries": queries},
                timeout=30,
            )
            results = response.json()["results"]

        for dependency, result in zip(chunk, results):
            vulnerabilities = [
                Vulnerability(osv_id=v.get("id")) for v in result.get("vulns", [])
            ]

            await update_dependency_vulnerabilities(db, dependency.id, vulnerabilities)
