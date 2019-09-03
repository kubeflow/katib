from logging import getLogger, StreamHandler, INFO, DEBUG
import itertools
import grpc
import numpy as np
from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from . import parsing_util

class GridService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        self.manager_addr = "katib-manager"
        self.manager_port = 6789
        self.default_grid = 10

    def _get_experiment(self, name):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            exp = client.GetExperiment(api_pb2.GetExperimentRequest(experiment_name=name), 10)
            return exp.experiment

    def _get_algorithm_settings(self, experiment_name):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            alg = client.GetAlgorithmExtraSettings(api_pb2.GetAlgorithmExtraSettingsRequest(
                experiment_name=experiment_name), 10)
            params = alg.extra_algorithm_settings
            alg_settings = {}
            for param in params:
                alg_settings[param.name] = param.value
            return alg_settings

    def _get_trials(self, experiment_name):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            trials = client.GetTrialList(api_pb2.GetTrialListRequest(
                experiment_name=experiment_name), 10)
            return trials.trials

    def _create_all_combinations(self, parameters, alg_settings):
        param_ranges = []
        cur_index = 0
        parameter_config = parsing_util.parse_parameter_configs(parameters)
        default_grid_size = alg_settings.get("DefaultGrid", self.default_grid)
        for idx, param_type in enumerate(parameter_config.parameter_types):
            param_name = parameter_config.names[idx]
            if param_type in [api_pb2.DOUBLE, api_pb2.INT]:
                num = alg_settings.get(param_name, default_grid_size)
                param_values = \
                    np.linspace(parameter_config.lower_bounds[0, cur_index],
                                parameter_config.upper_bounds[0, cur_index],
                                num=num)
                cur_index += 1
                if param_type == api_pb2.INT:
                    param_values = param_values.astype(np.int64)
            elif param_type == api_pb2.DISCRETE:
                for discrete_param in parameter_config.discrete_info:
                    if param_name == discrete_param["name"]:
                        param_values = discrete_param["values"]
                        break
                cur_index += 1
            elif param_type == api_pb2.CATEGORICAL:
                for categ_param in parameter_config.categorical_info:
                    if param_name == categ_param["name"]:
                        param_values = categ_param["values"]
                        break
                cur_index += categ_param["number"]
            param_ranges.append(param_values)
        all_combinations = [comb for comb in itertools.product(*param_ranges)]
        return all_combinations, parameter_config

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        experiment_name = request.experiment_name
        request_number = request.request_number
        experiment = self._get_experiment(experiment_name)
        parameters = experiment.spec.parameter_specs.parameters
        alg_settings = self._get_algorithm_settings(experiment_name)
        combinations, parameter_config = self._create_all_combinations(parameters, alg_settings)
        total_combinations = len(combinations)

        allocated_trials = self._get_trials(experiment_name)
        total_allocated_trials = len(allocated_trials)
        return_start_index = total_allocated_trials
        return_end_index = return_start_index + request_number

        if return_start_index > total_combinations:
            return_start_index = 0
            return_end_index = return_start_index + request_number
        elif return_start_index + request_number > total_combinations:
            return_start_index = total_combinations - request_number
            return_end_index = total_combinations
        if return_start_index < 0:
            return_start_index = 0

        trial_specs = []
        for elem in combinations[return_start_index:return_end_index]:
            suggestion = parsing_util.parse_x_next_tuple(elem, parameter_config.parameter_types,
                                                         parameter_config.names)
            trial_spec = api_pb2.TrialSpec()
            trial_spec.experiment_name = experiment_name
            for param in suggestion:
                trial_spec.parameter_assignments.assignments.add(name=param['name'],
                                                                 value=str(param['value']))
            trial_specs.append(trial_spec)
        reply = api_pb2.GetSuggestionsReply()
        for trial_spec in trial_specs:
            reply.trials.add(spec=trial_spec)
        return reply
