# Copyright 2026 The Kubeflow Authors.
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

import grpc
import grpc_testing
import utils

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.pbt.service import PbtService


class TestPbt(unittest.TestCase):
    def setUp(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name["Suggestion"]: PbtService()
        }
        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time()
        )

    def _spec(self, resample_probability=None):
        settings = [
            api_pb2.AlgorithmSetting(name="suggestion_trial_dir", value="/tmp"),
            api_pb2.AlgorithmSetting(name="n_population", value="10"),
            api_pb2.AlgorithmSetting(name="truncation_threshold", value="0.2"),
        ]
        if resample_probability is not None:
            settings.append(
                api_pb2.AlgorithmSetting(
                    name="resample_probability", value=resample_probability
                )
            )
        return api_pb2.ExperimentSpec(
            algorithm=api_pb2.AlgorithmSpec(
                algorithm_name="pbt", algorithm_settings=settings
            )
        )

    def test_validate_algorithm_settings(self):
        # Valid, without the optional resample_probability.
        _, _, code, _ = utils.call_validate(self.test_server, self._spec())
        self.assertEqual(code, grpc.StatusCode.OK)

        # Valid, with a resample_probability inside [0, 1]. Its value arrives as a
        # string, so the range check must cast it before comparing.
        _, _, code, _ = utils.call_validate(
            self.test_server, self._spec(resample_probability="0.5")
        )
        self.assertEqual(code, grpc.StatusCode.OK)

        # Invalid, resample_probability outside [0, 1].
        _, _, code, _ = utils.call_validate(
            self.test_server, self._spec(resample_probability="1.5")
        )
        self.assertEqual(code, grpc.StatusCode.INVALID_ARGUMENT)


if __name__ == "__main__":
    unittest.main()
