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


if __name__ == '__main__':
    unittest.main()
