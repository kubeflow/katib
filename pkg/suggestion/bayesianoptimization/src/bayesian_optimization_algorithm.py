""" module for bayesian optimization algorithm """
import numpy as np

from .global_optimizer import GlobalOptimizer


class BOAlgorithm:
    """ class for bayesian optimization """

    def __init__(self, parameter_config, suggestion_config, logger=None):
        # np.random.seed(0)
        self.parameter_config = parameter_config
        self.suggestion_config = suggestion_config
        self.l = np.zeros((1, parameter_config.dim))
        self.u = np.ones((1, parameter_config.dim))
        self.burn_in = int(suggestion_config.get("burn_in", 1))
        self.scaler = self.parameter_config.create_scaler()
        self.logger = logger

    def get_suggestion(self, X_train, y_train, request_num):
        """ main function to provide suggestion """
        if len(y_train) < self.burn_in:
            x_next_list = [self.parameter_config.random_sample() for _ in range(request_num)]
        else:
            current_optimal = max(y_train)

            # initialize the global optimizer
            optimizer = GlobalOptimizer(
                self.suggestion_config.get("N", 100),
                self.l,
                self.u,
                self.scaler,
                X_train,
                y_train,
                current_optimal,
                mode=self.suggestion_config.get("mode", "pi"),
                trade_off=self.suggestion_config.get("trade_off", 0.01),
                length_scale=self.suggestion_config.get("length_scale", 0.5),
                noise=self.suggestion_config.get("noise", 0.0005),
                nu=self.suggestion_config.get("nu", 1.5),
                kernel_type=self.suggestion_config.get("kernel_type", "matern"),
                n_estimators=self.suggestion_config.get("n_estimators", 50),
                max_features=self.suggestion_config.get("max_features", "auto"),
                model_type=self.suggestion_config.get("model_type", "gp"),
                logger=self.logger
            )
            x_next_list = []
            _, x_next_list_que = optimizer.direct(request_num)
            for xn in x_next_list_que:
                x = np.array(xn).reshape(1, self.parameter_config.dim)
                x = self.scaler.inverse_transform(x)
                x_next_list.append(x)
        return x_next_list
