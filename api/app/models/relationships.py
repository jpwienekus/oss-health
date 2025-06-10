
from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from . import Base


class RepositoryDependencyVersion(Base):
    __tablename__ = "repository_dependency_version"

    repository_id: Mapped[int] = mapped_column(ForeignKey("repositories.id"), primary_key=True)
    dependency_id: Mapped[int] = mapped_column(ForeignKey("dependencies.id"), primary_key=True)
    version_id: Mapped[int] = mapped_column(ForeignKey("versions.id"), primary_key=True)

    repository: Mapped["Repository"] = relationship(back_populates="dependency_versions") # type: ignore # noqa: F821
    dependency: Mapped["Dependency"] = relationship() # type: ignore # noqa: F821
    version: Mapped["Version"] = relationship() # type: ignore # noqa: F821


class VersionVulnerability(Base):
    __tablename__ = "version_vulnerability"

    version_id: Mapped[int] = mapped_column(ForeignKey("versions.id"), primary_key=True)
    vulnerability_id: Mapped[int] = mapped_column(
        ForeignKey("vulnerabilities.id"), primary_key=True
    )

