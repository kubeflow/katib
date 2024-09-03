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

"""
Module containing helper functions to translate objects that come
to/from the grpc API into the format accepted/returned by the different
suggestion generation algorithms.
"""
from collections.abc import Iterable

import numpy as np

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.suggestion.v1beta1.hyperband.parameter import ParameterConfig


def _deal_with_discrete(feasible_values, current_value):
    """function to embed the current values to the feasible discrete space"""
    diff = np.subtract(feasible_values, current_value)
    diff = np.absolute(diff)
    return feasible_values[np.argmin(diff)]


def _deal_with_categorical(feasible_values, one_hot_values):
    """function to do the one hot encoding of the categorical values"""
    index = np.argmax(one_hot_values)
    return feasible_values[int(index)]


def parse_parameter_configs(parameter_configs):
    name_ids = {}
    dim = 0
    lower_bounds = []
    upper_bounds = []
    parameter_types = []
    names = []
    discrete_info = []
    categorical_info = []
    for param_idx, param in enumerate(parameter_configs):
        name_ids[param.name] = param_idx
        parameter_types.append(param.parameter_type)
        names.append(param.name)
        if param.parameter_type == api_pb2.DOUBLE:
            new_lower = float(param.feasible_space.min)
            new_upper = float(param.feasible_space.max)
        elif param.parameter_type == api_pb2.INT:
            new_lower = int(param.feasible_space.min)
            new_upper = int(param.feasible_space.max)
        elif param.parameter_type == api_pb2.DISCRETE:
            discrete_values = [int(x) for x in param.feasible_space.list]
            new_lower = min(discrete_values)
            new_upper = max(discrete_values)
            discrete_info.append({"name": param.name, "values": discrete_values})
        elif param.parameter_type == api_pb2.CATEGORICAL:
            num_feasible = len(param.feasible_space.list)
            new_lower = [0 for _ in range(num_feasible)]
            new_upper = [1 for _ in range(num_feasible)]
            categorical_info.append(
                {
                    "name": param.name,
                    "values": param.feasible_space.list,
                    "number": num_feasible,
                }
            )
        if isinstance(new_lower, Iterable):  # handles categorical parameters
            lower_bounds.extend(new_lower)
            upper_bounds.extend(new_upper)
            dim += len(new_lower)
        else:  # handles ints, doubles, and discrete parameters
            lower_bounds.append(new_lower)
            upper_bounds.append(new_upper)
            dim += 1
    parsed_config = ParameterConfig(
        name_ids,
        dim,
        lower_bounds,
        upper_bounds,
        parameter_types,
        names,
        discrete_info,
        categorical_info,
    )
    return parsed_config


def parse_previous_observations(parameters_list, dim, name_id, types, categorical_info):
    parsed_X = np.zeros(shape=(len(parameters_list), dim))
    for row_idx, parameters in enumerate(parameters_list):
        offset = 0
        for p in parameters:
            map_id = name_id[p.name]
            if types[map_id] in [api_pb2.DOUBLE, api_pb2.INT, api_pb2.DISCRETE]:
                parsed_X[row_idx, offset] = float(p.value)
                offset += 1
            elif types[map_id] == api_pb2.CATEGORICAL:
                for ci in categorical_info:
                    if ci["name"] == p.name:
                        value_num = ci["values"].index(p.value)
                        parsed_X[row_idx, offset + value_num] = 1
                        offset += ci["number"]
    return parsed_X


def parse_metric(y_train, goal):
    """
    Parse the metric to the dictionary
    """
    y_array = np.array(y_train, dtype=np.float64)
    if goal == api_pb2.MINIMIZE:
        y_array *= -1
    return y_array


def parse_x_next_vector(
    x_next, param_types, param_names, discrete_info, categorical_info
):
    """parse the next suggestion to the proper format"""
    counter = 0
    result = []
    if isinstance(x_next, np.ndarray):
        x_next = x_next.squeeze(axis=0)
    for par_type, par_name in zip(param_types, param_names):
        if par_type == api_pb2.INT:
            value = int(round(x_next[counter], 0))
            counter = counter + 1
        elif par_type == api_pb2.DOUBLE:
            value = float(x_next[counter])
            counter = counter + 1
        elif par_type == api_pb2.DISCRETE:
            for param in discrete_info:
                if param["name"] == par_name:
                    value = _deal_with_discrete(param["values"], x_next[counter])
                    counter = counter + 1
                    break
        elif par_type == api_pb2.CATEGORICAL:
            for param in categorical_info:
                if param["name"] == par_name:
                    value = _deal_with_categorical(
                        feasible_values=param["values"],
                        one_hot_values=x_next[counter : counter + param["number"]],
                    )
                    counter = counter + param["number"]
                    break
        result.append({"name": par_name, "value": value, "type": par_type})
    return result


def parse_x_next_tuple(x_next, param_types, param_names):
    result = []
    for value, param_type, param_name in zip(x_next, param_types, param_names):
        result.append({"name": param_name, "type": param_type, "value": str(value)})
    return result
