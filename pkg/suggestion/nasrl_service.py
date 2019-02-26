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

MANAGER_ADDRESS = "vizier-core"
MANAGER_PORT = 6789

class NAS_RL_StudyJob(object):
    def __init__(self, request, logger):
        self.logger = logger
        self.study_id = request.study_id
        self.param_id = request.param_id
        self.study_name = None
        self.tf_graph = tf.Graph()
        self.prev_trial_id = None
        self.ctrl_cache_file = "ctrl_cache/{}/{}.ckpt".format(request.study_id, request.study_id)
        self.ctrl_step = 0
        self.is_first_run = True
        self.suggestion_config = None
        self.controller = None
        self.num_layers  = None
        self.input_size = None
        self.output_size = None
        self.num_operations = None
        self.search_space = None
        self.opt_direction = None
        self.objective_name = None
        
        self.logger.info("-" * 100 + "\nSetting Up Suggestion for StudyJob ID {}\n".format(request.study_id) + "-" * 100)
        self._get_study_param()
        self._get_suggestion_param()
        self._setup_controller()
        self.logger.info("Suggestion for StudyJob {} (ID: {}) has been initialized.\n".format(self.study_name, self.study_id))
        
    def _get_study_param(self):
        # this function need to
        # 1) get the number of layers
        # 2) get the I/O size
        # 3) get the available operations
        # 4) get the optimization direction (i.e. minimize or maximize)
        # 5) get the objective name
        # 6) get the study name

        channel = grpc.beta.implementations.insecure_channel(MANAGER_ADDRESS, MANAGER_PORT)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            api_study_param = client.GetStudy(api_pb2.GetStudyRequest(study_id=self.study_id), 10)

        self.study_name = api_study_param.study_config.name
        self.opt_direction = api_study_param.study_config.optimization_type
        self.objective_name = api_study_param.study_config.objective_value_name

        all_params = api_study_param.study_config.nas_config
        
        graph_config = all_params.graph_config
        self.num_layers = int(graph_config.num_layers)
        self.input_size = list(map(int, graph_config.input_size))
        self.output_size = list(map(int, graph_config.output_size))

        search_space_raw = all_params.operations
        search_space_object = SearchSpace(search_space_raw)
        self.search_space = search_space_object.search_space
        self.num_operations = search_space_object.num_operations
        
        self.print_search_space()
    
    def _get_suggestion_param(self):
        channel = grpc.beta.implementations.insecure_channel(MANAGER_ADDRESS, MANAGER_PORT)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            api_suggestion_param = client.GetSuggestionParameters(api_pb2.GetSuggestionParametersRequest(param_id=self.param_id), 10)
        
        params_raw = api_suggestion_param.suggestion_parameters
        self.suggestion_config = parseSuggestionParam(params_raw)

        self.print_suggestion_params()
    
    def _setup_controller(self):
        
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
                name="Ctrl_" + self.study_id,
                logger=self.logger)

            self.controller.build_trainer()

    def print_search_space(self):
        if self.search_space is None:
            self.logger.warning("Error! The Suggestion has not yet been initialized!")
            return
        
        self.logger.info("Search Space for StudyJob {} (ID: {}):".format(self.study_name, self.study_id))
        for opt in self.search_space:
            opt.print_op(self.logger)
        self.logger.info("There are {} operations in total.\n".format(self.num_operations))
    
    def print_suggestion_params(self):
        if self.suggestion_config is None:
            self.logger.warning("Error! The Suggestion has not yet been initialized!")
            return
        
        self.logger.info("Parameters of LSTM Controller for StudyJob {} (ID: {}):".format(self.study_name, self.study_id))
        for spec in self.suggestion_config:
            if len(spec) > 13:
                self.logger.info("{}: \t{}".format(spec, self.suggestion_config[spec]))
            else:
                self.logger.info("{}: \t\t{}".format(spec, self.suggestion_config[spec]))
        self.logger.info("")


