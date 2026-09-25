# Copyright 2026 The Kubeflow Authors.
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


def _set_algorithm_setting(experiment_manifest, name, value):
    algorithm_settings = (
        experiment_manifest["spec"]["algorithm"].setdefault("algorithmSettings", [])
    )
    for setting in algorithm_settings:
        if setting["name"] == name:
            setting["value"] = value
            return
    algorithm_settings.append({"name": name, "value": value})


def _append_trial_command_argument(experiment_manifest, argument):
    command = experiment_manifest["spec"]["trialTemplate"]["trialSpec"]["spec"][
        "template"
    ]["spec"]["containers"][0]["command"]
    if argument not in command:
        command.append(argument)


def apply_e2e_experiment_overrides(experiment_manifest):
    experiment_name = experiment_manifest["metadata"]["name"]

    if experiment_name == "darts-cpu":
        _set_algorithm_setting(experiment_manifest, "use_synthetic_data", "true")
        _set_algorithm_setting(experiment_manifest, "number_of_examples", "64")
        _set_algorithm_setting(experiment_manifest, "batch_size", "16")
        _set_algorithm_setting(experiment_manifest, "num_workers", "0")
    elif experiment_name == "enas-cpu":
        _append_trial_command_argument(experiment_manifest, "--use_synthetic_data")
        _append_trial_command_argument(experiment_manifest, "--number_of_examples=64")
