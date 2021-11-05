# Copyright 2021 The Kubeflow Authors.
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
from logging import getLogger, StreamHandler, INFO
import json

from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer
from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc


class DartsService(api_pb2_grpc.SuggestionServicer, HealthServicer):

    def __init__(self):
        super(DartsService, self).__init__()
        self.is_first_run = True

        self.logger = getLogger(__name__)
        FORMAT = '%(asctime)-15s Experiment %(experiment_name)s %(message)s'
        logging.basicConfig(format=FORMAT)
        handler = StreamHandler()
        handler.setLevel(INFO)
        self.logger.setLevel(INFO)
        self.logger.addHandler(handler)
        self.logger.propagate = False

    # TODO: Add validation
    def ValidateAlgorithmSettings(self, request, context):
        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        if self.is_first_run:
            nas_config = request.experiment.spec.nas_config
            num_layers = str(nas_config.graph_config.num_layers)

            search_space = get_search_space(nas_config.operations)

            settings_raw = request.experiment.spec.algorithm.algorithm_settings
            algorithm_settings = get_algorithm_settings(settings_raw)

            search_space_json = json.dumps(search_space)
            algorithm_settings_json = json.dumps(algorithm_settings)

            search_space_str = str(search_space_json).replace('\"', '\'')
            algorithm_settings_str = str(algorithm_settings_json).replace('\"', '\'')

            self.is_first_run = False

        parameter_assignments = []
        for i in range(request.current_request_number):

            self.logger.info(">>> Generate new Darts Trial Job")

            self.logger.info(">>> Number of layers {}\n".format(num_layers))

            self.logger.info(">>> Search Space")
            self.logger.info("{}\n".format(search_space_str))

            self.logger.info(">>> Algorithm Settings")
            self.logger.info("{}\n\n".format(algorithm_settings_str))

            parameter_assignments.append(
                api_pb2.GetSuggestionsReply.ParameterAssignments(
                    assignments=[
                        api_pb2.ParameterAssignment(
                            name="algorithm-settings",
                            value=algorithm_settings_str
                        ),
                        api_pb2.ParameterAssignment(
                            name="search-space",
                            value=search_space_str
                        ),
                        api_pb2.ParameterAssignment(
                            name="num-layers",
                            value=num_layers
                        )
                    ]
                )
            )

        return api_pb2.GetSuggestionsReply(parameter_assignments=parameter_assignments)


def get_search_space(operations):
    search_space = []

    for operation in list(operations.operation):
        opt_type = operation.operation_type

        if opt_type == "skip_connection":
            search_space.append(opt_type)
        else:
            # Currently support only one Categorical parameter - filter size
            opt_spec = list(operation.parameter_specs.parameters)[0]
            for filter_size in list(opt_spec.feasible_space.list):
                search_space.append(opt_type+"_{}x{}".format(filter_size, filter_size))
    return search_space


def get_algorithm_settings(settings_raw):

    algorithm_settings_default = {
        "num_epochs":           50,
        "w_lr":                 0.025,
        "w_lr_min":             0.001,
        "w_momentum":           0.9,
        "w_weight_decay":       3e-4,
        "w_grad_clip":          5.,
        "alpha_lr":             3e-4,
        "alpha_weight_decay":   1e-3,
        "batch_size":           128,
        "num_workers":          4,
        "init_channels":        16,
        "print_step":           50,
        "num_nodes":            4,
        "stem_multiplier":      3,
    }

    for setting in settings_raw:
        s_name = setting.name
        s_value = setting.value
        algorithm_settings_default[s_name] = s_value

    return algorithm_settings_default
