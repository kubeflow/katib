""" module for algorithm manager """
from collections.abc import Iterable

import numpy as np

from pkg.api.python import api_pb2
from .utils import get_logger


def deal_with_discrete(feasible_values, current_value):
    """ function to embed the current values to the feasible discrete space"""
    diff = np.subtract(feasible_values, current_value)
    diff = np.absolute(diff)
    return feasible_values[np.argmin(diff)]


def deal_with_categorical(feasible_values, one_hot_values):
    """ function to do the one hot encoding of the categorical values """
    index = np.argmax(one_hot_values)
    #index = one_hot_values.argmax()
    return feasible_values[int(index)]


class AlgorithmManager:
    """ class for the algorithm manager
    provide some helper functions
    """
    def __init__(self, study_id, study_config, X_train, y_train, logger=None):
        self.logger = logger if (logger is not None) else get_logger()
        self._study_id = study_id
        self._study_config = study_config
        self._goal = self._study_config.optimization_type
        name_ids, dim, lower_bounds, upper_bounds, parameter_types, names, discrete_info, categorical_info = \
            AlgorithmManager._parse_parameter_configs(self._study_config.parameter_configs.configs)
        self._dim = dim
        self._lowerbound = lower_bounds
        self._upperbound = upper_bounds
        self._types = parameter_types
        self._names = names
        self._discrete_info = discrete_info
        self._categorical_info = categorical_info
        self._name_id = name_ids

        self._X_train = AlgorithmManager._mapping_params(X_train, self._dim, self._name_id, self._types, self._categorical_info)
        self._y_train = AlgorithmManager._parse_metric(y_train, self.goal)

    @property
    def study_id(self):
        """ return the study id """
        return self._study_id

    @property
    def study_config(self):
        """ return the study configuration """
        return self._study_config

    @property
    def goal(self):
        """ return the optimization goal"""
        return self._goal

    @property
    def dim(self):
        """ return the dimension """
        return self._dim

    @property
    def lower_bound(self):
        """ return the lower bound of all the parameters """
        return self._lowerbound

    @property
    def upper_bound(self):
        """ return the upper bound of all the parameters """
        return self._upperbound

    @property
    def types(self):
        """ return the types of all the parameters """
        return self._types

    @property
    def names(self):
        """ return the names of all the parameters """
        return self._names

    @property
    def discrete_info(self):
        """ return the info of all the discrete parameters """
        return self._discrete_info

    @property
    def categorical_info(self):
        """ return the info of all the categorical parameters """
        return self._categorical_info

    @property
    def X_train(self):
        """ return the training data """
        return self._X_train

    @property
    def y_train(self):
        """ return the target of the training data"""
        return self._y_train

    @staticmethod
    def _parse_parameter_configs(parameter_configs):
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
            if param.parameter_type in [api_pb2.DOUBLE, api_pb2.INT]:
                new_lower = param.feasible.min
                new_upper = param.feasible.max
            elif param.parameter_type == api_pb2.DISCRETE:
                discrete_values = [int(x) for x in param.feasible.list]
                new_lower = min(discrete_values)
                new_upper = max(discrete_values)
                discrete_info.append(
                    {"name": param.name, "values": discrete_values})
            elif param.parameter_type == api_pb2.CATEGORICAL:
                num_feasible = len(param.feasible.list)
                new_lower = [0 for _ in range(num_feasible)]
                new_upper = [1 for _ in range(num_feasible)]
                categorical_info.append({
                    "name": param.name,
                    "values": param.feasible.list,
                    "number": num_feasible,
                })
            if isinstance(new_lower, Iterable):
                lower_bounds.extend(new_lower)
                upper_bounds.extend(new_upper)
                dim += len(new_lower)
            else:
                lower_bounds.append(new_lower)
                upper_bounds.append(new_upper)
                dim += 1
        return name_ids, dim, lower_bounds, upper_bounds, parameter_types, names, discrete_info, categorical_info

    @staticmethod
    def _mapping_params(parameters_list, dim, name_id, types, categorical_info):
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
                            parsed_X[row_idx, offset+value_num] = 1
                            offset += ci["number"]
        return parsed_X

    @staticmethod
    def _parse_metric(y_train, goal):
        """ parse the metric to the dictionary """
        y_array = np.array(y_train, dtype=np.float64)
        if goal == api_pb2.MINIMIZE:
                y_array *= -1
        return y_array

    def parse_x_next(self, x_next):
        """ parse the next suggestion to the proper format """
        counter = 0
        result = []
        for i in range(len(self._types)):
            if self._types[i] == api_pb2.INT:
                result.append(int(round(x_next[counter], 0)))
                counter = counter + 1
            elif self._types[i] == api_pb2.DISCRETE:
                for param in self._discrete_info:
                    if param["name"] == self._names[i]:
                        result.append(
                            deal_with_discrete(param["values"], x_next[counter])
                        )
                        counter = counter + 1
                        break
            elif self._types[i] == api_pb2.CATEGORICAL:
                for param in self._categorical_info:
                    if param["name"] == self._names[i]:
                        result.append(deal_with_categorical(
                            feasible_values=param["values"],
                            one_hot_values=x_next[counter:counter + param["number"]],
                        ))
                        counter = counter + param["number"]
                        break
            elif self._types[i] == api_pb2.DOUBLE:
                result.append(x_next[counter])
                counter = counter + 1
        return result

    def convert_to_dict(self, x_next):
        """ convert the next suggestion to the dictionary """
        result = []
        for i in range(len(x_next)):
            tmp = dict({
                "name": self._names[i],
                "value": x_next[i],
                "type": self._types[i],
            })
            result.append(tmp)
        return result
