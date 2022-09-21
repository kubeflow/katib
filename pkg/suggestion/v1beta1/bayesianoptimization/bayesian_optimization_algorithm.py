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

""" module for bayesian optimization algorithm """
import numpy as np
from sklearn.preprocessing import MinMaxScaler

from pkg.suggestion.v1beta1.bayesianoptimization.global_optimizer import GlobalOptimizer


class BOAlgorithm:
    """ class for bayesian optimization """

    def __init__(self, experiment_name, dim, N, lowerbound, upperbound, X_train, y_train, mode, trade_off,
                 length_scale, noise, nu, kernel_type, n_estimators, max_features, model_type, logger=None):
        # np.random.seed(0)
        self._experiment_name = experiment_name
        self.dim = dim
        self.N = N or 100
        self.l = np.zeros((1, dim))
        self.u = np.ones((1, dim))
        self.lowerbound = lowerbound.reshape(1, dim)
        self.upperbound = upperbound.reshape(1, dim)
        self.logger = logger

        # normalize the upperbound and lowerbound to [0, 1]
        self.scaler = MinMaxScaler()
        self.scaler.fit(np.append(self.lowerbound, self.upperbound, axis=0))

        self.X_train = X_train
        self.y_train = y_train
        if self.y_train is None:
            self.current_optimal = None
        else:
            self.current_optimal = max(self.y_train)
            self.logger.debug("create optimizer", extra={
                "Experiment": self._experiment_name})
            # initialize the global optimizer
            self.optimizer = GlobalOptimizer(
                N,
                self.l,
                self.u,
                self.scaler,
                self.X_train,
                self.y_train,
                self.current_optimal,
                experiment_name=self._experiment_name,
                mode=mode,
                trade_off=trade_off,
                length_scale=length_scale,
                noise=noise,
                nu=nu,
                kernel_type=kernel_type,
                n_estimators=n_estimators,
                max_features=max_features,
                model_type=model_type,
                logger=logger,
            )
            self.logger.debug("optimizer created", extra={
                "Experiment": self._experiment_name})

    def get_suggestion(self, request_num):
        """ main function to provide suggestion """
        x_next_list = []
        if self.X_train is None and self.y_train is None and self.current_optimal is None:
            # randomly pick a point as the first trial
            for _ in range(request_num):
                x_next_list.append(np.random.uniform(
                    self.lowerbound, self.upperbound, size=(1, self.dim)))
        else:
            _, x_next_list_que = self.optimizer.direct(request_num)
            for xn in x_next_list_que:
                x = np.array(xn).reshape(1, self.dim)
                x = self.scaler.inverse_transform(x)
                x_next_list.append(x)
        return x_next_list
