from pathlib import Path
from typing import List
from app.scanners.scanner import Dependency, register_parser


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
