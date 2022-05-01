# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

""" module for algorithm manager """
import numpy as np

from pkg.apis.manager.v1beta1.python import api_pb2

from pkg.suggestion.v1beta1.bayesianoptimization.utils import get_logger


def deal_with_discrete(feasible_values, current_value):
    """ function to embed the current values to the feasible discrete space"""
    diff = np.subtract(feasible_values, current_value)
    diff = np.absolute(diff)
    return feasible_values[np.argmin(diff)]


def deal_with_categorical(feasible_values, one_hot_values):
    """ function to do the one hot encoding of the categorical values """
    index = np.argmax(one_hot_values)
    #index = one_hot_values.argmax()
    return feasible_values[int(index)]


class AlgorithmManager:
    """ class for the algorithm manager
    provide some helper functions
    """

    def __init__(self, experiment_name, experiment, parameter_config, X_train, y_train, logger=None):
        self.logger = logger if (logger is not None) else get_logger()
        self._experiment_name = experiment_name
        self._experiment = experiment
        self._goal = self._experiment.spec.objective.type
        self._dim = parameter_config.dim
        self._lowerbound = parameter_config.lower_bounds
        self._upperbound = parameter_config.upper_bounds
        self._types = parameter_config.parameter_types
        self._names = parameter_config.names
        # record all the feasible values of discrete type variables
        self._discrete_info = parameter_config.discrete_info
        self._categorical_info = parameter_config.categorical_info
        self._name_id = parameter_config.name_ids

        self._X_train = self._mapping_params(X_train)
        self.parse_X()

        self._y_train = y_train
        self._parse_metric()

    @property
    def experiment_name(self):
        """ return the experiment_name """
        return self._experiment_name

    @property
    def experiment(self):
        """ return the experiment """
        return self._experiment

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
        """ return the upper bound of all the parameters """
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

    def _mapping_params(self, parameters_list):
        if len(parameters_list) == 0:
            return None
        ret = []
        for parameters in parameters_list:
            maplist = [np.zeros(1)]*len(self._names)
            for p in parameters:
                self.logger.debug("mapping: %r", p, extra={
                    "Experiment": self._experiment_name})
                map_id = self._name_id[p.name]
                if self._types[map_id] in [api_pb2.DOUBLE, api_pb2.INT, api_pb2.DISCRETE]:
                    maplist[map_id] = float(p.value)
                elif self._types[map_id] == api_pb2.CATEGORICAL:
                    for ci in self._categorical_info:
                        if ci["name"] == p.name:
                            maplist[map_id] = np.zeros(ci["number"])
                            for i, v in enumerate(ci["values"]):
                                if v == p.value:
                                    maplist[map_id][i] = 1
                                    break
            self.logger.debug("mapped: %r", maplist, extra={
                "Experiment": self._experiment_name})
            ret.append(np.hstack(maplist))
        return ret

    def _parse_metric(self):
        """ parse the metric to the dictionary """
        self.logger.info("Ytrain: %r", self._y_train, extra={
            "Experiment": self._experiment_name})
        if not self._y_train:
            self._y_train = None
            return
        y = []
        for metric in self._y_train:
            if self._goal == api_pb2.MAXIMIZE:
                y.append(float(metric))
            else:
                y.append(-float(metric))
        self.logger.debug("Ytrain: %r", y, extra={
            "Experiment": self._experiment_name})
        self._y_train = np.array(y)

    def parse_X(self):
        if not self._X_train:
            self._X_train = None
            return
        self.logger.debug("Xtrain: %r", self._X_train, extra={
            "Experiment": self._experiment_name})
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
                            deal_with_discrete(
                                param["values"], x_next[counter])
                        )
                        counter = counter + 1
                        break
            elif self._types[i] == api_pb2.CATEGORICAL:
                for param in self._categorical_info:
                    if param["name"] == self._names[i]:
                        result.append(deal_with_categorical(
                            feasible_values=param["values"],
                            one_hot_values=x_next[counter:counter +
                                                  param["number"]],
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
