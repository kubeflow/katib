""" module for acquisition function"""
import numpy as np
from scipy.stats import norm
import forestci as fci

from pkg.suggestion.bayesianoptimization.src.model.gp import GaussianProcessModel
from pkg.suggestion.bayesianoptimization.src.model.rf import RandomForestModel


class AcquisitionFunc:
    """ class for acquisition function
    expected improvement in this case
    """
    def __init__(self, X_train, y_train, current_optimal, mode, trade_off, length_scale,
                 noise, nu, kernel_type, n_estimators, max_features, model_type):
        """
        :param mode: pi: probability of improvement, ei: expected improvement, lcb: lower confident bound
        :param trade_off: a parameter to control the trade off between exploiting and exploring
        :param model_type: gp: gaussian process, rf: random forest
        """
        self.X_train = X_train
        self.y_train = y_train
        self.current_optimal = current_optimal
        self.mode = mode or "ei"
        self.trade_off = trade_off or 0.01
        self.model_type = model_type or "gp"
        if self.model_type == "gp":
            self.model = GaussianProcessModel(
                length_scale=length_scale,
                noise=noise,
                nu=nu,
                kernel_type=kernel_type,
            )
        else:
            self.model = RandomForestModel(
                n_estimators=n_estimators,
                max_features=max_features,
            )

    def compute(self, X_test):
        if self.model_type == "gp":
            self.model.gp.fit(self.X_train, self.y_train)
            y_mean, y_std = self.model.gp.predict(X_test, return_std=True)
            y_variance = y_std ** 2
        else:
            self.model.rf.fit(self.y_train, self.y_train)
            y_mean = self.model.rf.predict(X_test)
            y_variance = fci.random_forest_error(self.model.rf, self.X_train, X_test)
            y_std = np.sqrt(y_variance)

        z = (y_mean - self.current_optimal - self.trade_off) / y_std

        if self.mode == "ei":
            if y_std < 0.000001:
                return 0, y_mean, y_variance
            result = y_std * (z * norm.cdf(z) + norm.pdf(z))
        elif self.mode == "pi":
            result = norm.cdf(z)
        else:
            result = - (y_mean - self.trade_off * y_std)
        return np.squeeze(result), np.squeeze(y_mean), np.squeeze(y_variance)
