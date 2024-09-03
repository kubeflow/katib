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
import os
import shutil
import uuid

import grpc
import numpy as np

import pkg.suggestion.v1beta1.internal.constant as constant
from pkg.apis.manager.v1beta1.python import api_pb2, api_pb2_grpc
from pkg.suggestion.v1beta1.internal.base_health_service import HealthServicer
from pkg.suggestion.v1beta1.internal.search_space import (
    HyperParameter,
    HyperParameterSearchSpace,
)
from pkg.suggestion.v1beta1.internal.trial import Assignment

logger = logging.getLogger(__name__)


_REQUIRED_SETTINGS = ["suggestion_trial_dir", "n_population", "truncation_threshold"]
_DATA_PATH = "/opt/katib/data"


class PbtService(api_pb2_grpc.SuggestionServicer, HealthServicer):
    def __init__(self):
        super(PbtService, self).__init__()
        self.is_first_run = True
        self.settings = None
        self.job_queue = None

    def ValidateAlgorithmSettings(self, request, context):
        settings = {
            entry.name: entry.value
            for entry in request.experiment.spec.algorithm.algorithm_settings
        }
        missing_settings = [k for k in _REQUIRED_SETTINGS if k not in settings]
        if len(missing_settings) > 0:
            return self._set_validate_context_error(
                context,
                "Required params missing: {}".format(", ".join(missing_settings)),
            )

        if int(settings["n_population"]) < 5:
            return self._set_validate_context_error(
                context, "Param(n_population) should be >= 5"
            )
        if not 0 <= float(settings["truncation_threshold"]) <= 1:
            return self._set_validate_context_error(
                context,
                "Param(truncation_threshold) should be between 0 and 1, inclusive",
            )
        if (
            "resample_probability" in settings
            and not 0 <= settings["resample_probability"] <= 1
        ):
            return self._set_validate_context_error(
                context,
                "Param(resample_probability) should be null to perturb at 0.8 or 1.2, "
                "or be between 0 and 1, inclusive, to resample",
            )

        return api_pb2.ValidateAlgorithmSettingsReply()

    def GetSuggestions(self, request, context):
        if self.is_first_run:
            settings = {
                entry.name: entry.value
                for entry in request.experiment.spec.algorithm.algorithm_settings
            }  # all type:str
            search_space = HyperParameterSearchSpace.convert(request.experiment)
            search_space = [
                HyperParameterSampler.from_hyperparameter(i)
                for i in search_space.params
            ]
            # Always maximize the objective_metric
            objective_scale = (
                1 if request.experiment.spec.objective.type == api_pb2.MAXIMIZE else -1
            )
            objective_metric = request.experiment.spec.objective.objective_metric_name
            # Create Job Manager
            self.job_queue = PbtJobQueue(
                request.experiment.name,
                int(settings["n_population"]),
                float(settings["truncation_threshold"]),
                (
                    None
                    if "resample_probability" not in settings
                    else float(settings["resample_probability"])
                ),
                search_space,
                objective_metric,
                objective_scale,
            )

            self.is_first_run = False

        # Update states for completed trials
        for trial in request.trials:
            self.job_queue.update(trial)

        request_count = request.current_request_number
        if len(self.job_queue) < request_count:
            self.job_queue.generate(request_count)

        jobs = [self.job_queue.get() for _ in range(request_count)]
        parameter_assignments = Assignment.generate(
            [j[0] for j in jobs],
            trial_names=[j[2] for j in jobs],
            labels=[j[1] for j in jobs],
        )
        logger.info("Transmitting suggestion...")

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
            self.sample_list = np.arange(
                float(self.min),
                float(self.max) + float(self.step) / 2,
                float(self.step),
            ).astype(int if self.type == constant.INTEGER else float)
        else:
            self.sample_list = self.list

    def sample(self):
        return np.random.choice(self.sample_list, 1)[0]

    def perturb(self, value):
        if self.type == constant.INTEGER:
            new_value = int(int(value) * np.random.choice([0.8, 1.2], 1)[0])
            return max(float(self.min), min(float(self.max), new_value))
        elif self.type == constant.DOUBLE:
            new_value = float(value) * np.random.choice([0.8, 1.2], 1)[0]
            return max(float(self.min), min(float(self.max), new_value))
        else:
            sample_index = (
                self.sample_list.index(value) + np.random.choice([-1, 1], 1)[0]
            )
            if sample_index >= len(self.sample_list):
                return self.sample_list[0]
            else:
                return self.sample_list[sample_index]

    @staticmethod
    def from_hyperparameter(hp: HyperParameter):
        return HyperParameterSampler(hp.name, hp.type, hp.min, hp.max, hp.list, hp.step)


