from typing import List

from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class Version(Base):
    __tablename__ = "versions"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    version: Mapped[str] = mapped_column(nullable=False)

    vulnerabilities: Mapped[List["Vulnerability"]] = relationship( # type: ignore # noqa: F821
        secondary="version_vulnerability", back_populates="versions"
    )
