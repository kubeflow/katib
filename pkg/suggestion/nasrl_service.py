from pkg.suggestion.NAS_Reinforcement_Learning.Controller import Controller
from pkg.suggestion.NAS_Reinforcement_Learning.Operation import SearchSpace
from pkg.suggestion.NAS_Reinforcement_Learning.SuggestionParam import parseSuggestionParam
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


class NasrlService(api_pb2_grpc.SuggestionServicer):
    def __init__(self):
        self.manager_addr = "vizier-core"
        self.manager_port = 6789
        self.current_study_id = ""
        self.current_trial_id = ""
        self.ctrl_cache_file = ""

        self.tf_graph = tf.get_default_graph()
        self.is_first_run = True

    def reset_controller(self, request):
        print("Resetting Controller for StudyJob {}.".format(request.study_id))
        self.current_study_id = request.study_id
        self.ctrl_cache_file = "ctrl_cache/{}.ckpt".format(self.current_study_id)
        if not os.path.exists("ctrl_cache/"):
            os.makedirs("ctrl_cache/")
        self.current_trial_id = ""
        self._get_suggestion_param(request.param_id)
        self._get_search_space(request.study_id)

        with self.tf_graph.as_default():
            self.controller = Controller(
                num_layers=self.num_layers,
                num_operations=self.num_operations,
                lstm_size=self.suggestion_config['lstm_num_cells'],
                lstm_num_layers=self.suggestion_config['lstm_num_layers'],
                lstm_keep_prob=self.suggestion_config['lstm_keep_prob'],
                lr_init=self.suggestion_config['init_learning_rate'],
                lr_dec_start=self.suggestion_config['lr_decay_start'],
                lr_dec_every=self.suggestion_config['lr_decay_every'],
                lr_dec_rate=self.suggestion_config['lr_decay_rate'],
                l2_reg=self.suggestion_config['l2_reg'],
                entropy_weight=self.suggestion_config['entropy_weight'],
                bl_dec=self.suggestion_config['baseline_decay'],
                optim_algo=self.suggestion_config['optimizer'],
                skip_target=self.suggestion_config['skip-target'],
                skip_weight=self.suggestion_config['skip-weight'],
                name="Ctrl_"+self.current_study_id)

            self.controller.build_trainer()

        print("Controller for StudyJob {} has been initialized.".format(request.study_id))

    def GetSuggestions(self, request, context):
        if request.study_id != self.current_study_id:
            self.reset_controller(request)
            self.is_first_run = True

        with self.tf_graph.as_default():

            saver = tf.train.Saver()

            controller_ops = {
                  "train_step": self.controller.train_step,
                  "loss": self.controller.loss,
                  "train_op": self.controller.train_op,
                  "lr": self.controller.lr,
                  "grad_norm": self.controller.grad_norm,
                  "optimizer": self.controller.optimizer,
                  "baseline": self.controller.baseline,
                  "entropy": self.controller.sample_entropy,
                  "sample_arc": self.controller.sample_arc,
                  "skip_rate": self.controller.skip_rate}

            run_ops = [
                controller_ops["loss"],
                controller_ops["entropy"],
                controller_ops["lr"],
                controller_ops["grad_norm"],
                controller_ops["baseline"],
                controller_ops["skip_rate"],
                controller_ops["train_op"]]

            if self.is_first_run:
                with tf.Session() as sess:
                    sess.run(tf.global_variables_initializer())
                    arc = sess.run(controller_ops["sample_arc"])
                    saver.save(sess, self.ctrl_cache_file)

                self.is_first_run = False

            else:
                with tf.Session() as sess:
                    saver.restore(sess, self.ctrl_cache_file)

                    valid_acc = self.controller.reward
                    result = self.GetEvaluationResult(request.study_id)
                    loss, entropy, lr, gn, bl, skip, _ = sess.run(
                        fetches=run_ops,
                        feed_dict={valid_acc: result})
                    arc = sess.run(controller_ops["sample_arc"])

                    saver.save(sess, self.ctrl_cache_file)

        arc = arc.tolist()
        organized_arc = [0 for _ in range(self.num_layers)]
        record = 0
        for l in range(self.num_layers):
            organized_arc[l] = arc[record: record + l + 1]
            record += l + 1

        nn_config = dict()
        nn_config['num_layers'] = self.num_layers
        nn_config['input_size'] = self.input_size
        nn_config['output_size'] = self.output_size
        nn_config['embedding'] = dict()
        for l in range(self.num_layers):
            opt = organized_arc[l][0]
            nn_config['embedding'][opt] = self.search_space[opt].get_dict()

        organized_arc_json = json.dumps(organized_arc)
        nn_config_json = json.dumps(nn_config)

        organized_arc_str = str(organized_arc_json).replace('\"', '\'')
        nn_config_str = str(nn_config_json).replace('\"', '\'')

        print("Nerual Network Architecture json:")
        print(organized_arc_json)
        print("Seach Space Embedding json:")
        print(nn_config_json)

        trials = []
        trials.append(api_pb2.Trial(
                study_id=request.study_id,
                parameter_set=[
                    api_pb2.Parameter(
                        name="architecture",
                        value=organized_arc_str,
                        parameter_type= api_pb2.CATEGORICAL),
                    api_pb2.Parameter(
                        name="nn_config",
                        value=nn_config_str,
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
                self.current_trial_id = ctrep.trial_id
            print("TRIALS INSERTED")
            print(ctrep.trial_id)
            self.current_trial_id = ctrep.trial_id
        
        return api_pb2.GetSuggestionsReply(trials=trials)

    def GetEvaluationResult(self, studyID):
        worker_hist = []
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            print(">>> Calling GetEvaluationResult(). ")
            print(">>> The current study id is:", studyID)
            print(">>> The current trial id is:", self.current_trial_id)
            time.sleep(20)
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=studyID, trial_id=self.current_trial_id, only_latest_log=True), 10)
            worker_hist = gwfrep.worker_full_infos

        for w in worker_hist:
            if w.Worker.status == api_pb2.COMPLETED:
                for ml in w.metrics_logs:
                    print("Evaluation Result")
                    print(ml.values[-1].value)
                    print()

        return float(ml.values[-1].value)

    def _get_search_space(self, studyID):

        # this function need to
        # 1) get the number of layers
        # 2) get the I/O size
        # 2) get the available operations

        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID, job_type="NAS"), 10)

        all_params = gsrep.study_config.nas_config
        graph_config = all_params.graph_config
        search_space_raw = all_params.operations

        self.num_layers = int(graph_config.num_layers)
        self.input_size = list(map(int, graph_config.input_size))
        self.output_size = list(map(int, graph_config.output_size))
        search_space_object = SearchSpace(search_space_raw)


        print("\n==================== Search Space ====================")
        self.num_operations = search_space_object.num_operations
        print("There are", self.num_operations, "operations in total")
        self.search_space = search_space_object.search_space
        for opt in self.search_space:
            opt.print_op()    
            print()
            

    def _get_suggestion_param(self, paramID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
        
        params_raw = gsprep.suggestion_parameters

        suggestion_params = parseSuggestionParam(params_raw)

        print("\n=============== Parameters for LSTM Controller ===============")
        for spec in suggestion_params:
            print(spec, suggestion_params[spec])

        self.suggestion_config = suggestion_params
