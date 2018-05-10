import grpc

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
from pkg.suggestion.test_func import func
from pkg.suggestion.types import DEFAULT_PORT


def run():
    channel = grpc.insecure_channel(DEFAULT_PORT)
    stub = api_pb2_grpc.SuggestionStub(channel)
    set_param_response = stub.SetSuggestionParameters(api_pb2.SetSuggestionParametersRequest(
        study_id="1",
        suggestion_parameters=[
            api_pb2.SuggestionParameter(
                name="N",
                value="100",
            ),
            api_pb2.SuggestionParameter(
                name="kernel_type",
                value="matern",
            ),
            api_pb2.SuggestionParameter(
                name="mode",
                value="ei",
            ),
            api_pb2.SuggestionParameter(
                name="trade_off",
                value="0.01",
            ),
            api_pb2.SuggestionParameter(
                name="model_type",
                value="gp",
            ),
            api_pb2.SuggestionParameter(
                name="n_estimators",
                value="50",
            ),
        ]
    ))
    completed_trials = []
    maximum = -1
    iter = 0
    for i in range(30):
        response = stub.GenerateTrials(api_pb2.GenerateTrialsRequest(
            study_id="1",
            configs=api_pb2.StudyConfig(
                name="test_study",
                owner="me",
                optimization_type=api_pb2.MAXIMIZE,
                optimization_goal=0.2,
                parameter_configs=api_pb2.StudyConfig.ParameterConfigs(
                    configs=[
                        # api_pb2.ParameterConfig(
                        #     name="param1",
                        #     parameter_type=api_pb2.INT,
                        #     feasible=api_pb2.FeasibleSpace(max="5", min="1", list=[]),
                        # ),
                        # api_pb2.ParameterConfig(
                        #     name="param2",
                        #     parameter_type=api_pb2.CATEGORICAL,
                        #     feasible=api_pb2.FeasibleSpace(max=None, min=None, list=["cat1", "cat2", "cat3"])
                        # ),
                        # api_pb2.ParameterConfig(
                        #     name="param3",
                        #     parameter_type=api_pb2.DISCRETE,
                        #     feasible=api_pb2.FeasibleSpace(max=None, min=None, list=["3", "2", "6"])
                        # ),
                        # api_pb2.ParameterConfig(
                        #     name="param4",
                        #     parameter_type=api_pb2.DOUBLE,
                        #     feasible=api_pb2.FeasibleSpace(max="5", min="1", list=[])
                        # )
                        api_pb2.ParameterConfig(
                            name="param1",
                            parameter_type=api_pb2.DOUBLE,
                            feasible=api_pb2.FeasibleSpace(max="1", min="0", list=[]),
                        ),
                        api_pb2.ParameterConfig(
                            name="param2",
                            parameter_type=api_pb2.DOUBLE,
                            feasible=api_pb2.FeasibleSpace(max="1", min="0", list=[])
                        ),
                    ],
                ),
                access_permissions=[],
                suggest_algorithm="BO",
                autostop_algorithm="",
                study_task_name="task",
                suggestion_parameters=[],
                tags=[],
                objective_value_name="precision",
                metrics=[],
                image="",
                command=["", ""],
                gpu=0,
                scheduler="",
                mount=api_pb2.MountConf(
                    pvc="",
                    path="",
                ),
                pull_secret=""
            ),
            completed_trials=completed_trials,
            running_trials=[],)
        )
        x1 = response.trials[0].parameter_set[0].value
        x2 = response.trials[0].parameter_set[1].value
        objective_value = func(float(x1), float(x2))
        if objective_value > maximum:
            maximum = objective_value
            iter = i
        print(objective_value)
        completed_trials.append(api_pb2.Trial(
            trial_id=response.trials[0].trial_id,
            study_id="1",
            status=api_pb2.COMPLETED,
            eval_logs=[],
            objective_value=str(objective_value),
            parameter_set=[
                api_pb2.Parameter(
                    name="param1",
                    parameter_type=api_pb2.DOUBLE,
                    value=x1,
                ),
                api_pb2.Parameter(
                    name="param2",
                    parameter_type=api_pb2.DOUBLE,
                    value=x2,
                ),
            ]
        ))
        print(str(response.trials[0].parameter_set))
    stop_study_response = stub.StopSuggestion(api_pb2.StopStudyRequest(
        study_id="1"
    ))

    print("found the maximum: {} at {} iteration".format(maximum, iter))

if __name__ == "__main__":
    run()
