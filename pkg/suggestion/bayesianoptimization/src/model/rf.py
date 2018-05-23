from sklearn.ensemble import RandomForestRegressor


class RandomForestModel:
    def __init__(self, n_estimators, max_features):
        n_estimators = n_estimators or 50
        max_features = max_features or "auto"
        self.rf = RandomForestRegressor(
            n_estimators=n_estimators,
            max_features=max_features,
        )
