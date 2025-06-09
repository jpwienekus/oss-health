from pathlib import Path
from typing import List
from app.models.dependency import Dependency
from app.parsers.base import register_parser


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

                dependencies.append(
                    Dependency(name=name, version=version, ecosystem="PyPI")
                )

    return dependencies
