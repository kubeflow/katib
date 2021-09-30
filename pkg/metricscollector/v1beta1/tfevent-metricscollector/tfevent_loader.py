# Copyright 2021 The Kubeflow Authors.
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
# For example, in the Kubeflow tf-operator tutorial for mnist with summary:
# https://github.com/kubeflow/tf-operator/blob/master/examples/tensorflow/mnist_with_summaries/mnist_with_summaries.py.
# The "accuracy" metric is saved under "train" and "test" directories.
# So in the Metrics Collector specification, please specify name of "train" or "test" directory.
# Check TFJob example for more information:
# https://github.com/kubeflow/katib/blob/master/examples/v1beta1/kubeflow-training-operator/tfjob-mnist-with-summaries.yaml#L16-L22


import tensorflow as tf
import os
from datetime import datetime
import rfc3339
import api_pb2
from logging import getLogger, StreamHandler, INFO
import const


class TFEventFileParser:
    def find_all_files(self, directory):
        for root, dirs, files in os.walk(directory):
            yield root
            for f in files:
                yield os.path.join(root, f)

    def parse_summary(self, tfefile, metrics):
        metric_logs = []
        for summary in tf.train.summary_iterator(tfefile):
            paths = tfefile.split("/")
            for v in summary.summary.value:
                for m in metrics:
                    tag = str(v.tag)
                    if len(paths) >= 2 and len(m.split("/")) >= 2:
                        tag = str(paths[-2]+"/" + v.tag)
                    if tag.startswith(m):
                        ml = api_pb2.MetricLog(
                            time_stamp=rfc3339.rfc3339(datetime.fromtimestamp(summary.wall_time)),
                            metric=api_pb2.Metric(
                                name=m,
                                value=str(v.simple_value)
                            )
                        )
                        metric_logs.append(ml)
        return metric_logs


class MetricsCollector:
    def __init__(self, metric_names):
        self.logger = getLogger(__name__)
        handler = StreamHandler()
        handler.setLevel(INFO)
        self.logger.setLevel(INFO)
        self.logger.addHandler(handler)
        self.logger.propagate = False
        self.metrics = metric_names
        self.parser = TFEventFileParser()

    def parse_file(self, directory):
        mls = []
        for f in self.parser.find_all_files(directory):
            if os.path.isdir(f):
                continue
            try:
                self.logger.info(f + " will be parsed.")
                mls.extend(self.parser.parse_summary(f, self.metrics))
            except Exception as e:
                self.logger.warning("Unexpected error: " + str(e))
                continue

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
                        name=self.metrics[0],
                        value=const.UNAVAILABLE_METRIC_VALUE
                    )
                )
            ]
            self.logger.info("Objective metric {} is not found in training logs, {} value is reported".format(
                self.metrics[0], const.UNAVAILABLE_METRIC_VALUE))

        return api_pb2.ObservationLog(metric_logs=mls)
