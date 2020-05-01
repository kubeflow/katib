import logging
from logging import getLogger, StreamHandler, INFO
import json

from pkg.suggestion.v1alpha3.internal.base_health_service import HealthServicer
from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc

logger = logging.getLogger(__name__)


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

            search_space = get_search_space(nas_config.operations, logger)

            settings_raw = request.experiment.spec.algorithm.algorithm_setting
            algorithm_settings = get_algorithm_settings(settings_raw)

            search_space_json = json.dumps(search_space)
            algorithm_settings_json = json.dumps(algorithm_settings)

            search_space_str = str(search_space_json).replace('\"', '\'')
            algorithm_settings_str = str(algorithm_settings_json).replace('\"', '\'')

            self.is_first_run = False

        parameter_assignments = []
        for i in range(request.request_number):

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


def get_search_space(operations, logger):
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


# TODO: Add more algorithm settings
def get_algorithm_settings(settings_raw):

    algorithm_settings_default = {
        "num_epoch": 50
    }

    for setting in settings_raw:
        s_name = setting.name
        s_value = setting.value
        algorithm_settings_default[s_name] = s_value

    return algorithm_settings_default
