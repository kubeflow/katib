from pkg.suggestion.v1alpha2.NAS_Reinforcement_Learning.Controller import Controller
from pkg.suggestion.v1alpha2.NAS_Reinforcement_Learning.Operation import SearchSpace
from pkg.suggestion.v1alpha2.NAS_Reinforcement_Learning.AlgorithmSettings import parseAlgorithmSettings
import tensorflow as tf
import grpc
from pkg.api.v1alpha2.python import api_pb2
from pkg.api.v1alpha2.python import api_pb2_grpc

import logging
from logging import getLogger, StreamHandler, INFO, DEBUG
import json
import os
import time

MANAGER_ADDRESS = "katib-manager"
MANAGER_PORT = 6789
RESPAWN_SLEEP = 20
RESPAWN_LIMIT = 10


class NAS_RL_Experiment(object):
    def __init__(self, request, logger):
        self.logger = logger
        self.experiment_name = request.experiment_name
        self.num_trials = 1
        if request.request_number > 0:
            self.num_trials = request.request_number
        self.tf_graph = tf.Graph()
        # self.prev_trial_ids = list()
        # self.prev_trials = None
        self.ctrl_cache_file = "ctrl_cache/{}.ckpt".format(request.experiment_name)
        self.ctrl_step = 0
        self.is_first_run = True
        self.algorithm_settings = None
        self.controller = None
        self.num_layers  = None
        self.input_sizes = None
        self.output_sizes = None
        self.num_operations = None
        self.search_space = None
        self.opt_direction = None
        self.objective_name = None
        # self.respawn_count = 0
        
        self.logger.info("-" * 100 + "\nSetting Up Suggestion for Experiment {}\n".format(request.experiment_name) + "-" * 100)
        self._get_experiment_param()
        self._setup_controller()
        self.logger.info(">>> Suggestion for Experiment {} has been initialized.\n".format(self.experiment_name))
        
    def _get_experiment_param(self):
        # this function need to
        # 1) get the number of layers
        # 2) get the I/O size
        # 3) get the available operations
        # 4) get the optimization direction (i.e. minimize or maximize)
        # 5) get the objective name
        # 6) get the algorithm settings

        channel = grpc.beta.implementations.insecure_channel(MANAGER_ADDRESS, MANAGER_PORT)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            api_experiment_param = client.GetExperiment(api_pb2.GetExperimentRequest(experiment_name=self.experiment_name), 10)

        # Get Search Space
        self.experiment_name = api_experiment_param.experiment.name
        self.opt_direction = api_experiment_param.experiment.spec.objective.type
        self.objective_name = api_experiment_param.experiment.spec.objective.objective_metric_name

        nas_config = api_experiment_param.experiment.spec.nas_config
        
        graph_config = nas_config.graph_config
        self.num_layers = int(graph_config.num_layers)
        self.input_sizes = list(map(int, graph_config.input_sizes))
        self.output_sizes = list(map(int, graph_config.output_sizes))

        search_space_raw = nas_config.operations
        search_space_object = SearchSpace(search_space_raw)
        self.search_space = search_space_object.search_space
        self.num_operations = search_space_object.num_operations
        
        self.print_search_space()

        # Get Experiment Parameters
        params_raw = api_experiment_param.experiment.spec.algorithm.algorithm_setting
        self.algorithm_settings = parseAlgorithmSettings(params_raw)

        self.print_algorithm_settings()
    
    def _setup_controller(self):
        
        with self.tf_graph.as_default():

            self.controller = Controller(
                num_layers=self.num_layers,
                num_operations=self.num_operations,
                lstm_size=self.algorithm_settings['lstm_num_cells'],
                lstm_num_layers=self.algorithm_settings['lstm_num_layers'],
                lstm_keep_prob=self.algorithm_settings['lstm_keep_prob'],
                lr_init=self.algorithm_settings['init_learning_rate'],
                lr_dec_start=self.algorithm_settings['lr_decay_start'],
                lr_dec_every=self.algorithm_settings['lr_decay_every'],
                lr_dec_rate=self.algorithm_settings['lr_decay_rate'],
                l2_reg=self.algorithm_settings['l2_reg'],
                entropy_weight=self.algorithm_settings['entropy_weight'],
                bl_dec=self.algorithm_settings['baseline_decay'],
                optim_algo=self.algorithm_settings['optimizer'],
                skip_target=self.algorithm_settings['skip-target'],
                skip_weight=self.algorithm_settings['skip-weight'],
                name="Ctrl_" + self.experiment_name,
                logger=self.logger)

            self.controller.build_trainer()

    def print_search_space(self):
        if self.search_space is None:
            self.logger.warning("Error! The Suggestion has not yet been initialized!")
            return
        
        self.logger.info(">>> Search Space for Experiment {}".format(self.experiment_name))
        for opt in self.search_space:
            opt.print_op(self.logger)
        self.logger.info("There are {} operations in total.\n".format(self.num_operations))
    
    def print_algorithm_settings(self):
        if self.algorithm_settings is None:
            self.logger.warning("Error! The Suggestion has not yet been initialized!")
            return
        
        self.logger.info(">>> Parameters of LSTM Controller for Experiment {}".format(self.experiment_name))
        for spec in self.algorithm_settings:
            if len(spec) > 13:
                self.logger.info("{}: \t{}".format(spec, self.algorithm_settings[spec]))
            else:
                self.logger.info("{}: \t\t{}".format(spec, self.algorithm_settings[spec]))
        self.logger.info("RequestNumber:\t\t{}".format(self.num_trials))
        self.logger.info("")


