from typing import List

from sqlalchemy import ForeignKey, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Version(Base):
    __tablename__ = "versions"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    version: Mapped[str] = mapped_column(nullable=False)
    dependency_id: Mapped[int] = mapped_column(
        ForeignKey("dependencies.id"), nullable=False
    )

    dependency: Mapped["Dependency"] = relationship(back_populates="versions")  # type: ignore # noqa: F821
    vulnerabilities: Mapped[List["Vulnerability"]] = relationship(  # type: ignore # noqa: F821
        secondary="version_vulnerability", back_populates="versions"
    )

    __table_args__ = (UniqueConstraint("version", "dependency_id"),)
