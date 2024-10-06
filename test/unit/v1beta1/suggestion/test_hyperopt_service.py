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

import unittest

import grpc
import grpc_testing
import utils

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.hyperopt.service import HyperoptService
from pkg.suggestion.v1beta1.internal.constant import LOG_UNIFORM


class TestHyperopt(unittest.TestCase):
    def setUp(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name["Suggestion"]: HyperoptService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time())

    def test_get_suggestion(self):
        trials = [
            api_pb2.Trial(
                name="test-asfjh",
                spec=api_pb2.TrialSpec(
                    objective=api_pb2.ObjectiveSpec(
                        type=api_pb2.MAXIMIZE,
                        objective_metric_name="metric-2",
                        goal=0.9
                    ),
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=[
                            api_pb2.ParameterAssignment(
                                name="param-1",
                                value="2",
                            ),
                            api_pb2.ParameterAssignment(
                                name="param-2",
                                value="cat1",
                            ),
                            api_pb2.ParameterAssignment(
                                name="param-3",
                                value="2",
                            ),
                            api_pb2.ParameterAssignment(
                                name="param-4",
                                value="3.44",
                            )
                        ]
                    )
                ),
                status=api_pb2.TrialStatus(
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(
                                name="metric=1",
                                value="435"
                            ),
                            api_pb2.Metric(
                                name="metric=2",
                                value="5643"
                            ),
                        ]
                    )
                )
            ),
            api_pb2.Trial(
                name="test-234hs",
                spec=api_pb2.TrialSpec(
                    objective=api_pb2.ObjectiveSpec(
                        type=api_pb2.MAXIMIZE,
                        objective_metric_name="metric-2",
                        goal=0.9
                    ),
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=[
                            api_pb2.ParameterAssignment(
                                name="param-1",
                                value="3",
                            ),
                            api_pb2.ParameterAssignment(
                                name="param-2",
                                value="cat2",
                            ),
                            api_pb2.ParameterAssignment(
                                name="param-3",
                                value="6",
                            ),
                            api_pb2.ParameterAssignment(
                                name="param-4",
                                value="4.44",
                            )
                        ]
                    )
                ),
                status=api_pb2.TrialStatus(
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(
                                name="metric=1",
                                value="123"
                            ),
                            api_pb2.Metric(
                                name="metric=2",
                                value="3028"
                            ),
                        ]
                    )
                )
            )
        ]
        experiment = api_pb2.Experiment(
            name="test",
            spec=api_pb2.ExperimentSpec(
                algorithm=api_pb2.AlgorithmSpec(
                    algorithm_name="tpe",
                    algorithm_settings=[
                        api_pb2.AlgorithmSetting(
                            name="random_state",
                            value="10"
                        ),
                        api_pb2.AlgorithmSetting(
                            name="gamma",
                            value="0.25"
                        ),
                        api_pb2.AlgorithmSetting(
                            name="prior_weight",
                            value="1.0"
                        ),
                        api_pb2.AlgorithmSetting(
                            name="n_EI_candidates",
                            value="24"
                        ),
                    ],
                ),
                objective=api_pb2.ObjectiveSpec(
                    type=api_pb2.MAXIMIZE,
                    goal=0.9
                ),
                parameter_specs=api_pb2.ExperimentSpec.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="param-1",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[]),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-2",
                            parameter_type=api_pb2.CATEGORICAL,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["cat1", "cat2", "cat3"])
                        ),
                        api_pb2.ParameterSpec(
                            name="param-3",
                            parameter_type=api_pb2.DISCRETE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["3", "2", "6"])
                        ),
                        api_pb2.ParameterSpec(
                            name="param-4",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[])
                        ),
                        api_pb2.ParameterSpec(
                            name="param-5",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[], step="0.5", distribution=api_pb2.LOG_UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-6",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[], distribution=api_pb2.LOG_UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-7",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], step="0.8", distribution=api_pb2.UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-8",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], distribution=api_pb2.UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-9",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], step="0.8", distribution=api_pb2.NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-10",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], distribution=api_pb2.NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-11",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], step="0.8", distribution=api_pb2.LOG_NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-12",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], distribution=api_pb2.LOG_NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-13",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[], distribution=api_pb2.UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-14",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[], step="0.8", distribution=api_pb2.UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-15",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], distribution=api_pb2.LOG_UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-16",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="10", min="5", list=[], step="0.01", distribution=api_pb2.LOG_UNIFORM)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-17",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="100", min="5", list=[], distribution=api_pb2.NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-18",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="100", min="5", list=[], step="0.01", distribution=api_pb2.NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-19",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="64", min="32", distribution=api_pb2.LOG_NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-20",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="64", min="32", step="0.01", distribution=api_pb2.LOG_NORMAL)
                        ),
                        api_pb2.ParameterSpec(
                            name="param-21",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="64", min="32", step="0.01")
                        ),
                        api_pb2.ParameterSpec(
                            name="param-22",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="64", min="32", step="0.01")
                        )
                    ]
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
                               .services_by_name["Suggestion"]
                               .methods_by_name["GetSuggestions"]),
            invocation_metadata={},
            request=request, timeout=1)

        response, metadata, code, details = get_suggestion.termination()
        print(response.parameter_assignments)
        self.assertEqual(code, grpc.StatusCode.OK)
        self.assertEqual(2, len(response.parameter_assignments))

    def test_validate_algorithm_settings(self):
        # Valid cases.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="tpe",
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(
                        name="random_state",
                        value="10"
                    ),
                    api_pb2.AlgorithmSetting(
                        name="gamma",
                        value="0.25"
                    ),
                    api_pb2.AlgorithmSetting(
                        name="prior_weight",
                        value="1.0"
                    ),
                    api_pb2.AlgorithmSetting(
                        name="n_EI_candidates",
                        value="24"
                    ),
                ]
            )
        )

        _, _, code, _ = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.OK)

        # Invalid cases.
        # Unknown algorithm name.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="unknown"
            )
        )

        _, _, code, details = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "unknown algorithm name unknown")

        # Unknown algorithm setting name.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="random",
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(name="unknown_conf", value="1111")
                ]
            )
        )

        _, _, code, details = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "unknown setting unknown_conf for algorithm random")

        # Invalid gamma value.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="tpe",
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(name="gamma", value="1.5")
                ]
            )
        )

        _, _, code, details = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "gamma should be in the range of (0, 1)")

        # Invalid n_EI_candidates value.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="tpe",
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(name="n_EI_candidates", value="0")
                ]
            )
        )

        _, _, code, details = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "n_EI_candidates should be great than zero")

        # Invalid random_state value.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="tpe",
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(name="random_state", value="-1")
                ]
            )
        )

        _, _, code, details = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "random_state should be great or equal than zero")

        # Invalid prior_weight value.
        experiment_spec = api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="tpe",
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(name="prior_weight", value="aaa")
                ]
            )
        )

        _, _, code, details = utils.call_validate(self.test_server, experiment_spec)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertTrue(details.startswith("failed to validate prior_weight(aaa)"))


if __name__ == "__main__":
    unittest.main()
