""" module for bayesian optimization algorithm """
import numpy as np
from sklearn.preprocessing import MinMaxScaler

from .global_optimizer import GlobalOptimizer


class BOAlgorithm:
    """ class for bayesian optimization """

    def __init__(self, parameter_config, suggestion_config, X_train, y_train, logger=None):
        # np.random.seed(0)
        self.parameter_config = parameter_config
        self.suggestion_config = suggestion_config
        self.l = np.zeros((1, parameter_config.dim))
        self.u = np.ones((1, parameter_config.dim))
        self.lowerbound = np.array(parameter_config.lower_bounds).reshape((1, parameter_config.dim))
        self.upperbound = np.array(parameter_config.upper_bounds).reshape((1, parameter_config.dim))

        # normalize the upperbound and lowerbound to [0, 1]
        self.scaler = MinMaxScaler()
        self.scaler.fit(np.append(self.lowerbound, self.upperbound, axis=0))

        self.X_train = X_train
        self.y_train = y_train
        if len(self.y_train) > 0:
            self.current_optimal = max(self.y_train)

            # initialize the global optimizer
            self.optimizer = GlobalOptimizer(
                suggestion_config.get("N", 100),
                self.l,
                self.u,
                self.scaler,
                self.X_train,
                self.y_train,
                self.current_optimal,
                mode=suggestion_config.get("mode", "pi"),
                trade_off=suggestion_config.get("trade_off", 0.01),
                length_scale=suggestion_config.get("length_scale", 0.5),
                noise=suggestion_config.get("noise", 0.0005),
                nu=suggestion_config.get("nu", 1.5),
                kernel_type=suggestion_config.get("kernel_type", "matern"),
                n_estimators=suggestion_config.get("n_estimators", 50),
                max_features=suggestion_config.get("max_features", "auto"),
                model_type=suggestion_config.get("model_type", "gp"),
                logger=logger
            )

    def get_suggestion(self, request_num):
        """ main function to provide suggestion """
        x_next_list = []
        if len(self.y_train) == 0:
            # randomly pick a point as the first trial
            for _ in range(request_num):
                x_next_list.append(np.random.uniform(self.lowerbound, self.upperbound, size=(1, self.parameter_config.dim)))
        else:
            _, x_next_list_que = self.optimizer.direct(request_num)
            for xn in x_next_list_que:
                x = np.array(xn).reshape(1, self.parameter_config.dim)
                x = self.scaler.inverse_transform(x)
                x_next_list.append(x)
        return x_next_list
