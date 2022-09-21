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

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.nas.common.validation import validate_operations


class TestNasCommon(unittest.TestCase):

    def test_validate_operations(self):

        # Valid Case
        valid_operations = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="filter_size",
                            parameter_type=api_pb2.CATEGORICAL,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["3", "5"]),
                        ),
                        api_pb2.ParameterSpec(
                            name="pool_size",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="2", min="3", step="1", list=[]),
                        ),
                        api_pb2.ParameterSpec(
                            name="valid_type_double_example",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="1.0", min="3.0", step="0.1", list=[]),
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(valid_operations)
        self.assertEqual(is_valid, True)

        # Invalid OperationType
        invalid_operation_type = [
            api_pb2.Operation(
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="filter_size",
                            parameter_type=api_pb2.CATEGORICAL,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["3", "5"])
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(invalid_operation_type)
        self.assertEqual(is_valid, False)

        # Invalid ParameterConfigs
        invalid_parameter_configs = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(),
            ),
        ]
        is_valid, _ = validate_operations(invalid_parameter_configs)
        self.assertEqual(is_valid, False)

        # Invalid ParameterName
        invalid_parameter_name = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="",
                            parameter_type=api_pb2.CATEGORICAL,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["3", "5"]),
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(invalid_parameter_name)
        self.assertEqual(is_valid, False)

        # Invalid ParameterType
        invalid_parameter_type = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="filter_size",
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, list=["3", "5"]),
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(invalid_parameter_type)
        self.assertEqual(is_valid, False)

        # invalid List in Categorical
        invalid_categorical_list = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="filter_size",
                            parameter_type=api_pb2.CATEGORICAL,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="1", min="2", list=None),
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(invalid_categorical_list)
        self.assertEqual(is_valid, False)

        # invalid Min and Max
        invalid_min_max = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="pool_size",
                            parameter_type=api_pb2.INT,
                            feasible_space=api_pb2.FeasibleSpace(
                                max=None, min=None, step=None, list=["1", "2"]),
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(invalid_min_max)
        self.assertEqual(is_valid, False)

        # Invalid Double type parameter
        invalid_double_parameter = [
            api_pb2.Operation(
                operation_type="separable_convolution",
                parameter_specs=api_pb2.Operation.ParameterSpecs(
                    parameters=[
                        api_pb2.ParameterSpec(
                            name="invalid_type_double_example",
                            parameter_type=api_pb2.DOUBLE,
                            feasible_space=api_pb2.FeasibleSpace(
                                max="1.0", min="3.0", step=None, list=None),
                        ),
                    ],
                ),
            ),
        ]
        is_valid, _ = validate_operations(invalid_double_parameter)
        self.assertEqual(is_valid, False)


if __name__ == '__main__':
    unittest.main()
