from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from . import Base


class RepositoryDependency(Base):
    __tablename__ = "repository_dependency"
    
    repository_id: Mapped[int] = mapped_column(ForeignKey("repositories.id"), primary_key=True)
    dependency_id: Mapped[int] = mapped_column(ForeignKey("dependencies.id"), primary_key=True)

