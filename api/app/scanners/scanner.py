import fnmatch
import tempfile
import subprocess
from pathlib import Path
from dataclasses import dataclass
from typing import Callable, List, Tuple


@dataclass
class Dependency:
    name: str
    version: str
    ecosystem: str


dependency_parsers: List[Tuple[str, Callable[[Path], List[Dependency]], str]] = []


def register_parser(pattern: str, ecosystem: str):
    def decorator(func: Callable[[Path], List[Dependency]]):
        dependency_parsers.append((pattern, func, ecosystem))
        return func

    return decorator


@register_parser("requirements.txt", "pypi")
@register_parser("requirements-*.txt", "pypi")
@register_parser("requirements/*.txt", "pypi")
def parse_requirements_txt(file_path: Path) -> List[Dependency]:
    dependencies = []
    with file_path.open() as file:
        for line in file:
            line = line.strip()
            if line and not line.startswith("#"):
                if "==" in line:
                    name, version = line.split("==")
                else:
                    name, version = line, "unknown"

                dependencies.append(Dependency(name.strip(), version.strip(), "pypi"))

    return dependencies


def clone_repository(repository_url: str) -> Path:
    temp_dir = Path(tempfile.mkdtemp())
    subprocess.run(
        ["git", "clone", "--depth=1", repository_url, str(temp_dir)], check=True
    )

    return temp_dir


def extract_dependencies(repository_path: Path):
    all_dependencies = []

    for path in repository_path.rglob("*"):
        for pattern, parser, ecosystem in dependency_parsers:
            if fnmatch.fnmatch(path.name, pattern):
                try:
                    dependencies = parser(path)
                    all_dependencies.extend(dependencies)
                except Exception as e:
                    print(f"Failed to parse {path.name}: {e}")
    return all_dependencies


def get_repository_dependencies(repository_url: str) -> List[Dependency]:
    repository_path = None
    try:
        repository_path = clone_repository(repository_url)
        return extract_dependencies(repository_path)
    finally:
        if repository_path and repository_path.exists():
            subprocess.run(["rm", "-rf", str(repository_path)])
