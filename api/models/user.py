from sqlalchemy.orm import Mapped, mapped_column
from . import Base


class User(Base):
    __tablename__ = "user"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True, index=True)
    github_id: Mapped[str] = mapped_column(unique=True, nullable=False)
    github_username: Mapped[str] = mapped_column(nullable=False)
