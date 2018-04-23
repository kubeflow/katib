""" module for algorithm manager """

import numpy as np

from pkg.api.python import api_pb2


def deal_with_discrete(feasible_values, current_value):
    """ function to embed the current values to the feasible discrete space"""
    diff = np.subtract(feasible_values, current_value)
    diff = np.absolute(diff)
    return feasible_values[np.argmin(diff)]


def deal_with_categorical(feasible_values, one_hot_values):
    """ function to do the one hot encoding of the categorical values """
    index = np.argmax(one_hot_values)
    return feasible_values[index]


class AlgorithmManager:
    """ class for the algorithm manager
    provide some helper functions
    """
    def __init__(self, study_id, study_config, X_train, y_train):
        self._study_id = study_id
        self._study_config = study_config
        self._goal = self._study_config.optimization_type
        self._dim = 0
        self._lowerbound = []
        self._upperbound = []
        self._types = []
        self._names = []
        # record all the feasible values of discrete type variables
        self._discrete_info = []
        self._categorical_info = []

        self._parse_config()

        self._X_train = X_train
        self.parse_X()

        self._y_train = y_train
        self._parse_metric()

        # print(self._X_train)
        # print(self._y_train)

    @property
    def study_id(self):
        """ return the study id """
        return self._study_id

    @property
    def study_config(self):
        """ return the study configuration """
        return self._study_config

    @property
    def goal(self):
        """ return the optimization goal"""
        return self._goal

    @property
    def dim(self):
        """ return the dimension """
        return self._dim

    @property
    def lower_bound(self):
        """ return the lower bound of all the parameters """
        return self._lowerbound

    @property
    def upper_bound(self):
        """ return the ipper bound of all the parameters """
        return self._upperbound

    @property
    def types(self):
        """ return the types of all the parameters """
        return self._types

    @property
    def names(self):
        """ return the names of all the parameters """
        return self._names

    @property
    def discrete_info(self):
        """ return the info of all the discrete parameters """
        return self._discrete_info

    @property
    def categorical_info(self):
        """ return the info of all the categorical parameters """
        return self._categorical_info

    @property
    def X_train(self):
        """ return the training data """
        return self._X_train

    @property
    def y_train(self):
        """ return the target of the training data"""
        return self._y_train

    def _parse_config(self):
        """ extract info from the study configuration """
        for param in self._study_config.parameter_configs.configs:
            self._types.append(param.parameter_type)
            self._names.append(param.name)
            if param.parameter_type == api_pb2.DOUBLE or param.parameter_type == api_pb2.INT:
                self._dim = self._dim + 1
                self._lowerbound.append(float(param.feasible.min))
                self._upperbound.append(float(param.feasible.max))
            elif param.parameter_type == api_pb2.DISCRETE:
                self._dim = self._dim + 1
                discrete_values = [int(x) for x in param.feasible.list]
                min_value = min(discrete_values)
                max_value = max(discrete_values)
                self._lowerbound.append(min_value)
                self._upperbound.append(max_value)
                self._discrete_info.append(dict({
                    "name": param.name,
                    "values": discrete_values,
                }))
            # one hot encoding for categorical type
            elif param.parameter_type == api_pb2.CATEGORICAL:
                num_feasible = len(param.feasible.list)
                for i in range(num_feasible):
                    self._lowerbound.append(0)
                    self._upperbound.append(1)
                self._categorical_info.append(dict({
                    "name": param.name,
                    "values": param.feasible.list,
                    "number": num_feasible,
                }))
                self._dim += num_feasible

    def _parse_metric(self):
        """ parse the metric to the dictionary """
        if not self._y_train:
            self._y_train = None
            return
        y = []
        for metric in self._y_train:
            if self._goal == api_pb2.MAXIMIZE:
                y.append(float(metric))
            else:
                y.append(-float(metric))

        self._y_train = np.array(y)

    def parse_X(self):
        if not self._X_train:
            self._X_train = None
            return

        self._X_train = np.array(self._X_train)

    def parse_x_next(self, x_next):
        """ parse the next suggestion to the proper format """
        counter = 0
        result = []
        for i in range(len(self._types)):
            if self._types[i] == api_pb2.INT:
                result.append(int(round(x_next[counter], 0)))
                counter = counter + 1
            elif self._types[i] == api_pb2.DISCRETE:
                for param in self._discrete_info:
                    if param["name"] == self._names[i]:
                        result.append(
                            deal_with_discrete(param["values"], x_next[counter])
                        )
                        counter = counter + 1
                        break
            elif self._types[i] == api_pb2.CATEGORICAL:
                for param in self._categorical_info:
                    if param["name"] == self._names[i]:
                        result.append(deal_with_categorical(
                            feasible_values=param["values"],
                            one_hot_values=x_next[counter:counter + param["number"]],
                        ))
                        counter = counter + param["number"]
                        break
            elif self._types[i] == api_pb2.DOUBLE:
                result.append(x_next[counter])
                counter = counter + 1
        return result

    def convert_to_dict(self, x_next):
        """ convert the next suggestion to the dictionary """
        result = []
        for i in range(len(x_next)):
            tmp = dict({
                "name": self._names[i],
                "value": x_next[i],
                "type": self._types[i],
            })
            result.append(tmp)
        return result

