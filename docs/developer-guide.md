<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Developer Guide](#developer-guide)
  - [Requirements](#requirements)
  - [Build from source code](#build-from-source-code)
  - [Implement new suggestion algorithm](#implement-new-suggestion-algorithm)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer Guide

## Requirements

- [Go](https://golang.org/)
- [Dep](https://golang.github.io/dep/)
- [Docker](https://docs.docker.com/) (17.05 or later.)

## Build from source code

Check source code as follows:

```bash
make build
```

You can deploy katib v1alpha3 manifests into a k8s cluster as follows:

```bash
make deploy
```

You can undeploy katib v1alpha3 manifests from a k8s cluster as follows:

```bash
make undeploy
```

## Implement a new algorithm

The design of katib follows the [`ask-and-tell` pattern](https://scikit-optimize.github.io/notebooks/ask-and-tell.html):

> They often follow a pattern a bit like this: 1. ask for a new set of parameters 1. walk to the experiment and program in the new parameters 1. observe the outcome of running the experiment 1. walk back to your laptop and tell the optimizer about the outcome 1. go to step 1

When an experiment is created, one algorithm service will be created. Then katib asks for new sets of parameters via `GetSuggestions` GRPC call. After that, katib creates new trials according to the sets and observe the outcome. When the trials are finished, katib tells the metrics of the finished trials to the algorithm, and ask another new sets. One sample algorithm looks like:

```python
from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.suggestion.v1alpha3.internal.search_space import HyperParameter, HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.hyperopt.base_hyperopt_service import BaseHyperoptService

# Inherit SuggestionServicer and implement GetSuggestions
class HyperoptService(
        api_pb2_grpc.SuggestionServicer):
    def GetSuggestions(self, request, context):
        # Convert the experiment in GRPC request to the search space.
        # search_space example:
        #   HyperParameterSearchSpace(
        #       goal: MAXIMIZE, 
        #       params: [HyperParameter(name: param-1, type: INTEGER, min: 1, max: 5, step: 0), 
        #                HyperParameter(name: param-2, type: CATEGORICAL, list: cat1, cat2, cat3),
        #                HyperParameter(name: param-3, type: DISCRETE, list: 3, 2, 6),
        #                HyperParameter(name: param-4, type: DOUBLE, min: 1, max: 5, step: )]
        #   )
        search_space = HyperParameterSearchSpace.convert(request.experiment)
        # Convert the trials in GRPC request to the trials in algorithm side.
        # trials example:
        #   [Trial(
        #       assignment: [Assignment(name=param-1, value=2), 
        #                    Assignment(name=param-2, value=cat1), 
        #                    Assignment(name=param-3, value=2), 
        #                    Assignment(name=param-4, value=3.44)],
        #       target_metric: Metric(name="metric-2" value="5643"), 
        #       additional_metrics: [Metric(name=metric-1, value=435), 
        #                            Metric(name=metric-3, value=5643)],
        #   Trial(
        #       assignment: [Assignment(name=param-1, value=3),
        #                    Assignment(name=param-2, value=cat2),
        #                    Assignment(name=param-3, value=6),
        #                    Assignment(name=param-4, value=4.44)],
        #       target_metric: Metric(name="metric-2" value="3242"), 
        #       additional_metrics: [Metric(name=metric=1, value=123), 
        #                            Metric(name=metric-3, value=543)],
        trials = Trial.convert(request.trials)
        #--------------------------------------------------------------
        # Your code here
        # Implment the logic to generate new assignments for the given request number.
        # For example, if request.request_number is 2, you should return:
        # [
        #   [Assignment(name=param-1, value=3), 
        #    Assignment(name=param-2, value=cat2), 
        #    Assignment(name=param-3, value=3), 
        #    Assignment(name=param-4, value=3.22)
        #   ],
        #   [Assignment(name=param-1, value=4), 
        #    Assignment(name=param-2, value=cat4), 
        #    Assignment(name=param-3, value=2), 
        #    Assignment(name=param-4, value=4.32)
        #   ],
        # ]
        list_of_assignments = your_logic(search_space, trials, request.request_number)
        #--------------------------------------------------------------
        # Convert list_of_assignments to 
        return api_pb2.GetSuggestionsReply(
            trials=Assignment.generate(list_of_assignments)
        )
```