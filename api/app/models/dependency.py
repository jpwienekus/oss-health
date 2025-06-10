from sqlalchemy.orm import Mapped, mapped_column

from . import Base


class Dependency(Base):
    __tablename__ = "dependencies"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    name: Mapped[str] = mapped_column(nullable=True)
    ecosystem: Mapped[str] = mapped_column(nullable=True)
