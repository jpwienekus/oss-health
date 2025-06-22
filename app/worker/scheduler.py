# import asyncio
# from typing import Optional
#
# from apscheduler.schedulers.asyncio import AsyncIOScheduler
# from apscheduler.triggers.cron import CronTrigger
# from sqlalchemy.ext.asyncio import AsyncSession
#
# from config.settings import settings
# from core.utils.loggin import configure_logging
# from core.crud.dependency import resolve_pending_dependencies
# from worker.db_helpers import with_db_session
# from worker.tasks.dependency_tasks import resolve_github_urls
#
# configure_logging(log_level_str="INFO")
# scheduler = AsyncIOScheduler(timezone="UTC")
#
# # Schedule each ecosystem task
# scheduler.add_job(
#     resolve_github_urls,
#     trigger=CronTrigger(minute="*/1"),  # every minute, like your Celery config
#     kwargs={"ecosystem": "npm", "batch_size": 1, "offset": 0},
#     id="resolve_npm_github_urls",
#     replace_existing=True,
# )
#
# scheduler.add_job(
#     resolve_github_urls,
#     trigger=CronTrigger(minute="*/1"),  # every minute, like your Celery config
#     kwargs={"ecosystem": "pypi", "batch_size": 1, "offset": 0},
#     id="resolve_npm_github_urls",
#     replace_existing=True,
# )
#
# if __name__ == "__main__":
#
#     async def main():
#         scheduler.start()
#         # keep the loop running forever
#         while True:
#             await asyncio.sleep(3600)
#
#     asyncio.run(main())