class NasrlService(api_pb2_grpc.SuggestionServicer):
    def __init__(self, logger=None):

        self.registered_studies = dict()

        if logger == None:
            self.logger = getLogger(__name__)
            FORMAT = '%(asctime)-15s studyID %(study_id)s %(message)s'
            logging.basicConfig(format=FORMAT)
            handler = StreamHandler()
            handler.setLevel(INFO)
            self.logger.setLevel(INFO)
            self.logger.addHandler(handler)
            self.logger.propagate = False
        else:
            self.logger = logger

        if not os.path.exists("ctrl_cache/"):
            os.makedirs("ctrl_cache/")  

    def GetSuggestions(self, request, context):
        if request.study_id not in self.registered_studies:
            self.registered_studies[request.study_id] = NAS_RL_StudyJob(request, self.logger)
        
        study = self.registered_studies[request.study_id]

        self.logger.info("-" * 100 + "\nSuggestion Step {} for StudyJob {} (ID: {})\n".format(study.ctrl_step, study.study_name, study.study_id) + "-" * 100)

        with study.tf_graph.as_default():

            saver = tf.train.Saver()
            ctrl = study.controller

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

            if study.is_first_run:
                self.logger.info("First time running suggestion for {}. Random architecture will be given.".format(study.study_name))
                with tf.Session() as sess:
                    sess.run(tf.global_variables_initializer())
                    arc = sess.run(controller_ops["sample_arc"])
                    # TODO: will use PVC to store the checkpoint to protect against unexpected suggestion pod restart
                    saver.save(sess, study.ctrl_cache_file)

                study.is_first_run = False

            else:
                with tf.Session() as sess:
                    saver.restore(sess, study.ctrl_cache_file)

                    valid_acc = ctrl.reward
                    result = self.GetEvaluationResult(study)

                    # This lstm cell is designed to maximize the metrics
                    # However, if the user want to minimize the metrics, we can take the negative of the result
                    if study.opt_direction == api_pb2.MINIMIZE:
                        result = -result

                    loss, entropy, lr, gn, bl, skip, _ = sess.run(
                        fetches=run_ops,
                        feed_dict={valid_acc: result})
                    self.logger.info("Suggetion updated. LSTM Controller Reward: {}".format(loss))

                    arc = sess.run(controller_ops["sample_arc"])

                    saver.save(sess, study.ctrl_cache_file)

        arc = arc.tolist()
        organized_arc = [0 for _ in range(study.num_layers)]
        record = 0
        for l in range(study.num_layers):
            organized_arc[l] = arc[record: record + l + 1]
            record += l + 1

        nn_config = dict()
        nn_config['num_layers'] = study.num_layers
        nn_config['input_size'] = study.input_size
        nn_config['output_size'] = study.output_size
        nn_config['embedding'] = dict()
        for l in range(study.num_layers):
            opt = organized_arc[l][0]
            nn_config['embedding'][opt] = study.search_space[opt].get_dict()

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

        channel = grpc.beta.implementations.insecure_channel(MANAGER_ADDRESS, MANAGER_PORT)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            for i, t in enumerate(trials):
                ctrep = client.CreateTrial(api_pb2.CreateTrialRequest(trial=t), 10)
                trials[i].trial_id = ctrep.trial_id
            self.logger.info("Trial {} Created\n".format(ctrep.trial_id))
            study.prev_trial_id = ctrep.trial_id
        
        study.ctrl_step += 1

        return api_pb2.GetSuggestionsReply(trials=trials)

    def GetEvaluationResult(self, study):
        worker_list = []
        channel = grpc.beta.implementations.insecure_channel(MANAGER_ADDRESS, MANAGER_PORT)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            gwfrep = client.GetWorkerFullInfo(api_pb2.GetWorkerFullInfoRequest(study_id=study.study_id, trial_id=study.prev_trial_id, only_latest_log=True), 10)
            worker_list = gwfrep.worker_full_infos

        for w in worker_list:
            if w.Worker.status == api_pb2.COMPLETED:
                for ml in w.metrics_logs:
                    if ml.name == study.objective_name:
                        self.logger.info("Evaluation result of previous candidate: {}".format(ml.values[-1].value))
                        return float(ml.values[-1].value)

        # TODO: add support for multiple trials
