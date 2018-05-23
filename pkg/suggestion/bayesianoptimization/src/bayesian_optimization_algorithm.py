""" module for bayesian optimization algorithm """
import numpy as np
from sklearn.preprocessing import MinMaxScaler

from pkg.suggestion.bayesianoptimization.src.global_optimizer.global_optimizer import GlobalOptimizer


class BOAlgorithm:
    """ class for bayesian optimization """
    def __init__(self, dim, N, lowerbound, upperbound, X_train, y_train, mode, trade_off,
                 length_scale, noise, nu, kernel_type, n_estimators, max_features, model_type):
        # np.random.seed(0)
        self.dim = dim
        self.N = N or 100
        self.l = np.zeros((1, dim))
        self.u = np.ones((1, dim))
        self.lowerbound = lowerbound.reshape(1, dim)
        self.upperbound = upperbound.reshape(1, dim)

        # normalize the upperbound and lowerbound to [0, 1]
        self.scaler = MinMaxScaler()
        self.scaler.fit(np.append(self.lowerbound, self.upperbound, axis=0))

        self.x_next = None
        self.X_train = X_train
        self.y_train = y_train
        if self.y_train is None:
            self.current_optimal = None
        else:
            self.current_optimal = max(self.y_train)

        # initialize the global optimizer
        self.optimizer = GlobalOptimizer(
            N,
            self.l,
            self.u,
            self.scaler,
            self.X_train,
            self.y_train,
            self.current_optimal,
            mode=mode,
            trade_off=trade_off,
            length_scale=length_scale,
            noise=noise,
            nu=nu,
            kernel_type=kernel_type,
            n_estimators=n_estimators,
            max_features=max_features,
            model_type=model_type,
        )

    def get_suggestion(self):
        """ main function to provide suggestion """
        if self.X_train is None and self.y_train is None and self.current_optimal is None:
            # randomly pick a point as the first trial
            self.x_next = np.random.uniform(self.lowerbound, self.upperbound, size=(1, self.dim))
            return self.x_next
        _, self.x_next = self.optimizer.direct()
        self.x_next = np.array(self.x_next).reshape(1, self.dim)
        self.x_next = self.scaler.inverse_transform(self.x_next)
        return self.x_next
