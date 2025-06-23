from datetime import datetime, timedelta

from fastapi import HTTPException
from jose import JWTError, jwt

from config.settings import settings


def create_token_pair(user_id: int):
    return jwt.encode(
        {
            "sub": str(user_id),
            "exp": datetime.now() + timedelta(minutes=15),
        },
        settings.secret_key,
        settings.algorithm,
    )


def decode_token(token: str):
    try:
        payload = jwt.decode(
            token, settings.secret_key, algorithms=[settings.algorithm]
        )
        return int(payload["sub"])
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")
