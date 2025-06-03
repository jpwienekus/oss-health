from jose import jwt, JWTError
from datetime import datetime, timedelta
from fastapi import HTTPException
from app.configuration import settings


def create_token_pair(user_id: int, access_token: str):
    return jwt.encode(
        {
            "sub": str(user_id),
            "access_token": access_token,
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
        return payload
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")
