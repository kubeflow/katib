"""
Module for random search algorithm.
"""
from .parsing_utils import parse_x_next_vector


class RandomSearch:

    def __init__(self, parameter_config, suggestion_config, logger=None):
        self.parameter_config = parameter_config
        self.suggestion_config = suggestion_config
        self.logger = logger

    def get_suggestion(self, request_num):
        """
        Main function to provide suggestion.
        """
        x_next_list = [self.parameter_config.random_sample() for _ in range(request_num)]
        new_suggestions = [parse_x_next_vector(x_next,
                                               self.parameter_config.parameter_types,
                                               self.parameter_config.names,
                                               self.parameter_config.discrete_info,
                                               self.parameter_config.categorical_info)
                           for x_next in x_next_list]
        return new_suggestions
