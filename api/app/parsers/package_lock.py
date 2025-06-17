import json
from pathlib import Path
from typing import List

from app.parsers.base import register_parser


@register_parser("package-lock.json", "npm")
def parse_package_lock(file_path: Path) -> List[tuple[str, str, str]]:
    dependencies = []
    with file_path.open() as file:
        data = json.load(file)
        for name, info in data.get("dependencies", {}).items():
            version = info.get("version", "unknown")
            dependencies.append((name, version, "npm"))

    return dependencies
