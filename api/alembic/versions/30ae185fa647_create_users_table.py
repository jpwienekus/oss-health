"""create users table

Revision ID: 30ae185fa647
Revises:
Create Date: 2025-05-30 16:17:12.548933+00:00

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = "30ae185fa647"
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "users",
        sa.Column("id", sa.Integer, primary_key=True),
        sa.Column("github_id", sa.String(50), nullable=False),
        sa.Column("gihub_username", sa.String(200)),
    )


def downgrade() -> None:
    op.drop_table("users")
