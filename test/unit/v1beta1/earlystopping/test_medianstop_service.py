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
from unittest.mock import patch

import grpc
import grpc_testing
import utils

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.earlystopping.v1beta1.medianstop.service import MedianStopService


class TestMedianStop(unittest.TestCase):
    def setUp(self):
        # Mock load Kubernetes config.
        patcher = patch('pkg.earlystopping.v1beta1.medianstop.service.config.load_kube_config')
        self.mock_sum = patcher.start()
        self.addCleanup(patcher.stop)

        servicers = {
            api_pb2.DESCRIPTOR.services_by_name['EarlyStopping']: MedianStopService(
            )
        }

        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time())

    def test_validate_early_stopping_settings(self):
        # Valid cases
        early_stopping = api_pb2.EarlyStoppingSpec(
            algorithm_name="medianstop",
            algorithm_settings=[
                api_pb2.EarlyStoppingSetting(
                    name="min_trials_required",
                    value="2",
                ),
                api_pb2.EarlyStoppingSetting(
                    name="start_step",
                    value="5",
                ),
            ],
        )

        _, _, code, _ = utils.call_validate(self.test_server, early_stopping)
        self.assertEqual(code, grpc.StatusCode.OK)

        # Invalid cases
        # Unknown algorithm name
        early_stopping = api_pb2.EarlyStoppingSpec(algorithm_name="unknown")

        _, _, code, details = utils.call_validate(self.test_server, early_stopping)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "unknown algorithm name unknown")

        # Unknown config name
        early_stopping = api_pb2.EarlyStoppingSpec(
            algorithm_name="medianstop",
            algorithm_settings=[
                api_pb2.EarlyStoppingSetting(
                    name="unknown_conf",
                    value="100",
                ),
            ],
        )

        _, _, code, details = utils.call_validate(self.test_server, early_stopping)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "unknown setting unknown_conf for algorithm medianstop")

        # Wrong min_trials_required
        early_stopping = api_pb2.EarlyStoppingSpec(
            algorithm_name="medianstop",
            algorithm_settings=[
                api_pb2.EarlyStoppingSetting(
                    name="min_trials_required",
                    value="0",
                ),
            ],
        )

        _, _, code, details = utils.call_validate(self.test_server, early_stopping)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "min_trials_required must be greater than zero (>0)")

        # Wrong start_step
        early_stopping = api_pb2.EarlyStoppingSpec(
            algorithm_name="medianstop",
            algorithm_settings=[
                api_pb2.EarlyStoppingSetting(
                    name="start_step",
                    value="0",
                ),
            ],
        )
        _, _, code, details = utils.call_validate(self.test_server, early_stopping)
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)
        self.assertEqual(details, "start_step must be greater or equal than one (>=1)")

    def test_get_earlystopping_rules(self):
        # TODO (andreyvelich): Add more informative tests.
        trials = [
            api_pb2.Trial(
                name="test-asfjh",
            ),
            api_pb2.Trial(
                name="test-234hs",
            )
        ]

        experiment = api_pb2.Experiment(
            name="test",
        )

        request = api_pb2.GetEarlyStoppingRulesRequest(
            experiment=experiment,
            trials=trials,
            db_manager_address="katib-db-manager.kubeflow:6789"
        )

        get_earlystopping_rules = self.test_server.invoke_unary_unary(
            method_descriptor=(api_pb2.DESCRIPTOR
                .services_by_name['EarlyStopping']
                .methods_by_name['GetEarlyStoppingRules']),
            invocation_metadata={},
            request=request, timeout=1)

        _, _, code, _ = get_earlystopping_rules.termination()

        self.assertEqual(code, grpc.StatusCode.OK)


    def test_get_earlystopping_rules_uses_median(self):
        # Average of the first start_step metric values for each Trial.
        metric_logs = {
            "trial-1": ["0.1", "0.1", "0.5"],
            "trial-2": ["0.2", "0.2", "0.5"],
            "trial-3": ["0.9", "0.9", "0.5"],
        }

        def get_observation_log(request, timeout):
            return api_pb2.GetObservationLogReply(
                observation_log=api_pb2.ObservationLog(
                    metric_logs=[
                        api_pb2.MetricLog(
                            metric=api_pb2.Metric(name=request.metric_name, value=value),
                        )
                        for value in metric_logs[request.trial_name]
                    ]
                )
            )

        experiment = api_pb2.Experiment(
            name="test",
            spec=api_pb2.ExperimentSpec(
                objective=api_pb2.ObjectiveSpec(
                    type=api_pb2.MAXIMIZE,
                    objective_metric_name="accuracy",
                ),
                early_stopping=api_pb2.EarlyStoppingSpec(
                    algorithm_name="medianstop",
                    algorithm_settings=[
                        api_pb2.EarlyStoppingSetting(name="min_trials_required", value="3"),
                        api_pb2.EarlyStoppingSetting(name="start_step", value="2"),
                    ],
                ),
            ),
        )
        trials = [
            api_pb2.Trial(
                name=name,
                status=api_pb2.TrialStatus(condition=api_pb2.TrialStatus.SUCCEEDED),
            )
            for name in metric_logs
        ]
        request = api_pb2.GetEarlyStoppingRulesRequest(
            experiment=experiment,
            trials=trials,
            db_manager_address="katib-db-manager.kubeflow:6789"
        )

        with patch("pkg.earlystopping.v1beta1.medianstop.service.api_pb2_grpc.DBManagerStub") as stub:
            stub.return_value.GetObservationLog.side_effect = get_observation_log
            get_earlystopping_rules = self.test_server.invoke_unary_unary(
                method_descriptor=(api_pb2.DESCRIPTOR
                    .services_by_name['EarlyStopping']
                    .methods_by_name['GetEarlyStoppingRules']),
                invocation_metadata={},
                request=request, timeout=1)

            response, _, code, _ = get_earlystopping_rules.termination()

        self.assertEqual(code, grpc.StatusCode.OK)
        self.assertEqual(len(response.early_stopping_rules), 1)
        rule = response.early_stopping_rules[0]
        self.assertEqual(rule.name, "accuracy")
        self.assertEqual(rule.comparison, api_pb2.LESS)
        self.assertEqual(rule.start_step, 2)
        # Median of the Trial averages (0.1, 0.2, 0.9), not their mean (0.4).
        self.assertAlmostEqual(float(rule.value), 0.2)


if __name__ == '__main__':
    unittest.main()
