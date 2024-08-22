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

import logging
import time
from concurrent import futures

import grpc

from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.earlystopping.v1beta1.medianstop.service import MedianStopService

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
DEFAULT_PORT = "0.0.0.0:6788"

logger = logging.getLogger()
logging.basicConfig(level=logging.INFO)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service = MedianStopService()
    api_pb2_grpc.add_EarlyStoppingServicer_to_server(service, server)

    server.add_insecure_port(DEFAULT_PORT)
    logger.info("Start Median Stop service at address {}".format(DEFAULT_PORT))
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()
