import tensorflow as tf
import numpy as np
import random
import grpc
from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
import logging
from logging import getLogger, StreamHandler, INFO, DEBUG
import json
import os
import time
from pkg.suggestion.NASenvelopenet.Operation import SearchSpace
from pkg.suggestion.NASenvelopenet.suggestion_param import parseSuggestionParam
from pkg.suggestion.NASenvelopenet.nac_gen import NAC

class EnvelopenetService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        self.manager_addr = "vizier-core"
        self.manager_port = 6789
        self.current_study_id = ""
        self.current_trial_id = ""
        self.ctrl_cache_file = ""
        self.is_first_run = True
        self.current_itr=0

    def generate_arch(self, request):
        self.current_study_id = request.study_id
        self.current_trial_id = ""
        self._get_suggestion_param(request.param_id)
        self._get_search_space(request.study_id)
        self.restruct_itr=self.suggestion_config["iterations"]
        self.generator = NAC(
            self.search_space,
            algorithm = self.suggestion_config["algorithm"],
            stages=self.stages)
    
    
    def GetSuggestions(self, request, context):
        if request.study_id != self.current_study_id:
            self.generate_arch(request)
        
        if self.current_itr==0:
            self.arch=self.generator.get_init_arch()
            self.is_first_run = False
        elif self.current_itr<=self.restruct_itr:
            result = self.GetEvaluationResult(request.study_id, self.prev_trial_id)
            self.arch=self.generator.get_arch(self.arch, result)
 
        print("Architecture at itr=", self.current_itr)
        print(self.arch)
        arch_json=json.dumps(self.arch)
        config_json=json.dumps(self.suggestion_config)
        arch=str(arch_json).replace('\"', '\'')
        config=str(config_json).replace('\"', '\'')    
        
        trials = []
        trials.append(api_pb2.Trial(
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name="architecture",
                        value=arch,
                        parameter_type= api_pb2.CATEGORICAL),
                    api_pb2.Parameter(
                        name="parameters",
                        value=config,
                        parameter_type= api_pb2.CATEGORICAL),
                    api_pb2.Parameter(
                        name="current_itr",
                        value=str(self.current_itr),
                        parameter_type= api_pb2.CATEGORICAL)
                ], 
            )
        )

        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            print("INSIDE CLIENT")
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
            print("TRIALS INSERTED")
            print(ctrep.trial_id)
            self.prev_trial_id = ctrep.trial_id

        self.current_itr+=1

        return api_pb2.GetSuggestionsReply(trials=trials)

    def GetEvaluationResult(self, studyID, trialID):
        worker_list = []
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            print(">>> Calling GetEvaluationResult(). ")
            print(">>> The current study id is:", studyID)
            print(">>> The current trial id is:", trialID)
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=studyID, trial_id=trialID, only_latest_log=False), 10)
            worker_list = gwfrep.worker_full_infos
            print(worker_list)
        for w in worker_list:
            print(w)
            if w.Worker.status == api_pb2.COMPLETED:
                print(w.metrics_logs)
                for ml in w.metrics_logs:
                    print(ml.name)
                    if ml.name == self.objective_name:
                        samples=self.get_msss(ml)   
                        print(samples)
                        return samples

    
    def _get_search_space(self, studyID):

        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID), 10)

        self.opt_direction = gsrep.study_config.optimization_type
        self.objective_name = gsrep.study_config.objective_value_name
        all_params = gsrep.study_config.nas_config
        graph_config = all_params.graph_config
        search_space_raw = all_params.operations

        self.stages = int(graph_config.num_layers)
        self.input_size = list(map(int, graph_config.input_size))
        self.output_size = list(map(int, graph_config.output_size))
        search_space_object = SearchSpace(search_space_raw)
        self.search_space = search_space_object.search_space
        self.search_space.update({"input_size":self.input_size[0]})
        self.search_space.update({"output_size":self.output_size[0]})
        print(self.search_space)
        
    def _get_suggestion_param(self, paramID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
        
            params_raw = gsprep.suggestion_parameters
            suggestion_params = parseSuggestionParam(params_raw)
            self.suggestion_config = suggestion_params

    def get_msss(self, metric_object):
        samples=list()
        i=0
        for i in range(len(metric_object.values)):
            samples.append(metric_object.values[i].value)
        return samples
