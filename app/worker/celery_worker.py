from celery import Celery

# from core.utils.loggin import configure_logging

# configure_logging(log_level_str="INFO")

celery_app = Celery(
    "dependency_monitor",
    broker="redis://localhost:6379/0",
    backend="redis://localhost:6379/1",
)

celery_app.conf.task_routes = {
    "worker.tasks.resolve_npm_github_urls": {"queue": "npm"},
    "worker.tasks.resolve_pypi_github_urls": {"queue": "pypi"}
}

