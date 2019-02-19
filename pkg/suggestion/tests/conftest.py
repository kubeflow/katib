# pylint: disable=missing-docstring
import pytest
import numpy as np


def booth_function(X):
    f = (X[:, 0] + 2 * X[:, 1] - 7) ** 2 + (2 * X[:, 0] + X[:, 1] - 5) ** 2
    return f


@pytest.fixture
def X_train():
    x_range, y_range = np.arange(-5, 6), np.arange(-5, 6)
    X_train = np.array([(x, y) for x in x_range for y in y_range])
    return X_train


@pytest.fixture
def X_test():
    x_range, y_range = np.arange(-1.5, 2.5), np.arange(-1.5, 2.5)
    X_test = np.array([(x, y) for x in x_range for y in y_range])
    return X_test


@pytest.fixture
def y_train(X_train):
    y_train = booth_function(X_train)
    return y_train
