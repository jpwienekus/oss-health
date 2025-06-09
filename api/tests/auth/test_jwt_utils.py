
import pytest
from fastapi import HTTPException

from app.auth.jwt_utils import create_token_pair, decode_token


def test_create_and_decode_token():
    user_id = 123
    token = create_token_pair(user_id)

    # Should return the original user_id
    assert decode_token(token) == user_id


def test_invalid_token_raises():
    # An obviously invalid token
    with pytest.raises(HTTPException) as exc_info:
        decode_token("not.a.valid.token")

    assert exc_info.value.status_code == 401
    assert "Invalid token" in str(exc_info.value.detail)


def test_tampered_token_raises():
    user_id = 123
    token = create_token_pair(user_id)

    # Tamper with token (e.g., remove last character)
    tampered_token = token[:-1]

    with pytest.raises(HTTPException):
        decode_token(tampered_token)
