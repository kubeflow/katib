import pytest

from ..bayesianoptimization.src.acquisition_func import AcquisitionFunc


@pytest.mark.parametrize("aq_mode", ["ei", "pi", "lcb"])
def test_ei(aq_mode, X_train, y_train, X_test):
    aq = AcquisitionFunc(X_train,
                         y_train,
                         current_optimal=1.0,
                         mode=aq_mode,
                         trade_off=0.01,
                         length_scale=0.5,
                         noise=0.00005,
                         nu=1.5,
                         kernel_type="matern",
                         n_estimators=None,
                         max_features=None,
                         model_type="gp")
    results, y_mean, y_variance = aq.compute(X_test)
