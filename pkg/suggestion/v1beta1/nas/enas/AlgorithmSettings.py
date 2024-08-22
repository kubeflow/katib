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


algorithmSettingsValidator = {
    "controller_hidden_size": [int, [1, "inf"]],
    "controller_temperature": [float, [0, "inf"]],
    "controller_tanh_const": [float, [0, "inf"]],
    "controller_entropy_weight": [float, [0.0, "inf"]],
    "controller_baseline_decay": [float, [0.0, 1.0]],
    "controller_learning_rate": [float, [0.0, 1.0]],
    "controller_skip_target": [float, [0.0, 1.0]],
    "controller_skip_weight": [float, [0.0, "inf"]],
    "controller_train_steps": [int, [1, "inf"]],
    "controller_log_every_steps": [int, [1, "inf"]],
}
enableNoneSettingsList = [
    "controller_temperature",
    "controller_tanh_const",
    "controller_entropy_weight",
    "controller_skip_weight",
]


def parseAlgorithmSettings(settings_raw):

    algorithm_settings_default = {
        "controller_hidden_size": 64,
        "controller_temperature": 5.0,
        "controller_tanh_const": 2.25,
        "controller_entropy_weight": 1e-5,
        "controller_baseline_decay": 0.999,
        "controller_learning_rate": 5e-5,
        "controller_skip_target": 0.4,
        "controller_skip_weight": 0.8,
        "controller_train_steps": 50,
        "controller_log_every_steps": 10,
    }

    for setting in settings_raw:
        s_name = setting.name
        s_value = setting.value
        if s_value == "None":
            algorithm_settings_default[s_name] = None
        else:
            s_type = algorithmSettingsValidator[s_name][0]
            algorithm_settings_default[s_name] = s_type(s_value)

    return algorithm_settings_default
