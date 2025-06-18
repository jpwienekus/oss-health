from typing import List

import httpx

from core.utils.chunking import chunk_list


async def get_dependency_version_vulnerability(
    dependency_versions: List[tuple[int, str, str, str]],
) -> List[tuple[int, List[str]]]:
    vulnerabilities: List[tuple[int, List[str]]] = []

    for chunk in chunk_list(dependency_versions, 500):
        queries = [
            {
                "package": {"name": name, "ecosystem": ecosystem},
                "version": version,
            }
            for _, name, version, ecosystem in chunk
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

        for (version_id, name, version, ecosystem), result in zip(chunk, results):
            vulnerability_ids = [v["id"] for v in result.get("vulns", [])]
            vulnerabilities.append((version_id, vulnerability_ids))

            if len(vulnerability_ids) > 0:
                # ruff: noqa: E501
                print(
                    f"Package: {name}, Version: {version}, Ecosystem: {ecosystem} -> Vulnerabilities: {vulnerability_ids}"
                )

    return vulnerabilities
