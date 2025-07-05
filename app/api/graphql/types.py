from datetime import datetime
from typing import List, Optional

import strawberry

@strawberry.type
class DependencyType:
    id: int
    name: str
    ecosystem: str
    status: str
    repository_url_checked_at: Optional[datetime]
    repository_url_resolve_failed_reason: Optional[str]

    @classmethod
    def from_model(
        cls,
        model,
    ) -> "DependencyType":
        return cls(
            id=model.id,
            name=model.name,
            ecosystem=model.ecosystem,
            status=model.status,
            repository_url_checked_at=model.repository_url_checked_at,
            repository_url_resolve_failed_reason=model.repository_url_resolve_failed_reason
        )

@strawberry.type
class DependencyPaginatedResponse:
    dependencies: List[DependencyType]
    total: int = 0



# @strawberry.type
# class DependencyEdge:
#     node: DependencyType
#     cursor: int

# @strawberry.type
# class PageInfo:
#     has_next_page: bool
#     has_previous_page: bool
#     start_cursor: Optional[str]
#     end_cursor: Optional[str]

# @strawberry.type
# class DependencyConnection:
#     edges: List[DependencyEdge]
#     page_info: PageInfo


@strawberry.type
class Dependency:
    name: str
    version: str
    ecosystem: str


@strawberry.type
class GitHubRepository:
    id: int | None
    name: str
    description: str | None
    github_id: int
    stars: int
    watchers: int
    forks: int
    private: bool
    score: int | None
    vulnerabilities: int | None
    dependencies: int | None
    url: str
    last_scanned_at: datetime | None
    updated_at: datetime | None

    @classmethod
    def from_model(
        cls,
        model,
        id: int | None = None,
        score: int | None = None,
        number_of_dependencies: int | None = None,
        number_of_vulnerabilities: int | None = None,
        last_scanned_at: datetime | None = None,
    ) -> "GitHubRepository":
        return cls(
            id=id,
            name=model.get("name"),
            description=model.get("description"),
            github_id=model.get("id"),
            private=model.get("private"),
            score=score,
            stars=model.get("stargazers_count"),
            watchers=model.get("watchers_count"),
            forks=model.get("forks_count"),
            vulnerabilities=number_of_vulnerabilities,
            dependencies=number_of_dependencies,
            url=model.get("url"),
            last_scanned_at=last_scanned_at,
            updated_at=datetime.fromisoformat(
                model.get("updated_at").replace("Z", "+00:00")
            ),
        )
