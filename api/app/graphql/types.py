from datetime import datetime

import strawberry


@strawberry.type
class Dependency:
    name: str
    version: str
    ecosystem: str


# @strawberry.type
# class Vulnerability:
#     id: str
#     summary: str
#     severity: str
#     # affected_versions: List[str]


@strawberry.type
class GitHubRepository:
    name: str
    description: str | None
    github_id: int
    stars: int
    watchers: int
    updated_at: datetime | None
    private: bool
    forks: int
    score: int
    vulnerabilities: int
    dependencies: int
    clone_url: str

    @classmethod
    def from_model(cls, model, score: int) -> "GitHubRepository":
        date = model.get("updated_at")
        updated_at = datetime.strptime(date, "%Y-%m-%dT%H:%M:%SZ") if date else None
        return cls(
            name=model.get("name"),
            description=model.get("description"),
            github_id=model.get("id"),
            stars=model.get("stargazers_count"),
            watchers=model.get("watchers_count"),
            updated_at=updated_at,
            private=model.get("private"),
            forks=model.get("forks_count"),
            score=score,
            vulnerabilities=0,
            dependencies=0,
            clone_url=model.get("clone_url"),
        )
