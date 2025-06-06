import strawberry
from datetime import datetime


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

    @classmethod
    def from_model(cls, model) -> "GitHubRepository":
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
        )
