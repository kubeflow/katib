from logging import getLogger, StreamHandler, INFO, DEBUG
from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc
import grpc
from . import parsing_util

class RandomService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        self.manager_addr = "katib-manager"
        self.manager_port = 6789

    def _get_experiment(self, name):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            exp = client.GetExperiment(api_pb2.GetExperimentRequest(experiment_name=name), 10)
            return exp.experiment

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        experiment = self._get_experiment(request.experiment_name)
        parameter_config = parsing_util.parse_parameter_configs(
            experiment.spec.parameter_specs.parameters)
        trial_specs = []
        for _ in range(request.request_number):
            sample = parameter_config.random_sample()
            suggestion = parsing_util.parse_x_next_vector(sample,
                                                           parameter_config.parameter_types,
                                                           parameter_config.names,
                                                           parameter_config.discrete_info,
                                                           parameter_config.categorical_info)
            trial_spec = api_pb2.TrialSpec()
            trial_spec.experiment_name = request.experiment_name
            for param in suggestion:
                trial_spec.parameter_assignments.assignments.add(name=param['name'],
                                                                 value=str(param['value']))
            trial_specs.append(trial_spec)

        reply = api_pb2.GetSuggestionsReply()
        for trial_spec in trial_specs:
            reply.trials.add(spec=trial_spec)

        return reply
