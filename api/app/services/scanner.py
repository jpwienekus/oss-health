import fnmatch
import subprocess
import tempfile
from pathlib import Path
from typing import List

from app.parsers import dependency_parsers


def clone_repository(repository_url: str) -> Path:
    temp_dir = Path(tempfile.mkdtemp())
    subprocess.run(
        ["git", "clone", "--depth=1", repository_url, str(temp_dir)], check=True
    )

    return temp_dir


def extract_dependencies(repository_path: Path):
    all_dependencies: List[tuple[str, str, str]] = []

    for path in repository_path.rglob("*"):
        for pattern, parser, ecosystem in dependency_parsers:
            if fnmatch.fnmatch(path.name, pattern):
                try:
                    dependencies = parser(path)
                    all_dependencies.extend(dependencies)
                except Exception as e:
                    print(f"Failed to parse {path.name}: {e}")

    return all_dependencies


def get_repository_dependencies(repository_url: str) -> List[tuple[str, str, str]]:
    repository_path = None
    try:
        repository_path = clone_repository(repository_url)
        dependencies = extract_dependencies(repository_path)
        return dependencies
    finally:
        if repository_path and repository_path.exists():
            subprocess.run(["rm", "-rf", str(repository_path)])
