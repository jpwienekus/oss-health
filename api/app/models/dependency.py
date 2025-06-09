from typing import List
from sqlalchemy.orm import Mapped, mapped_column, relationship
from . import Base


class Dependency(Base):
    __tablename__ = "dependencies"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    name: Mapped[str] = mapped_column(nullable=True)
    version: Mapped[str] = mapped_column(nullable=True)
    ecosystem: Mapped[str] = mapped_column(nullable=True)

    repositories: Mapped[List["Repository"]] = relationship( # type: ignore
        secondary="repository_dependency",
        back_populates="dependencies"
    )
