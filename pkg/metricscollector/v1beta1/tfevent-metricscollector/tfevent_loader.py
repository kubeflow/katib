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

import tensorflow as tf
from tensorboard.backend.event_processing.event_accumulator import EventAccumulator, TensorEvent
from tensorboard.backend.event_processing.tag_types import TENSORS
import os
import rfc3339
from datetime import datetime
import api_pb2
from logging import getLogger, StreamHandler, INFO
from pkg.metricscollector.v1beta1.common import const


class TFEventFileParser:
    def __init__(self, metric_names):
        self.metric_names = metric_names

    @staticmethod
    def find_all_files(directory):
        for root, dirs, files in os.walk(directory):
            for f in files:
                yield os.path.join(root, f)

    def parse_summary(self, tfefile):
        metric_logs = []
        event_accumulator = EventAccumulator(tfefile, size_guidance={TENSORS: 0})
        event_accumulator.Reload()
        for tag in event_accumulator.Tags()[TENSORS]:
            for m in self.metric_names:

                tfefile_parent_dir = os.path.dirname(m) if len(m.split("/")) >= 2 else os.path.dirname(tfefile)
                basedir_name = os.path.dirname(tfefile)
                if not tag.startswith(m.split("/")[-1]) or not basedir_name.endswith(tfefile_parent_dir):
                    continue

                for tensor in event_accumulator.Tensors(tag):
                    ml = api_pb2.MetricLog(
                        time_stamp=rfc3339.rfc3339(datetime.fromtimestamp(tensor.wall_time)),
                        metric=api_pb2.Metric(
                            name=m,
                            value=str(tf.make_ndarray(tensor.tensor_proto))
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
        self.parser = TFEventFileParser(self.metrics)

    def parse_file(self, directory):
        mls = []
        for f in self.parser.find_all_files(directory):
            if os.path.isdir(f):
                continue
            try:
                self.logger.info(f + " will be parsed.")
                mls.extend(self.parser.parse_summary(f))
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
