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
