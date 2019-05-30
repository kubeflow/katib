import random
import string

import grpc
import numpy as np

from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc
from pkg.suggestion.v1alpha2.bayesianoptimization.src.bayesian_optimization_algorithm import BOAlgorithm
from pkg.suggestion.v1alpha2.bayesianoptimization.src.algorithm_manager import AlgorithmManager
import logging
from logging import getLogger, StreamHandler, INFO, DEBUG


class BayesianService(api_pb2_grpc.SuggestionServicer):
    def __init__(self, logger=None):
        self.manager_addr = "katib-manager"
        self.manager_port = 6789
        if logger == None:
            self.logger = getLogger(__name__)
            FORMAT = '%(asctime)-15s Experiment %(experiment_name)s %(message)s'
            logging.basicConfig(format=FORMAT)
            handler = StreamHandler()
            handler.setLevel(INFO)
            self.logger.setLevel(INFO)
            self.logger.addHandler(handler)
            self.logger.propagate = False
        else:
            self.logger = logger

    def _get_experiment(self, name):
        channel = grpc.beta.implementations.insecure_channel(
            self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            exp = client.GetExperiment(
                api_pb2.GetExperimentRequest(experiment_name=name), 10)
            return exp.experiment

    def ValidateAlgorithmSettings(self, request, context):
        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        """
        Main function to provide suggestion.
        """
        service_params = self.parseParameters(request.experiment_name)
        experiment = self._get_experiment(request.experiment_name)
        X_train, y_train = self.getEvalHistory(
            request.experiment_name, experiment.spec.objective.objective_metric_name, service_params["burn_in"])

        algo_manager = AlgorithmManager(
            experiment_name=request.experiment_name,
            experiment=experiment,
            X_train=X_train,
            y_train=y_train,
            logger=self.logger,
        )

        lowerbound = np.array(algo_manager.lower_bound)
        upperbound = np.array(algo_manager.upper_bound)
        self.logger.debug("lowerbound: %r", lowerbound,
                          extra={"experiment_name": request.experiment_name})
        self.logger.debug("upperbound: %r", upperbound,
                          extra={"experiment_name": request.experiment_name})
        alg = BOAlgorithm(
            dim=algo_manager.dim,
            N=int(service_params["N"]),
            lowerbound=lowerbound,
            upperbound=upperbound,
            X_train=algo_manager.X_train,
            y_train=algo_manager.y_train,
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
            self.logger.debug("xnext: %r ", x_next, extra={
                              "experiment_name": request.experiment_name})
            x_next = algo_manager.parse_x_next(x_next)
            x_next = algo_manager.convert_to_dict(x_next)
            trials.append(api_pb2.Trial(
                spec=api_pb2.TrialSpec(
                    experiment_name=request.experiment_name,
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=[
                            api_pb2.ParameterAssignment(
                                name=x["name"],
                                value=str(x["value"]),
                            ) for x in x_next
                        ]
                    )
                )
            ))
        return api_pb2.GetSuggestionsReply(
            trials=trials
        )

    def getEvalHistory(self, experiment_name, obj_name, burn_in):
        worker_hist = []
        x_train = []
        y_train = []
        channel = grpc.beta.implementations.insecure_channel(
            self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            trialsrep = client.GetTrialList(api_pb2.GetTrialListRequest(
                experiment_name=experiment_name
            ))
            for t in trialsrep.trials:
                if t.status.condition == 2:
                    gwfrep = client.GetObservationLog(
                        api_pb2.GetObservationLogRequest(
                            trial_name=t.name,
                            metric_name=obj_name))
                    w = gwfrep.observation_log
                    for ml in w.metrics_logs:
                        if ml.name == obj_name:
                            y_train.append(float(ml.values[-1].value))
                            x_train.append(w.parameter_set)
                            break
        self.logger.info("%d completed trials are found.",
                         len(x_train), extra={"Experiment": experiment_name})
        if len(x_train) <= burn_in:
            x_train = []
            y_train = []
            self.logger.info("Trials will be sampled until %d trials for burn-in are completed.",
                             burn_in, extra={"experiment_name": experiment_name})
        else:
            self.logger.debug("Completed trials: %r", x_train,
                              extra={"experiment_name": experiment_name})

        return x_train, y_train

    def parseParameters(self, experiment_name):
        channel = grpc.beta.implementations.insecure_channel(
            self.manager_addr, self.manager_port)
        params = []
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetAlgorithmExtraSettings(
                api_pb2.GetAlgorithmExtraSettingsRequest(param_id=experiment_name), 10)
            params = gsprep.extra_algorithm_settings

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
        kernel_types = ["matern", "rbf"]

        for param in params:
            if param.name in parsed_service_params.keys():
                if param.name == "length_scale" or param.name == "noise" or param.name == "nu" or param.name == "trade_off":
                    try:
                        float(param.value)
                    except ValueError:
                        self.logger.warning(
                            "Parameter must be float for %s: %s back to default value", param.name, param.value)
                    else:
                        parsed_service_params[param.name] = float(param.value)

                elif param.name == "N" or param.name == "n_estimators" or param.name == "burn_in":
                    try:
                        int(param.value)
                    except ValueError:
                        self.logger.warning(
                            "Parameter must be int for %s: %s back to default value", param.name, param.value)
                    else:
                        parsed_service_params[param.name] = int(param.value)

                elif param.name == "kernel_type":
                    if param.value != "rbf" and param.value != "matern":
                        parsed_service_params[param.name] = param.value
                    else:
                        self.logger.warning(
                            "Unknown Parameter for %s: %s back to default value", param.name, param.value)
                elif param.name == "mode" and param.value in modes:
                    if param.value != "lcb" and param.value != "ei" and param.value != "pi":
                        parsed_service_params[param.name] = param.value
                    else:
                        self.logger.warning(
                            "Unknown Parameter for %s: %s back to default value", param.name, param.value)
                elif param.name == "model_type" and param.value in model_types:
                    if param.value != "rf" and param.value != "gp":
                        parsed_service_params[param.name] = param.value
                    else:
                        self.logger.warning(
                            "Unknown Parameter for %s: %s back to default value", param.name, param.value)
            else:
                self.logger.warning("Unknown Parameter name: %s ", param.name)

        return parsed_service_params
