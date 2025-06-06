from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from . import Base


class Repository(Base):
    __tablename__ = "repositories"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    github_id: Mapped[int] = mapped_column(unique=True, nullable=False)
    user_id: Mapped[int] = mapped_column(ForeignKey("user.id"))
    score: Mapped[int] = mapped_column(nullable=True)
