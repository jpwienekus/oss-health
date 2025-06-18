from datetime import datetime

import strawberry


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
    clone_url: str
    scanned_date: datetime | None
    updated_at: datetime | None

    @classmethod
    def from_model(cls, model, id: int| None = None, score: int | None = None, number_of_dependencies: int | None = None, number_of_vulnerabilities: int | None = None, scanned_date: datetime | None = None) -> "GitHubRepository":
        updated_at = datetime.fromisoformat(model.get("updated_at").replace("Z", "+00:00"))
        print(updated_at)
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
            clone_url=model.get("clone_url"),
            scanned_date=scanned_date,
            updated_at=datetime.fromisoformat(model.get("updated_at").replace("Z", "+00:00"))
        )
