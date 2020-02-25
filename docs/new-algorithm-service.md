# Document about how to add a new algorithm in Katib

## Implement a new algorithm and use it in Katib

### Implement the algorithm

The design of Katib follows the `ask-and-tell` pattern:

> They often follow a pattern a bit like this: 1. ask for a new set of parameters 1. walk to the experiment and program in the new parameters 1. observe the outcome of running the experiment 1. walk back to your laptop and tell the optimizer about the outcome 1. go to step 1

When an experiment is created, one algorithm service will be created. Then Katib asks for new sets of parameters via `GetSuggestions` GRPC call. After that, Katib creates new trials according to the sets and observe the outcome. When the trials are finished, Katib tells the metrics of the finished trials to the algorithm, and ask another new sets. 

The new algorithm needs to implement `Suggestion` service defined in [api.proto](../pkg/apis/manager/v1alpha3/api.proto). One sample algorithm looks like:

```python
from pkg.apis.manager.v1alpha3.python import api_pb2
from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.suggestion.v1alpha3.internal.search_space import HyperParameter, HyperParameterSearchSpace
from pkg.suggestion.v1alpha3.internal.trial import Trial, Assignment
from pkg.suggestion.v1alpha3.hyperopt.base_hyperopt_service import BaseHyperoptService
from pkg.suggestion.v1alpha3.base_health_service import HealthServicer


# Inherit SuggestionServicer and implement GetSuggestions
class HyperoptService(
        api_pb2_grpc.SuggestionServicer, HealthServicer):
    def ValidateAlgorithmSettings(self, request, context):
        # Optional, it is used to validate algorithm settings defined by users.
        pass
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

### Make the algorithm a GRPC server

Create a package under [cmd/suggestion](../cmd/suggestion). Then create the main function and Dockerfile. The new GRPC server should serve in port 6789.

Here is an example: [cmd/suggestion/hyperopt](../cmd/suggestion/hyperopt). Then build the Docker image.

### Use the algorithm in Katib.

Update the [katib-config](../manifests/v1alpha3/katib-controller/katib-config.yaml), add a new object:

```json
  suggestion: |-
    {
      "tpe": {
        "image": "gcr.io/kubeflow-images-public/katib/v1alpha3/suggestion-hyperopt"
      },
      "random": {
        "image": "gcr.io/kubeflow-images-public/katib/v1alpha3/suggestion-hyperopt"
      },
      "<new-algorithm-name>": {
        "image": "image built in the previous stage"
      }
    }
```

### Contribute the algorithm to Katib

If you want to contribute the algorithm to Katib, you could add unit test or e2e test for it in CI and submit a PR.

#### Unit Test

Here is an example [test_hyperopt_service.py](../test/suggestion/v1alpha3/test_hyperopt_service.py):

```python
import grpc
import grpc_testing
import unittest

from pkg.apis.manager.v1alpha3.python import api_pb2_grpc
from pkg.apis.manager.v1alpha3.python import api_pb2

from pkg.suggestion.v1alpha3.hyperopt_service import HyperoptService

class TestHyperopt(unittest.TestCase):
    def setUp(self):
        servicers = {
            api_pb2.DESCRIPTOR.services_by_name['Suggestion']: HyperoptService()
        }

        self.test_server = grpc_testing.server_from_dictionary(
            servicers, grpc_testing.strict_real_time())


if __name__ == '__main__':
    unittest.main()
```

You can setup the GRPC server using `grpc_testing`, then define you own test cases.

#### E2E Test (Optional)

E2e tests help Katib verify that the algorithm works well.
To add a e2e test for the new algorithm, in [test/scripts/v1alpha3](../test/scripts/v1alpha3) you need to:

1. Create a new Experiment yaml file in [examples/v1alpha3](../examples/v1alpha3) with the new algorithm.

2. Create a new script `build-suggestion-xxx.sh` to build new suggestion. Here is an example [test/scripts/v1alpha3/build-suggestion-hyperopt.sh](../test/scripts/v1alpha3/build-suggestion-hyperopt.sh).

3. Create a new script `run-suggestion-xxx.sh` to run new suggestion. Below is an example (Replace `<name>` with the new algorithm name):

```bash
#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This shell script is used to build a cluster and create a namespace from our
# argo workflow

set -o errexit
set -o nounset
set -o pipefail

CLUSTER_NAME="${CLUSTER_NAME}"
ZONE="${GCP_ZONE}"
PROJECT="${GCP_PROJECT}"
NAMESPACE="${DEPLOY_NAMESPACE}"
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Configuring kubectl"

echo "CLUSTER_NAME: ${CLUSTER_NAME}"
echo "ZONE: ${GCP_ZONE}"
echo "PROJECT: ${GCP_PROJECT}"

gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
  --zone ${ZONE}
kubectl config set-context $(kubectl config current-context) --namespace=default
USER=`gcloud config get-value account`

echo "All Katib components are running."
kubectl version
kubectl cluster-info
echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod

mkdir -p ${GO_DIR}
cp -r . ${GO_DIR}/
cp -r pkg/apis/manager/v1alpha3/python/* ${GO_DIR}/test/e2e/v1alpha3
cd ${GO_DIR}/test/e2e/v1alpha3

echo "Running e2e <name> experiment"
export KUBECONFIG=$HOME/.kube/config
go run run-e2e-experiment.go ../../../examples/v1alpha3/<name>-example.yaml
kubectl -n kubeflow describe suggestion
kubectl -n kubeflow delete experiment <name>-example
exit 0
```

Then add a new step in our CI to run the new e2e test case in [test/workflows/components/workflows-v1alpha3.libsonnet](../test/workflows/components/workflows-v1alpha3.libsonnet) (Replace `<name>` with the new algorithm name):

```diff
// ...
                  {
                    name: "run-nasrl-e2e-tests",
                    template: "run-nasrl-e2e-tests",
                  },
                  {
                    name: "run-hyperband-e2e-tests",
                    template: "run-hyperband-e2e-tests",
                  },
                  {
                    name: "run-tpe-e2e-tests",
                    template: "run-tpe-e2e-tests",
                  },
+                  {
+                    name: "run-<name>-e2e-tests",
+                    template: "run-<name>-e2e-tests",
+                  },
// ...
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tpe-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-tpe.sh",
            ]),  // run tpe algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-hyperband-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-hyperband.sh",
            ]),  // run hyperband algorithm
+            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-<name>-e2e-tests", testWorkerImage, [
+              "test/scripts/v1alpha3/run-suggestion-<name>.sh",
+            ]),  // run <name> algorithm
```
