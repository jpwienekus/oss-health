from pathlib import Path

from core.parsers.package_lock import parse_package_lock


def test_parse_package_lock(tmp_path: Path):
    file = tmp_path / "package-lock.json"
    file.write_text('{"dependencies": {"lodash": {"version": "4.17.21"}}}')

    result = parse_package_lock(file)

    assert len(result) == 1
    name, version, ecosystem = result[0]
    assert name == "lodash"
    assert version == "4.17.21"
    assert ecosystem == "npm"
