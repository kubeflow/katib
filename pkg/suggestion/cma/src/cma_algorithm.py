import numpy as np
from scipy.special import gamma
from sklearn.preprocessing import MinMaxScaler

from pkg.suggestion.cma.src.handle_boundary import cal_boundary_param, cal_penalty


def sum_of_sign(array, sign):
    if sign == "positive":
        flag = (array >= 0).astype("int32")
    else:
        flag = (array < 0).astype("int32")
    return np.sum((array * flag))


class CMAES:
    def __init__(self, dim, lowerbound, upperbound):
        self.dim = dim
        self.scaler = MinMaxScaler()
        lowerbound = lowerbound.reshape(1, dim)
        upperbound = upperbound.reshape(1, dim)
        self.scaler.fit(np.append(lowerbound, upperbound, axis=0))

        # 1. set parameters
        self.popsize = 4 + int(3 * np.log(self.dim))
        self.select_size = int(self.popsize / 2)
        weights_dash = self.cal_weights_dash()
        self.mu_eff = np.sum(weights_dash[:self.select_size]) ** 2 / np.sum(weights_dash[:self.select_size] ** 2)
        mu_eff_bar = np.sum(weights_dash[self.select_size:]) ** 2 / np.sum(weights_dash[self.select_size:] ** 2)

        # 1.1 parameters for covariance matrix adaptation
        self.c1 = 2 / ((self.popsize + 1.3) ** 2 + self.mu_eff)
        self.cc = (4 + self.mu_eff / self.popsize) / (self.popsize + 4 + 2 * self.mu_eff / self.popsize)
        self.c_mu = min(1 - self.c1, 2 * (self.mu_eff - 2 + 1 / self.mu_eff) / ((self.popsize + 2) ** 2 + self.mu_eff))

        # 1.2 parameters for step size control
        self.c_sigma = (self.mu_eff + 2) / (self.popsize + self.mu_eff + 5)
        self.d_sigma = 1 + 2 * max(0, np.sqrt((self.mu_eff - 1) / (self.popsize + 1)) - 1) + self.c_sigma

        # 1.3 parameters for selection and recombination
        self.weights = self.cal_weights(self.c1, self.c_mu, mu_eff_bar, weights_dash)
        self.cm = 1

        # 2. init parameters
        self.l = np.zeros((self.dim, 1))
        self.u = np.ones((self.dim, 1))

        self.boundary_w = np.zeros((self.dim, 1))

    def init_params(self):
        path_sigma = np.zeros((self.dim, 1))
        path_c = np.zeros((self.dim, 1))
        C = np.eye(self.dim)
        sigma = 0.3 * (self.u[0] - self.l[0])
        mean = np.random.uniform(self.l, self.u, size=(self.dim, 1))

        return path_sigma, path_c, C, sigma, mean

    def get_suggestion(self, mean, sigma, C):
        # sample new population of search points

        y = np.random.multivariate_normal(
            mean=np.zeros((self.dim,)),
            cov=C,
            size=self.popsize
        )
        x = mean.T + sigma * y

        # if the mean is out of bounds, set the boundary weights
        if np.sum(np.less(mean, self.l)) > 0 or np.sum(np.greater(mean, self.u)) > 0:
            self.boundary_w = set_boundary_weights(sigma, C)

        # increase the boundary weights if the mean transcend the boundary too much
        threshold = 3 * sigma * np.diag(C) * max(1, np.sqrt(self.dim) / self.mu_eff)
        lower = np.greater(self.l - mean, threshold)
        upper = np.greater(mean - self.u, threshold)
        merge = np.logical_or(lower, upper)
        self.boundary_w = self.boundary_w * (1.1 ** (max(1, self.mu_eff / (10 * self.dim)) * merge))

        # selection and recombination
        suggestions = []
        for i in range(y.shape[0]):
            x_temp = x[i, :].reshape(x[i, :].shape[0], 1)
            lower_violate = np.less(x_temp, self.l)
            upper_violate = np.greater(x_temp, self.u)
            if np.sum(lower_violate) > 0 or np.sum(upper_violate) > 0:
                boundary_param = cal_boundary_param(C)
                # get the closest feasible suggestion
                x_feasible = x_temp + lower_violate*(self.l-x_temp)
                x_feasible = x_feasible - upper_violate*(x_feasible-self.u)

                # calculate penalty term
                penalty = cal_penalty(x_temp, x_feasible, self.boundary_w, boundary_param)

                suggestions.append(dict(
                    suggestion=np.squeeze(self.scaler.inverse_transform(x_feasible.T)),
                    penalty=penalty,
                ))
            else:
                suggestions.append(dict(
                    suggestion=np.squeeze(self.scaler.inverse_transform(x[i:i + 1], )),
                    penalty=0,
                ))
        return suggestions

    def report_metric(self, objective_dict, mean, sigma, C, path_sigma, path_c):
        for i in range(len(objective_dict)):
            objective_dict[i]["x"] = np.squeeze(self.scaler.transform(objective_dict[i]["x"].reshape(1, self.dim)))
            objective_dict[i]["x"] = (objective_dict[i]["x"] - mean.T) / sigma
            objective_dict[i]["x"] = np.squeeze(objective_dict[i]["x"])

            objective_dict[i]["y"] += objective_dict[i]["penalty"]
        objective_values = sorted(objective_dict, key=lambda k: k["y"])
        sorted_y = []
        for i in range(self.popsize):
            sorted_y.append(objective_values[i]["x"])
        sorted_y = np.array(sorted_y)

        y_w = np.sum(self.weights[:self.select_size, ] * sorted_y[:self.select_size, ], axis=0)
        y_w = y_w.reshape((y_w.shape[0], 1))
        next_mean = mean + self.cm * sigma * y_w

        eigenvalue, B = np.linalg.eig(C)
        D_inverse = np.diag(1 / np.sqrt(eigenvalue))

        # step size control
        next_path_sigma = (1 - self.c_sigma) * path_sigma + np.sqrt(
            self.c_sigma * (2 - self.c_sigma) * self.mu_eff) * np.dot(
            np.dot(np.dot(B, D_inverse), B.T), y_w)
        expectation = np.sqrt(2) * gamma((self.dim + 1) / 2) / gamma(self.dim / 2)
        next_sigma = sigma * np.exp(self.c_sigma / self.d_sigma * (np.sqrt(np.sum(next_path_sigma ** 2)) / expectation - 1))

        # covariance matrix adaptation
        next_path_c = (1 - self.cc) * path_c + np.sqrt(self.cc * (2 - self.cc) / self.mu_eff) * y_w
        weight_node = []
        for i in range(self.popsize):
            if self.weights[i] >= 0:
                weight_node.append(self.weights[i])
            else:
                temp = np.dot(np.dot(np.dot(B, D_inverse), B.T), sorted_y[i,].reshape(sorted_y[i,].shape[0], 1))
                norm = self.dim / np.sum(temp ** 2)
                weight_node.append(norm * self.weights[i])
        weight_sum = np.zeros((self.dim, self.dim))
        for i in range(self.popsize):
            vec = sorted_y[i, :].reshape(sorted_y[i, :].shape[0], 1)
            weight_sum = weight_sum + weight_node[i] * np.dot(vec, vec.T)
        next_C = (1 - self.c1 - self.c_mu * np.sum(self.weights)) * C + self.c1 * np.dot(next_path_c,
                                                                                    next_path_c.T) + self.c_mu * weight_sum

        return next_path_sigma, next_path_c, next_C, next_sigma, next_mean

    def cal_weights_dash(self):
        weights_dash = []
        for i in range(self.popsize):
            weights_dash.append(np.log((self.popsize + 1) / 2) - np.log(i + 1))
        return np.array(weights_dash)

    def cal_weights(self, c1, c_mu, mu_eff_bar, weights_dash):
        temp_mu = 1 + c1 / c_mu
        temp_mu_eff = 1 + 2 * mu_eff_bar / (self.mu_eff + 2)
        temp_posdef = (1 - c1 - c_mu) / (self.popsize * c_mu)

        positive_sum = sum_of_sign(weights_dash, "positive")
        negative_sum = sum_of_sign(weights_dash, "negative")
        weights = []
        for i in range(self.popsize):
            if weights_dash[i] >= 0:
                weights.append(weights_dash[i] / positive_sum)
            else:
                weights.append(min(temp_mu, temp_mu_eff, temp_posdef) * weights_dash[i] / negative_sum)

        return np.array(weights).reshape((self.popsize, 1))
