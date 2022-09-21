# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

""" module for acquisition function"""
import numpy as np
from scipy.stats import norm


class AcquisitionFunc:
    """
    Class for acquisition function with options for expected improvement,
    probability of improvement, or lower confident bound.
    """

    def __init__(self, model, current_optimal, mode="ei", trade_off=0.01):
        """
        :param mode: pi: probability of improvement, ei: expected improvement, lcb: lower confident bound
        :param trade_off: a parameter to control the trade off between exploiting and exploring
        :param model_type: gp: gaussian process, rf: random forest
        """
        self.model = model
        self.current_optimal = current_optimal
        self.mode = mode
        self.trade_off = trade_off

    def compute(self, X_test):
        y_mean, y_std, y_variance = self.model.predict(X_test)

        z = (y_mean - self.current_optimal - self.trade_off) / y_std

        if self.mode == "ei":
            if y_std.any() < 0.000001:
                return 0, y_mean, y_variance
            result = y_std * (z * norm.cdf(z) + norm.pdf(z))
        elif self.mode == "pi":
            result = norm.cdf(z)
        else:
            result = - (y_mean - self.trade_off * y_std)
        return np.squeeze(result), np.squeeze(y_mean), np.squeeze(y_variance)
