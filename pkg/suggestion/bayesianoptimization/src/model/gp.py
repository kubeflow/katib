""" module for gaussian process prior """
from sklearn.gaussian_process.kernels import RBF, Matern
from sklearn.gaussian_process import GaussianProcessRegressor


class GaussianProcessModel:
    """ use the gaussian process as a prior """
    def __init__(self, length_scale, noise, nu, kernel_type):
        """
        :param length_scale: the larger the length_scale is, the smoother the gaussian prior is. If a float,
        an isotropic kernel is used. If an array, an anisotropic kernel is used where each dimension of it defines
        the length-scale of the respective feature dimension.
        :param noise:
        :param nu: control the smoothness of the prior using Matern kernel. The larger nu is, the smoother the
        approximate function is.
        :param kernel_type: "rbf": squared exponential kernel, "matern": Matern kernel.
        """

        length_scale = length_scale or 0.5
        noise = noise or 0.00005
        nu = nu or 1.5
        kernel_type = kernel_type or "matern"

        if kernel_type == "rbf":
            kernel = RBF(length_scale=length_scale)
        else:
            kernel = Matern(length_scale=length_scale, nu=nu)

        self.gp = GaussianProcessRegressor(
            kernel=kernel,
            alpha=noise,
            random_state=0,
            optimizer=None,
        )
