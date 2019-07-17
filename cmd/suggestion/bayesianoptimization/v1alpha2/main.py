import grpc
from concurrent import futures

import time

from pkg.api.v1alpha2.python import api_pb2_grpc
from pkg.suggestion.v1alpha2.bayesian_service import BayesianService

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
DEFAULT_PORT = "0.0.0.0:6789"

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    api_pb2_grpc.add_SuggestionServicer_to_server(BayesianService(), server)
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
