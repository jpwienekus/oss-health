from enum import Enum
import strawberry


@strawberry.enum
class DependencySortField(Enum):
    ID = "id"
    NAME = "name"
    ECOSYSTEM = "ecosystem"
    CHECKED_AT = "github_url_checked_at"

@strawberry.enum
class SortDirection(Enum):
    ASC = "asc"
    DESC = "desc"

@strawberry.input
class DependencyFilter:
    name: str = ""
    ecosystem: str = ""
    github_url_resolve_failed: bool | None = None


@strawberry.input
class DependencySortInput:
    field: DependencySortField = DependencySortField.ID
    direction: SortDirection = SortDirection.ASC


@strawberry.input
class PaginationInput:
    page: int = 0
    page_size: int = 25
