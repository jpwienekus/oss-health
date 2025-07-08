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
    scan_status: Mapped[str] = mapped_column(nullable=False)
    scanned_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=True)
    error_message: Mapped[str] = mapped_column(nullable=True)
    dependency_repository_id: Mapped[Optional[int]] = mapped_column(
        ForeignKey("dependency_repository.id"), nullable=True
    )

    versions: Mapped[List["Version"]] = relationship(back_populates="dependency")  # type: ignore # noqa: F821
    dependency_repository: Mapped[Optional["DependencyRepository"]] = relationship(back_populates="dependencies")  # type: ignore # noqa: F821, E501
