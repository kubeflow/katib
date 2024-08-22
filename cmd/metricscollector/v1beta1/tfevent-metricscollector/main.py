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

import argparse
from logging import INFO, StreamHandler, getLogger

import api_pb2
import api_pb2_grpc
import const
import grpc
from pns import WaitMainProcesses
from tfevent_loader import MetricsCollector

timeout_in_seconds = 60


def parse_options():
    parser = argparse.ArgumentParser(
        description='TF-Event MetricsCollector',
        add_help=True
    )

    # TODO (andreyvelich): Add early stopping flags.
    parser.add_argument("-s-db", "--db_manager_server_addr", type=str, default="")
    parser.add_argument("-t", "--trial_name", type=str, default="")
    parser.add_argument("-path", "--metrics_file_dir", type=str, default=const.DEFAULT_METRICS_FILE_DIR)
    parser.add_argument("-m", "--metric_names", type=str, default="")
    parser.add_argument("-o-type", "--objective_type", type=str, default="")
    parser.add_argument("-f", "--metric_filters", type=str, default="")
    parser.add_argument("-p", "--poll_interval", type=int, default=const.DEFAULT_POLL_INTERVAL)
    parser.add_argument("-timeout", "--timeout", type=int, default=const.DEFAULT_TIMEOUT)
    parser.add_argument("-w", "--wait_all_processes", type=str, default=const.DEFAULT_WAIT_ALL_PROCESSES)

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
    wait_all_processes = opt.wait_all_processes.lower() == "true"
    db_manager_server = opt.db_manager_server_addr.split(':')
    if len(db_manager_server) != 2:
        raise Exception(
            f"Invalid Katib DB manager service address: {opt.db_manager_server_addr}"
        )

    WaitMainProcesses(
        pool_interval=opt.poll_interval,
        timout=opt.timeout,
        wait_all=wait_all_processes,
        completed_marked_dir=opt.metrics_file_dir,
    )

    mc = MetricsCollector(opt.metric_names.split(";"))
    observation_log = mc.parse_file(opt.metrics_file_dir)

    with grpc.insecure_channel(opt.db_manager_server_addr) as channel:
        stub = api_pb2_grpc.DBManagerStub(channel)
        logger.info(
            f"In {opt.trial_name} {str(len(observation_log.metric_logs))} metrics will be reported."
        )
        stub.ReportObservationLog(
            api_pb2.ReportObservationLogRequest(
                trial_name=opt.trial_name, observation_log=observation_log
            ),
            timeout=timeout_in_seconds,
        )
