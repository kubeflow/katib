import math
import traceback

import logging
from logging import getLogger, StreamHandler, DEBUG
from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc
import grpc
from . import parsing_util

class HyperbandService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        self.manager_addr = "katib-manager"
        self.manager_port = 6789
        FORMAT = '%(asctime)-15s Experiment %(experiment_name)s %(message)s'
        self.logger = getLogger(__name__)
        logging.basicConfig(format=FORMAT)
        handler = StreamHandler()
        handler.setLevel(DEBUG)
        self.logger.setLevel(DEBUG)
        self.logger.addHandler(handler)

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
            reply = client.GetTrialList(api_pb2.GetTrialListRequest(experiment_name=experiment_name), 10)
            return reply.trials

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        try:
          reply = api_pb2.GetSuggestionsReply()
          experiment_name = request.experiment_name
          experiment = self._get_experiment(experiment_name)
          alg_settings = self._get_algorithm_settings(experiment_name)

          sParams = self._parse_suggestionParameters(experiment, alg_settings)
          if sParams["current_s"] < 0:
            return reply

          trials = self._make_bracket(experiment, sParams)
          self._update_algorithm_extrasettings(experiment_name, sParams)
          for trial in trials:
            reply.trials.add(spec=trial)
          return reply
        except Exception as e:
            self.logger.error("Fail to generate trials: \n%s",
                    traceback.format_exc(), extra={"experiment_name": experiment_name})
            raise e

    def _update_hbParameters(self, sParams):
        sParams["current_i"] += 1
        if sParams["current_i"] > sParams["current_s"]:
            self._new_hbParameters(sParams)

    def _new_hbParameters(self, sParams):
        sParams["current_s"] -= 1
        sParams["current_i"] = 0
        if sParams["current_s"] >= 0:
        # when sParams["current_s"] < 0, hyperband algorithm reaches the end
            sParams["n"] = int(math.ceil(float(sParams["sMax"] + 1) * (
                float(sParams["eta"]**sParams["current_s"]) / float(sParams["current_s"]+1))))
            sParams["r"] = sParams["r_l"]*sParams["eta"]**(-sParams["current_s"])

    def _parse_suggestionParameters(self, experiment, alg_settings):
        sParams = {
            "eta":              3,
            "sMax":             -1,
            "r_l":              -1,
            "b_l":              -1,
            "r":                -1,
            "n":                -1,
            "current_s":        -2,
            "current_i":        -1,
            "resourceName":     "",
            "evaluatingTrials": 0,
        }

        for k, v in alg_settings.items():
            if k in ["eta", "r_l", "b_l"]:
                sParams[k] = float(v)
            elif k in ["n", "r", "current_s", "current_i", "sMax", "evaluatingTrials"]:
                sParams[k] = int(float(v))
            elif k == "resourceName":
                sParams[k] = v
            else:
                self.logger.info("Unknown HyperBand Param %s, ignore it",
                        k, extra={"experiment_name": experiment.name})
        if sParams["current_s"] == -1:
            # Hyperband outlerloop has finished
            self.logger.info("HyperBand outlerloop has finished.",
                    extra={"experiment_name": experiment.name})
            return sParams

        if sParams["eta"] <= 0:
            sParams["eta"] = 3
        if sParams["sMax"] < 0:
            sParams["sMax"] = int(math.log(sParams["r_l"]) / math.log(sParams["eta"]))
        if sParams["b_l"] < 0:
            sParams["b_l"] = (sParams["sMax"] + 1) * sParams["r_l"]
        if sParams["current_s"] < 0:
            sParams["current_s"] = sParams["sMax"]
        if sParams["current_i"] < 0:
            sParams["current_i"] = 0
        if sParams["n"] < 0:
            sParams["n"] = int(math.ceil(float(sParams["sMax"] + 1) * (
                float(sParams["eta"]**sParams["current_s"]) / float(sParams["current_s"]+1))))
        if sParams["r"] < 0:
            sParams["r"] = sParams["r_l"]*sParams["eta"]**(-sParams["current_s"])

        return sParams

    def _make_bracket(self, experiment, sParams):
        if sParams["evaluatingTrials"] == 0:
            trialSpecs = self._make_master_bracket(experiment, sParams)
        else:
            trialSpecs = self._make_child_bracket(experiment, sParams)
        if sParams["current_i"] < sParams["current_s"]:
            sParams["evaluatingTrials"] = len(trialSpecs)
        else:
            sParams["evaluatingTrials"] = 0

        self.logger.info("HyperBand Param eta %d.",
                sParams["eta"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param R %d.",
                sParams["r_l"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param sMax %d.",
                sParams["sMax"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param B %d.",
                sParams["b_l"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param n %d.",
                sParams["n"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param r %d.",
                sParams["r"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param s %d.",
                sParams["current_s"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param i %d.",
                sParams["current_i"], extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand evaluating trials count %d.",
                sParams["evaluatingTrials"], extra={"experiment_name": experiment.name})

        if sParams["evaluatingTrials"] == 0:
            self._new_hbParameters(sParams)

        return trialSpecs

    def _make_child_bracket(self, experiment, sParams):
        n_i = math.ceil(sParams["n"] * sParams["eta"]**(-sParams["current_i"]))
        top_trials_num = int(math.ceil(n_i / sParams["eta"]))
        self._update_hbParameters(sParams)
        r_i = int(sParams["r"] * sParams["eta"]**sParams["current_i"])
        last_trials = self._get_top_trial(sParams["evaluatingTrials"], top_trials_num, experiment)
        trialSpecs = self._copy_trials(last_trials, r_i, sParams["resourceName"])

        self.logger.info("Generate %d trials by child bracket.",
                top_trials_num, extra={"experiment_name": experiment.name})
        return trialSpecs

    def _get_last_trials(self, all_trials, latest_trials_num):
        sorted_trials = sorted(all_trials, key=lambda trial: trial.status.start_time)
        if len(sorted_trials) > latest_trials_num:
            return sorted_trials[-latest_trials_num:]
        else:
            return sorted_trials

    def _get_top_trial(self, latest_trials_num, top_trials_num, experiment):
        objective_metric = experiment.spec.objective.objective_metric_name
        objective_type = experiment.spec.objective.type

        def get_objective_value(t):
            for m in t.status.observation.metrics:
                if m.name == objective_metric:
                    return float(m.value)

        top_trials = []
        all_trials = self._get_trials(experiment.name)
        latest_trials = self._get_last_trials(all_trials, latest_trials_num)

        for t in latest_trials:
            if t.status.condition != api_pb2.TrialStatus.TrialConditionType.SUCCEEDED:
                raise Exception(
                    "There are some trials which are not completed yet for experiment %s." % experiment.name)

        if objective_type == api_pb2.MAXIMIZE:
            top_trials.extend(sorted(latest_trials, key=get_objective_value, reverse=True))
        else:
            top_trials.extend(sorted(latest_trials, key=get_objective_value))
        return top_trials[:top_trials_num]

    def _copy_trials(self, trials, r_i, resourceName):
        trialSpecs = []
        for t in trials:
            trial_spec = api_pb2.TrialSpec()
            for assignment in t.spec.parameter_assignments.assignments:
                if assignment.name == resourceName:
                    value = str(r_i)
                else:
                    value = assignment.value
                trial_spec.parameter_assignments.assignments.add(name=assignment.name,
                                                                 value=value)
            trialSpecs.append(trial_spec)
        return trialSpecs

    def _make_master_bracket(self, experiment, sParams):
        n = sParams["n"]
        r = int(sParams["r"])
        parameter_config = parsing_util.parse_parameter_configs(
            experiment.spec.parameter_specs.parameters)
        trial_specs = []
        for _ in range(n):
            sample = parameter_config.random_sample()
            suggestion = parsing_util.parse_x_next_vector(
                                          sample,
                                          parameter_config.parameter_types,
                                          parameter_config.names,
                                          parameter_config.discrete_info,
                                          parameter_config.categorical_info)
            trial_spec = api_pb2.TrialSpec()
            trial_spec.experiment_name = experiment.name
            for param in suggestion:
                if param['name'] == sParams["resourceName"]:
                    param['value'] = str(r)
                trial_spec.parameter_assignments.assignments.add(name=param['name'],
                                                                 value=str(param['value']))
            trial_specs.append(trial_spec)
        self.logger.info("Generate %d trials by master bracket.",
                n, extra={"experiment_name": experiment.name})
        return trial_specs

    def _update_algorithm_extrasettings(self, experiment_name, sParams):
        as_list = []
        for k, v in sParams.items():
            as_list.append(api_pb2.AlgorithmSetting(name=k, value=str(v)))
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            client.UpdateAlgorithmExtraSettings(api_pb2.UpdateAlgorithmExtraSettingsRequest(
                experiment_name=experiment_name, extra_algorithm_settings=as_list), 10)

    def _set_validate_context_error(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        self.logger.info(error_message)
        return api_pb2.ValidateAlgorithmSettingsReply()

    def ValidateAlgorithmSettings(self, request, context):
        params = request.experiment_spec.parameter_specs.parameters
        settings = request.experiment_spec.algorithm.algorithm_setting
        setting_dict = {}
        for setting in settings:
            setting_dict[setting.name] = setting.value
        if "r_l" not in setting_dict or "resourceName" not in setting_dict:
            return self._set_validate_context_error(context, "r_l and resourceName must be set.")
        try:
            rl = float(setting_dict["r_l"])
        except:
            return self._set_validate_context_error(context, "r_l must be a positive float number.")
        else:
            if rl < 0:
                return self._set_validate_context_error(context, "r_l must be a positive float number.")

        if "eta" in setting_dict:
            eta = int(float(setting_dict["eta"]))
            if eta <= 0:
                eta = 3
        else:
            eta = 3

        smax = int(math.log(rl)/math.log(eta))
        max_parallel = int(math.ceil(eta**smax))
        if request.experiment_spec.parallel_trial_count < max_parallel:
            return self._set_validate_context_error(context,
                            "parallelTrialCount must be not less than %d." % max_parallel)

        valid_resourceName = False
        for param in params:
            if param.name == setting_dict["resourceName"]:
                valid_resourceName = True
                break
        if not valid_resourceName:
            return self._set_validate_context_error(context,
                            "value of resourceName setting must be in parameters.")

        return api_pb2.ValidateAlgorithmSettingsReply()
