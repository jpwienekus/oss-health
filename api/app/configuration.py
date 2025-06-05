from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    database_user: str = "dev-user"
    database_password: str = "password"
    database_url: str = "localhost"
    database_name: str = "dev_db"
    echo_sql: bool = True
    test: bool = False
    project_name: str = "My FastAPI project"
    oauth_token_secret: str = "my_dev_secret"
    log_level: str = "DEBUG"
    github_client_id: str = ""
    github_client_secret: str = ""
    secret_key: str = ""
    algorithm: str = "HS256"

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()  # type: ignore
