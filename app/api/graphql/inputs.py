from enum import Enum
from typing import List

import strawberry


@strawberry.enum
class DependencySortField(Enum):
    ID = "id"
    NAME = "name"
    ECOSYSTEM = "ecosystem"
    SCANNED_AT = "scanned_at"
    SCAN_STATUS = "scan_status"
    ERROR_MESSAGE = "error_message"


@strawberry.enum
class SortDirection(Enum):
    ASC = "asc"
    DESC = "desc"


@strawberry.input
class DependencyFilter:
    name: str = ""
    statuses: List[str]


@strawberry.input
class DependencySortInput:
    field: DependencySortField = DependencySortField.ID
    direction: SortDirection = SortDirection.ASC


@strawberry.input
class PaginationInput:
    page: int = 0
    page_size: int = 25
