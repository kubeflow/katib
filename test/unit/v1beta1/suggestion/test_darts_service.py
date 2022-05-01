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
import unittest
import json

from pkg.apis.manager.v1beta1.python import api_pb2

from pkg.suggestion.v1beta1.nas.darts.service import DartsService


class TestDarts(unittest.TestCase):
    def setUp(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name['Suggestion']: DartsService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time())

    def test_get_suggestion(self):
        experiment = api_pb2.Experiment(
            name="darts-experiment",
            spec=api_pb2.ExperimentSpec(
                algorithm=api_pb2.AlgorithmSpec(
                    algorithm_name="darts",
                    algorithm_settings=[
                        api_pb2.AlgorithmSetting(
                            name="num_epoch",
                            value="10"
                        )
                    ],
                ),
                objective=api_pb2.ObjectiveSpec(
                    type=api_pb2.MAXIMIZE,
                    objective_metric_name="Best-Genotype"
                ),
                parallel_trial_count=1,
                max_trial_count=1,
                nas_config=api_pb2.NasConfig(
                    graph_config=api_pb2.GraphConfig(
                        num_layers=3,
                    ),
                    operations=api_pb2.NasConfig.Operations(
                        operation=[
                            api_pb2.Operation(
                                operation_type="separable_convolution",
                                parameter_specs=api_pb2.Operation.ParameterSpecs(
                                    parameters=[
                                        api_pb2.ParameterSpec(
                                            name="filter_size",
                                            parameter_type=api_pb2.CATEGORICAL,
                                            feasible_space=api_pb2.FeasibleSpace(
                                                max=None, min=None, list=["3", "5"])
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
            current_request_number=1,
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
        self.assertEqual(1, len(response.parameter_assignments))

        exp_algorithm_settings = {}
        for setting in experiment.spec.algorithm.algorithm_settings:
            exp_algorithm_settings[setting.name] = setting.value

        exp_num_layers = experiment.spec.nas_config.graph_config.num_layers

        exp_search_space = ["separable_convolution_3x3", "separable_convolution_5x5"]
        for pa in response.parameter_assignments[0].assignments:
            if (pa.name == "algorithm-settings"):
                algorithm_settings = pa.value.replace("\'", "\"")
                algorithm_settings = json.loads(algorithm_settings)
                self.assertDictContainsSubset(exp_algorithm_settings, algorithm_settings)
            elif (pa.name == "num-layers"):
                self.assertEqual(exp_num_layers, int(pa.value))
            elif (pa.name == "search-space"):
                search_space = pa.value.replace("\'", "\"")
                search_space = json.loads(search_space)
                self.assertEqual(exp_search_space, search_space)


if __name__ == '__main__':
    unittest.main()
