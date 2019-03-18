import pytest

from katib_suggestion.model.rf import RandomForestModel
from katib_suggestion.model.gp import GaussianProcessModel


MODELS = [RandomForestModel(n_estimators=5),
          GaussianProcessModel(),
          GaussianProcessModel(kernel_type="rbf")]


@pytest.mark.parametrize("model", MODELS)
def test_fit_predict(model, X_train, y_train):
    model.fit(X_train, y_train)
    y_mean, y_std, y_variance = model.predict(X_train)
    assert y_mean.shape == (2,)
    assert y_std.shape == (2,)
    assert y_variance.shape == (2,)


def test_gp_kernel_type_exception():
    with pytest.raises(Exception):
        _ = GaussianProcessModel(kernel_type="different_kernel")
