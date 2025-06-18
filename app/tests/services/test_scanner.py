import tempfile
from pathlib import Path
from typing import List
from unittest.mock import patch

from core.parsers import dependency_parsers
from core.services.scanner import extract_dependencies, get_repository_dependencies


def test_extract_dependencies_with_mock_parser(tmp_path: Path):
    dummy_file = tmp_path / "test.txt"
    dummy_file.write_text("dummy")

    def dummy_parser(_) -> List[tuple[str, str, str]]:
        return [("test", "1.0.0", "dummy")]

    dependency_parsers.append(("test.txt", dummy_parser, "dummy"))

    result = extract_dependencies(tmp_path)
    assert len(result) == 1
    name, version, ecosystem = result[0]
    assert name == "test"
    assert version == "1.0.0"
    assert ecosystem == "dummy"


@patch("core.services.scanner.clone_repository")
@patch("core.services.scanner.extract_dependencies")
@patch("core.services.scanner.subprocess.run")
def test_get_repository_dependencies(mock_rm, mock_extract, mock_clone):
    mock_path = tempfile.mkdtemp()
    mock_clone.return_value = Path(mock_path)
    mock_extract.return_value = []

    dependencies = get_repository_dependencies("https://fake.repo")

    assert dependencies == []
    mock_rm.assert_called()
