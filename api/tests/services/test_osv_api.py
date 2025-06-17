from unittest.mock import patch

import pytest

from app.services.osv_api import get_dependency_version_vulnerability


@pytest.mark.asyncio
@patch("app.services.osv_api.httpx.AsyncClient.post")
async def test_update_dependency_vulnerability(mock_post):
    dependency_versions = [
        (1, "requests", "2.25.1", "PyPI"),
        (2, "flask", "2.0.0", "PyPI"),
    ]

    mock_post.return_value.json = lambda: {
        "results": [
            {"vulns": [{"id": "OSV-123"}, {"id": "OSV-456"}]},
            {"vulns": []},
        ]
    }

    result = await get_dependency_version_vulnerability(dependency_versions)
    assert result == [
        (1, ["OSV-123", "OSV-456"]),
        (2, []),
    ]

    expected_queries = [
        {"package": {"name": "requests", "ecosystem": "PyPI"}, "version": "2.25.1"},
        {"package": {"name": "flask", "ecosystem": "PyPI"}, "version": "2.0.0"},
    ]
    mock_post.assert_awaited_once_with(
        "https://api.osv.dev/v1/querybatch",
        json={"queries": expected_queries},
        timeout=30,
    )
