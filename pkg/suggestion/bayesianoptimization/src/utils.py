import os
import logging

from logging import getLogger, StreamHandler


FORMAT = '%(asctime)-15s StudyID %(studyid)s %(message)s'
LOG_LEVEL = os.environ.get("LOG_LEVEL", "INFO")


def get_logger(name=__name__):
    logger = getLogger(name)
    logging.basicConfig(format=FORMAT)
    handler = StreamHandler()
    logger.setLevel(LOG_LEVEL)
    logger.addHandler(handler)
    logger.propagate = False
    return logger
