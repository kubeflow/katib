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

import os
import unittest
import utils


class TestTFEventMetricsCollector(unittest.TestCase):
    def test_parse_file(self):

        current_dir = os.path.dirname(os.path.abspath(__file__))
        logs_dir = os.path.join(current_dir, "testdata/tfevent-metricscollector/logs")

        # Metric format is "{{dirname}}/{{metrics name}}"
        metric_names = ["train/accuracy", "train/loss", "test/loss", "test/accuracy"]
        metric_logs = utils.get_metric_logs(logs_dir, metric_names)
        self.assertEqual(20, len(metric_logs))

        for log in metric_logs:
            actual = log["metric"]["name"]
            self.assertIn(actual, metric_names)

        # Metric format is "{{metrics name}}"
        metric_names = ["accuracy", "loss"]
        metrics_file_dir = os.path.join(logs_dir, "train")
        metric_logs = utils.get_metric_logs(metrics_file_dir, metric_names)
        self.assertEqual(10, len(metric_logs))

        for log in metric_logs:
            actual = log["metric"]["name"]
            self.assertIn(actual, metric_names)


if __name__ == '__main__':
    unittest.main()
