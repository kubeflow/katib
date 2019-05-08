import grpc
from concurrent import futures

import time

from pkg.api.python import api_pb2_grpc
from pkg.suggestion.nasenvelopenet_service import EnvelopenetService
from pkg.suggestion.types import DEFAULT_PORT
from logging import getLogger, StreamHandler, INFO, DEBUG


_ONE_DAY_IN_SECONDS = 60 * 60 * 24


def serve():
    print("NAS Envelopenet Suggestion Service")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    api_pb2_grpc.add_SuggestionServicer_to_server(EnvelopenetService(), server)
    server.add_insecure_port(DEFAULT_PORT)
    print("Listening...")
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
