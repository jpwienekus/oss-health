# import os
# from dotenv import load_dotenv
from pydantic_settings import BaseSettings

# load_dotenv()

# GITHUB_CLIENT_ID = os.getenv("GITHUB_CLIENT_ID")
# GITHUB_CLIENT_SECRET = os.getenv("GITHUB_CLIENT_SECRET")
# SECRET_KEY = os.getenv("SECRET_KEY")
# ALGORITHM = "HS256"


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


settings = Settings()  # type: ignore
