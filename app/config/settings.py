import os

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    database_url: str = "postgresql+asyncpg://dev-user:password@localhost:5432/dev_db"
    echo_sql: bool = False
    test: bool = False
    project_name: str = "My FastAPI project"
    oauth_token_secret: str = "my_dev_secret"
    log_level: str = "INFO"
    github_client_id: str = ""
    github_client_secret: str = ""
    secret_key: str = ""
    algorithm: str = "HS256"
    environment: str = "development"
    allowed_origins: str = "http://localhost:5173"

    model_config = {
        "env_file": os.getenv("ENV_FILE", ".env"),
        "env_file_encoding": "utf-8",
    }


settings = Settings()  # type: ignore
