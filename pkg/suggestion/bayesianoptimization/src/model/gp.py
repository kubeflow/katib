""" module for gaussian process prior """
from sklearn.gaussian_process.kernels import RBF, Matern
from sklearn.gaussian_process import GaussianProcessRegressor


class GaussianProcessModel:
    """ use the gaussian process as a prior """
    def __init__(self, length_scale=0.5, noise=0.00005,
                 nu=1.5, kernel_type="matern"):
        """
        :param length_scale: the larger the length_scale is, the smoother the gaussian prior is. If a float,
        an isotropic kernel is used. If an array, an anisotropic kernel is used where each dimension of it defines
        the length-scale of the respective feature dimension.
        :param noise:
        :param nu: control the smoothness of the prior using Matern kernel. The larger nu is, the smoother the
        approximate function is.
        :param kernel_type: "rbf": squared exponential kernel, "matern": Matern kernel.
        """
        if kernel_type == "rbf":
            kernel = RBF(length_scale=length_scale)
        elif kernel_type == "matern":
            kernel = Matern(length_scale=length_scale, nu=nu)
        else:
            raise Exception("kernel_type must be 'rbf' or 'matern'")
        self.gp = GaussianProcessRegressor(
            kernel=kernel,
            alpha=noise,
            random_state=0,
            optimizer=None,
        )

    def fit(self, X_train, y_train):
        self.gp.fit(X_train, y_train)

    def predict(self, X_test):
        y_mean, y_std = self.gp.predict(X_test, return_std=True)
        y_variance = y_std ** 2
        return y_mean, y_std, y_variance
