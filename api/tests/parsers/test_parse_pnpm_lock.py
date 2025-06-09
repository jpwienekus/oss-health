from pathlib import Path

from app.parsers.pnpm_lock import parse_pnpm_lock


def test_parse_pnpm_lock(tmp_path: Path):
    file = tmp_path / "pnpm-lock.yaml"
    file.write_text(
        """
packages:
  axios@1.0.0:
    resolution: {integrity: sha512}
  node_modules/ignored@1.0.0:
    resolution: {integrity: sha512}
"""
    )

    result = parse_pnpm_lock(file)

    assert len(result) == 1
    assert result[0].name == "axios"
    assert result[0].version == "1.0.0"
