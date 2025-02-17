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

import os
import tempfile
import unittest

import tensorboardX
import utils

METRIC_DIR_NAMES = ("train", "test")
METRIC_NAMES = ("accuracy", "loss")
QUALIFIED_METRIC_NAMES = tuple(
    f"{dir}/{metric}"
    for dir in METRIC_DIR_NAMES
    for metric in METRIC_NAMES
)

class TestTFEventMetricsCollector(unittest.TestCase):
    def test_parse_file(self):

        current_dir = os.path.dirname(os.path.abspath(__file__))
        logs_dir = os.path.join(current_dir, "testdata/tfevent-metricscollector/logs")


        metric_logs = utils.get_metric_logs(logs_dir, QUALIFIED_METRIC_NAMES)
        self.assertEqual(20, len(metric_logs))

        for log in metric_logs:
            actual = log["metric"]["name"]
            self.assertIn(actual, QUALIFIED_METRIC_NAMES)

        train_metric_logs = utils.get_metric_logs(
            os.path.join(logs_dir, "train"), METRIC_NAMES)
        self.assertEqual(10, len(train_metric_logs))

        for log in train_metric_logs:
            actual = log["metric"]["name"]
            self.assertIn(actual, METRIC_NAMES)

    def test_parse_file_with_tensorboardX(self):
        logs_dir = tempfile.mkdtemp()
        num_iters = 3

        for dir_name in METRIC_DIR_NAMES:
            with tensorboardX.SummaryWriter(os.path.join(logs_dir, dir_name)) as writer:
                for metric_name in METRIC_NAMES:
                    for iter in range(num_iters):
                        writer.add_scalar(metric_name, 0.1, iter)


        metric_logs = utils.get_metric_logs(logs_dir, QUALIFIED_METRIC_NAMES)
        self.assertEqual(num_iters * len(QUALIFIED_METRIC_NAMES), len(metric_logs))

        for log in metric_logs:
            actual = log["metric"]["name"]
            self.assertIn(actual, QUALIFIED_METRIC_NAMES)

        train_metric_logs = utils.get_metric_logs(
            os.path.join(logs_dir, "train"), METRIC_NAMES)
        self.assertEqual(num_iters * len(METRIC_NAMES), len(train_metric_logs))

        for log in train_metric_logs:
            actual = log["metric"]["name"]
            self.assertIn(actual, METRIC_NAMES)


if __name__ == '__main__':
    unittest.main()
