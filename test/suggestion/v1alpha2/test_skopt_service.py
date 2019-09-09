import grpc

from pkg.api.v1alpha2.python import api_pb2_grpc
from pkg.api.v1alpha2.python import api_pb2

DEFAULT_PORT = "0.0.0.0:6789"


def run():
    channel = grpc.insecure_channel(DEFAULT_PORT)
    stub = api_pb2.SuggestionStub(channel)
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
                algorithm_name="skopt-bayesian-optimization",
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

    response = stub.GetSuggestions(api_pb2.GetSuggestionsRequest(
        experiment=experiment,
        trials=trials,
        request_number=2,
    ))
    print(response)


if __name__ == "__main__":
    run()
