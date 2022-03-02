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

import grpc
import grpc_testing

import pytest

from pkg.apis.manager.v1beta1.python import api_pb2

from pkg.suggestion.v1beta1.optuna.service import OptunaService


class TestOptuna:
    def setup_method(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name['Suggestion']: OptunaService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time())

    @pytest.mark.parametrize(
        ["algorithm_name", "algorithm_settings"],
        [
            ["tpe", {"n_startup_trials": "20", "n_ei_candidates": "10", "random_state": "71"}],
            ["multivariate-tpe", {"n_startup_trials": "20", "n_ei_candidates": "10", "random_state": "71"}],
            ["cmaes", {"restart_strategy": "ipop", "sigma": "2", "random_state": "71"}],
            ["random", {"random_state": "71"}],
        ],
    )
    def test_get_suggestion(self, algorithm_name, algorithm_settings):
        experiment = api_pb2.Experiment(
            name="test",
            spec=api_pb2.ExperimentSpec(
                algorithm=api_pb2.AlgorithmSpec(
                    algorithm_name=algorithm_name,
                    algorithm_settings=[
                        api_pb2.AlgorithmSetting(
                            name=name,
                            value=value
                        ) for name, value in algorithm_settings.items()
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
                        )
                    ]
                )
            )
        )

        # Run the first suggestion with no previous trials in the request
        request = api_pb2.GetSuggestionsRequest(
            experiment=experiment,
            trials=[],
            current_request_number=2,
        )

        get_suggestion = self.test_server.invoke_unary_unary(
            method_descriptor=(api_pb2.DESCRIPTOR
                               .services_by_name['Suggestion']
                               .methods_by_name['GetSuggestions']),
            invocation_metadata={},
            request=request, timeout=1)

        response, metadata, code, details = get_suggestion.termination()
        assert code == grpc.StatusCode.OK
        assert 2 == len(response.parameter_assignments)

        # Run the second suggestion with trials whose parameters are assigned in the first request
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
                        assignments=response.parameter_assignments[0].assignments
                    )
                ),
                status=api_pb2.TrialStatus(
                    condition=api_pb2.TrialStatus.TrialConditionType.SUCCEEDED,
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(
                                name="metric-1",
                                value="435"
                            ),
                            api_pb2.Metric(
                                name="metric-2",
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
                        assignments=response.parameter_assignments[1].assignments
                    )
                ),
                status=api_pb2.TrialStatus(
                    condition=api_pb2.TrialStatus.TrialConditionType.SUCCEEDED,
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(
                                name="metric-1",
                                value="123"
                            ),
                            api_pb2.Metric(
                                name="metric-2",
                                value="3028"
                            ),
                        ]
                    )
                )
            )
        ]

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
            request=request, timeout=1)

        response, metadata, code, details = get_suggestion.termination()
        assert code == grpc.StatusCode.OK
        assert 2 == len(response.parameter_assignments)


if __name__ == '__main__':
    pytest.main()
