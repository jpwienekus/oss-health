from pathlib import Path

from core.parsers.pyproject_toml import parse_pyproject_toml


def test_parse_pyproject_toml(tmp_path: Path):
    file = tmp_path / "pyproject.toml"
    file.write_text(
        """
[tool.poetry]
name = "example"
version = "0.1.0"
description = ""
authors = ["Jane Doe <jane@example.com>"]

[tool.poetry.dependencies]
python = "^3.10"
requests = "^2.31.0"
httpx = { version = "^0.27.0", extras = ["http2"] }
custom = {}

[tool.poetry.group.dev.dependencies]
pytest = "^8.0.0"
black = { version = "^24.3.0" }
mypy = { some_other_field = "irrelevant" }
"""
    )

    result = parse_pyproject_toml(file)

    assert ("requests", "2.31.0", "PyPI") in result
    assert ("httpx", "0.27.0", "PyPI") in result
    assert ("custom", "unknown", "PyPI") in result
    assert ("pytest", "8.0.0", "PyPI") in result
    assert ("black", "24.3.0", "PyPI") in result
    assert ("mypy", "unknown", "PyPI") in result

    # Should not include 'python'
    assert all(name.lower() != "python" for name, _, _ in result)
