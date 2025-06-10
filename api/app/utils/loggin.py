import logging
import sys


def configure_logging(log_level_str="INFO"):
    log_level = logging.DEBUG if log_level_str.upper() == "DEBUG" else logging.INFO

    formatter = logging.Formatter('%(levelname)-9s %(message)s')
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(formatter)

    root_logger = logging.getLogger()
    root_logger.handlers = []
    root_logger.addHandler(handler)
    root_logger.setLevel(log_level)
