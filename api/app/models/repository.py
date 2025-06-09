from typing import List

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Repository(Base):
    __tablename__ = "repositories"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    github_id: Mapped[int] = mapped_column(unique=True, nullable=False)
    user_id: Mapped[int] = mapped_column(ForeignKey("user.id"))
    score: Mapped[int] = mapped_column(nullable=True)
    clone_url: Mapped[str] = mapped_column(nullable=True)

    dependencies: Mapped[List["Dependency"]] = relationship(  # type: ignore
        secondary="repository_dependency", back_populates="repositories"
    )
