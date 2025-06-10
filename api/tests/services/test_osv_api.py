# from unittest.mock import AsyncMock, patch
#
# import pytest
#
# from app.models.dependency import Dependency
# from app.services.osv_api import update_dependency_vulnerability
#
#
# @pytest.mark.asyncio
# @patch("app.services.osv_api.httpx.AsyncClient.post")
# @patch("app.services.osv_api.update_dependency_vulnerabilities")
# async def test_update_dependency_vulnerability(mock_update_db, mock_post):
#     dep = Dependency(id=1, name="requests", version="2.25.1", ecosystem="PyPI")
#     mock_post.return_value.json = lambda: {
#         "results": [{"vulns": [{"id": "OSV-2023-001"}]}]
#     }
#
#     db = AsyncMock()
#     await update_dependency_vulnerability(db, [dep])
#
#     mock_post.assert_called_once()
#     mock_update_db.assert_called_once()
