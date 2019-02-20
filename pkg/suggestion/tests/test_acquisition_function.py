import pytest

from ..bayesianoptimization.src.acquisition_func import AcquisitionFunc
from ..bayesianoptimization.src.model.gp import GaussianProcessModel


@pytest.fixture
def model(X_train, y_train):
    model = GaussianProcessModel()
    model.fit(X_train, y_train)
    return model


@pytest.mark.parametrize("aq_mode", ["ei", "pi", "lcb"])
def test_ei(aq_mode, model, X_test):
    aq = AcquisitionFunc(model,
                         current_optimal=1.0,
                         mode=aq_mode,
                         trade_off=0.01)
    results, y_mean, y_variance = aq.compute(X_test)
    assert results.shape == (16,)
    assert y_mean.shape == (16,)
    assert y_variance.shape == (16,)
