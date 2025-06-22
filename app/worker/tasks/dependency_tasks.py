# import asyncio
# from typing import Optional
#
# from sqlalchemy.ext.asyncio import AsyncSession
#
# from core.crud.dependency import resolve_pending_dependencies
# from worker.celery_worker import celery_app
# from worker.db_helpers import with_db_session
#
#
# @celery_app.task(name="worker.tasks.resolve_npm_github_urls")
# def resolve_npm_github_urls():
#     asyncio.run(resolve_github_urls("npm", 1, 0))
#
#
# @celery_app.task(name="worker.tasks.resolve_pypi_github_urls")
# def resolve_pypi_github_urls():
#     asyncio.run(resolve_github_urls("pypi", 1, 0))
#
#
# @with_db_session
# async def resolve_github_urls(
#     ecosystem: str,
#     batch_size: int,
#     offset: int,
#     db_session: Optional[AsyncSession] = None,
# ):
#     if not db_session:
#         return
#
#     await resolve_pending_dependencies(db_session, batch_size, offset, ecosystem)
