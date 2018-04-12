""" module for gaussian process prior """
from sklearn.gaussian_process.kernels import RBF
from sklearn.gaussian_process import GaussianProcessRegressor


class Model:
    """ use the gaussian process as a prior """
    def __init__(self, length_scale=0.5, noise=0.00005):
        se_kernel = RBF(length_scale)
        self.gp = GaussianProcessRegressor(
            kernel=se_kernel,
            alpha=noise,
            random_state=0,
            optimizer=None,
        )
