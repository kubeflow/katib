import grpc
import numpy as np

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc


def func(x1, x2):
    return 0.75 * np.exp(-(9 * x1 - 2) ** 2 / 4 - (9 * x2 - 2) ** 2 / 4) + 0.75 * np.exp(
        -(9 * x1 + 1) ** 2 / 49 - (9 * x2 + 1) / 10) + \
           0.5 * np.exp(-(9 * x1 - 7) ** 2 / 4 - (9 * x2 - 3) ** 2 / 4) - 0.2 * np.exp(
        -(9 * x1 - 4) ** 2 - (9 * x2 - 7) ** 2)


def run():
    channel = grpc.insecure_channel("localhost:50051")
    stub = api_pb2_grpc.ManagerStub(channel)
    study_configs = api_pb2.StudyConfig(
        name="test_study",
        owner="me",
        optimization_type=api_pb2.MAXIMIZE,
        optimization_goal=0.2,
        parameter_configs=api_pb2.StudyConfig.ParameterConfigs(
            configs=[
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
        default_suggestion_algorithm="BO",
        default_early_stopping_algorithm="",
        tags=[],
        objective_value_name="precision",
        metrics=[],
    )
    create_study_response = stub.CreateStudy(api_pb2.CreateStudyRequest(
        study_config=study_configs,
    ))
    study_id = create_study_response.study_id
    get_study_response = stub.GetStudy(api_pb2.GetStudyRequest(
        study_id=study_id
    ))
    maximum = -1
    iter = 0
    for i in range(20):
        get_suggestion_response = stub.GetSuggestions(api_pb2.GetSuggestionsRequest(
            study_id=study_id,
            suggestion_algorithm="cma",
        ))
        for trial in get_suggestion_response.trials:
            x1 = trial.parameter_set[0].value
            x2 = trial.parameter_set[1].value

            objective_value = func(float(x1), float(x2))
            if objective_value > maximum:
                maximum = objective_value
                iter = i
            stub.UpdateTrial(api_pb2.UpdateTrialRequest(
                trial_id=trial.trial_id,
                objective_value=str(objective_value),
                status=api_pb2.COMPLETED,
            ))
    print("find max {} in {} iteration".format(maximum, iter))


if __name__ == "__main__":
    run()
