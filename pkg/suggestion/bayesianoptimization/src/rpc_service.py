import grpc

from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
from .config import MANAGER_ADDRESS, MANAGER_PORT, SEARCH_ALGORITHM
from .algorithm_register import ALGORITHM_REGISTER
from .utils import get_logger
from . import parsing_utils


class SuggestionService(api_pb2_grpc.SuggestionServicer):

    def __init__(self, logger=None):
        self.manager_addr = MANAGER_ADDRESS
        self.manager_port = MANAGER_PORT
        self.logger = logger if (logger is not None) else get_logger()

    def GetSuggestions(self, request, context):
        suggestion_config = self._parse_suggestion_parameters(request.param_id)
        suggestion_config = {param.name: param.value for param in suggestion_config}
        study_conf = self._get_study_config(request.study_id)
        past_suggestions, past_metrics = self._get_eval_history(request.study_id, study_conf.objective_value_name)
        parameter_config = parsing_utils.parse_parameter_configs(study_conf.parameter_configs.configs)

        self.logger.debug("lowerbound: %r", parameter_config.lower_bounds, extra={"StudyID": request.study_id})
        self.logger.debug("upperbound: %r", parameter_config.upper_bounds, extra={"StudyID": request.study_id})
        X_train = parsing_utils.parse_previous_observations(
            past_suggestions,
            parameter_config.dim,
            parameter_config.name_ids,
            parameter_config.parameter_types,
            parameter_config.categorical_info
        )
        y_train = parsing_utils.parse_metric(past_metrics,
                                             study_conf.optimization_type)
        alg = ALGORITHM_REGISTER[SEARCH_ALGORITHM](parameter_config, suggestion_config, X_train, y_train, logger=self.logger)
        trials = []
        x_next_list = alg.get_suggestion(request.request_number)
        for x_next in x_next_list:
            x_next = x_next.squeeze()
            self.logger.debug("xnext: %r ", x_next, extra={"StudyID": request.study_id})
            x_next = parsing_utils.parse_x_next(x_next,
                                                parameter_config.parameter_types,
                                                parameter_config.names,
                                                parameter_config.discrete_info,
                                                parameter_config.categorical_info)
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

    def _get_study_config(self, studyID): # pragma: no cover
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID), 10)
            return gsrep.study_config

    def _get_eval_history(self, studyID, obj_name): # pragma: no cover
        x_train = []
        y_train = []
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=studyID, only_latest_log=True), 10)
            worker_hist = gwfrep.worker_full_infos
        self.logger.info("Eval Trials Log: %r", worker_hist, extra={"StudyID": studyID})
        for w in worker_hist:
            if w.Worker.status == api_pb2.COMPLETED:
                for ml in w.metrics_logs:
                    if ml.name == obj_name:
                        y_train.append(float(ml.values[-1].value))
                        x_train.append(w.parameter_set)
                        break
        self.logger.info("%d completed trials are found.", len(x_train), extra={"StudyID": studyID})

        return x_train, y_train

    def _register_trials(self, trials): # pragma: no cover
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
        return trials

    def _parse_suggestion_parameters(self, paramID): # pragma: no cover
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
            params = gsprep.suggestion_parameters
        return params