class NasrlService(api_pb2_grpc.SuggestionServicer):
    def __init__(self, logger=None):

        self.registered_experiments = dict()

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

        if not os.path.exists("ctrl_cache/"):
            os.makedirs("ctrl_cache/")

    def ValidateAlgorithmSettings(self, request, context):
        self.logger.info("Validate Algorithm Settings start")
        graph_config = request.experiment_spec.nas_config.graph_config

        # Validate GraphConfig
        # Check InputSize
        if not graph_config.input_sizes:
            return self.SetValidateContextError(context, "Missing InputSizes in GraphConfig:\n{}".format(graph_config))

        # Check OutputSize
        if not graph_config.output_sizes:
            return self.SetValidateContextError(context, "Missing OutputSizes in GraphConfig:\n{}".format(graph_config))

        # Check NumLayers
        if not graph_config.num_layers:
            return self.SetValidateContextError(context, "Missing NumLayers in GraphConfig:\n{}".format(graph_config))

        # Validate each operation
        operations_list = list(request.experiment_spec.nas_config.operations.operation)
        for operation in operations_list:

            # Check OperationType
            if not operation.operation_type:
                return self.SetValidateContextError(context, "Missing operationType in Operation:\n{}".format(operation))

            # Check ParameterConfigs
            if not operation.parameter_specs.parameters:
                return self.SetValidateContextError(context, "Missing ParameterConfigs in Operation:\n{}".format(operation))
            
            # Validate each ParameterConfig in Operation
            parameters_list = list(operation.parameter_specs.parameters)
            for parameter in parameters_list:

                # Check Name
                if not parameter.name:
                    return self.SetValidateContextError(context, "Missing Name in ParameterConfig:\n{}".format(parameter))

                # Check ParameterType
                if not parameter.parameter_type:
                    return self.SetValidateContextError(context, "Missing ParameterType in ParameterConfig:\n{}".format(parameter))

                # Check List in Categorical or Discrete Type
                if parameter.parameter_type == api_pb2.CATEGORICAL or parameter.parameter_type == api_pb2.DISCRETE:
                    if not parameter.feasible_space.list:
                        return self.SetValidateContextError(context, "Missing List in ParameterConfig.feasibleSpace:\n{}".format(parameter) )

                # Check Max, Min, Step in Int or Double Type
                elif parameter.parameter_type == api_pb2.INT or parameter.parameter_type == api_pb2.DOUBLE:
                    if not parameter.feasible_space.min and not parameter.feasible_space.max:
                        return self.SetValidateContextError(context, "Missing Max and Min in ParameterConfig.feasibleSpace:\n{}".format(parameter))

                    if parameter.parameter_type == api_pb2.DOUBLE and (not parameter.feasible_space.step or float(parameter.feasible_space.step) <= 0):
                        return self.SetValidateContextError(context, "Step parameter should be > 0 in ParameterConfig.feasibleSpace:\n{}".format(parameter))

        self.logger.info("All Experiment Settings are Valid")
        return api_pb2.ValidateAlgorithmSettingsReply()

    def SetValidateContextError(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        self.logger.info(error_message)
        return api_pb2.ValidateSuggestionParametersReply()

    def GetSuggestions(self, request, context):

        if request.experiment_name not in self.registered_experiments:
            self.registered_experiments[request.experiment_name] = NAS_RL_Experiment(request, self.logger)
        
        experiment = self.registered_experiments[request.experiment_name]

        self.logger.info("-" * 100 + "\nSuggestion Step {} for Experiment {}\n".format(experiment.ctrl_step, experiment.experiment_name) + "-" * 100)

        with experiment.tf_graph.as_default():
            saver = tf.train.Saver()
            ctrl = experiment.controller

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

            if experiment.is_first_run:
                self.logger.info(">>> First time running suggestion for {}. Random architecture will be given.".format(experiment.experiment_name))
                with tf.Session() as sess:
                    sess.run(tf.global_variables_initializer())
                    candidates = list()
                    for _ in range(experiment.num_trials):
                        candidates.append(sess.run(controller_ops["sample_arc"]))
                    
                    # TODO: will use PVC to store the checkpoint to protect against unexpected suggestion pod restart
                    saver.save(sess, experiment.ctrl_cache_file)

                experiment.is_first_run = False

            else:
                with tf.Session() as sess:
                    saver.restore(sess, experiment.ctrl_cache_file)

                    valid_acc = ctrl.reward
                    result = self.GetEvaluationResult(experiment)

                    # TODO: (andreyvelich) I deleted this part, should it be handle by controller?
                    # Sometimes training container may fail and GetEvaluationResult() will return None
                    # In this case, the Suggestion will:
                    # 1. Firstly try to respawn the previous trials after waiting for RESPAWN_SLEEP seconds
                    # 2. If respawning the trials for RESPAWN_LIMIT times still cannot collect valid results,
                    #    then fail the task because it may indicate that the training container has errors.
                    if result is None:
                        self.logger.warning(">>> Suggestion has spawned trials, but they all failed.")
                        self.logger.warning(">>> Please check whether the training container is correctly implemented")
                        self.logger.info(">>> Experiment {} failed".format(experiment.experiment_name))
                        return []

                    # This LSTM network is designed to maximize the metrics
                    # However, if the user wants to minimize the metrics, we can take the negative of the result

                    if experiment.opt_direction == api_pb2.MINIMIZE:
                        result = -result

                    loss, entropy, lr, gn, bl, skip, _ = sess.run(
                        fetches=run_ops,
                        feed_dict={valid_acc: result})
                    
                    self.logger.info(">>> Suggestion updated. LSTM Controller Reward: {}".format(loss))

                    candidates = list()
                    for _ in range(experiment.num_trials):
                        candidates.append(sess.run(controller_ops["sample_arc"]))

                    saver.save(sess, experiment.ctrl_cache_file)
        
        organized_candidates = list()
        trials = list()

        for i in range(experiment.num_trials):
            arc = candidates[i].tolist()
            organized_arc = [0 for _ in range(experiment.num_layers)]
            record = 0
            for l in range(experiment.num_layers):
                organized_arc[l] = arc[record: record + l + 1]
                record += l + 1
            organized_candidates.append(organized_arc)

            nn_config = dict()
            nn_config['num_layers'] = experiment.num_layers
            nn_config['input_sizes'] = experiment.input_sizes
            nn_config['output_sizes'] = experiment.output_sizes
            nn_config['embedding'] = dict()
            for l in range(experiment.num_layers):
                opt = organized_arc[l][0]
                nn_config['embedding'][opt] = experiment.search_space[opt].get_dict()

            organized_arc_json = json.dumps(organized_arc)
            nn_config_json = json.dumps(nn_config)

            organized_arc_str = str(organized_arc_json).replace('\"', '\'')
            nn_config_str = str(nn_config_json).replace('\"', '\'')

            self.logger.info("\n>>> New Neural Network Architecture Candidate #{} (internal representation):".format(i))
            self.logger.info(organized_arc_json)
            self.logger.info("\n>>> Corresponding Seach Space Description:")
            self.logger.info(nn_config_str)

            trials.append(api_pb2.Trial(
                spec=api_pb2.TrialSpec(
                    experiment_name=request.experiment_name,
                    parameter_assignments=api_pb2.TrialSpec.ParameterAssignments(
                        assignments=[
                            api_pb2.ParameterAssignment(
                                name="architecture",
                                value=organized_arc_str
                            ),
                            api_pb2.ParameterAssignment(
                                name="nn_config",
                                value=nn_config_str
                            )
                        ]
                    )
                )
            ))
        
        self.logger.info("")
        self.logger.info(">>> {} Trials were created for Experiment {}".format(experiment.num_trials, experiment.experiment_name))
        self.logger.info("")

        experiment.ctrl_step += 1

        return api_pb2.GetSuggestionsReply(trials=trials)        

    def GetEvaluationResult(self, experiment):
        channel = grpc.beta.implementations.insecure_channel(MANAGER_ADDRESS, MANAGER_PORT)
        with api_pb2.beta_create_Manager_stub(channel) as client:
            trials_resp = client.GetTrialList(api_pb2.GetTrialListRequest(experiment_name=experiment.experiment_name), 10)
            trials_list = trials_resp.trials
        
        completed_trials = dict()
        failed_trials = []
        for t in trials_list:
            if t.status.condition == api_pb2.TrialStatus.TrialConditionType.SUCCEEDED:
                obslog_resp = client.GetObservationLog(
                    api_pb2.GetObservationLogRequest(
                        trial_name=t.name,
                        metric_name=t.spec.objective.objective_metric_name
                    ), 10
                )

                # Take only the latest metric value
                completed_trials[t.name] = float(obslog_resp.observation_log.metric_logs[-1].metric.value)

            if t.status.condition == api_pb2.TrialStatus.TrialConditionType.FAILED:
                failed_trials.append(t.name)

        n_completed = len(completed_trials)
        self.logger.info(">>> By now: {} Trials succeeded, {} Trials failed".format(n_completed, len(failed_trials)))
        for tname in completed_trials:
            self.logger.info("Trial: {}, Value: {}".format(tname, completed_trials[tname]))
        for tname in failed_trials:
            self.logger.info("Trial: {} was failed".format(tname))
       
        if n_completed > 0:
            avg_metrics = sum(completed_trials.values()) / n_completed
            self.logger.info("The average is {}\n".format(avg_metrics))

            return avg_metrics
