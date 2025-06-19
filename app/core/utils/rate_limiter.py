import logging
import time
import redis
from config.settings import settings

logger = logging.getLogger(__name__)

class RedisRateLimiter:
    def __init__(self, name: str, limit: int, period: int = 60) -> None:
        self.name = name
        self.limit = limit
        self.period = period
        self.redis = redis.Redis.from_url(settings.broker_url)

    def _key(self):
        return f"ratelimit:{self.name}"

    def allow(self) -> bool:
        key = self._key()
        now = int(time.time())

        pipeline = self.redis.pipeline()
        pipeline.zremrangebyscore(key, 0, now - self.period)
        pipeline.zadd(key, { str(now): now})
        pipeline.zcard(key)
        pipeline.expire(key, self.period)
        _, _, count, _ = pipeline.execute()

        allowed = count <= self.limit
        # TODO: change to debug
        logger.info(f"[{self.name}] {count}/{self.limit} used")

        return allowed

    def wail_until_allowed(self):
        while not self.allow():
            # TODO: change to debug
            logger.info(f"[{self.name}] rate limit hit")
            time.sleep(1)

