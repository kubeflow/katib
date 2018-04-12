""" module for acquisition function"""
import numpy as np
from scipy.stats import norm

from suggestion.BO.model.Model import Model


class AcquisitionFunc:
    """ class for acquisition function
    expected improvement in this case
    """
    def __init__(self, X_train, y_train, current_optimal):
        self.X_train = X_train
        self.y_train = y_train
        self.current_optimal = current_optimal

    def get_expected_improvement(self, X_test):
        """ method to calculate expected improvement """
        model = Model()
        model.gp.fit(self.X_train, self.y_train)
        y_mean, y_std = model.gp.predict(X_test, return_std=True)
        y_variance = y_std**2

        if y_variance < 0.000001:
            return 0, y_mean, y_variance
        z = (y_mean - self.current_optimal) / y_variance
        result = (y_mean - self.current_optimal) * norm.cdf(z) + y_variance * norm.pdf(z)
        return np.squeeze(result), np.squeeze(y_mean), np.squeeze(y_variance)
