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

from pkg.apis.manager.v1beta1.python import api_pb2


def validate_operations(operations: list[api_pb2.Operation]) -> (bool, str):

    # Validate each operation
    for operation in operations:

        # Check OperationType
        if not operation.operation_type:
            return False, "Missing operationType in Operation:\n{}".format(operation)

        # Check ParameterConfigs
        if not operation.parameter_specs.parameters:
            return False, "Missing ParameterConfigs in Operation:\n{}".format(operation)

        # Validate each ParameterConfig in Operation
        parameters_list = list(operation.parameter_specs.parameters)
        for parameter in parameters_list:

            # Check Name
            if not parameter.name:
                return False, "Missing Name in ParameterConfig:\n{}".format(parameter)

            # Check ParameterType
            if not parameter.parameter_type:
                return False, "Missing ParameterType in ParameterConfig:\n{}".format(parameter)

            # Check List in Categorical or Discrete Type
            if parameter.parameter_type == api_pb2.CATEGORICAL or parameter.parameter_type == api_pb2.DISCRETE:
                if not parameter.feasible_space.list:
                    return False, "Missing List in ParameterConfig.feasibleSpace:\n{}".format(parameter)

            # Check Max, Min, Step in Int or Double Type
            elif parameter.parameter_type == api_pb2.INT or parameter.parameter_type == api_pb2.DOUBLE:
                if not parameter.feasible_space.min and not parameter.feasible_space.max:
                    return False, "Missing Max and Min in ParameterConfig.feasibleSpace:\n{}".format(parameter)

                try:
                    if (parameter.parameter_type == api_pb2.DOUBLE and
                            (not parameter.feasible_space.step or float(parameter.feasible_space.step) <= 0)):
                        return False, \
                               "Step parameter should be > 0 in ParameterConfig.feasibleSpace:\n{}".format(parameter)
                except Exception as e:
                    return False, \
                           "failed to validate ParameterConfig.feasibleSpace \n{parameter}):\n{exception}".format(
                               parameter=parameter, exception=e)

    return True, ""
