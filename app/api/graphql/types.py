from datetime import datetime
from typing import List, Optional

import strawberry

@strawberry.type
class DependencyType:
    id: int
    name: str
    ecosystem: str
    scan_status: str
    scanned_at: Optional[datetime]
    error_message: Optional[str]
    repository_url: Optional[str]

    @classmethod
    def from_model(
        cls,
        model,
    ) -> "DependencyType":
        return cls(
            id=model.id,
            name=model.name,
            ecosystem=model.ecosystem,
            scan_status=model.scan_status,
            repository_url=model.dependency_repository.repository_url if model.dependency_repository else None,
            scanned_at=model.scanned_at,
            error_message=model.error_message
        )

@strawberry.type
class DependencyPaginatedResponse:
    dependencies: List[DependencyType]
    total_pages: int = 0
    completed: int = 0
    pending: int = 0
    failed: int = 0


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
