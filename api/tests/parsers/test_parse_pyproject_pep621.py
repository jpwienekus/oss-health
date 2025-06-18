from pathlib import Path

from app.parsers.pyproject_pep621 import parse_pep621_pyproject


def test_parse_pep621_pyproject(tmp_path: Path):
    file = tmp_path / "pyproject.toml"
    file.write_text(
        """
[project]
name = "example"
version = "0.1.0"
dependencies = [
    "requests (>=2.25)",
    "httpx (==0.27.0)",
    "custom-lib",
    "fastapi[standard] (==0.115.13)",
]
"""
    )

    result = parse_pep621_pyproject(file)

    assert len(result) == 4
    assert ("requests", "2.25", "PyPI") in result
    assert ("httpx", "0.27.0", "PyPI") in result
    assert ("custom-lib", "unknown", "PyPI") in result
    assert ("fastapi", "0.115.13", "PyPI") in result
