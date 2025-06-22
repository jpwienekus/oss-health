import logging
import sys
from datetime import datetime

RESET = "\x1b[0m"
BOLD = "\x1b[1m"
DIM = "\x1b[2m"
COLORS = {
    "DEBUG": "\x1b[34m",
    "INFO": "\x1b[32m",
    "WARNING": "\x1b[33m",
    "ERROR": "\x1b[31m",
    "CRITICAL": "\x1b[31m",
}


class ColoredLoggerFormatter(logging.Formatter):
    def format(self, record):
        level_color = COLORS.get(record.levelname, "")
        timestamp = datetime.fromtimestamp(record.created).strftime("%Y-%m-%d %H:%M:%S")
        level = f"{level_color}{record.levelname:<8}{RESET}"

        return f"{level}{DIM}{timestamp}{RESET} {BOLD} {record.name}:{RESET} {record.getMessage()}"


def configure_logging(log_level_str="INFO"):
    level_map = {
        "DEBUG": logging.DEBUG,
        "INFO": logging.INFO,
        "WARNING": logging.WARNING,
        "ERROR": logging.ERROR,
    }
    log_level = level_map.get(log_level_str.upper(), logging.INFO)

    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(ColoredLoggerFormatter())

    root_logger = logging.getLogger()
    root_logger.handlers.clear()
    root_logger.propagate = False
    root_logger.setLevel(log_level)
    root_logger.addHandler(handler)

    for name in ("", "uvicorn", "uvicorn.error", "uvicorn.access", "fastapi", "celery"):
        logger = logging.getLogger(name)
        logger.handlers.clear()
        logger.propagate = False
        logger.setLevel(log_level)
        logger.addHandler(handler)
