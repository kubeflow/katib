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

import grpc
import logging
import numpy as np
import os
import shutil

from pkg.apis.manager.v1beta1.python import api_pb2
from pkg.apis.manager.v1beta1.python import api_pb2_grpc
from pkg.suggestion.v1beta1.internal.search_space import HyperParameter, HyperParameterSearchSpace
import pkg.suggestion.v1beta1.internal.constant as constant
from pkg.suggestion.v1beta1.internal.trial import Trial, Assignment
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer

logger = logging.getLogger(__name__)


_REQUIRED_SETTINGS = ["checkpoint_dir", "n_population", "truncation_threshold"]
_DATA_PATH = "/opt/katib/data"


class PbtService(api_pb2_grpc.SuggestionServicer, HealthServicer):
    def __init__(self):
        super(PbtService, self).__init__()
        self.is_first_run = True
        self.settings = None
        self.job_queue = None

    def ValidateAlgorithmSettings(self, request, context):
        settings = {entry.name: entry.value for entry in request.experiment.spec.algorithm.algorithm_settings}
        missing_settings = [k for k in _REQUIRED_SETTINGS if k not in settings]
        if len(missing_settings) > 0:
            return self._set_validate_context_error(context, "Required params missing: {}".format(", ".join(missing_settings)))

        if int(settings["n_population"]) < 5:
            return self._set_validate_context_error(context, "Param(n_population) should be >= 5")
        if not 0 <= float(settings["truncation_threshold"]) <= 1:
            return self._set_validate_context_error(context, "Param(truncation_threshold) should be between 0 and 1, inclusive")
        if "resample_probability" in settings and not 0 <= settings["resample_probability"] <= 1:
            return self._set_validate_context_error(
                context,
                "Param(resample_probability) should be null to perturb at 0.8 or 1.2, or be between 0 and 1, inclusive, to resample",
            )

        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        if self.is_first_run:
            settings = {entry.name: entry.value for entry in request.experiment.spec.algorithm.algorithm_settings}  # all type:str
            search_space = HyperParameterSearchSpace.convert(request.experiment)
            search_space = [HyperParameterSampler.from_hyperparameter(i) for i in search_space.params]
            # Always maximize the objective_metric
            objective_scale = 1 if request.experiment.spec.objective.type == api_pb2.MAXIMIZE else -1
            objective_metric = request.experiment.spec.objective.objective_metric_name
            # Create Job Manager
            self.job_queue = PbtJobQueue(
                request.experiment.name,
                int(settings["n_population"]),
                float(settings["truncation_threshold"]),
                None if not "resample_probability" in settings else float(settings["resample_probability"]),
                search_space,
                objective_metric,
                objective_scale,
            )

            self.is_first_run = False

        # Update states for completed trials
        for trial in request.trials:
            if not trial.status.condition in (
                api_pb2.TrialStatus.TrialConditionType.CREATED,
                api_pb2.TrialStatus.TrialConditionType.RUNNING,
            ) and not self.job_queue.is_processed(trial.name):
                self.job_queue.update(trial)

        # Synchronize generation
        request_count = request.current_request_number
        if len(self.job_queue) < request_count:
            if len(self.job_queue) > 0:
                request_count = len(self.job_queue)
            elif len(self.job_queue.running) > 0:
                logger.info("Waiting for trials to complete before spawning next generation...")
                return api_pb2.GetSuggestionsReply(parameter_assignments=[])
            else:
                logger.info("Spawning next generation...")
                self.job_queue.generate()

        parameter_assignments = Assignment.generate([self.job_queue.get() for _ in range(request_count)])
        logger.info("Transmiting suggestion...")
        return api_pb2.GetSuggestionsReply(parameter_assignments=parameter_assignments)

    def _set_validate_context_error(self, context, error_message):
        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
        context.set_details(error_message)
        logger.info(error_message)
        return api_pb2.ValidateAlgorithmSettingsReply()


class HyperParameterSampler(HyperParameter):
    def __init__(self, *args, **kwargs):
        super(HyperParameterSampler, self).__init__(*args, **kwargs)
        if self.type in (constant.INTEGER, constant.DOUBLE):
            self.sample_list = np.arange(float(self.min), float(self.max) + float(self.step) / 2, float(self.step)).astype(
                int if self.type == constant.INTEGER else float
            )
        else:
            self.sample_list = self.list

    def sample(self):
        return np.random.choice(self.sample_list, 1)[0]

    def perterb(self, value):
        if self.type == constant.INTEGER:
            new_value = int(int(value) * np.random.choice([0.8, 1.2], 1)[0])
            return max(float(self.min), min(float(self.max), new_value))
        elif self.type == constant.DOUBLE:
            new_value = float(value) * np.random.choice([0.8, 1.2], 1)[0]
            return max(float(self.min), min(float(self.max), new_value))
        else:
            sample_index = self.sample_list.index(value) + np.random.choice([-1, 1], 1)[0]
            if sample_index >= len(self.sample_list):
                return self.sample_list[0]
            else:
                return self.sample_list[sample_index]

    @staticmethod
    def from_hyperparameter(hp: HyperParameter):
        return HyperParameterSampler(hp.name, hp.type, hp.min, hp.max, hp.list, hp.step)


class PbtJob(object):
    def __init__(self, assignment_list, index=None, generation=0, previous_index=None):
        self.params = {a.name: str(a.value) for a in assignment_list}  # all assignments transmitted as str
        self.index = index
        self.generation = generation
        self.previous_index = previous_index
        self.metric_value = None

    def assignments(self):
        return [Assignment(k, v) for k, v in self.params.items()]

    def __repr__(self):
        return "({}|{})<-({}) : {}".format(self.index, self.generation, self.previous_index, self.params)


