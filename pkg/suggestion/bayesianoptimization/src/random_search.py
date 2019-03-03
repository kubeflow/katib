"""
Module for random search algorithm.
"""
import numpy as np
from sklearn.preprocessing import MinMaxScaler


class RandomSearch:

    def __init__(self, parameter_config, suggestion_config, X_train, y_train, logger=None):
        # np.random.seed(0)
        self.parameter_config = parameter_config
        self.suggestion_config = suggestion_config
        self.logger = logger
        self.l = np.zeros((1, parameter_config.dim))
        self.u = np.ones((1, parameter_config.dim))
        self.lowerbound = np.array(parameter_config.lower_bounds).reshape((1, parameter_config.dim))
        self.upperbound = np.array(parameter_config.upper_bounds).reshape((1, parameter_config.dim))

        # normalize the upperbound and lowerbound to [0, 1]
        self.scaler = MinMaxScaler()
        self.scaler.fit(np.append(self.lowerbound, self.upperbound, axis=0))

        self.X_train = X_train
        self.y_train = y_train

    def get_suggestion(self, request_num):
        """
        Main function to provide suggestion.
        """
        x_next_list = []
        for _ in range(request_num):
            x_next_list.append(np.random.uniform(self.lowerbound, self.upperbound, size=(1, self.parameter_config.dim)))
        return x_next_list
