import time
from concurrent import futures

import grpc

from pkg.api.python import api_pb2_grpc
from pkg.suggestion.katib_suggestion.rpc_service import SuggestionService
from pkg.suggestion.types import DEFAULT_PORT

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    api_pb2_grpc.add_SuggestionServicer_to_server(SuggestionService(), server)
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
