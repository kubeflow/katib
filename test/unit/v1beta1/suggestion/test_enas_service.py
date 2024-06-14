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
import shutil
import unittest

import grpc
import grpc_testing
import pytest

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.nas.enas.service import EnasService


class TestEnas(unittest.TestCase):
    def setUp(self):
        services = {
            api_pb2.DESCRIPTOR.services_by_name['Suggestion']: EnasService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            services, grpc_testing.strict_real_time())

    def test_get_suggestion(self):
        trials = [
            api_pb2.Trial(
                name="first-trial",
                spec=api_pb2.TrialSpec(
                    objective=api_pb2.ObjectiveSpec(
                     type=api_pb2.MAXIMIZE,
                     objective_metric_name="Validation-Accuracy",
                     goal=0.99
                    ),
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=[
                            api_pb2.ParameterAssignment(
                                name="architecture",
                                value="[[3], [0, 1], [0, 0, 1], [2, 1, 0, 0]]",
                            ),
                            api_pb2.ParameterAssignment(
                                name="nn_config",
                                value="{'num_layers': 4}",
                            ),
                        ]
                    )
                ),
                status=api_pb2.TrialStatus(
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(
                                name="Validation-Accuracy",
                                value="0.88"
                            ),
                        ]
                    ),
                    condition=api_pb2.TrialStatus.TrialConditionType.SUCCEEDED,

                )
            ),
            api_pb2.Trial(
                name="second-trial",
                spec=api_pb2.TrialSpec(
                    objective=api_pb2.ObjectiveSpec(
                     type=api_pb2.MAXIMIZE,
                     objective_metric_name="Validation-Accuracy",
                     goal=0.99
                    ),
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=[
                            api_pb2.ParameterAssignment(
                                name="architecture",
                                value="[[1], [0, 1], [2, 1, 1], [2, 1, 1, 0]]",
                            ),
                            api_pb2.ParameterAssignment(
                                name="nn_config",
                                value="{'num_layers': 4}",
                            ),
                        ],
                    )
                ),
                status=api_pb2.TrialStatus(
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(
                                name="Validation-Accuracy",
                                value="0.84"
                            ),
                        ]
                    ),
                    condition=api_pb2.TrialStatus.TrialConditionType.SUCCEEDED,
                )
            )
        ]
        experiment = api_pb2.Experiment(
            name="enas-experiment",
            spec=api_pb2.ExperimentSpec(
                algorithm=api_pb2.AlgorithmSpec(
                    algorithm_name="enas",
                ),
                objective=api_pb2.ObjectiveSpec(
                    type=api_pb2.MAXIMIZE,
                    goal=0.9,
                    objective_metric_name="Validation-Accuracy"
                ),
                parallel_trial_count=2,
                max_trial_count=10,
                nas_config=api_pb2.NasConfig(
                    graph_config=api_pb2.GraphConfig(
                        num_layers=4,
                        input_sizes=[32, 32, 8],
                        output_sizes=[10]
                    ),
                    operations=api_pb2.NasConfig.Operations(
                        operation=[
                            api_pb2.Operation(
                                operation_type="convolution",
                                parameter_specs=api_pb2.Operation.ParameterSpecs(
                                    parameters=[
                                        api_pb2.ParameterSpec(
                                            name="filter_size",
                                            parameter_type=api_pb2.CATEGORICAL,
                                            feasible_space=api_pb2.FeasibleSpace(
                                                max=None, min=None, list=["5"])
                                        ),
                                        api_pb2.ParameterSpec(
                                            name="num_filter",
                                            parameter_type=api_pb2.CATEGORICAL,
                                            feasible_space=api_pb2.FeasibleSpace(
                                                max=None, min=None, list=["128"])
                                        ),
                                        api_pb2.ParameterSpec(
                                            name="stride",
                                            parameter_type=api_pb2.CATEGORICAL,
                                            feasible_space=api_pb2.FeasibleSpace(
                                                max=None, min=None, list=["1", "2"])
                                        ),
                                    ]
                                )
                            ),
                            api_pb2.Operation(
                                operation_type="reduction",
                                parameter_specs=api_pb2.Operation.ParameterSpecs(
                                    parameters=[
                                        api_pb2.ParameterSpec(
                                            name="reduction_type",
                                            parameter_type=api_pb2.CATEGORICAL,
                                            feasible_space=api_pb2.FeasibleSpace(
                                                max=None, min=None, list=["max_pooling"])
                                        ),
                                        api_pb2.ParameterSpec(
                                            name="pool_size",
                                            parameter_type=api_pb2.INT,
                                            feasible_space=api_pb2.FeasibleSpace(
                                                min="2", max="3", step="1", list=[])
                                        ),
                                    ]
                                )
                            ),
                        ],
                    )
                )
            )
        )

        request = api_pb2.GetSuggestionsRequest(
            experiment=experiment,
            trials=trials,
            current_request_number=2,
        )

        get_suggestion = self.test_server.invoke_unary_unary(
            method_descriptor=(api_pb2.DESCRIPTOR
                               .services_by_name['Suggestion']
                               .methods_by_name['GetSuggestions']),
            invocation_metadata={},
            request=request, timeout=100)

        response, metadata, code, details = get_suggestion.termination()
        print(response.parameter_assignments)
        self.assertEqual(code, grpc.StatusCode.OK)
        self.assertEqual(2, len(response.parameter_assignments))


@pytest.fixture(scope='function', autouse=True)
def tear_down():
    yield
    working_dir = os.getcwd()
    target_path = os.path.join(working_dir, "ctrl_cache")
    if os.path.isdir(target_path):
        shutil.rmtree(target_path)


if __name__ == '__main__':
    unittest.main()
