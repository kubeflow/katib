import numpy as np


def set_boundary_weights(sigma, C):
    dim = C.shape[0]
    temp = (np.sum(np.diag(C)) / dim)
    weight = 2 / (sigma ** 2 * temp)
    return np.full((dim, 1), weight)


def cal_boundary_param(C):
    dim = C.shape[0]
    temp = np.log(np.diag(C)) - 1 / dim * np.sum(np.log(np.diag(C)))
    return np.exp(0.9 * temp)


def cal_penalty(x, x_feas, boundary_w, param):
    dim = x.shape[0]
    return np.sum((x_feas - x) ** 2 / param * boundary_w) * (1 / dim)
