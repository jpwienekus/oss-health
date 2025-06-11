# from sqlalchemy import select
# from sqlalchemy.ext.asyncio import AsyncSession
# from sqlalchemy.orm import selectinload
#
# from app.models import Dependency as DependencyDBModel
# from app.models import Version as VerionDBModel


# async def get_dependency_version(
#     db_session: AsyncSession, dependency_id: int, version: str
# ) -> VerionDBModel | None:
#     return (
#         await db_session.execute(
#             select(VerionDBModel)
#             .options(selectinload(VerionDBModel.dependencies))
#             .options(selectinload(VerionDBModel.vulnerabilities))
#             .where(VerionDBModel.version == version)
#             .where(DependencyDBModel.id == dependency_id)
#         )
#     ).scalar_one_or_none()
