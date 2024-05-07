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

# This script creates Argo Workflow for the e2e Katib tests.

from kubeflow.testing import argo_build_util

# Main worker image to execute Workflow.
IMAGE_WORKER = "public.ecr.aws/j1r0q0g6/kubeflow-testing:latest"
# Kaniko image to build Katib images.
IMAGE_KANIKO = "gcr.io/kaniko-project/executor:v1.0.0"

# Volume to store test data among the Workflow tasks.
VOLUME_TEST_DATA = "kubeflow-test-volume"
# Volume mount path to store test data among the Workflow tasks.
MOUNT_PATH = "/mnt/test-data-volume"
# Volume to store GitHub token to clone repos.
VOLUME_GITHUB_TOKEN = "github-token"
# Volume to store AWS secret for the Kaniko build.
VOLUME_AWS_SECRET = "aws-secret"
# Volume to store Docker config for Kaniko build.
VOLUME_DOCKER_CONFIG = "docker-config"

# Entrypoint for the Argo Workflow.
ENTRYPOINT = "e2e"
# The template that should always run when the Workflow is complete.
EXIT_HANDLER = "exit-handler"

# Dict with all Katib images.
# Key - image name, Value - dockerfile location.
KATIB_IMAGES = {
    "katib-controller":              "cmd/katib-controller/v1beta1/Dockerfile",
    "katib-db-manager":              "cmd/db-manager/v1beta1/Dockerfile",
    "katib-ui":                      "cmd/ui/v1beta1/Dockerfile",
    "file-metrics-collector":        "cmd/metricscollector/v1beta1/file-metricscollector/Dockerfile",
    "tfevent-metrics-collector":     "cmd/metricscollector/v1beta1/tfevent-metricscollector/Dockerfile",
    "suggestion-hyperopt":           "cmd/suggestion/hyperopt/v1beta1/Dockerfile",
    "suggestion-skopt":              "cmd/suggestion/skopt/v1beta1/Dockerfile",
    "suggestion-hyperband":          "cmd/suggestion/hyperband/v1beta1/Dockerfile",
    "suggestion-goptuna":            "cmd/suggestion/goptuna/v1beta1/Dockerfile",
    "suggestion-optuna":             "cmd/suggestion/optuna/v1beta1/Dockerfile",
    "suggestion-pbt":                "cmd/suggestion/pbt/v1beta1/Dockerfile",
    "suggestion-enas":               "cmd/suggestion/nas/enas/v1beta1/Dockerfile",
    "suggestion-darts":              "cmd/suggestion/nas/darts/v1beta1/Dockerfile",
    "earlystopping-medianstop":      "cmd/earlystopping/medianstop/v1beta1/Dockerfile",
    "trial-pytorch-mnist":           "examples/v1beta1/trial-images/pytorch-mnist/Dockerfile",
    "trial-tf-mnist-with-summaries": "examples/v1beta1/trial-images/tf-mnist-with-summaries/Dockerfile",
    "trial-enas-cnn-cifar10-gpu":    "examples/v1beta1/trial-images/enas-cnn-cifar10/Dockerfile.gpu",
    "trial-enas-cnn-cifar10-cpu":    "examples/v1beta1/trial-images/enas-cnn-cifar10/Dockerfile.cpu",
    "trial-darts-cnn-cifar10":       "examples/v1beta1/trial-images/darts-cnn-cifar10/Dockerfile",
    "trial-simple-pbt":              "examples/v1beta1/trial-images/simple-pbt/Dockerfile",
}

# Dict with Katib Experiments to run during the test.
# Key - image name, Value - dockerfile location.
KATIB_EXPERIMENTS = {
    "random":                                 "examples/v1beta1/hp-tuning/random.yaml",
    "grid":                                   "examples/v1beta1/hp-tuning/grid.yaml",
    "bayesianoptimization":                   "examples/v1beta1/hp-tuning/bayesian-optimization.yaml",
    "tpe":                                    "examples/v1beta1/hp-tuning/tpe.yaml",
    "multivariate-tpe":                       "examples/v1beta1/hp-tuning/multivariate-tpe.yaml",
    "cmaes":                                  "examples/v1beta1/hp-tuning/cma-es.yaml",
    "hyperband":                              "examples/v1beta1/hp-tuning/hyperband.yaml",
    "pbt":                                    "examples/v1beta1/hp-tuning/simple-pbt.yaml",
    "enas":                                   "examples/v1beta1/nas/enas-cpu.yaml",
    "darts":                                  "examples/v1beta1/nas/darts-cpu.yaml",
    "pytorchjob":                             "examples/v1beta1/kubeflow-training-operator/pytorchjob-mnist.yaml",
    "tfjob":                                  "examples/v1beta1/kubeflow-training-operator/tfjob-mnist-with-summaries.yaml",
    "file-metricscollector":                  "examples/v1beta1/metrics-collector/file-metrics-collector.yaml",
    "file-metricscollector-with-json-format": "examples/v1beta1/metrics-collector/file-metrics-collector-with-json-format.yaml",
    "never-resume":                           "examples/v1beta1/resume-experiment/never-resume.yaml",
    "from-volume-resume":                     "examples/v1beta1/resume-experiment/from-volume-resume.yaml",
    "median-stop":                            "examples/v1beta1/early-stopping/median-stop.yaml",
    "median-stop-with-json-format":           "examples/v1beta1/early-stopping/median-stop-with-json-format.yaml",
}
# How many Experiments are running in parallel.
PARALLEL_EXECUTION = 5


