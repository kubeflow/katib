import grpc
import numpy as np

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
from pkg.suggestion.config import MANAGER_ADDRESS, MANAGER_PORT
from pkg.suggestion.bayesianoptimization.src.bayesian_optimization_algorithm import BOAlgorithm
from pkg.suggestion.bayesianoptimization.src.utils import get_logger
from pkg.suggestion.bayesianoptimization.src import parsing_utils


class BayesianService(api_pb2_grpc.SuggestionServicer):

    def __init__(self, logger=None):
        self.manager_addr = MANAGER_ADDRESS
        self.manager_port = MANAGER_PORT
        self.logger = logger if (logger is not None) else get_logger()

    def GetSuggestions(self, request, context):
        service_params = self._parse_suggestion_parameters(request.param_id)
        study_conf = self._get_study_config(request.study_id)
        past_suggestions, past_metrics = self._get_eval_history(request.study_id, study_conf.objective_value_name, service_params["burn_in"])
        name_ids, dim, lower_bounds, upper_bounds, parameter_types, names, discrete_info, categorical_info = \
            parsing_utils.parse_parameter_configs(study_conf.parameter_configs.configs)


        lower_bounds = np.array(lower_bounds)
        upper_bounds = np.array(upper_bounds)
        self.logger.debug("lowerbound: %r", lower_bounds, extra={"StudyID": request.study_id})
        self.logger.debug("upperbound: %r", upper_bounds, extra={"StudyID": request.study_id})
        X_train = parsing_utils.parse_previous_observations(past_suggestions)
        y_train = parsing_utils.parse_metric(past_metrics)
        alg = BOAlgorithm(
            dim=dim,
            N=int(service_params["N"]),
            lowerbound=lower_bounds,
            upperbound=upper_bounds,
            X_train=X_train,
            y_train=y_train,
            mode=service_params["mode"],
            trade_off=service_params["trade_off"],
            # todo: support length_scale with array type
            length_scale=service_params["length_scale"],
            noise=service_params["noise"],
            nu=service_params["nu"],
            kernel_type=service_params["kernel_type"],
            n_estimators=service_params["n_estimators"],
            max_features=service_params["max_features"],
            model_type=service_params["model_type"],
            logger=self.logger,
        )
        trials = []
        x_next_list = alg.get_suggestion(request.request_number)
        for x_next in x_next_list:
            x_next = x_next.squeeze()
            self.logger.debug("xnext: %r ", x_next, extra={"StudyID": request.study_id})
            x_next = parsing_utils.parse_x_next(x_next, parameter_types, names, discrete_info, categorical_info)
            trials.append(api_pb2.Trial(
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name=x["name"],
                        value=str(x["value"]),
                        parameter_type=x["type"],
                    ) for x in x_next
                    ]
            ))
        trials = self._register_trials(trials)
        return api_pb2.GetSuggestionsReply(
            trials=trials
        )

    def _get_study_config(self, studyID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID), 10)
            return gsrep.study_config

    def _get_eval_history(self, studyID, obj_name, burn_in):
        x_train = []
        y_train = []
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=studyID, only_latest_log=True), 10)
            worker_hist = gwfrep.worker_full_infos
        #self.logger.debug("Eval Trials Log: %r", worker_hist, extra={"StudyID": studyID})
        for w in worker_hist:
            if w.Worker.status == api_pb2.COMPLETED:
                for ml in w.metrics_logs:
                    if ml.name == obj_name:
                        y_train.append(float(ml.values[-1].value))
                        x_train.append(w.parameter_set)
                        break
        self.logger.info("%d completed trials are found.", len(x_train), extra={"StudyID": studyID})
        if len(x_train) <= burn_in:
            x_train = []
            y_train = []
            self.logger.info("Trials will be sampled until %d trials for burn-in are completed.", burn_in, extra={"StudyID": studyID})
        else:
            self.logger.debug("Completed trials: %r", x_train, extra={"StudyID": studyID})

        return x_train, y_train

    def _register_trials(self, trials):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
        return trials

    def _parse_suggestion_parameters(self, paramID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
            params = gsprep.suggestion_parameters

        parsed_service_params = {
            "N":            100,
            "model_type":   "gp",
            "max_features": "auto",
            "length_scale": 0.5,
            "noise":        0.0005,
            "nu":           1.5,
            "kernel_type":  "matern",
            "n_estimators": 50,
            "mode":         "pi",
            "trade_off":    0.01,
            "trial_hist":   "",
            "burn_in":      10,
        }
        modes = ["pi", "ei"]
        model_types = ["gp", "rf"]

        for param in params:
            if param.name in parsed_service_params.keys():
                if param.name == "length_scale" or param.name == "noise" or param.name == "nu" or param.name == "trade_off":
                    try:
                        float(param.value)
                    except ValueError:
                        self.logger.warning("Parameter must be float for %s: %s back to default value", param.name, param.value)
                    else:
                        parsed_service_params[param.name] = float(param.value)

                elif param.name == "N" or param.name == "n_estimators" or param.name == "burn_in":
                    try:
                        int(param.value)
                    except ValueError:
                        self.logger.warning("Parameter must be int for %s: %s back to default value", param.name, param.value)
                    else:
                        parsed_service_params[param.name] = int(param.value)

                elif param.name == "kernel_type":
                    if param.value != "rbf" and param.value != "matern":
                        parsed_service_params[param.name] = param.value
                    else:
                        self.logger.warning("Unknown Parameter for %s: %s back to default value", param.name, param.value)
                elif param.name == "mode"  and param.value in modes:
                    if param.value != "lcb" and param.value != "ei" and param.value != "pi":
                        parsed_service_params[param.name] = param.value
                    else:
                        self.logger.warning("Unknown Parameter for %s: %s back to default value", param.name, param.value)
                elif param.name == "model_type" and param.value in model_types:
                    if param.value != "rf" and param.value != "gp":
                        parsed_service_params[param.name] = param.value
                    else:
                        self.logger.warning("Unknown Parameter for %s: %s back to default value", param.name, param.value)
            else:
                self.logger.warning("Unknown Parameter name: %s ", param.name)

        return parsed_service_params
