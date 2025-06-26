from datetime import datetime
from typing import List

from sqlalchemy import DateTime, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Repository(Base):
    __tablename__ = "repositories"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    github_id: Mapped[int] = mapped_column(unique=True, nullable=False)
    user_id: Mapped[int] = mapped_column(ForeignKey("user.id"))
    score: Mapped[int] = mapped_column(nullable=True)
    url: Mapped[str] = mapped_column(nullable=True)
    last_scanned_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), nullable=True
    )
    scan_day: Mapped[int] = mapped_column(nullable=False)
    scan_hour: Mapped[int] = mapped_column(nullable=False)
    scan_status: Mapped[str] = mapped_column(nullable=False)
    error_message: Mapped[str] = mapped_column(nullable=True)

    dependency_versions: Mapped[List["RepositoryDependencyVersion"]] = relationship(  # type: ignore # noqa: F821
        back_populates="repository", cascade="all, delete-orphan"
    )
