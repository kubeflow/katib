import pytest

from ..bayesianoptimization.src.model.rf import RandomForestModel
from ..bayesianoptimization.src.model.gp import GaussianProcessModel


MODELS = [RandomForestModel(n_estimators=5),
          GaussianProcessModel(),
          GaussianProcessModel(kernel_type="rbf")]


@pytest.mark.parametrize("model", MODELS)
def test_fit_predict(model, X_train, y_train, X_test):
    model.fit(X_train, y_train)
    y_mean, y_std, y_variance = model.predict(X_test)
    assert y_mean.shape == (16,)
    assert y_std.shape == (16,)
    assert y_variance.shape == (16,)


def test_gp_kernel_type_exception():
    with pytest.raises(Exception):
        _ = GaussianProcessModel(kernel_type="different_kernel")
