# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import logging
from logging import getLogger, StreamHandler, INFO
import json
import os
import tensorflow as tf
import grpc

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.suggestion.v1beta1.nas.enas.Controller import Controller
from pkg.suggestion.v1beta1.nas.enas.Operation import SearchSpace
from pkg.suggestion.v1beta1.nas.enas.AlgorithmSettings import (
    parseAlgorithmSettings, algorithmSettingsValidator, enableNoneSettingsList)
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer
from pkg.suggestion.v1beta1.nas.common.validation import validate_operations


class EnasExperiment:
    def __init__(self, request, logger):
        self.logger = logger
        self.experiment_name = request.experiment.name
        self.experiment = request.experiment
        self.num_trials = 1
        self.tf_graph = tf.Graph()
        self.ctrl_cache_file = "ctrl_cache/{}.ckpt".format(
            self.experiment_name)
        self.suggestion_step = 0
        self.algorithm_settings = None
        self.controller = None
        self.num_layers = None
        self.input_sizes = None
        self.output_sizes = None
        self.num_operations = None
        self.search_space = None
        self.opt_direction = None
        self.objective_name = None
        self.logger.info("-" * 100 + "\nSetting Up Suggestion for Experiment {}\n".format(
            self.experiment_name) + "-" * 100)
        self._get_experiment_param()
        self._setup_controller()
        self.logger.info(">>> Suggestion for Experiment {} has been initialized.\n".format(
            self.experiment_name))

    def _get_experiment_param(self):
        # this function need to
        # 1) get the number of layers
        # 2) get the I/O size
        # 3) get the available operations
        # 4) get the optimization direction (i.e. minimize or maximize)
        # 5) get the objective name
        # 6) get the algorithm settings

        # Get Search Space
        self.opt_direction = self.experiment.spec.objective.type
        self.objective_name = self.experiment.spec.objective.objective_metric_name

        nas_config = self.experiment.spec.nas_config

        graph_config = nas_config.graph_config
        self.num_layers = int(graph_config.num_layers)
        self.input_sizes = list(map(int, graph_config.input_sizes))
        self.output_sizes = list(map(int, graph_config.output_sizes))

        search_space_raw = nas_config.operations
        search_space_object = SearchSpace(search_space_raw)
        self.search_space = search_space_object.search_space
        self.num_operations = search_space_object.num_operations

        self.print_search_space()

        # Get Experiment Algorithm Settings
        settings_raw = self.experiment.spec.algorithm.algorithm_settings
        self.algorithm_settings = parseAlgorithmSettings(settings_raw)

        self.print_algorithm_settings()

    def _setup_controller(self):

        with self.tf_graph.as_default():

            self.controller = Controller(
                num_layers=self.num_layers,
                num_operations=self.num_operations,
                controller_hidden_size=self.algorithm_settings['controller_hidden_size'],
                controller_temperature=self.algorithm_settings['controller_temperature'],
                controller_tanh_const=self.algorithm_settings['controller_tanh_const'],
                controller_entropy_weight=self.algorithm_settings['controller_entropy_weight'],
                controller_baseline_decay=self.algorithm_settings['controller_baseline_decay'],
                controller_learning_rate=self.algorithm_settings["controller_learning_rate"],
                controller_skip_target=self.algorithm_settings['controller_skip_target'],
                controller_skip_weight=self.algorithm_settings['controller_skip_weight'],
                controller_name="Ctrl_" + self.experiment_name,
                logger=self.logger)

            self.controller.build_trainer()

    def print_search_space(self):
        if self.search_space is None:
            self.logger.warning(
                "Error! The Suggestion has not yet been initialized!")
            return

        self.logger.info(
            ">>> Search Space for Experiment {}".format(self.experiment_name))
        for opt in self.search_space:
            opt.print_op(self.logger)
        self.logger.info(
            "There are {} operations in total.\n".format(self.num_operations))

    def print_algorithm_settings(self):
        if self.algorithm_settings is None:
            self.logger.warning(
                "Error! The Suggestion has not yet been initialized!")
            return

        self.logger.info(">>> Parameters of LSTM Controller for Experiment {}\n".format(
            self.experiment_name))
        for spec in self.algorithm_settings:
            if len(spec) > 22:
                self.logger.info("{}:\t{}".format(
                    spec, self.algorithm_settings[spec]))
            else:
                self.logger.info("{}:\t\t{}".format(
                    spec, self.algorithm_settings[spec]))

        self.logger.info("")


