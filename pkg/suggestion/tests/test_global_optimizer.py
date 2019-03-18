import numpy as np

from katib_suggestion.global_optimizer import GlobalOptimizer


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
