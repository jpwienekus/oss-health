from celery import Celery
from config.settings import settings

# from core.utils.loggin import configure_logging

# configure_logging(log_level_str="INFO")

celery_app = Celery(
    "dependency_monitor",
    broker=f"{settings.broker_url}/0",
    backend=f"{settings.broker_url}/1",
)

celery_app.conf.task_routes = {
    "worker.tasks.resolve_npm_github_urls": {"queue": "npm"},
    "worker.tasks.resolve_pypi_github_urls": {"queue": "pypi"}
}

