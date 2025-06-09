from pathlib import Path
from typing import Callable, List, Tuple

from app.models.dependency import Dependency

dependency_parsers: List[Tuple[str, Callable[[Path], List[Dependency]], str]] = []


def register_parser(pattern: str, ecosystem: str):
    def decorator(func: Callable[[Path], List[Dependency]]):
        dependency_parsers.append((pattern, func, ecosystem))
        return func

    return decorator
