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

from pkg.apis.manager.v1beta1.python import api_pb2


def call_validate(test_server, experiment_spec):
    experiment = api_pb2.Experiment(name="validation-test", spec=experiment_spec)

    request = api_pb2.ValidateAlgorithmSettingsRequest(experiment=experiment)
    validate_algorithm_settings = test_server.invoke_unary_unary(
        method_descriptor=(api_pb2.DESCRIPTOR
                           .services_by_name['Suggestion']
                           .methods_by_name['ValidateAlgorithmSettings']),
        invocation_metadata={},
        request=request, timeout=1)

    response, metadata, code, details = validate_algorithm_settings.termination()

    return response, metadata, code, details