class EnasService(api_pb2_grpc.SuggestionServicer, HealthServicer):
    def __init__(self, logger=None):
        super(EnasService, self).__init__()
        self.is_first_run = True
        self.experiment = None
        if logger is None:
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
        nas_config = request.experiment.spec.nas_config
        graph_config = nas_config.graph_config

        # Validate GraphConfig
        # Check InputSize
        if not graph_config.input_sizes:
            return self.set_validate_context_error(context,
                                                   "Missing InputSizes in GraphConfig:\n{}".format(graph_config))

        # Check OutputSize
        if not graph_config.output_sizes:
            return self.set_validate_context_error(context,
                                                   "Missing OutputSizes in GraphConfig:\n{}".format(graph_config))

        # Check NumLayers
        if not graph_config.num_layers:
            return self.set_validate_context_error(context,
                                                   "Missing NumLayers in GraphConfig:\n{}".format(graph_config))

        # Validate Operations
        is_valid, message = validate_operations(nas_config.operations.operation)
        if not is_valid:
            return self.set_validate_context_error(context, message)

        # Validate Algorithm Settings
        settings_raw = request.experiment.spec.algorithm.algorithm_settings
        for setting in settings_raw:
            if setting.name in algorithmSettingsValidator.keys():
                if setting.name in enableNoneSettingsList and setting.value == "None":
                    continue
                setting_type = algorithmSettingsValidator[setting.name][0]
                setting_range = algorithmSettingsValidator[setting.name][1]
                try:
                    converted_value = setting_type(setting.value)
                except Exception as e:
                    return self.set_validate_context_error(context,
                                                           "Algorithm Setting {} must be {} type: exception {}".format(
                                                               setting.name, setting_type.__name__, e))

                if setting_type == float:
                    if (converted_value <= setting_range[0] or
                            (setting_range[1] != 'inf' and converted_value > setting_range[1])):
                        return self.set_validate_context_error(
                            context, "Algorithm Setting {}: {} with {} type must be in range ({}, {}]".format(
                                setting.name,
                                converted_value,
                                setting_type.__name__,
                                setting_range[0],
                                setting_range[1])
                        )

                elif converted_value < setting_range[0]:
                    return self.set_validate_context_error(
                        context, "Algorithm Setting {}: {} with {} type must be in range [{}, {})".format(
                            setting.name,
                            converted_value,
                            setting_type.__name__,
                            setting_range[0],
                            setting_range[1])
                    )
            else:
                return self.set_validate_context_error(context,
                                                       "Unknown Algorithm Setting name: {}".format(setting.name))

        self.logger.info("All Experiment Settings are Valid")
        return api_pb2.ValidateAlgorithmSettingsReply()

    def set_validate_context_error(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        self.logger.info(error_message)
        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        if self.is_first_run:
            self.experiment = EnasExperiment(request, self.logger)
        experiment = self.experiment
        if request.current_request_number > 0:
            experiment.num_trials = request.current_request_number
        self.logger.info("-" * 100 + "\nSuggestion Step {} for Experiment {}\n".format(
            experiment.suggestion_step, experiment.experiment_name) + "-" * 100)

        self.logger.info("")
        self.logger.info(">>> Current Request Number:\t\t{}".format(experiment.num_trials))
        self.logger.info("")

        with experiment.tf_graph.as_default():
            saver = tf.compat.v1.train.Saver()
            ctrl = experiment.controller

            controller_ops = {
                "loss": ctrl.loss,
                "entropy": ctrl.sample_entropy,
                "grad_norm": ctrl.grad_norm,
                "baseline": ctrl.baseline,
                "skip_rate": ctrl.skip_rate,
                "train_op": ctrl.train_op,
                "train_step": ctrl.train_step,
                "sample_arc": ctrl.sample_arc,
                "child_val_accuracy": ctrl.child_val_accuracy,
            }

            if self.is_first_run:
                self.logger.info(">>> First time running suggestion for {}. Random architecture will be given.".format(
                    experiment.experiment_name))
                with tf.compat.v1.Session() as sess:
                    sess.run(tf.compat.v1.global_variables_initializer())
                    candidates = list()
                    for _ in range(experiment.num_trials):
                        candidates.append(
                            sess.run(controller_ops["sample_arc"]))

                    # TODO: will use PVC to store the checkpoint to protect against unexpected suggestion pod restart
                    saver.save(sess, experiment.ctrl_cache_file)

                self.is_first_run = False

            else:
                with tf.compat.v1.Session() as sess:
                    saver.restore(sess, experiment.ctrl_cache_file)

                    result = self.GetEvaluationResult(request.trials)

                    # TODO: (andreyvelich) I deleted this part, should it be handle by controller?
                    # Sometimes training container may fail and GetEvaluationResult() will return None
                    # In this case, the Suggestion will:
                    # 1. Firstly try to respawn the previous trials after waiting for RESPAWN_SLEEP seconds
                    # 2. If respawning the trials for RESPAWN_LIMIT times still cannot collect valid results,
                    #    then fail the task because it may indicate that the training container has errors.
                    if result is None:
                        self.logger.warning(
                            ">>> Suggestion has spawned trials, but they all failed.")
                        self.logger.warning(
                            ">>> Please check whether the training container is correctly implemented")
                        self.logger.info(">>> Experiment {} failed".format(
                            experiment.experiment_name))
                        return []

                    # This LSTM network is designed to maximize the metrics
                    # However, if the user wants to minimize the metrics, we can take the negative of the result

                    if experiment.opt_direction == api_pb2.MINIMIZE:
                        result = -result

                    self.logger.info(">>> Suggestion updated. LSTM Controller Training\n")
                    log_every = experiment.algorithm_settings["controller_log_every_steps"]
                    for ctrl_step in range(1, experiment.algorithm_settings["controller_train_steps"]+1):
                        run_ops = [
                            controller_ops["loss"],
                            controller_ops["entropy"],
                            controller_ops["grad_norm"],
                            controller_ops["baseline"],
                            controller_ops["skip_rate"],
                            controller_ops["train_op"]
                        ]

                        loss, entropy, grad_norm, baseline, skip_rate, _ = sess.run(
                            fetches=run_ops,
                            feed_dict={controller_ops["child_val_accuracy"]: result})

                        controller_step = sess.run(controller_ops["train_step"])
                        if ctrl_step % log_every == 0:
                            log_string = ""
                            log_string += "Controller Step: {} - ".format(controller_step)
                            log_string += "Loss: {:.4f} - ".format(loss)
                            log_string += "Entropy: {:.9} - ".format(entropy)
                            log_string += "Gradient Norm: {:.7f} - ".format(grad_norm)
                            log_string += "Baseline={:.4f} - ".format(baseline)
                            log_string += "Skip Rate={:.4f}".format(skip_rate)
                            self.logger.info(log_string)

                    candidates = list()
                    for _ in range(experiment.num_trials):
                        candidates.append(
                            sess.run(controller_ops["sample_arc"]))

                    saver.save(sess, experiment.ctrl_cache_file)

        organized_candidates = list()
        parameter_assignments = list()

        for i in range(experiment.num_trials):
            arc = candidates[i].tolist()
            organized_arc = [0 for _ in range(experiment.num_layers)]
            record = 0
            for layer in range(experiment.num_layers):
                organized_arc[layer] = arc[record: record + layer + 1]
                record += layer + 1
            organized_candidates.append(organized_arc)

            nn_config = dict()
            nn_config['num_layers'] = experiment.num_layers
            nn_config['input_sizes'] = experiment.input_sizes
            nn_config['output_sizes'] = experiment.output_sizes
            nn_config['embedding'] = dict()
            for layer in range(experiment.num_layers):
                opt = organized_arc[layer][0]
                nn_config['embedding'][opt] = experiment.search_space[opt].get_dict()

            organized_arc_json = json.dumps(organized_arc)
            nn_config_json = json.dumps(nn_config)

            organized_arc_str = str(organized_arc_json).replace('\"', '\'')
            nn_config_str = str(nn_config_json).replace('\"', '\'')

            self.logger.info(
                "\n>>> New Neural Network Architecture Candidate #{} (internal representation):".format(i))
            self.logger.info(organized_arc_json)
            self.logger.info("\n>>> Corresponding Seach Space Description:")
            self.logger.info(nn_config_str)

            parameter_assignments.append(
                api_pb2.GetSuggestionsReply.ParameterAssignments(
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

        self.logger.info("")
        self.logger.info(">>> {} Trials were created for Experiment {}".format(
            experiment.num_trials, experiment.experiment_name))
        self.logger.info("")

        experiment.suggestion_step += 1

        return api_pb2.GetSuggestionsReply(parameter_assignments=parameter_assignments)

    def GetEvaluationResult(self, trials_list):
        completed_trials = dict()
        failed_trials = []
        for t in trials_list:
            if t.status.condition == api_pb2.TrialStatus.TrialConditionType.SUCCEEDED:
                target_value = None
                for metric in t.status.observation.metrics:
                    if metric.name == t.spec.objective.objective_metric_name:
                        target_value = metric.value
                        break

                # Take only the first metric value
                # In current cifar-10 training container this value is the latest
                completed_trials[t.name] = float(target_value)

            if t.status.condition == api_pb2.TrialStatus.TrialConditionType.FAILED:
                failed_trials.append(t.name)

        n_completed = len(completed_trials)
        self.logger.info(">>> By now: {} Trials succeeded, {} Trials failed".format(
            n_completed, len(failed_trials)))
        for tname in completed_trials:
            self.logger.info("Trial: {}, Value: {}".format(
                tname, completed_trials[tname]))
        for tname in failed_trials:
            self.logger.info("Trial: {} was failed".format(tname))

        if n_completed > 0:
            avg_metrics = sum(completed_trials.values()) / n_completed
            self.logger.info("The average is {}\n".format(avg_metrics))

            return avg_metrics
