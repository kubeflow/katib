from pkg.suggestion.NAS_Reinforcement_Learning.Controller import Controller
from pkg.suggestion.NAS_Reinforcement_Learning.Operation import SearchSpace
from pkg.suggestion.NAS_Reinforcement_Learning.SuggestionParam import parseSuggestionParam
import tensorflow as tf
import grpc
from pkg.api.python import api_pb2
from pkg.api.python import api_pb2_grpc
import logging
from logging import getLogger, StreamHandler, INFO, DEBUG
import json
import os
import time


class NasrlService(api_pb2_grpc.SuggestionServicer):
    def __init__(self, logger=None):
        self.manager_addr = "vizier-core"
        self.manager_port = 6789
        self.registered_studies = list()

        self.ctrl_cache_file = ""
        self.ctrl_step = 0
        self.is_first_run = True

        if not os.path.exists("ctrl_cache/"):
            os.makedirs("ctrl_cache/")

        if logger == None:
            self.logger = getLogger(__name__)
            FORMAT = '%(asctime)-15s StudyID %(studyid)s %(message)s'
            logging.basicConfig(format=FORMAT)
            handler = StreamHandler()
            handler.setLevel(INFO)
            self.logger.setLevel(INFO)
            self.logger.addHandler(handler)
            self.logger.propagate = False
        else:
            self.logger = logger

    def setup_controller(self, request):
        self.logger.info("-" * 80 + "\nSetting Up Suggestion for StudyJob {}\n".format(request.study_id) + "-" * 80)
        self.tf_graph = tf.Graph()
        self.ctrl_step = 0
        self.ctrl_cache_file = "ctrl_cache/{}.ckpt".format(request.study_id)
        self._get_suggestion_param(request.param_id, request.study_id)
        self._get_search_space(request.study_id)

        with self.tf_graph.as_default():
            ctrl_param = self.suggestion_config
            self.controllers = Controller(
                num_layers=self.num_layers,
                num_operations=self.num_operations,
                lstm_size=ctrl_param['lstm_num_cells'],
                lstm_num_layers=ctrl_param['lstm_num_layers'],
                lstm_keep_prob=ctrl_param['lstm_keep_prob'],
                lr_init=ctrl_param['init_learning_rate'],
                lr_dec_start=ctrl_param['lr_decay_start'],
                lr_dec_every=ctrl_param['lr_decay_every'],
                lr_dec_rate=ctrl_param['lr_decay_rate'],
                l2_reg=ctrl_param['l2_reg'],
                entropy_weight=ctrl_param['entropy_weight'],
                bl_dec=ctrl_param['baseline_decay'],
                optim_algo=ctrl_param['optimizer'],
                skip_target=ctrl_param['skip-target'],
                skip_weight=ctrl_param['skip-weight'],
                name="Ctrl_"+request.study_id)

            self.controllers.build_trainer()

        self.logger.info("Suggestion for StudyJob {} has been initialized.".format(request.study_id))

    def GetSuggestions(self, request, context):
        if request.study_id not in self.registered_studies:
            self.setup_controller(request)
            self.is_first_run = True
            self.registered_studies.append(request.study_id)

        self.logger.info("-" * 80 + "\nSuggestion Step {} for Study {}\n".format(self.ctrl_step, request.study_id) + "-" * 80)

        with self.tf_graph.as_default():

            saver = tf.train.Saver()
            ctrl = self.controllers

            controller_ops = {
                  "train_step": ctrl.train_step,
                  "loss": ctrl.loss,
                  "train_op": ctrl.train_op,
                  "lr": ctrl.lr,
                  "grad_norm": ctrl.grad_norm,
                  "optimizer": ctrl.optimizer,
                  "baseline": ctrl.baseline,
                  "entropy": ctrl.sample_entropy,
                  "sample_arc": ctrl.sample_arc,
                  "skip_rate": ctrl.skip_rate}

            run_ops = [
                controller_ops["loss"],
                controller_ops["entropy"],
                controller_ops["lr"],
                controller_ops["grad_norm"],
                controller_ops["baseline"],
                controller_ops["skip_rate"],
                controller_ops["train_op"]]

            if self.is_first_run:
                self.logger.info("First time running suggestion for {}. Random architecture will be given.".format(request.study_id))
                with tf.Session() as sess:
                    sess.run(tf.global_variables_initializer())
                    arc = sess.run(controller_ops["sample_arc"])
                    # TODO: will use PVC to store the checkpoint to protect against unexpected suggestion pod restart
                    saver.save(sess, self.ctrl_cache_file)

                self.is_first_run = False

            else:
                with tf.Session() as sess:
                    saver.restore(sess, self.ctrl_cache_file)

                    valid_acc = ctrl.reward
                    result = self.GetEvaluationResult(request.study_id)

                    # handle the special case where there is no completed trials (though not very likely)
                    while result is None:
                        time.sleep(20)
                        result = self.GetEvaluationResult(request.study_id)

                    # This lstm cell is designed to maximize the metrics
                    # However, if the user want to minimize the metrics, we can take the negative of the result
                    if self.opt_direction == api_pb2.MINIMIZE:
                        result = -result

                    loss, entropy, lr, gn, bl, skip, _ = sess.run(
                        fetches=run_ops,
                        feed_dict={valid_acc: result})
                    self.logger.info("Suggetion updated. LSTM Controller Loss: {}".format(loss))
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

        self.logger.info("\nNew Neural Network Architecture (internal representation):")
        self.logger.info(organized_arc_json)
        self.logger.info("\nCorresponding Seach Space Description:")
        self.logger.info(nn_config_str)
        self.logger.info("")

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
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
            self.logger.info("Trial {} Created\n".format(ctrep.trial_id))
            self.prev_trial_id = ctrep.trial_id
        
        self.ctrl_step += 1
        return api_pb2.GetSuggestionsReply(trials=trials)

    def GetEvaluationResult(self, studyID):
        worker_list = []
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=studyID, trial_id=self.prev_trial_id, only_latest_log=True), 10)
            worker_list = gwfrep.worker_full_infos

        completed_count = 0
        for w in worker_list:
            if w.Worker.status == api_pb2.COMPLETED:
                for ml in w.metrics_logs:
                    if ml.name != self.objective_name:
                        continue
                    completed_count += 1
                    self.logger.info("Evaluation result of previous candidate: {}".format(ml.values[-1].value))

        # TODO: add support for multiple trials
        if completed_count >= 1:
            return float(ml.values[-1].value)
        else:
            self.logger.warning("Error. No trials have been completed.")
            return None


    def _get_search_space(self, studyID):

        # this function need to
        # 1) get the number of layers
        # 2) get the I/O size
        # 3) get the available operations
        # 4) get the optimization direction (i.e. minimize or maximize)
        # 5) get the objective name
        
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsrep = client.GetStudy(api_pb2.GetStudyRequest(study_id=studyID), 10)
        
        self.opt_direction = gsrep.study_config.optimization_type
        self.objective_name = gsrep.study_config.objective_value_name

        all_params = gsrep.study_config.nas_config
        graph_config = all_params.graph_config
        search_space_raw = all_params.operations

        self.num_layers = int(graph_config.num_layers)
        self.input_size = list(map(int, graph_config.input_size))
        self.output_size = list(map(int, graph_config.output_size))
        search_space_object = SearchSpace(search_space_raw)

        self.logger.info("Search Space for Study {}:".format(studyID))

        self.search_space = search_space_object.search_space
        for opt in self.search_space:
            opt.print_op(self.logger)
        
        self.num_operations = search_space_object.num_operations
        self.logger.info("There are {} operations in total.\n".format(self.num_operations))
            

    def _get_suggestion_param(self, paramID, studyID):
        channel = grpc.beta.implementations.insecure_channel(self.manager_addr, self.manager_port)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gsprep = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=paramID), 10)
        
        params_raw = gsprep.suggestion_parameters

        suggestion_params = parseSuggestionParam(params_raw)

        self.logger.info("Parameters of LSTM Controller for Study {}:".format(studyID))
        for spec in suggestion_params:
            if len(spec) > 13:
                self.logger.info("{}: \t{}".format(spec, suggestion_params[spec]))
            else:
                self.logger.info("{}: \t\t{}".format(spec, suggestion_params[spec]))

        self.suggestion_config = suggestion_params
