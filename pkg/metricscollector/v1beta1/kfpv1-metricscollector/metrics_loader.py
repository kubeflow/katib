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

# TFEventFileParser parses tfevent files and returns an ObservationLog of the metrics specified.
# When the event file is under a directory(e.g. test dir), please specify "{{dirname}}/{{metrics name}}"
# For example, in the Tensorflow MNIST Classification With Summaries:
# https://github.com/kubeflow/katib/blob/master/examples/v1beta1/trial-images/tf-mnist-with-summaries/mnist.py.
# The "accuracy" and "loss" metric is saved under "train" and "test" directories.
# So in the Metrics Collector specification, please specify name of "train" or "test" directory.
# Check TFJob example for more information:
# https://github.com/kubeflow/katib/blob/master/examples/v1beta1/kubeflow-training-operator/tfjob-mnist-with-summaries.yaml#L16-L22

from datetime import datetime
from logging import getLogger, StreamHandler, INFO
from typing import List
import json

import rfc3339
import api_pb2
from pkg.metricscollector.v1beta1.common import const


def parse_metrics(fn: str) -> List[api_pb2.MetricLog]:
    """Parse a kubeflow pipeline metrics file

    Args:
        fn (function): path to metrics file

    Returns:
        List[api_pb2.MetricLog]: A list of logged metrics
    """
    metrics = []
    with open(fn, "r") as f:
        metrics_dict = json.load(f)
        for m in metrics_dict["metrics"]:
            name = m["name"]
            value = m["numberValue"]
            ml = api_pb2.MetricLog(
                time_stamp=rfc3339.rfc3339(datetime.now()),
                metric=api_pb2.Metric(name=name, value=str(value)),
            )
            metrics.append(ml)
    return metrics


class MetricsCollector:
    def __init__(self, metric_names):
        self.logger = getLogger(__name__)
        handler = StreamHandler()
        handler.setLevel(INFO)
        self.logger.setLevel(INFO)
        self.logger.addHandler(handler)
        self.logger.propagate = False
        self.metrics = metric_names

    def parse_file(self, filename):
        self.logger.info(filename + " will be parsed.")
        mls = parse_metrics(filename)

        # Metrics logs must contain at least one objective metric value
        # Objective metric is located at first index
        is_objective_metric_reported = False
        for ml in mls:
            if ml.metric.name == self.metrics[0]:
                is_objective_metric_reported = True
                break
        # If objective metrics were not reported, insert unavailable value in the DB
        if not is_objective_metric_reported:
            mls = [
                api_pb2.MetricLog(
                    time_stamp=rfc3339.rfc3339(datetime.now()),
                    metric=api_pb2.Metric(
                        name=self.metrics[0], value=const.UNAVAILABLE_METRIC_VALUE
                    ),
                )
            ]
            self.logger.info(
                "Objective metric {} is not found in metrics file, {} value is reported".format(
                    self.metrics[0], const.UNAVAILABLE_METRIC_VALUE
                )
            )

        return api_pb2.ObservationLog(metric_logs=mls)
