import toml
from pathlib import Path
from typing import List
import re

from app.parsers.base import register_parser


@register_parser("pyproject.toml", "pypi")
def parse_pep621_pyproject(file_path: Path) -> List[tuple[str, str, str]]:
    dependencies = []
    data = toml.load(file_path)

    project = data.get("project", {})

    for dep in project.get("dependencies", []):
        if " (" in dep and dep.endswith(")"):
            name, version = dep[:-1].split(" (", 1)
        else:
            name, version = dep, "unknown"

        name = re.sub(r"\[.*?\]", "", name).strip()
        version = re.sub(r"^[~^<>=!]+", "", version).strip()

        dependencies.append((name, version, "PyPI"))

    return dependencies

