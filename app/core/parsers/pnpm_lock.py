from pathlib import Path
from typing import List

import yaml

from core.parsers.base import register_parser


@register_parser("pnpm-lock.yaml", "npm")
def parse_pnpm_lock(file_path: Path) -> List[tuple[str, str, str]]:
    dependencies = []
    with file_path.open() as file:
        data = yaml.safe_load(file)

        for package_ref in data.get("packages").keys():
            parts = package_ref.split("/")

            if not parts or "node_modules" in parts:
                continue

            if "@" in package_ref:
                name_version = package_ref.lstrip("/").rsplit("@", 1)
                if len(name_version) == 2:
                    name, version = name_version
                    if name:
                        dependencies.append((name, version, "npm"))

    return dependencies