class WorkflowBuilder(object):
    def __init__(self, workflow_name, workflow_namespace, test_dir, ecr_registry):
        """WorkflowBuilder constructor.

        :param workflow_name: Argo Workflow name.
        :param workflow_namespace: Argo Workflow namespace.
        :param test_dir: Root directory to store all data for a particular test run.
        :param ecr_registry: ECR registry to push the test images.
        """

        self.workflow_name = workflow_name
        self.workflow_namespace = workflow_namespace
        self.test_dir = test_dir
        self.katib_dir = test_dir + "/src/github.com/kubeflow/katib"
        self.manifest_dir = test_dir + "/src/github.com/kubeflow/manifests"
        self.ecr_registry = ecr_registry

    def create_task_template(self, task_name, exec_image, command):
        """Creates template for all the Workflow tasks.

        :param task_name: Template name for the task.
        :param exec_image: Container image to execute the task.
        :param command: List of container commands.

        :return: Created task template.
        """

        # Container environment variables.
        # TODO (andreyvelich): Add PYTHONPATH ?
        env = [
            {
                "name": "AWS_ACCESS_KEY_ID",
                "valueFrom": {
                    "secretKeyRef": {
                        "name": "aws-credentials",
                        "key": "AWS_ACCESS_KEY_ID"
                    }
                }
            },
            {
                "name": "AWS_SECRET_ACCESS_KEY",
                "valueFrom": {
                    "secretKeyRef": {
                        "name": "aws-credentials",
                        "key": "AWS_SECRET_ACCESS_KEY"
                    }
                }
            },
            {
                "name": "AWS_REGION",
                "value": "us-west-2"
            },
            {
                "name": "CLUSTER_NAME",
                "value": self.workflow_name
            },
            {
                "name": "EKS_CLUSTER_VERSION",
                "value": "1.19"
            },
            {
                "name": "ECR_REGISTRY",
                "value": self.ecr_registry
            },
            {
                "name": "GIT_TOKEN",
                "valueFrom": {
                    "secretKeyRef": {
                        "name": "github-token",
                        "key": "github_token"
                    }
                }
            },
            {
                "name": "MANIFESTS_DIR",
                "value": self.manifest_dir
            },
            {
                "name": "EXTRA_REPOS",
                "value": "kubeflow/testing@HEAD;kubeflow/manifests@v1.5-branch"
            },
            # Set GOPATH to test_dir because Katib repo is located under /src/github.com/kubeflow/katib
            {
                "name": "GOPATH",
                "value": self.test_dir
            }
        ]

        # Container volume mounts.
        volume_mounts = [
            {
                "name": VOLUME_TEST_DATA,
                "mountPath": MOUNT_PATH
            },
            {
                "name": VOLUME_GITHUB_TOKEN,
                "mountPath": "/secret/github-token"
            },
            {
                "name": VOLUME_AWS_SECRET,
                "mountPath": "/root/.aws/"
            },
            {
                "name": VOLUME_DOCKER_CONFIG,
                "mountPath": "/kaniko/.docker/"
            },
        ]

        task_template = {
            "name": task_name,
            # Each container can be alive for 40 minutes.
            "retryStrategy": {
                "limit": "3",
                "retryPolicy": "Always",
                "backoff": {
                    "duration": "1",
                    "factor": "2",
                    "maxDuration": "1m",
                },
            },
            "container": {
                "command": command,
                "image": exec_image,
                "workingDir": self.katib_dir,
                "env": env,
                "volumeMounts": volume_mounts,
            }
        }

        # Add prow env to the task template.
        prow_env_dict = argo_build_util.get_prow_dict()
        for k, v in prow_env_dict.items():
            task_template["container"]["env"].append({"name": k, "value": v})

        return task_template

    def create_init_workflow(self):
        """Creates initial structure for the Argo Workflow.

        :return: Initial Argo Workflow.
        """

        # Volumes which are used in Argo Workflow.
        volumes = [
            {
                "name": VOLUME_TEST_DATA,
                "persistentVolumeClaim": {
                    "claimName": "nfs-external"
                },
            },
            {
                "name": VOLUME_GITHUB_TOKEN,
                "secret": {
                    "secretName": VOLUME_GITHUB_TOKEN
                },
            },
            {
                "name": VOLUME_AWS_SECRET,
                "secret": {
                    "secretName": VOLUME_AWS_SECRET
                },
            },
            {
                "name": VOLUME_DOCKER_CONFIG,
                "configMap": {
                    "name": VOLUME_DOCKER_CONFIG
                },
            },
        ]

        workflow = {
            "apiVersion": "argoproj.io/v1alpha1",
            "kind": "Workflow",
            "metadata": {
                "name": self.workflow_name,
                "namespace": self.workflow_namespace,
            },
            "spec": {
                "entrypoint": ENTRYPOINT,
                "volumes": volumes,
                "templates": [
                    {
                        "name": ENTRYPOINT,
                        "dag": {
                            "tasks": []
                        }
                    },
                    {
                        "name": EXIT_HANDLER,
                        "dag": {
                            "tasks": []
                        }
                    }
                ],
                "onExit": EXIT_HANDLER
            },
        }

        return workflow


