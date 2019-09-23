import grpc
import grpc_testing
import unittest

from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.apis.manager.v1alpha3.python import api_pb2

from pkg.suggestion.v1alpha3.nasrl_service import NasrlService


class TestNasRL(unittest.TestCase):
    def setUp(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name['Suggestion']: NasrlService(
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
                    algorithm_name="hyperopt-tpe",
                    algorithm_setting=[
                        api_pb2.AlgorithmSetting(
                            name="random_state",
                            value="10"
                        )
                    ],
                ),
                objective=api_pb2.ObjectiveSpec(
                    type=api_pb2.MAXIMIZE,
                    goal=0.9
                ),
                nas_config=api_pb2.NasConfig(
                    graph_config=api_pb2.GraphConfig(
                        num_layers=8,
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
                                                max=None, min=None, list=["3", "5", "7"])
                                        )
                                    ]
                                )
                            )
                        ]
                    )
                )
            )
        )

        request = api_pb2.GetSuggestionsRequest(
            experiment=experiment,
            trials=trials,
            request_number=2,
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


if __name__ == '__main__':
    unittest.main()
