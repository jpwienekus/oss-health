from typing import List

from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Dependency(Base):
    __tablename__ = "dependencies"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    name: Mapped[str] = mapped_column(nullable=True)
    ecosystem: Mapped[str] = mapped_column(nullable=True)

    versions: Mapped[List["Version"]] = relationship(back_populates="dependency")  # type: ignore # noqa: F821
