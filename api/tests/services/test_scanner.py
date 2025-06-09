import tempfile
from pathlib import Path
from unittest.mock import patch

from app.models.dependency import Dependency
from app.parsers import dependency_parsers
from app.services.scanner import extract_dependencies, get_repository_dependencies


def test_extract_dependencies_with_mock_parser(tmp_path: Path):
    dummy_file = tmp_path / "test.txt"
    dummy_file.write_text("dummy")

    def dummy_parser(_):
        return [Dependency(name="test", version="1.0.0", ecosystem="dummy")]

    dependency_parsers.append(("test.txt", dummy_parser, "dummy"))

    result = extract_dependencies(tmp_path)
    assert len(result) == 1
    assert result[0].name == "test"


@patch("app.services.scanner.clone_repository")
@patch("app.services.scanner.extract_dependencies")
@patch("app.services.scanner.subprocess.run")
def test_get_repository_dependencies(mock_rm, mock_extract, mock_clone):
    mock_path = tempfile.mkdtemp()
    mock_clone.return_value = Path(mock_path)
    mock_extract.return_value = []

    result = get_repository_dependencies("https://fake.repo")

    assert result == []
    mock_rm.assert_called()
