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

import grpc
import grpc_testing
import pytest
import utils

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.optuna.service import OptunaService


class TestOptuna:
    def setup_method(self):
        services = {api_pb2.DESCRIPTOR.services_by_name["Suggestion"]: OptunaService()}

        self.test_server = grpc_testing.server_from_dictionary(
            services, grpc_testing.strict_real_time()
        )

    @pytest.mark.parametrize(
        ["algorithm_name", "algorithm_settings"],
        [
            [
                "tpe",
                {
                    "n_startup_trials": "20",
                    "n_ei_candidates": "10",
                    "random_state": "71",
                },
            ],
            [
                "multivariate-tpe",
                {
                    "n_startup_trials": "20",
                    "n_ei_candidates": "10",
                    "random_state": "71",
                },
            ],
            ["cmaes", {"restart_strategy": "ipop", "sigma": "2", "random_state": "71"}],
            ["random", {"random_state": "71"}],
            # ["grid", {"random_state": "71"}],
        ],
    )
    def test_get_suggestion(self, algorithm_name, algorithm_settings):
        experiment = api_pb2.Experiment(
            name="test",
            spec=api_pb2.ExperimentSpec(
                algorithm=api_pb2.AlgorithmSpec(
                    algorithm_name=algorithm_name,
                    algorithm_settings=[
                        api_pb2.AlgorithmSetting(name=name, value=value)
                        for name, value in algorithm_settings.items()
                    ],
                ),
                objective=api_pb2.ObjectiveSpec(type=api_pb2.MAXIMIZE, goal=0.9),
                parameter_specs=api_pb2.ExperimentSpec.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="param-1",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", list=[]
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-2",
                            parameter_type=api_pb2.CATEGORICAL,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["cat1", "cat2", "cat3"]
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-3",
                            parameter_type=api_pb2.DISCRETE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["3", "2", "6"]
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-4",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", step="1", list=[]
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-5",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", step="2", distribution=api_pb2.UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-6",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", distribution=api_pb2.UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-7",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", step="2", distribution=api_pb2.LOG_UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-8",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", distribution=api_pb2.LOG_UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-9",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="11", min="1", step="2.5", distribution=api_pb2.UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-10",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="11", min="1", step="2.5", distribution=api_pb2.LOG_UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-11",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", distribution=api_pb2.UNIFORM
                            ),
                        ),
                        api_pb2.ParameterSpec(
                            name="param-12",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="5", min="1", distribution=api_pb2.LOG_UNIFORM
                            ),
                        ),
                    ]
                ),
            ),
        )

        # Run the first suggestion with no previous trials in the request
        request = api_pb2.GetSuggestionsRequest(
            experiment=experiment,
            trials=[],
            current_request_number=2,
        )

        get_suggestion = self.test_server.invoke_unary_unary(
            method_descriptor=(
                api_pb2.DESCRIPTOR.services_by_name["Suggestion"].methods_by_name[
                    "GetSuggestions"
                ]
            ),
            invocation_metadata={},
            request=request,
            timeout=1,
        )

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
                        goal=0.9,
                    ),
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=response.parameter_assignments[0].assignments
                    ),
                ),
                status=api_pb2.TrialStatus(
                    condition=api_pb2.TrialStatus.TrialConditionType.SUCCEEDED,
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(name="metric-1", value="435"),
                            api_pb2.Metric(name="metric-2", value="5643"),
                        ]
                    ),
                ),
            ),
            api_pb2.Trial(
                name="test-234hs",
                spec=api_pb2.TrialSpec(
                    objective=api_pb2.ObjectiveSpec(
                        type=api_pb2.MAXIMIZE,
                        objective_metric_name="metric-2",
                        goal=0.9,
                    ),
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=response.parameter_assignments[1].assignments
                    ),
                ),
                status=api_pb2.TrialStatus(
                    condition=api_pb2.TrialStatus.TrialConditionType.SUCCEEDED,
                    observation=api_pb2.Observation(
                        metrics=[
                            api_pb2.Metric(name="metric-1", value="123"),
                            api_pb2.Metric(name="metric-2", value="3028"),
                        ]
                    ),
                ),
            ),
        ]

        request = api_pb2.GetSuggestionsRequest(
            experiment=experiment,
            trials=trials,
            current_request_number=2,
        )

        get_suggestion = self.test_server.invoke_unary_unary(
            method_descriptor=(
                api_pb2.DESCRIPTOR.services_by_name["Suggestion"].methods_by_name[
                    "GetSuggestions"
                ]
            ),
            invocation_metadata={},
            request=request,
            timeout=1,
        )

        response, metadata, code, details = get_suggestion.termination()
        assert code == grpc.StatusCode.OK
        assert 2 == len(response.parameter_assignments)

    @pytest.mark.parametrize(
        [
            "algorithm_name",
            "algorithm_settings",
            "max_trial_count",
            "parameters",
            "result",
        ],
        [
            # Invalid algorithm name
            ["invalid", {}, 1, [], grpc.StatusCode.INVALID_ARGUMENT],
            # [TPE] Valid case
            [
                "tpe",
                {"n_startup_trials": "5", "n_ei_candidates": "24", "random_state": "1"},
                100,
                [],
                grpc.StatusCode.OK,
            ],
            # [TPE] Invalid parameter name
            ["tpe", {"invalid": "5"}, 100, [], grpc.StatusCode.INVALID_ARGUMENT],
            # [TPE] Invalid n_startup_trials
            [
                "tpe",
                {"n_startup_trials": "-1"},
                100,
                [],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [TPE] Invalid n_ei_candidate
            [
                "tpe",
                {"n_ei_candidate": "-1"},
                100,
                [],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [TPE] Invalid random_state
            ["tpe", {"random_state": "-1"}, 100, [], grpc.StatusCode.INVALID_ARGUMENT],
            # [Multivariate-TPE] Valid case
            [
                "multivariate-tpe",
                {"n_startup_trials": "5", "n_ei_candidates": "24", "random_state": "1"},
                100,
                [],
                grpc.StatusCode.OK,
            ],
            # [CMAES] Valid case
            [
                "cmaes",
                {"restart_strategy": "ipop", "sigma": "0.1", "random_state": "10"},
                20,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="10", min="9", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.OK,
            ],
            # [CMAES] Invalid parameter name
            [
                "cmaes",
                {"invalid": "invalid", "sigma": "0.1"},
                100,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="10", min="9", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [CMAES] Invalid restart_strategy
            [
                "cmaes",
                {"restart_strategy": "invalid", "sigma": "0.1"},
                15,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="10", min="9", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [CMAES] Invalid sigma
            [
                "cmaes",
                {"restart_strategy": "None", "sigma": "-10"},
                55,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="10", min="9", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [CMAES] Invalid random_state
            [
                "cmaes",
                {"sigma": "0.2", "random_state": "-20"},
                25,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="10", min="9", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [CMAES] Invalid number of parameters
            [
                "cmaes",
                {"sigma": "0.2"},
                5,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    }
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [RANDOM] Valid Case
            [
                "random",
                {"random_state": "10"},
                23,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="10", min="9", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.OK,
            ],
            # [RANDOM] Invalid parameter name
            [
                "random",
                {"invalid": "invalid"},
                33,
                [],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [RANDOM] Invalid random_state
            [
                "random",
                {"random_state": "-1"},
                33,
                [],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [GRID] Valid Case
            [
                "grid",
                {"random_state": "10"},
                5,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.OK,
            ],
            # [GRID] Invalid parameter name
            ["grid", {"invalid": "invalid"}, 33, [], grpc.StatusCode.INVALID_ARGUMENT],
            # [GRID] Invalid random_state
            ["grid", {"random_state": "-1"}, 10, [], grpc.StatusCode.INVALID_ARGUMENT],
            # [GRID] Invalid feasible_space
            [
                "grid",
                {"random_state": "1"},
                26,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.DOUBLE,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
            # [GRID] Invalid max_trial_count
            [
                "grid",
                {"random_state": "1"},
                26,
                [
                    {
                        "name": "param-1",
                        "type": api_pb2.INT,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", list=[]
                        ),
                    },
                    {
                        "name": "param-2",
                        "type": api_pb2.DOUBLE,
                        "feasible_space": api_pb2.FeasibleSpace(
                            max="5", min="1", step="1", list=[]
                        ),
                    },
                ],
                grpc.StatusCode.INVALID_ARGUMENT,
            ],
        ],
    )
    def test_validate_algorithm_settings(
        self, algorithm_name, algorithm_settings, max_trial_count, parameters, result
    ):
        experiment_spec = api_pb2.ExperimentSpec(
            max_trial_count=max_trial_count,
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name=algorithm_name,
                algorithm_settings=[
                    api_pb2.AlgorithmSetting(name=name, value=value)
                    for name, value in algorithm_settings.items()
                ],
            ),
            parameter_specs=api_pb2.ExperimentSpec.ParameterSpecs(
                parameters=[
                    api_pb2.ParameterSpec(
                        name=param["name"],
                        parameter_type=param["type"],
                        feasible_space=param["feasible_space"],
                    )
                    for param in parameters
                ]
            ),
        )
        _, _, code, _ = utils.call_validate(self.test_server, experiment_spec)
        assert code == result


if __name__ == "__main__":
    pytest.main()
