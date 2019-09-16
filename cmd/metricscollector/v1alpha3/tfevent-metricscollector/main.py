import grpc
import argparse
import api_pb2
import api_pb2_grpc
from tfevent_loader import MetricsCollector
from logging import getLogger, StreamHandler, INFO

def parse_options():
    parser = argparse.ArgumentParser(
            description='TF-Event MetricsCollector',
            add_help = True
            )
    parser.add_argument("-a", "--manager_addr", type = str, default = "katib-manager")
    parser.add_argument("-p", "--manager_port", type = int, default = 6789 )
    parser.add_argument("-t", "--trial_name", type = str, default = "")
    parser.add_argument("-d", "--log_dir", type = str, default = "/log")
    parser.add_argument("-m", "--metric_names", type = str, default = "")
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
    
    mc = MetricsCollector(opt.metric_names.split(','))
    observation_log = mc.parse_file(opt.log_dir)

    channel = grpc.beta.implementations.insecure_channel(opt.manager_addr, opt.manager_port)
    with api_pb2.beta_create_Manager_stub(channel) as client:
        logger.info("In " + opt.trial_name + " " + str(len(observation_log.metric_logs)) + " metrics will be reported.")
        client.ReportObservationLog(api_pb2.ReportObservationLogRequest(
            trial_name=opt.trial_name,
            observation_log=observation_log
            ), 10)
