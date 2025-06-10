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
    name, version, ecosystem = result[0]
    assert name == "axios"
    assert version == "1.0.0"
    assert ecosystem == "npm"
