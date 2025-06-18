from typing import Awaitable, Callable, Dict, Optional

package_repository_resolvers: Dict[str, Callable[[str], Awaitable[Optional[str]]]] = {}


def register_package_repository_resolver(ecosystem: str):
    def decorator(func: Callable[[str], Awaitable[Optional[str]]]):
        package_repository_resolvers[ecosystem.lower()] = func
        return func

    return decorator
