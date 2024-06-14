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

from google.protobuf import json_format
from tfevent_loader import MetricsCollector


def get_metric_logs(logs_dir, metric_names):
    mc = MetricsCollector(metric_names)
    observation_log = mc.parse_file(logs_dir)
    dict_observation_log = json_format.MessageToDict(observation_log)
    return dict_observation_log["metricLogs"]
