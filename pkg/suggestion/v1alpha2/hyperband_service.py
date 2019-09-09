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
        FORMAT = '%(asctime)-15s Experiment %(experiment_name)s %(message)s'
        self.logger = getLogger(__name__)
        logging.basicConfig(format=FORMAT)
        handler = StreamHandler()
        handler.setLevel(DEBUG)
        self.logger.setLevel(DEBUG)
        self.logger.addHandler(handler)

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        try:
            reply = api_pb2.GetSuggestionsReply()
            experiment = request.experiment
            alg_settings = experiment.spec.algorithm.algorithm_setting

            param = HyperBandParam.convert(alg_settings)
            if param.current_s < 0:
                # Hyperband outlerloop has finished
                return reply

            trials = self._make_bracket(experiment, param)
            for trial in trials:
                reply.trials.add(spec=trial)
            reply.algorithm = HyperBandParam.generate(param)
            return reply
        except Exception as e:
            self.logger.error("Fail to generate trials: \n%s",
                              traceback.format_exc(), extra={"experiment_name": experiment_name})
            raise e

    def _update_hbParameters(self, param):
        param.current_i += 1
        if param.current_i > param.current_s:
            self._new_hbParameters(param)

    def _new_hbParameters(self, param):
        param.current_s -= 1
        param.current_i = 0
        if param.current_s >= 0:
            # when param.current_s < 0, hyperband algorithm reaches the end
            param.n = int(math.ceil(float(param.s_max + 1) * (
                float(param.eta**param.current_s) / float(param.current_s+1))))
            param.r = param.r_l * \
                param.eta**(-param.current_s)

    def _make_bracket(self, experiment, param):
        if param.evaluating_trials == 0:
            trialSpecs = self._make_master_bracket(experiment, param)
        else:
            trialSpecs = self._make_child_bracket(experiment, param)
        if param.current_i < param.current_s:
            param.evaluating_trials = len(trialSpecs)
        else:
            param.evaluating_trials = 0

        self.logger.info("HyperBand Param eta %d.",
                         param.eta, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param R %d.",
                         param.r_l, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param sMax %d.",
                         param.s_max, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param B %d.",
                         param.b_l, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param n %d.",
                         param.n, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param r %d.",
                         param.r, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param s %d.",
                         param.current_s, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand Param i %d.",
                         param.current_i, extra={"experiment_name": experiment.name})
        self.logger.info("HyperBand evaluating trials count %d.",
                         param.evaluating_trials, extra={"experiment_name": experiment.name})
        # TODO: Log resource name

        if param.evaluating_trials == 0:
            self._new_hbParameters(param)

        return trialSpecs

    def _make_child_bracket(self, experiment, param):
        n_i = math.ceil(param.n * param.eta**(-param.current_i))
        top_trials_num = int(math.ceil(n_i / param.eta))
        self._update_hbParameters(param)
        r_i = int(param.r * param.eta**param.current_i)
        last_trials = self._get_top_trial(
            param.evaluating_trials, top_trials_num, experiment)
        trialSpecs = self._copy_trials(
            last_trials, r_i, param.resource_name)

        self.logger.info("Generate %d trials by child bracket.",
                         top_trials_num, extra={"experiment_name": experiment.name})
        return trialSpecs

    def _get_last_trials(self, all_trials, latest_trials_num):
        sorted_trials = sorted(
            all_trials, key=lambda trial: trial.status.start_time)
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
            top_trials.extend(
                sorted(latest_trials, key=get_objective_value, reverse=True))
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

    def _make_master_bracket(self, experiment, param):
        n = param.n
        r = int(param.r)
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
                if param['name'] == param.resource_name:
                    param['value'] = str(r)
                trial_spec.parameter_assignments.assignments.add(name=param['name'],
                                                                 value=str(param['value']))
            trial_specs.append(trial_spec)
        self.logger.info("Generate %d trials by master bracket.",
                         n, extra={"experiment_name": experiment.name})
        return trial_specs

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


class HyperBandParam(object):
    def __init__(self, eta=3, s_max=-1, r_l=-1,
                 b_l=-1, r=-1, n=-1, current_s=-2, current_i=-1, resource_name="",
                 evaluating_trials=0)
    self.eta = eta
    self.s_max = s_max
    self.r_l = r_l
    self.b_l = b_l
    self.r = r
    self.n = n
    self.current_s = current_s
    self.current_i = current_i
    self.resource_name = resource_name
    self.evaluating_trials = evaluating_trials

    @staticmethod
    def generate(param):
        algorithm_settings = [
            api_pb2.AlgorithmSetting(
            name="eta",
            value=str(param.eta)
        ), api_pb2.AlgorithmSetting(
            name="s_max",
            value=str(param.s_max)
        ), api_pb2.AlgorithmSetting(
            name="r_l",
            value=str(param.r_l)
        ), api_pb2.AlgorithmSetting(
            name="b_l",
            value=str(param.b_l)
        ), api_pb2.AlgorithmSetting(
            name="r",
            value=str(param.r)
        ), api_pb2.AlgorithmSetting(
            name="n",
            value=str(param.n)
        ), api_pb2.AlgorithmSetting(
            name="current_s",
            value=str(param.current_s)
        ), api_pb2.AlgorithmSetting(
            name="current_i",
            value=str(param.current_i)
        ), api_pb2.AlgorithmSetting(
            name="resource_name",
            value=param.resource_name
        ), api_pb2.AlgorithmSetting(
            name="evaluating_trials",
            value=str(param.evaluating_trials)
        )]
        return api_pb2.AlgorithmSpec(
            algorithm_setting=algorithm_settings
        )

    @staticmethod
    def convert(alg_settings):
        """Convert the algorithm settings to HyperBandParam.
        """
        param = HyperBandParam()
        # Set the param from the algorithm settings.
        for k, v in alg_settings.items():
            if k == "eta":
                param.eta = float(v)
            elif k == "r_l":
                param.r_l = float(v)
            elif k == "b_l":
                param.b_l = float(v)
            elif k == "n"
            param.n = int(float(v))
            elif k == "r"
            param.r = int(float(v))
            elif k == "current_s":
                param.current_s = int(float(v))
            elif k == "current_i":
                param.current_i = int(float(v))
            elif k == "s_max":
                param.s_max = int(float(v))
            elif k == "evaluating_trials":
                param.evaluating_trials = int(float(v))
            elif k == "resource_name":
                param.resource_name = v
            else:
                self.logger.info("Unknown HyperBand Param %s, ignore it", k)
        if param.current_s == -1:
            # Hyperband outlerloop has finished
            self.logger.info("HyperBand outlerloop has finished.")
            return param

        # Deal with illegal parameter values.
        if param.eta <= 0:
            param.eta = 3
        if param.s_max < 0:
            param.s_max = int(
                math.log(param.r_l) / math.log(param.eta))
        if param.b_l < 0:
            param.b_l = (param.s_max + 1) * param.r_l
        if param.current_s < 0:
            param.current_s = param.s_max
        if param.current_i < 0:
            param.current_i = 0
        if param.n < 0:
            param.n = int(math.ceil(float(param.s_max + 1) * (
                float(param.eta**param.current_s) / float(param.current_s+1))))
        if param.r < 0:
            param.r = param.r_l * \
                param.eta**(-param.current_s)

        return param
