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
import grpc

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc

from pkg.suggestion.v1beta1.internal.constant import INTEGER, DOUBLE, CATEGORICAL, DISCRETE
from pkg.suggestion.v1beta1.internal.search_space import HyperParameterSearchSpace
from pkg.suggestion.v1beta1.internal.trial import Trial, Assignment
from pkg.suggestion.v1beta1.chocolate.base_service import BaseChocolateService
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer

import numpy as np
import itertools

logger = logging.getLogger(__name__)


class ChocolateService(api_pb2_grpc.SuggestionServicer, HealthServicer):
    def __init__(self):
        super(ChocolateService, self).__init__()
        self.base_service = None
        self.is_first_run = True

    def ValidateAlgorithmSettings(self, request, context):
        algorithm_name = request.experiment.spec.algorithm.algorithm_name
        if algorithm_name == "grid":
            search_space = HyperParameterSearchSpace.convert(
                request.experiment)
            available_space = {}
            for param in search_space.params:
                if param.type == INTEGER:
                    available_space[param.name] = range(int(param.min), int(param.max)+1, int(param.step))

                elif param.type == DOUBLE:
                    if param.step == "" or param.step is None:
                        return self._set_validate_context_error(
                            context, "Param: {} step is nil".format(param.name))
                    double_list = np.arange(float(param.min), float(param.max)+float(param.step), float(param.step))
                    if double_list[-1] > float(param.max):
                        double_list = double_list[:-1]
                    available_space[param.name] = double_list

                elif param.type == CATEGORICAL or param.type == DISCRETE:
                    available_space[param.name] = param.list

            num_combinations = len(list(itertools.product(*available_space.values())))
            max_trial_count = request.experiment.spec.max_trial_count

            if max_trial_count > num_combinations:
                return self._set_validate_context_error(
                    context, "Max Trial Count: {} > all possible search space combinations: {}".format(
                        max_trial_count, num_combinations)
                )

        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """

        if self.is_first_run:
            search_space = HyperParameterSearchSpace.convert(
                request.experiment)
            self.base_service = BaseChocolateService(
                algorithm_name=request.experiment.spec.algorithm.algorithm_name,
                search_space=search_space)
            self.is_first_run = False

        trials = Trial.convert(request.trials)
        new_assignments = self.base_service.getSuggestions(
            trials, request.current_request_number, request.total_request_number)
        return api_pb2.GetSuggestionsReply(
            parameter_assignments=Assignment.generate(new_assignments)
        )

    def _set_validate_context_error(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        logger.info(error_message)
        return api_pb2.ValidateAlgorithmSettingsReply()
