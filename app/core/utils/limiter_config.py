
from core.utils.rate_limiter import RedisRateLimiter


LIMITERS = {
    "pypi": RedisRateLimiter("pypi.org", limit=60, period=60),
    "npm" : RedisRateLimiter("npmjs.org", limit=30, period=60)
}
