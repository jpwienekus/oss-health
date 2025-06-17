from typing import List

from datetime import datetime
from sqlalchemy import ForeignKey, DateTime
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Repository(Base):
    __tablename__ = "repositories"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    github_id: Mapped[int] = mapped_column(unique=True, nullable=False)
    user_id: Mapped[int] = mapped_column(ForeignKey("user.id"))
    score: Mapped[int] = mapped_column(nullable=True)
    clone_url: Mapped[str] = mapped_column(nullable=True)
    scanned_date: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        nullable=True
    )

    dependency_versions: Mapped[List["RepositoryDependencyVersion"]] = relationship(  # type: ignore # noqa: F821
        back_populates="repository", cascade="all, delete-orphan"
    )
