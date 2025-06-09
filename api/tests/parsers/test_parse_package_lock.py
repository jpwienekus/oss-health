from pathlib import Path
from app.parsers.package_lock import parse_package_lock


def test_parse_package_lock(tmp_path: Path):
    file = tmp_path / "package-lock.json"
    file.write_text('{"dependencies": {"lodash": {"version": "4.17.21"}}}')

    result = parse_package_lock(file)

    assert len(result) == 1
    assert result[0].name == "lodash"
    assert result[0].version == "4.17.21"
