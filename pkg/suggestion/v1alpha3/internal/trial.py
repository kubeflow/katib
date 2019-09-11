import logging
from pkg.apis.manager.v1alpha3.python import api_pb2 as api


logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger("Trial")


class Trial(object):
    def __init__(self, name, assignments, target_metric, metric_name, additional_metrics):
        self.name = name
        self.assignments = assignments
        self.target_metric = target_metric
        self.metric_name = metric_name
        self.additional_metrics = additional_metrics

    @staticmethod
    def convert(trials):
        res = []
        for trial in trials:
            res.append(Trial.convertTrial(trial))
        return res
    
    @staticmethod
    def generate(trials):
        res = []
        for trial in trials:
            assignments = []
            for assignment in trial.assignments:
                assignments.append(
                    api.ParameterAssignment(name=assignment.name, value=str(assignment.value)))
            rt = api.Trial(
                spec=api.TrialSpec(
                    parameter_assignments=api.TrialSpec.ParameterAssignments(
                        assignments=assignments
                    )
                )
            )
            res.append(rt)
        return res

    @staticmethod
    def convertTrial(trial):
        assignments = []
        for assignment in trial.spec.parameter_assignments.assignments:
            assignments.append(Assignment.convert(assignment))
        metric_name = trial.spec.objective.objective_metric_name
        target_metric, additional_metrics = Metric.convert(
            trial.status.observation, metric_name)
        trial = Trial(trial.name, assignments, target_metric,
                      metric_name, additional_metrics)
        return trial

    def __str__(self):
        if self.name == None:
            return "Trial(assignment: {})".format(", ".join([str(e) for e in self.assignments]))
        else:
            return "Trial(assignment: {}, metric_name: {}, metric: {}, additional_metrics: {})".format(
                ", ".join([str(e) for e in self.assignments]),
                self.metric_name, self.target_metric,
                ", ".join(str(e) for e in self.additional_metrics))


class Assignment(object):
    def __init__(self, name, value):
        self.name = name
        self.value = value

    @staticmethod
    def convert(assignment):
        return Assignment(assignment.name, assignment.value)

    def __str__(self):
        return "Assignment(name={}, value={})".format(self.name, self.value)


class Metric(object):
    def __init__(self, name, value):
        self.name = name
        self.value = value

    @staticmethod
    def convert(observation, target):
        metric = ""
        additional_metrics = []
        for metric in observation.metrics:
            if metric.name == target:
                metric = Metric(metric.name, metric.value)
            else:
                additional_metrics.append(Metric(metric.name, metric.value))
        return metric, additional_metrics

    def __str__(self):
        return "Metric(name={}, value={})".format(self.name, self.value)