def create_workflow(name, namespace, **kwargs):
    """Main function which returns Argo Workflow.

    :param name: Argo Workflow name.
    :param namespace: Argo Workflow namespace.
    :param kwargs: Argo Workflow additional arguments.

    :return: Created Argo Workflow.
    """

    test_dir = MOUNT_PATH + "/" + name
    ecr_registry = kwargs["registry"]
    builder = WorkflowBuilder(name, namespace, test_dir, ecr_registry)

    # Build initial structure for the Workflow.
    workflow = builder.create_init_workflow()

    # Delete AWS Cluster in the exit handler step.
    delete_cluster = builder.create_task_template(
        task_name="delete-cluster",
        exec_image=IMAGE_WORKER,
        command=[
            "/usr/local/bin/delete-eks-cluster.sh",
        ]
    )
    argo_build_util.add_task_to_dag(workflow, EXIT_HANDLER, delete_cluster, [])

    # Step 1. Checkout GitHub repositories.
    checkout = builder.create_task_template(
        task_name="checkout",
        exec_image=IMAGE_WORKER,
        command=[
            "/usr/local/bin/checkout.sh",
            test_dir + "/src/github.com"
        ]
    )
    argo_build_util.add_task_to_dag(workflow, ENTRYPOINT, checkout, [])

    # Step 2.1 Build all Katib images.
    depends = []
    for image, dockerfile in KATIB_IMAGES.items():
        build_image = builder.create_task_template(
            task_name="build-"+image,
            exec_image=IMAGE_KANIKO,
            command=[
                "/kaniko/executor",
                "--dockerfile={}/{}".format(builder.katib_dir, dockerfile),
                "--context=dir://" + builder.katib_dir,
                "--destination={}/katib/v1beta1/{}:$(PULL_PULL_SHA)".format(ecr_registry, image)
            ]
        )
        argo_build_util.add_task_to_dag(workflow, ENTRYPOINT, build_image, [checkout["name"]])
        depends.append(build_image["name"])

    # Step 2.2 Create AWS cluster.
    create_cluster = builder.create_task_template(
        task_name="create-cluster",
        exec_image=IMAGE_WORKER,
        command=[
            "/usr/local/bin/create-eks-cluster.sh",
        ]
    )
    argo_build_util.add_task_to_dag(workflow, ENTRYPOINT, create_cluster, [checkout["name"]])
    depends.append(create_cluster["name"])

    # Step 3. Setup Katib on AWS cluster.
    setup_katib = builder.create_task_template(
        task_name="setup-katib",
        exec_image=IMAGE_WORKER,
        command=[
            "test/e2e/v1beta1/scripts/setup-katib.sh"
        ]
    )

    # Installing Katib after cluster is created and images are built.
    argo_build_util.add_task_to_dag(workflow, ENTRYPOINT, setup_katib, depends)

    # Step 4. Run Katib Experiments.
    depends = [setup_katib["name"]]
    tmp_depends = []
    for index, (experiment, location) in enumerate(KATIB_EXPERIMENTS.items()):
        run_experiment = builder.create_task_template(
            task_name="run-e2e-experiment-"+experiment,
            exec_image=IMAGE_WORKER,
            command=[
                "test/e2e/v1beta1/scripts/run-e2e-experiment.sh",
                location
            ]
        )
        argo_build_util.add_task_to_dag(workflow, ENTRYPOINT, run_experiment, depends)
        tmp_depends.append(run_experiment["name"])
        # We run only X number of Experiments at the same time. index starts with 0
        if (index+1) % PARALLEL_EXECUTION == 0:
            depends, tmp_depends = tmp_depends, []

    return workflow
