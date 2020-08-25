import grpc
import argparse
import api_pb2
from pns import WaitMainProcesses
import const
from tfevent_loader import MetricsCollector
from logging import getLogger, StreamHandler, INFO

timeout_in_seconds = 60


def parse_options():
    parser = argparse.ArgumentParser(
        description='TF-Event MetricsCollector',
        add_help=True
    )

    parser.add_argument("-s", "--manager_server_addr", type=str, default="")
    parser.add_argument("-t", "--trial_name", type=str, default="")
    parser.add_argument("-path", "--metrics_file_dir", type=str, default=const.DEFAULT_METRICS_FILE_DIR)
    parser.add_argument("-m", "--metric_names", type=str, default="")
    parser.add_argument("-f", "--metric_filters", type=str, default="")
    parser.add_argument("-p", "--poll_interval", type=int, default=const.DEFAULT_POLL_INTERVAL)
    parser.add_argument("-timeout", "--timeout", type=int, default=const.DEFAULT_TIMEOUT)
    parser.add_argument("-w", "--wait_all", type=bool, default=const.DEFAULT_WAIT_ALL)

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

    WaitMainProcesses(
        pool_interval=opt.poll_interval,
        timout=opt.timeout,
        wait_all=opt.wait_all,
        completed_marked_dir=opt.metrics_file_dir)

    mc = MetricsCollector(opt.metric_names.split(';'))
    observation_log = mc.parse_file(opt.metrics_file_dir)

    channel = grpc.beta.implementations.insecure_channel(
        manager_server[0], int(manager_server[1]))

    with api_pb2.beta_create_DBManager_stub(channel) as client:
        logger.info("In " + opt.trial_name + " " +
                    str(len(observation_log.metric_logs)) + " metrics will be reported.")
        client.ReportObservationLog(api_pb2.ReportObservationLogRequest(
            trial_name=opt.trial_name,
            observation_log=observation_log
        ), timeout=timeout_in_seconds)
