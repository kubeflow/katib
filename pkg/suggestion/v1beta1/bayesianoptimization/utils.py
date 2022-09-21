# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import logging
from logging import getLogger, StreamHandler


FORMAT = '%(asctime)-15s Experiment %(experiment_name)s %(message)s'
LOG_LEVEL = os.environ.get("LOG_LEVEL", "INFO")


def get_logger(name=__name__):
    logger = getLogger(name)
    logging.basicConfig(format=FORMAT)
    handler = StreamHandler()
    logger.setLevel(LOG_LEVEL)
    logger.addHandler(handler)
    logger.propagate = False
    return logger
