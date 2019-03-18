import pytest

from katib_suggestion.acquisition_func import AcquisitionFunc


@pytest.mark.parametrize("aq_mode", ["ei", "pi", "lcb"])
def test_ei(aq_mode, model, X_train):
    aq = AcquisitionFunc(model,
                         current_optimal=1.0,
                         mode=aq_mode,
                         trade_off=0.01)
    results, y_mean, y_variance = aq.compute(X_train)
    assert results.shape == (2,)
    assert y_mean.shape == (2,)
    assert y_variance.shape == (2,)
