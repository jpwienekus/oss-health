from datetime import datetime
from typing import List

from sqlalchemy import DateTime
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Dependency(Base):
    __tablename__ = "dependencies"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    name: Mapped[str] = mapped_column(nullable=True)
    ecosystem: Mapped[str] = mapped_column(nullable=True)
    github_url: Mapped[str | None] = mapped_column(nullable=True)
    github_url_resolved: Mapped[bool] = mapped_column(nullable=False, default=False)
    github_url_checked_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), nullable=True
    )

    versions: Mapped[List["Version"]] = relationship(back_populates="dependency")  # type: ignore # noqa: F821
