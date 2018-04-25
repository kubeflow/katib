import unittest

import grpc
import numpy as np

from pkg.suggestion.types import DEFAULT_PORT
from pkg.api.python import api_pb2_grpc, api_pb2


def func(x1, x2):
    return 0.75 * np.exp(-(9 * x1 - 2) ** 2 / 4 - (9 * x2 - 2) ** 2 / 4) + 0.75 * np.exp(
        -(9 * x1 + 1) ** 2 / 49 - (9 * x2 + 1) / 10) + \
           0.5 * np.exp(-(9 * x1 - 7) ** 2 / 4 - (9 * x2 - 3) ** 2 / 4) - 0.2 * np.exp(
        -(9 * x1 - 4) ** 2 - (9 * x2 - 7) ** 2)


class TestCMA(unittest.TestCase):
    def __init__(self, *args, **kwargs):
        super(TestCMA, self).__init__(*args, **kwargs)
        self.channel = grpc.insecure_channel(DEFAULT_PORT)
        self.stub = api_pb2_grpc.SuggestionStub(self.channel)

    def test_set_suggestion_parameters(self):
        self.stub.SetSuggestionParameters(api_pb2.SetSuggestionParametersRequest(
            study_id="1",
            suggestion_parameters=[
                api_pb2.SuggestionParameter(
                )
            ]
        ))
        self.stub.SetSuggestionParameters(api_pb2.SetSuggestionParametersRequest(
            study_id="2",
            suggestion_parameters=[
                api_pb2.SuggestionParameter(
                )
            ]
        ))

    def test_generate_trials_maximize(self):
        completed_trials = []
        maximum = -1
        for i in range(20):
            response = self.stub.GenerateTrials(api_pb2.GenerateTrialsRequest(
                study_id="1",
                configs=api_pb2.StudyConfig(
                    optimization_type=api_pb2.MAXIMIZE,
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
                ),
                completed_trials=completed_trials,
                running_trials=[], )
            )
            for trial in response.trials:
                x1 = trial.parameter_set[0].value
                x2 = trial.parameter_set[1].value
                objective_value = func(float(x1), float(x2))
                if objective_value > maximum:
                    maximum = objective_value
                completed_trials.append(api_pb2.Trial(
                    trial_id=trial.trial_id,
                    study_id="1",
                    status=api_pb2.COMPLETED,
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
        self.assertTrue(maximum >= 1.2)

    def test_generate_trials_minimize(self):
        completed_trials = []
        minimum = float("inf")
        for i in range(20):
            response = self.stub.GenerateTrials(api_pb2.GenerateTrialsRequest(
                study_id="2",
                configs=api_pb2.StudyConfig(
                    optimization_type=api_pb2.MAXIMIZE,
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
                ),
                completed_trials=completed_trials,
                running_trials=[], )
            )
            for trial in response.trials:
                x1 = trial.parameter_set[0].value
                x2 = trial.parameter_set[1].value
                objective_value = func(float(x1), float(x2))
                if objective_value < minimum:
                    minimum = objective_value
                completed_trials.append(api_pb2.Trial(
                    trial_id=trial.trial_id,
                    study_id="1",
                    status=api_pb2.COMPLETED,
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
        self.assertTrue(minimum <= 0.1)

    def test_stop_study(self):
        self.stub.StopSuggestion(api_pb2.StopStudyRequest(
            study_id="1"
        ))
        self.stub.StopSuggestion(api_pb2.StopStudyRequest(
            study_id="2"
        ))


if __name__ == "__main__":
    unittest.main()
