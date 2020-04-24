import grpc
import time
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.apis.manager.health.python import health_pb2_grpc
from pkg.suggestion.v1alpha3.chocolate.chocolate_service import ChocolateService
from concurrent import futures

_ONE_DAY_IN_SECONDS = 60 * 60 * 24
DEFAULT_PORT = "0.0.0.0:6789"


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service = ChocolateService()
    api_pb2_grpc.add_SuggestionServicer_to_server(service, server)
    health_pb2_grpc.add_HealthServicer_to_server(service, server)
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
