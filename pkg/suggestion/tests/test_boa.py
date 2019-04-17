import numpy as np

from ..bayesianoptimization.src.bayesian_optimization_algorithm import BOAlgorithm


def test_boa(dim, request_num, lower_bounds, upper_bounds, X_train, y_train):
    boa = BOAlgorithm(dim=dim,
                      N=200,
                      lowerbound=np.array(lower_bounds, dtype=np.float64),
                      upperbound=np.array(upper_bounds, dtype=np.float64),
                      X_train=X_train,
                      y_train=y_train,
                      mode="ei",
                      trade_off=0.01,
                      length_scale=0.5,
                      noise=0.00005,
                      nu=1.5,
                      kernel_type="matern",
                      n_estimators=None,
                      max_features=None,
                      model_type="gp")
    response = boa.get_suggestion(request_num)
    assert len(response) == request_num
