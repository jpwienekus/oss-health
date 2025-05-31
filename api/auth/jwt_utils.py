from jose import jwt, JWTError
from datetime import datetime, timedelta
from configuration import settings
from fastapi import HTTPException

def create_token_pair(user_id):
    return jwt.encode({
        "sub": str(user_id),
        "exp": datetime.now() + timedelta(minutes=15),
    }, settings.secret_key, settings.algorithm)

def decode_token(token: str):
    try:
        payload = jwt.decode(token, settings.secret_key, algorithms=[settings.algorithm])
        return payload["sub"]
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")
