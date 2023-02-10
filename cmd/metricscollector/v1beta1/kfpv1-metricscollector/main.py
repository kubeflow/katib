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
import os
from logging import INFO, StreamHandler, getLogger

import api_pb2
import const
import grpc
from metrics_loader import MetricsCollector
from pns import WaitMainProcesses

timeout_in_seconds = 60

# Next steps:
#
# - check is it is possible to mount the argo share
# - read the metrics from the tgz archive
# -
def parse_options():
    parser = argparse.ArgumentParser(
        description="KFP V1 MetricsCollector", add_help=True
    )

    # TODO (andreyvelich): Add early stopping flags.
    parser.add_argument("-s-db", "--db_manager_server_addr", type=str, default="")
    parser.add_argument("-t", "--pod_name", type=str, default="")
    parser.add_argument(
        "-path",
        "--metrics_file_dir",
        type=str,
        default=const.DEFAULT_METRICS_FILE_KFPV1_DIR,
    )
    parser.add_argument("-m", "--metric_names", type=str, default="")
    parser.add_argument("-o-type", "--objective_type", type=str, default="")
    parser.add_argument("-f", "--metric_filters", type=str, default="")
    parser.add_argument(
        "-p", "--poll_interval", type=int, default=const.DEFAULT_POLL_INTERVAL
    )
    parser.add_argument(
        "-timeout", "--timeout", type=int, default=const.DEFAULT_TIMEOUT
    )
    parser.add_argument(
        "-w", "--wait_all_processes", type=str, default=const.DEFAULT_WAIT_ALL_PROCESSES
    )
    parser.add_argument(
        "-fn",
        "--metrics_file_name",
        type=str,
        default=const.DEFAULT_METRICS_FILE_KFPV1_FILE,
    )

    opt = parser.parse_args()
    return opt


if __name__ == "__main__":
    logger = getLogger(__name__)
    handler = StreamHandler()
    handler.setLevel(INFO)
    logger.setLevel(INFO)
    logger.addHandler(handler)
    logger.propagate = False
    opt = parse_options()
    wait_all_processes = opt.wait_all_processes.lower() == "true"
    db_manager_server = opt.db_manager_server_addr.split(":")
    trial_name = '-'.join(opt.pod_name.split('-')[:-1])
    if len(db_manager_server) != 2:
        raise Exception(
            "Invalid Katib DB manager service address: %s" % opt.db_manager_server_addr
        )

    WaitMainProcesses(
        pool_interval=opt.poll_interval,
        timout=opt.timeout,
        wait_all=wait_all_processes,
        completed_marked_dir=None,
    )

    mc = MetricsCollector(opt.metric_names.split(";"))
    metrics_file = os.path.join(opt.metrics_file_dir, opt.metrics_file_name)
    observation_log = mc.parse_file(metrics_file)

    channel = grpc.beta.implementations.insecure_channel(
        db_manager_server[0], int(db_manager_server[1])
    )

    with api_pb2.beta_create_DBManager_stub(channel) as client:
        logger.info(
            "In "
            + trial_name
            + " "
            + str(len(observation_log.metric_logs))
            + " metrics will be reported."
        )
        client.ReportObservationLog(
            api_pb2.ReportObservationLogRequest(
                trial_name=trial_name, observation_log=observation_log
            ),
            timeout=timeout_in_seconds,
        )
