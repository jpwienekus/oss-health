import re
from pathlib import Path
from typing import List

import toml

from core.parsers.base import register_parser


@register_parser("pyproject.toml", "pypi")
def parse_pyproject_toml(file_path: Path) -> List[tuple[str, str, str]]:
    dependencies = []

    data = toml.load(file_path)

    poetry_section = data.get("tool", {}).get("poetry", {})
    deps_section = poetry_section.get("dependencies", {})

    for name, version in deps_section.items():
        # Skip Python itself
        if name.lower() == "python":
            continue

        name, version = get_name_and_version(name, version)
        dependencies.append((name, version, "PyPI"))

    dev_deps = (
        data.get("tool", {})
        .get("poetry", {})
        .get("group", {})
        .get("dev", {})
        .get("dependencies", {})
    )
    for name, version in dev_deps.items():
        name, version = get_name_and_version(name, version)
        dependencies.append((name, version, "PyPI"))

    return dependencies


def get_name_and_version(name: str, version):
    parsed_version = "unknown"

    if isinstance(version, str):
        parsed_version = version
    elif isinstance(version, dict) and "version" in version:
        parsed_version = version["version"]

    name = re.sub(r"\[.*?\]", "", name).strip()
    parsed_version = re.sub(r"^[~^<>=!]+", "", parsed_version).strip()

    return name, parsed_version
