from typing import List

from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class DependencyRepository(Base):
    __tablename__ = "dependency_repository"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    github_url: Mapped[str | None] = mapped_column(nullable=False, unique=True)

    dependencies: Mapped[List["Dependency"]] = relationship(back_populates="dependency_repository")  # type: ignore # noqa: F821, E501
