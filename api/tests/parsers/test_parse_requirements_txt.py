from pathlib import Path

from app.parsers.requirements_txt import parse_requirements_txt


def test_parse_requirements_txt(tmp_path: Path):
    file = tmp_path / "requirements.txt"
    file.write_text("requests==2.25.1\nnumpy\n# comment\n")

    result = parse_requirements_txt(file)

    assert len(result) == 2
    name, version, ecosystem = result[0]
    assert name == "requests"
    assert version == "2.25.1"
    assert ecosystem == "PyPI"
