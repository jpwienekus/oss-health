from unittest.mock import AsyncMock, MagicMock, patch

import pytest

from app.graphql.schema import schema


@pytest.fixture
def mock_context():
    request = MagicMock()
    request.headers = {"Authorization": "Bearer fake_token"}

    return {"request": request, "db": MagicMock()}


@pytest.mark.asyncio
@patch("app.graphql.resolvers.get_user", new_callable=AsyncMock)
@patch("app.graphql.resolvers.decode_token", return_value=123)
async def test_username(mock_decode_token, mock_get_user, mock_context):
    mock_get_user.return_value.github_username = "testuser"

    query = """
        query {
            username
        }
    """
    result = await schema.execute(
        query, variable_values={"title": "The Great Gatsby"}, context_value=mock_context
    )

    assert result.errors is None
    assert result.data is not None
    assert result.data["username"] == "testuser"
    mock_decode_token.assert_called_once()
    mock_get_user.assert_called_once()


@pytest.mark.asyncio
@patch("app.graphql.resolvers.decode_token", return_value=123)
@patch("app.graphql.resolvers.get_access_token", new_callable=AsyncMock)
@patch("app.graphql.resolvers.httpx.AsyncClient.get", new_callable=AsyncMock)
async def test_github_repositories(
    mock_get, mock_get_access_token, mock_decode_token, mock_context
):
    mock_get_access_token.return_value = "fake_github_token"
    mock_get.return_value.json = lambda: [
        {"id": 1, "name": "repo1", "updated_at": "2025-03-17T17:49:00Z"},
        {"id": 2, "name": "repo2", "updated_at": "2025-03-17T17:49:00Z"},
    ]
    query = """
  query {
    githubRepositories {
      name
      githubId
    }
  }
    """

    result = await schema.execute(query, context_value=mock_context)

    assert result.errors is None
    assert result.data is not None
    assert len(result.data["githubRepositories"]) == 2
    mock_get_access_token.assert_called_once()
    mock_get.assert_called_once()


@pytest.mark.asyncio
@patch("app.graphql.resolvers.decode_token", return_value=123)
@patch("app.graphql.resolvers.get_access_token", new_callable=AsyncMock)
@patch("app.graphql.resolvers.httpx.AsyncClient.get", new_callable=AsyncMock)
@patch("app.graphql.resolvers.get_repositories", new_callable=AsyncMock)
async def test_repositories_only_returns_tracked(
    mock_get_repositories,
    mock_get,
    mock_get_access_token,
    mock_decode_token,
    mock_context,
):
    mock_get_access_token.return_value = "fake_token"
    mock_get.return_value.json = lambda: [
        {"id": 1, "name": "repo1", "updated_at": "2025-03-17T17:49:00Z"},
        {"id": 2, "name": "repo2", "updated_at": "2025-03-17T17:49:00Z"},
    ]
    mock_get_repositories.return_value = [
        MagicMock(github_id=1, score=10),
    ]

    query = """
  query {
    repositories {
      name
      githubId
    }
  }
    """

    result = await schema.execute(query, context_value=mock_context)

    assert result.errors is None
    assert result.data is not None
    assert len(result.data["repositories"]) == 1


@pytest.mark.asyncio
@patch("app.graphql.resolvers.decode_token", return_value=123)
@patch("app.graphql.resolvers.get_access_token", new_callable=AsyncMock)
@patch("app.graphql.resolvers.httpx.AsyncClient.get", new_callable=AsyncMock)
@patch("app.graphql.resolvers.add_repository_ids", new_callable=AsyncMock)
@patch("app.graphql.resolvers.get_repositories", new_callable=AsyncMock)
async def test_save_selected_repositories(
    mock_get_repositories,
    mock_sync,
    mock_get,
    mock_get_access_token,
    mock_decode_token,
    mock_context,
):
    mock_get_access_token.return_value = "gh_token"
    mock_get.return_value.json = lambda: [
        {"id": 1, "name": "repo1", "updated_at": "2025-03-17T17:49:00Z"},
        {"id": 2, "name": "repo2", "updated_at": "2025-03-17T17:49:00Z"},
    ]
    mock_get_repositories.return_value = [MagicMock(github_id=1)]
    mutation = """
  mutation SaveSelectedRepositories($selectedGithubRepositoryIds: [Int!]!) {
    saveSelectedRepositories(
      selectedGithubRepositoryIds: $selectedGithubRepositoryIds
    ) {
      name
      githubId
    }
  }
    """

    result = await schema.execute(
        mutation,
        context_value=mock_context,
        variable_values={"selectedGithubRepositoryIds": [1]},
    )

    assert result.errors is None
    assert result.data is not None
    assert len(result.data["saveSelectedRepositories"]) == 1
    assert result.data["saveSelectedRepositories"][0]["name"] == "repo1"
