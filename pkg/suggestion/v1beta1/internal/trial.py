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

from pkg.apis.manager.v1beta1.python import api_pb2 as api

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


class Trial(object):
    def __init__(
        self,
        name,
        assignments,
        target_metric,
        metric_name,
        additional_metrics,
        labels,
    ):
        self.name = name
        self.assignments = assignments
        self.target_metric = target_metric
        self.metric_name = metric_name
        self.additional_metrics = additional_metrics
        self.labels = labels

    @staticmethod
    def convert(trials):
        res = []
        for trial in trials:
            if (
                trial.status.condition == api.TrialStatus.TrialConditionType.SUCCEEDED
                or trial.status.condition
                == api.TrialStatus.TrialConditionType.EARLYSTOPPED
            ):
                new_trial = Trial.convertTrial(trial)
                if new_trial is not None:
                    res.append(Trial.convertTrial(trial))
        return res

    @staticmethod
    def convertTrial(trial):
        assignments = []
        for assignment in trial.spec.parameter_assignments.assignments:
            assignments.append(Assignment.convert(assignment))
        metric_name = trial.spec.objective.objective_metric_name
        target_metric, additional_metrics = Metric.convert(
            trial.status.observation, metric_name
        )
        labels = trial.spec.labels
        # If the target_metric is none, ignore the trial.
        if target_metric is not None:
            trial = Trial(
                trial.name,
                assignments,
                target_metric,
                metric_name,
                additional_metrics,
                labels,
            )
            return trial
        return None

    def __str__(self):
        if self.name is None:
            return "Trial(assignment: {})".format(
                ", ".join([str(e) for e in self.assignments])
            )
        else:
            return (
                "Trial(assignment: {}, metric_name: {}, metric: {}, "
                "additional_metrics: {})".format(
                    ", ".join(str(e) for e in self.assignments),
                    self.metric_name,
                    self.target_metric,
                    ", ".join(str(e) for e in self.additional_metrics),
                )
            )


class Assignment(object):
    def __init__(self, name, value):
        self.name = name
        self.value = value

    @staticmethod
    def convert(assignment):
        return Assignment(assignment.name, assignment.value)

    @staticmethod
    def generate(list_of_assignments, trial_names=None, labels=None):
        if trial_names is not None and len(list_of_assignments) != len(trial_names):
            raise RuntimeError("Assignment and trial list length mismatch")
        res = []
        for n, assignments in enumerate(list_of_assignments):
            buf = []
            for assignment in assignments:
                buf.append(
                    api.ParameterAssignment(
                        name=assignment.name, value=str(assignment.value)
                    )
                )
            kwargs = {"assignments": buf}
            if trial_names is not None:
                kwargs["trial_name"] = trial_names[n]
            if labels is not None:
                kwargs["labels"] = {
                    k: str(v) for k, v in labels[n].items()
                }  # force string encoding

            rt = api.GetSuggestionsReply.ParameterAssignments(**kwargs)
            res.append(rt)
        return res

    def __str__(self):
        return "Assignment(name={}, value={})".format(self.name, self.value)


class Metric(object):
    def __init__(self, name, value):
        self.name = name
        self.value = value

    @staticmethod
    def convert(observation, target):
        metric = None
        additional_metrics = []
        for m in observation.metrics:
            if m.name == target:
                metric = Metric(m.name, m.value)
            else:
                additional_metrics.append(Metric(m.name, m.value))
        return metric, additional_metrics

    def __str__(self):
        return "Metric(name={}, value={})".format(self.name, self.value)
