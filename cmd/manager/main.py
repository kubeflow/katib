import grpc
from concurrent import futures

import time

from pkg.api.python import api_pb2, api_pb2_grpc
from pkg.db.db_init import connect_db
from pkg.db.internface import create_study, get_study_config, get_suggestion_param, set_suggestion_param, \
    update_suggestion_param, get_trials, update_trial_status, update_trial_value, create_trial, \
    get_suggestion_param_list

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class ManagerService(api_pb2_grpc.ManagerServicer):
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

    def UpdateTrial(self, request, context):
        update_trial_status(self.cnx, request.trial_id, request.status)
        update_trial_value(self.cnx, request.trial_id, request.objective_value)

        ret = api_pb2.UpdateTrialReply(
            trial_id=request.trial_id,
        )

        return ret

    def GetWorkers(self, request, context):
        pass

    def GetEarlyStoppingParameters(self, request, context):
        pass

    def CreateStudy(self, request, context):
        study_id = create_study(self.cnx, request.study_config)
        ret = api_pb2.CreateStudyReply(
            study_id=study_id,
        )

        return ret

    def GetMetrics(self, request, context):
        pass

    def RunTrial(self, request, context):
        pass

    def GetTrial(self, request, context):
        trials = get_trials(self.cnx, request.trial_id, "")
        ret = api_pb2.GetTrialReply(
            trial=trials[0]
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
            suggestion_parameter_set=param_set,
        )

        return ret

    def GetTrials(self, request, context):
        pass

    def SaveStudy(self, request, context):
        pass

    def GetSuggestions(self, request, context):
        channel = grpc.insecure_channel("localhost:50052")
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
    server.add_insecure_port("localhost:50051")
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
