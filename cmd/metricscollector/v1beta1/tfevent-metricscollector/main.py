import grpc
import argparse
import api_pb2
import api_pb2_grpc
from pns import WaitOtherMainProcesses
from tfevent_loader import MetricsCollector
from logging import getLogger, StreamHandler, INFO

timeout_in_seconds = 60


def parse_options():
    parser = argparse.ArgumentParser(
        description='TF-Event MetricsCollector',
        add_help=True
    )
    parser.add_argument("-s", "--manager_server_addr",
                        type=str, default="katib-db-manager:6789")
    parser.add_argument("-t", "--trial_name", type=str, default="")
    parser.add_argument("-path", "--dir_path", type=str, default="/log")
    parser.add_argument("-m", "--metric_names", type=str, default="")
    parser.add_argument("-f", "--metric_filters", type=str, default="")
    opt = parser.parse_args()
    return opt


if __name__ == '__main__':
    logger = getLogger(__name__)
    handler = StreamHandler()
    handler.setLevel(INFO)
    logger.setLevel(INFO)
    logger.addHandler(handler)
    logger.propagate = False
    opt = parse_options()
    manager_server = opt.manager_server_addr.split(':')
    if len(manager_server) != 2:
        raise Exception("Invalid katib manager service address: %s" %
                        opt.manager_server_addr)

    WaitOtherMainProcesses(completed_marked_dir=opt.dir_path)

    mc = MetricsCollector(opt.metric_names.split(';'))
    observation_log = mc.parse_file(opt.dir_path)

    channel = grpc.beta.implementations.insecure_channel(
        manager_server[0], int(manager_server[1]))

    with api_pb2.beta_create_Manager_stub(channel) as client:
        logger.info("In " + opt.trial_name + " " +
                    str(len(observation_log.metric_logs)) + " metrics will be reported.")
        client.ReportObservationLog(api_pb2.ReportObservationLogRequest(
            trial_name=opt.trial_name,
            observation_log=observation_log
        ), timeout=timeout_in_seconds)
