from pathlib import Path
from typing import Callable, List, Tuple

# ruff: noqa: E501
dependency_parsers: List[
    Tuple[str, Callable[[Path], List[tuple[str, str, str]]], str]
] = []


def register_parser(pattern: str, ecosystem: str):
    def decorator(func: Callable[[Path], List[tuple[str, str, str]]]):
        dependency_parsers.append((pattern, func, ecosystem))
        return func

    return decorator
