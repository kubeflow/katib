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

import logging

from pkg.apis.manager.v1beta1.python import api_pb2 as api
import pkg.suggestion.v1beta1.internal.constant as constant


logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


class HyperParameterSearchSpace(object):
    def __init__(self):
        self.goal = ""
        self.params = []

    @staticmethod
    def convert(experiment):
        search_space = HyperParameterSearchSpace()
        if experiment.spec.objective.type == api.MAXIMIZE:
            search_space.goal = constant.MAX_GOAL
        elif experiment.spec.objective.type == api.MINIMIZE:
            search_space.goal = constant.MIN_GOAL
        for p in experiment.spec.parameter_specs.parameters:
            search_space.params.append(
                HyperParameterSearchSpace.convertParameter(p))
        return search_space

    def __str__(self):
        return "HyperParameterSearchSpace(goal: {}, ".format(self.goal) + \
            "params: {})".format(", ".join([element.__str__() for element in self.params]))

    @staticmethod
    def convertParameter(p):
        if p.parameter_type == api.INT:
            # Default value for INT parameter step is 1
            step = 1
            if p.feasible_space.step is not None and p.feasible_space.step != "":
                step = p.feasible_space.step
            return HyperParameter.int(p.name, p.feasible_space.min, p.feasible_space.max, step)
        elif p.parameter_type == api.DOUBLE:
            return HyperParameter.double(p.name, p.feasible_space.min, p.feasible_space.max, p.feasible_space.step)
        elif p.parameter_type == api.CATEGORICAL:
            return HyperParameter.categorical(p.name, p.feasible_space.list)
        elif p.parameter_type == api.DISCRETE:
            return HyperParameter.discrete(p.name, p.feasible_space.list)
        else:
            logger.error(
                "Cannot get the type for the parameter: %s (%s)", p.name, p.parameter_type)


class HyperParameter(object):
    def __init__(self, name, type_, min_, max_, list_, step):
        self.name = name
        self.type = type_
        self.min = min_
        self.max = max_
        self.list = list_
        self.step = step

    def __str__(self):
        if self.type == constant.INTEGER or self.type == constant.DOUBLE:
            return "HyperParameter(name: {}, type: {}, min: {}, max: {}, step: {})".format(
                self.name, self.type, self.min, self.max, self.step)
        else:
            return "HyperParameter(name: {}, type: {}, list: {})".format(
                self.name, self.type, ", ".join(self.list))

    @staticmethod
    def int(name, min_, max_, step):
        return HyperParameter(name, constant.INTEGER, min_, max_, [], step)

    @staticmethod
    def double(name, min_, max_, step):
        return HyperParameter(name, constant.DOUBLE, min_, max_, [], step)

    @staticmethod
    def categorical(name, lst):
        return HyperParameter(name, constant.CATEGORICAL, 0, 0, [str(e) for e in lst], 0)

    @staticmethod
    def discrete(name, lst):
        return HyperParameter(name, constant.DISCRETE, 0, 0, [str(e) for e in lst], 0)
