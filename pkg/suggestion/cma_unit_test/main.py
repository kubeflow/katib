import grpc
from concurrent import futures

import time

from pkg.api.python import api_pb2, api_pb2_grpc
from pkg.suggestion.cma_unit_test.db_init import connect_db
from pkg.suggestion.cma_unit_test.internface import create_study, get_study_config, get_suggestion_param, set_suggestion_param, \
    update_suggestion_param, create_trial, \
    get_suggestion_param_list, get_worker_list, get_worker, get_worker_logs, create_worker
from pkg.suggestion.cma_unit_test.worker import spawn_worker

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class ManagerService(api_pb2_grpc.ManagerServicer):
    def GetEarlyStoppingParameterList(self, request, context):
        pass

    def SaveStudy(self, request, context):
        pass

    def __init__(self):
        self.cnx = connect_db()

    def CreateTrial(self, request, context):
        trial_id = create_trial(self.cnx, request.trial)
        ret = api_pb2.CreateTrialReply(
            trial_id=trial_id,
        )
        return ret

    def SetSuggestionParameters(self, request, context):
        if request.param_id == "":
            id = set_suggestion_param(self.cnx, request.suggestion_algorithm, request.study_id,
                                      request.suggestion_parameters)
        else:
            id = request.param_id
            update_suggestion_param(self.cnx, request.param_id, request.suggestion_parameters)
        ret = api_pb2.SetSuggestionParametersReply(
            param_id=id,
        )
        return ret

    def GetWorkers(self, request, context):
        if not request.worker_id:
            workers = get_worker_list(self.cnx, request.study_id, request.trial_id)
        else:
            worker = get_worker(self.cnx, request.worker_id)
            workers = [worker]
        ret = api_pb2.GetWorkersReply(
            workers=workers,
        )
        return ret

    def GetEarlyStoppingParameters(self, request, context):
        pass

    def CreateStudy(self, request, context):
        study_id = create_study(self.cnx, request.study_config)
        ret = api_pb2.CreateStudyReply(
            study_id=study_id,
        )

        return ret

    def GetMetrics(self, request, context):
        if len(request.metrics_names) > 0:
            m_names = request.metrics_names
        else:
            study_config = get_study_config(self.cnx, request.study_id)
            m_names = study_config.metrics

        metric_log_sets = []
        for worker_id in request.worker_ids:
            metric_logs = []
            for name in m_names:
                worker_logs = get_worker_logs(self.cnx, worker_id, name)
                values = []
                for log in worker_logs:
                    values.append(log["value"])
                metric_logs.append(api_pb2.MetricsLog(
                    name=name,
                    values=values
                ))

            metric_log_sets.append(api_pb2.MetricsLogSet(
                worker_id=worker_id,
                metrics_logs=metric_logs,
            ))

        ret = api_pb2.GetMetricsReply(
            metrics_log_sets=metric_log_sets
        )
        return ret

    def RunTrial(self, request, context):
        worker_id = create_worker(self.cnx, api_pb2.Worker(
            study_id=request.study_id,
            trial_id=request.trial_id,
            runtime=request.runtime,
            config=request.worker_config,
        ))
        spawn_worker(worker_id, request.worker_config)

        ret = api_pb2.RunTrialReply(
            worker_id=worker_id
        )
        return ret

    def SaveModel(self, request, context):
        pass

    def GetSuggestionParameters(self, request, context):
        params = get_suggestion_param(self.cnx, request.param_id)
        ret = api_pb2.GetSuggestionParametersReply(
            suggestion_parameters=params,
        )

        return ret

    def GetSuggestionParameterList(self, request, context):
        param_set = get_suggestion_param_list(self.cnx, request.study_id)
        ret = api_pb2.GetSuggestionParameterListReply(
            suggestion_parameter_sets=param_set,
        )

        return ret

    def GetTrials(self, request, context):
        pass

    def GetSuggestions(self, request, context):
        channel = grpc.insecure_channel("0.0.0.0:6789")
        stub = api_pb2_grpc.SuggestionStub(channel)
        ret = stub.GetSuggestions(request)
        # for trial in ret.trials:
        #     create_trial(self.cnx, trial)
        return ret

    def GetShouldStopWorkers(self, request, context):
        pass

    def GetSavedModels(self, request, context):
        pass

    def SetEarlyStoppingParameters(self, request, context):
        pass

    def StopWorkers(self, request, context):
        pass

    def GetStudy(self, request, context):
        study_config = get_study_config(self.cnx, request.study_id)
        ret = api_pb2.GetStudyReply(
            study_config=study_config,
        )

        return ret

    def GetSavedStudies(self, request, context):
        pass

    def GetStudyList(self, request, context):
        pass

    def StopStudy(self, request, context):
        pass


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    api_pb2_grpc.add_ManagerServicer_to_server(ManagerService(), server)
    server.add_insecure_port("0.0.0.0:6788")
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()
