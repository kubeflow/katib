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

import numpy as np
import forestci as fci
from sklearn.ensemble import RandomForestRegressor


class RandomForestModel:

    def __init__(self, n_estimators=50, max_features="auto"):
        self.rf = RandomForestRegressor(
            n_estimators=n_estimators,
            max_features=max_features,
        )
        self.X_train = None

    def fit(self, X_train, y_train):
        print(X_train.shape, y_train.shape)
        self.X_train = X_train
        self.rf.fit(X_train, y_train)

    def predict(self, X_test):
        y_mean = self.rf.predict(X_test)
        y_variance = fci.random_forest_error(self.rf, self.X_train, X_test)
        y_std = np.sqrt(y_variance)
        return y_mean, y_std, y_variance