class PbtJobQueue(object):
    def __init__(
        self, experiment_name, population_size, truncation_threshold, resample_probability, search_space, metric_name, metric_scaler
    ):
        self.name = experiment_name

        self.population_size = population_size
        self.truncation_threshold = truncation_threshold
        self.resample_probability = resample_probability

        self.search_space = search_space
        self.metric_name = metric_name
        self.metric_scaler = metric_scaler

        self.pending = []
        self.running = []
        self.generation = 0
        self.completed = [{}]

        self.total_trial_count = 0

        # Seed the initial population
        for n in range(self.population_size):
            self.append([Assignment(param.name, param.sample()) for param in self.search_space], self.generation)

    def __len__(self):
        return len(self.pending)

    def _min_generation_indices(self):
        # find minimum generation
        min_gen = min([i.generation for i in self.pending if not i is None])
        # find indices of minimum generation
        return [n for n, i in enumerate(self.pending) if i.generation == min_gen]

    def _get_objective_value(self, trial):
        for m in trial.status.observation.metrics:
            if m.name == self.metric_name:
                return self.metric_scaler * float(m.value)
        logger.error("Objective metric value could not be found for {}".format(trial.name))

    def is_processed(self, trial_name):
        for i in range(self.generation + 1):
            if trial_name in self.completed[i]:
                return True
        return False

    def append(self, assignments, generation, link_idx=None):
        if generation > 0 and link_idx is None:
            logger.warning("Trial generation >0 but no previous checkpoint trial provided")
        obj = PbtJob(assignments, self.total_trial_count, generation, link_idx)
        self.pending.append(obj)
        new_trial_dir = os.path.join(_DATA_PATH, "{}-{}".format(self.name, self.total_trial_count))
        if link_idx is None:
            os.makedirs(new_trial_dir, exist_ok=True)
        else:
            previous_trial_dir = os.path.join(_DATA_PATH, "{}-{}".format(self.name, link_idx))
            shutil.copytree(previous_trial_dir, new_trial_dir)

        logger.info("PbtJob enqueued: {}".format(obj))
        self.total_trial_count += 1

    def get(self):
        if len(self.pending) == 0:
            raise RuntimeError("Pending queue is empty!")

        # Attempt to minimze collisions on update()
        pop_idx = None
        min_gen_indices = self._min_generation_indices()
        for n in min_gen_indices:
            compare_job = self.pending[n]
            if sum([r.params == compare_job.params for r in self.running]) == 0:
                pop_idx = n
                break
        if pop_idx is None:
            logger.warn("Multiple jobs running with same assignment list; collision possible")
            pop_idx = min_gen_indices[0]

        # Move job to running queue
        obj = self.pending.pop(pop_idx)
        logger.info("PbtJob submitted: {}".format(obj))
        self.running.append(obj)

        return obj.assignments()

    def update(self, trial):
        trial_assignments = [Assignment.convert(assignment) for assignment in trial.spec.parameter_assignments.assignments]
        trial_job = PbtJob(trial_assignments)

        # Find first match in running and move to completed with metrics
        pop_idx = None
        for n, i in enumerate(self.running):
            if i.params == trial_job.params:
                pop_idx = n
                break
        if pop_idx is None:
            logger.error("Unable to find trial match in PbtJobQueue.running!")
            return
        obj = self.running.pop(pop_idx)
        obj.metric_value = self._get_objective_value(trial)
        self.completed[self.generation][trial.name] = obj

        # Generate next trial if necessary
        if (
            trial.status.condition
            in (api_pb2.TrialStatus.TrialConditionType.SUCCEEDED, api_pb2.TrialStatus.TrialConditionType.EARLYSTOPPED)
            and Trial.convertTrial(trial) is None
        ):
            # Trial invalid, retry
            logger.error("Unable to convert trial {} (status: {}), re-adding to task queue".format(trial.name, trial.status.condition))
            self.append(trial_assignments, obj.generation, obj.previous_index)
        elif trial.status.condition in (api_pb2.TrialStatus.TrialConditionType.KILLED, api_pb2.TrialStatus.TrialConditionType.FAILED):
            # Trial error, retry
            logger.warning("Trial failed {} (status: {}), re-adding to task queue".format(trial.name, trial.status.condition))
            self.append(trial_assignments, obj.generation, obj.previous_index)

    def generate(self):
        values = [job.metric_value for _, job in self.completed[self.generation].items() if not job.metric_value is None]
        trunc_bounds = np.quantile(values, (self.truncation_threshold, 1 - self.truncation_threshold))
        exploit_names, explore_names, upper_names = [], [], []
        for trial_name, job in self.completed[self.generation].items():
            if job.metric_value is None:
                continue
            if job.metric_value < trunc_bounds[0]:
                exploit_names.append(trial_name)
            else:
                explore_names.append(trial_name)
                if job.metric_value >= trunc_bounds[1]:
                    upper_names.append(trial_name)

        # Generate exploit jobs
        exploit_replacements = np.random.choice(upper_names, len(exploit_names))
        for n, trial_name in enumerate(exploit_names):
            job = self.completed[self.generation][trial_name]
            self.append(self.completed[self.generation][exploit_replacements[n]].assignments(), job.generation + 1, job.index)

        # Generate perterbed trials
        for trial_name in explore_names:
            job = self.completed[self.generation][trial_name]
            assignment_list = []
            for param in self.search_space:
                if self.resample_probability is None:
                    value = job.params[param.name]
                    new_value = param.perterb(value)
                elif np.random.random() < self.resample_probability:
                    new_value = param.sample()
                else:
                    new_value = job.params[param.name]
                assignment_list.append(Assignment(param.name, new_value))
            self.append(assignment_list, job.generation + 1, job.index)

        self.completed.append({})
        self.generation += 1
