from datetime import datetime
from typing import List, Optional

from sqlalchemy import DateTime, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Dependency(Base):
    __tablename__ = "dependencies"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    name: Mapped[str] = mapped_column(nullable=True)
    ecosystem: Mapped[str] = mapped_column(nullable=True)
    github_url_resolved: Mapped[bool] = mapped_column(nullable=False, default=False)
    github_url_checked_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), nullable=True
    )
    github_url_resolve_failed: Mapped[bool] = mapped_column(
        nullable=False, default=False, index=True
    )
    github_url_resolve_failed_reason: Mapped[str] = mapped_column(nullable=True)
    dependency_repository_id: Mapped[Optional[int]] = mapped_column(
        ForeignKey("dependency_repository.id"), nullable=True
    )

    versions: Mapped[List["Version"]] = relationship(back_populates="dependency")  # type: ignore # noqa: F821
    dependency_repository: Mapped[Optional["DependencyRepository"]] = relationship(back_populates="dependencies")  # type: ignore # noqa: F821, E501
