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

import json
import unittest

import grpc
import grpc_testing

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.nas.darts.service import (
    DartsService,
    validate_algorithm_settings,
)


class TestDarts(unittest.TestCase):
    def setUp(self):
        services = {
            api_pb2.DESCRIPTOR.services_by_name['Suggestion']: DartsService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            services, grpc_testing.strict_real_time())

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
            if pa.name == "algorithm-settings":
                algorithm_settings = pa.value.replace("\'", "\"")
                algorithm_settings = json.loads(algorithm_settings)
                self.assertDictContainsSubset(exp_algorithm_settings, algorithm_settings)
            elif pa.name == "num-layers":
                self.assertEqual(exp_num_layers, int(pa.value))
            elif pa.name == "search-space":
                search_space = pa.value.replace("\'", "\"")
                search_space = json.loads(search_space)
                self.assertEqual(exp_search_space, search_space)

    def test_validate_algorithm_spec(self):

        # Valid Case
        valid = [
            api_pb2.AlgorithmSetting(name="num_epoch", value="10"),
            api_pb2.AlgorithmSetting(name="w_lr", value="0.01"),
            api_pb2.AlgorithmSetting(name="w_lr_min", value="0.01"),
            api_pb2.AlgorithmSetting(name="alpha_lr", value="0.01"),
            api_pb2.AlgorithmSetting(name="w_weight_decay", value="0.25"),
            api_pb2.AlgorithmSetting(name="alpha_weight_decay", value="0.25"),
            api_pb2.AlgorithmSetting(name="w_momentum", value="0.9"),
            api_pb2.AlgorithmSetting(name="w_grad_clip", value="5.0"),
            api_pb2.AlgorithmSetting(name="batch_size", value="100"),
            api_pb2.AlgorithmSetting(name="num_workers", value="0"),
            api_pb2.AlgorithmSetting(name="init_channels", value="1"),
            api_pb2.AlgorithmSetting(name="print_step", value="100"),
            api_pb2.AlgorithmSetting(name="num_nodes", value="4"),
            api_pb2.AlgorithmSetting(name="stem_multiplier", value="3"),
        ]
        is_valid, _ = validate_algorithm_settings(valid)
        self.assertEqual(is_valid, True)

        # Invalid num_epochs
        invalid = [api_pb2.AlgorithmSetting(name="num_epochs", value="0")]
        is_valid, _ = validate_algorithm_settings(invalid)
        self.assertEqual(is_valid, False)

        # Invalid w_lr
        invalid = [api_pb2.AlgorithmSetting(name="w_lr", value="-0.1")]
        is_valid, _ = validate_algorithm_settings(invalid)
        self.assertEqual(is_valid, False)

        # Invalid alpha_weight_decay
        invalid = [api_pb2.AlgorithmSetting(name="alpha_weight_decay", value="-0.02")]
        is_valid, _ = validate_algorithm_settings(invalid)
        self.assertEqual(is_valid, False)

        # Invalid w_momentum
        invalid = [api_pb2.AlgorithmSetting(name="w_momentum", value="-0.8")]
        is_valid, _ = validate_algorithm_settings(invalid)
        self.assertEqual(is_valid, False)

        # Invalid batch_size
        invalid = [api_pb2.AlgorithmSetting(name="batch_size", value="0")]
        is_valid, _ = validate_algorithm_settings(invalid)
        self.assertEqual(is_valid, False)

        # Valid batch_size
        valid = [api_pb2.AlgorithmSetting(name="batch_size", value="None")]
        is_valid, _ = validate_algorithm_settings(valid)
        self.assertEqual(is_valid, True)

        # Invalid print_step
        invalid = [api_pb2.AlgorithmSetting(name="print_step", value="0")]
        is_valid, _ = validate_algorithm_settings(invalid)
        self.assertEqual(is_valid, False)


if __name__ == '__main__':
    unittest.main()
