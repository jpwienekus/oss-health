from celery import Celery
from config.settings import settings

from core.utils.loggin import configure_logging

configure_logging(log_level_str="INFO")

celery_app = Celery(
    "dependency_monitor",
)

celery_app.config_from_object({
    "broker_url": f"{settings.broker_url}/0",
    "result_backend": f"{settings.broker_url}/1",
    "worker_hijack_root_logger": False,

})

celery_app.conf.task_routes = {
    "worker.tasks.resolve_npm_github_urls": {"queue": "npm"},
    "worker.tasks.resolve_pypi_github_urls": {"queue": "pypi"}
}

