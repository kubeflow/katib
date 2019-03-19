"""
Module for grid search algorithm.
"""
import itertools

import numpy as np

import api_pb2
from .parsing_utils import parse_x_next_tuple


class GridSearch:

    def __init__(self, parameter_config, suggestion_config, logger=None):
        self.parameter_config = parameter_config
        self.suggestion_config = suggestion_config
        self.logger = logger

    def _create_all_combinations(self):
        param_ranges = []
        location = 0
        default_grid_size = self.suggestion_config.get("DefaultGrid", 10)
        for idx, param_type in enumerate(self.parameter_config.parameter_types):
            param_name = self.parameter_config.names[idx]
            if param_type in [api_pb2.DOUBLE, api_pb2.INT]:
                num = self.suggestion_config.get(param_name, default_grid_size)
                param_values = \
                    np.linspace(self.parameter_config.lower_bounds[0, location],
                                self.parameter_config.upper_bounds[0, location],
                                num=num)
                location += 1
                if param_type == api_pb2.INT:
                    param_values = param_values.astype(np.int64)
            elif param_type == api_pb2.DISCRETE:
                for discrete_param in self.parameter_config.discrete_info:
                    if param_name == discrete_param["name"]:
                        param_values = discrete_param["values"]
                        break
                location += 1
            elif param_type == api_pb2.CATEGORICAL:
                for categ_param in self.parameter_config.categorical_info:
                    if param_name == categ_param["name"]:
                        param_values = categ_param["values"]
                        break
                location += categ_param["number"]
            param_ranges.append(param_values)
        all_combinations = [comb for comb in itertools.product(*param_ranges)]
        return all_combinations

    # TODO: get the number of previous observations without fetching X_train, y_train
    def get_suggestion(self, X_train, y_train, request_num):
        combinations = self._create_all_combinations()
        if request_num > 0:
            assert X_train.shape[0] == y_train.shape[0]
            start_idx = y_train.shape[0]
            end_idx = start_idx + request_num
        else:
            start_idx = 0
            end_idx = len(combinations)
        x_next_list = combinations[start_idx:end_idx]
        new_suggestions = [parse_x_next_tuple(x_next,
                                              self.parameter_config.parameter_types,
                                              self.parameter_config.names)
                           for x_next in x_next_list]
        return new_suggestions
