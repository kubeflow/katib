import pytest
import numpy as np
from sklearn.preprocessing import MinMaxScaler

from ..bayesianoptimization.src.global_optimizer import GlobalOptimizer


@pytest.fixture
def scaler(lower_bounds, upper_bounds, dim):
    scaler = MinMaxScaler()
    lower_bounds = np.array(lower_bounds, dtype=np.float64).reshape(1, dim)
    upper_bounds = np.array(upper_bounds, dtype=np.float64).reshape(1, dim)
    scaler.fit(np.append(lower_bounds, upper_bounds, axis=0))
    return scaler


def test_global_optimizer(dim, request_num, scaler, X_train, y_train):
    optimizer = GlobalOptimizer(N=200,
                                l=np.zeros((1, dim)),
                                u=np.ones((1, dim)),
                                scaler=scaler,
                                X_train=X_train,
                                y_train=y_train,
                                current_optimal=1.0,
                                mode="ei",
                                trade_off=0.01,
                                length_scale=0.5,
                                noise=0.00005,
                                nu=1.5,
                                kernel_type="matern",
                                n_estimators=None,
                                max_features=None,
                                model_type="gp")
    f_min, x_next_candidate = optimizer.direct(request_num)
    assert isinstance(f_min, float)
    assert np.array(x_next_candidate).shape == (request_num, 1, dim)