class PbtJob(object):
    def __init__(self, uid, assignment_list, generation, parent=None):
        self.uid = uid
        self.params = {
            a.name: str(a.value) for a in assignment_list
        }  # all assignments transmitted as str
        self.generation = generation
        self.parent = parent
        self.metric_value = None

    def get(self):
        assignments = [Assignment(k, v) for k, v in self.params.items()]
        labels = {
            "pbt.suggestion.katib.kubeflow.org/generation": self.generation,
        }
        if self.parent is not None:
            labels["pbt.suggestion.katib.kubeflow.org/parent"] = self.parent
        return assignments, labels, self.uid

    def __repr__(self):
        return "generation: {}, uid: {}, parent: {}".format(
            self.generation, self.uid, self.parent
        )


class PbtJobQueue(object):
    def __init__(
        self,
        experiment_name,
        population_size,
        truncation_threshold,
        resample_probability,
        search_space,
        metric_name,
        metric_scaler,
    ):
        self.experiment_name = experiment_name
        self.suggestion_dir = os.path.join(_DATA_PATH, self.experiment_name)

        self.population_size = population_size
        self.truncation_threshold = truncation_threshold
        self.resample_probability = resample_probability

        self.search_space = search_space
        self.metric_name = metric_name
        self.metric_scaler = metric_scaler

        self.pending = []
        self.running = {}
        self.completed = {}
        self.sample_pool = {"previous": [], "current": []}

        # Seed the initial population
        self._seed_from_base(self.population_size)

    def __len__(self):
        return len(self.pending)

    def _get_objective_value(self, trial):
        for m in trial.status.observation.metrics:
            if m.name == self.metric_name:
                return self.metric_scaler * float(m.value)
        logger.error(
            "Objective metric value could not be found for {}".format(trial.name)
        )

    def _seed_from_base(self, count):
        for n in range(count):
            self.append(
                assignments=[
                    Assignment(param.name, param.sample())
                    for param in self.search_space
                ],
                generation=0,
            )

    def append(self, assignments, generation, parent=None):
        if generation > 0 and parent is None:
            logger.warning(
                "Trial generation >0 but no previous checkpoint trial provided"
            )

        obj = PbtJob(
            uid="{}-{}".format(self.experiment_name, uuid.uuid4()),
            assignment_list=assignments,
            generation=generation,
            parent=parent,
        )
        self.pending.append(obj)

        new_trial_dir = os.path.join(self.suggestion_dir, obj.uid)
        if os.path.isdir(new_trial_dir):
            shutil.rmtree(new_trial_dir)

        if parent is None:
            os.makedirs(new_trial_dir, exist_ok=True)
        else:
            previous_trial_dir = os.path.join(self.suggestion_dir, obj.parent)
            shutil.copytree(previous_trial_dir, new_trial_dir)

        logger.info("PbtJob enqueued: {}".format(obj))
        return obj.uid

    def get(self):
        if len(self.pending) == 0:
            raise RuntimeError("Pending queue is empty!")

        # Move job to running queue
        obj = self.pending.pop(0)
        logger.info("PbtJob submitted: {}".format(obj))
        self.running[obj.uid] = obj

        return obj.get()

    def update(self, trial):
        uid = trial.name

        # Do not update active/pending trials
        if trial.status.condition in (
            api_pb2.TrialStatus.TrialConditionType.CREATED,
            api_pb2.TrialStatus.TrialConditionType.RUNNING,
        ):
            return
        # Do not update previously processed trials
        if uid in self.completed:
            return

        obj = self.running.pop(uid)
        obj.metric_value = self._get_objective_value(trial)
        self.completed[obj.uid] = obj

        # Re-queue killed or failed jobs
        if trial.status.condition in (
            api_pb2.TrialStatus.TrialConditionType.KILLED,
            api_pb2.TrialStatus.TrialConditionType.FAILED,
        ):
            # Trial error, retry
            logger.warning(
                "Trial failed {} (status: {}), re-queuing".format(
                    trial.name, trial.status.condition
                )
            )
            trial_assignments = [
                Assignment.convert(assignment)
                for assignment in trial.spec.parameter_assignments.assignments
            ]
            self.append(
                assignments=trial_assignments,
                generation=obj.generation,
                parent=obj.parent,
            )
            return

        self.sample_pool["current"].append(obj.uid)

    def _segment_sample_pool(self, pool, count):
        # Keep the first population_size samples to construct the new generation
        trial_pool = [self.completed[uid] for uid in self.sample_pool[pool]]

        values = [job.metric_value for job in trial_pool]
        trunc_bounds = np.quantile(
            values, (self.truncation_threshold, 1 - self.truncation_threshold)
        )

        exploit_names, explore_names, upper_names = [], [], []
        for job in trial_pool:
            if job.metric_value < trunc_bounds[0]:
                exploit_names.append(job.uid)
            else:
                explore_names.append(job.uid)
                if job.metric_value >= trunc_bounds[1]:
                    upper_names.append(job.uid)

        # Keep count samples in (exploit + explore)
        np.random.shuffle(exploit_names)
        np.random.shuffle(explore_names)
        exploit_names = list(exploit_names[: int(count * self.truncation_threshold)])
        explore_names = list(explore_names[: (count - len(exploit_names))])
        return exploit_names, explore_names, upper_names

    def generate(self, min_count):
        # Check if new generation can be created
        if len(self.sample_pool["current"]) <= self.population_size:
            # Check if first generation
            if len(self.sample_pool["previous"]) == 0:
                logger.info(
                    "Spawning {} additional samples from original search space...".format(
                        min_count
                    )
                )
                self._seed_from_base(min_count)
                return

            # Sample from previous generation
            logger.info(
                "Spawning {} additional samples from previous generation...".format(
                    min_count
                )
            )
            exploit_names, explore_names, upper_names = self._segment_sample_pool(
                "previous", min_count
            )
        else:
            logger.info("Spawning next generation...")
            exploit_names, explore_names, upper_names = self._segment_sample_pool(
                "current", self.population_size
            )
            self.sample_pool["previous"] = self.sample_pool["current"]
            self.sample_pool["current"] = []

        # Generate exploit jobs
        exploit_replacements = np.random.choice(upper_names, len(exploit_names))
        for n, trial_uid in enumerate(exploit_names):
            job = self.completed[trial_uid]
            self.append(
                assignments=self.completed[exploit_replacements[n]].get()[0],
                generation=job.generation + 1,
                parent=job.uid,
            )

        # Generate perturbed trials
        for trial_uid in explore_names:
            job = self.completed[trial_uid]
            assignment_list = []
            for param in self.search_space:
                if self.resample_probability is None:
                    value = job.params[param.name]
                    new_value = param.perturb(value)
                elif np.random.random() < self.resample_probability:
                    new_value = param.sample()
                else:
                    new_value = job.params[param.name]
                assignment_list.append(Assignment(param.name, new_value))
            self.append(
                assignments=assignment_list,
                generation=job.generation + 1,
                parent=job.uid,
            )
