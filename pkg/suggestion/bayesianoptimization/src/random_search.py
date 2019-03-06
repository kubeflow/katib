"""
Module for random search algorithm.
"""
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
        return x_next_list
