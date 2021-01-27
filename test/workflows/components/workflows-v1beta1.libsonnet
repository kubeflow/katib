{
  // TODO(https://github.com/ksonnet/ksonnet/issues/222): Taking namespace as an argument is a work around for the fact that ksonnet
  // doesn't support automatically piping in the namespace from the environment to prototypes.

  // convert a list of two items into a map representing an environment variable
  // TODO(jlewi): Should we move this into kubeflow/core/util.libsonnet
  listToMap:: function(v)
    {
      name: v[0],
      value: v[1],
    },

  // Function to turn comma separated list of prow environment variables into a dictionary.
  parseEnv:: function(v)
    local pieces = std.split(v, ",");
    if v != "" && std.length(pieces) > 0 then
      std.map(
        function(i) $.listToMap(std.split(i, "=")),
        std.split(v, ",")
      )
    else [],


  // default parameters.
  defaultParams:: {
  },

  // overrides is a dictionary of parameters to provide in addition to defaults.
  parts(namespace, name, overrides):: {
    // Workflow to run the e2e test.
    e2e(prow_env, bucket):
      local params = $.defaultParams + overrides;
      local registry = params.registry;

      // mountPath is the directory where the volume to store the test data should be mounted.
      local mountPath = "/mnt/test-data-volume";
      // testDir is the root directory for all data for a particular test run.
      local testDir = mountPath + "/" + name;
      // goDir is the directory to run Go e2e tests.
      // Since Katib repo is located under /src/github.com/kubeflow/katib, we can set GOPATH = testDir.
      // In fact, we don't need to create additional folder to execute Go e2e test.
      local goDir = testDir;

      // srcRootDir is the directory where all repos should be checked out.
      local srcRootDir = testDir + "/src/github.com";
      // katibDir is the directory containing the kubeflow/katib repo.
      local katibDir = srcRootDir + "/kubeflow/katib";
      // manifestsDir is the directory containing the kubeflow/manifests repo.
      local manifestsDir = srcRootDir + "/kubeflow/manifests";
      // kubeflowTestingPy is the directory within the kubeflow_testing submodule containing py scripts to use.
      local kubeflowTestingPy = srcRootDir + "/kubeflow/testing/py";

      // testWorkerImage is the main worker image to execute workflow.
      local testWorkerImage = "public.ecr.aws/j1r0q0g6/kubeflow-testing:latest";
      // kanikoExecutorImage is the image for Kaniko to build Katib images.
      local kanikoExecutorImage = "gcr.io/kaniko-project/executor:v1.0.0";
      // pythonImage is the image to run Katib Python unit test.
      local pythonImage = "python:3.6-jessie";

      // The name of the NFS volume claim to use for test files.
      local nfsVolumeClaim = "nfs-external";
      // The name to use for the volume to use to contain test data.
      local dataVolume = "kubeflow-test-volume";

      // We need to truncate the cluster to no more than 40 characters because
      // cluster names can be a max of 40 characters.
      // We expect the suffix of the cluster name to be unique salt.
      // We prepend a z because cluster name must start with an alphanumeric character
      // and if we cut the prefix we might end up starting with "-" or other invalid
      // character for first character.
      local cluster =
        if std.length(name) > 80 then
          std.substr(name, std.length(name) - 79, 79)
        else
          name;
      {
        // Build an Argo template to execute a particular command.
        // step_name: Name for the template.
        // command: List to pass as the container command.
        buildTemplate(step_name, image, command):: {
          name: step_name,
          // Each container can be alive for 40 minutes.
          retryStrategy: {
            limit: "3",
            retryPolicy: "Always",
            backoff: {
              duration: "1",
              factor: "2",
              maxDuration: "1m",
            },
          },
          container: {
            command: command,
            image: image,
            workingDir: katibDir,
            env: [
              {
                // Add the source directories to the python path.
                name: "PYTHONPATH",
                value: katibDir + ":" + kubeflowTestingPy,
              },
              {
                name: "MANIFESTS_DIR",
                value: manifestsDir,
              },
              {
                // Set the GOPATH
                name: "GOPATH",
                value: goDir,
              },
              {
                name: "CLUSTER_NAME",
                value: cluster,
              },
              {
                name: "AWS_REGION",
                value: "us-west-2",
              },
              {
                name: "AWS_ACCESS_KEY_ID",
                valueFrom: {
                  secretKeyRef: {
                    name: "aws-credentials",
                    key: "AWS_ACCESS_KEY_ID",
                  },
                },
              },
              {
                name: "AWS_SECRET_ACCESS_KEY",
                valueFrom: {
                  secretKeyRef: {
                    name: "aws-credentials",
                    key: "AWS_SECRET_ACCESS_KEY",
                  },
                },
              },
              {
                name: "ECR_REGISTRY",
                value: registry,
              },
              {
                name: "GIT_TOKEN",
                valueFrom: {
                  secretKeyRef: {
                    name: "github-token",
                    key: "github_token",
                  },
                },
              },
            ] + prow_env,
            volumeMounts: [
              {
                name: dataVolume,
                mountPath: mountPath,
              },
              {
                name: "github-token",
                mountPath: "/secret/github-token",
              },
              {
                name: "aws-secret",
                mountPath: "/root/.aws/",
              },
              {
                name: "docker-config",
                mountPath: "/kaniko/.docker/",
              },
            ],
          },
        },  // buildTemplate

        apiVersion: "argoproj.io/v1alpha1",
        kind: "Workflow",
        metadata: {
          name: name,
          namespace: namespace,
        },
        spec: {
          entrypoint: "e2e",
          volumes: [
            {
              name: "github-token",
              secret: {
                secretName: "github-token",
              },
            },
            {
              name: dataVolume,
              persistentVolumeClaim: {
                claimName: nfsVolumeClaim,
              },
            },
            // Attach aws-secret and docker-config for Kaniko build
            {
              name: "docker-config",
              configMap: {
                name: "docker-config",
              },
            },
            {
              name: "aws-secret",
              secret: {
                secretName: "aws-secret",
              },
            },
          ],  // volumes
          // onExit specifies the template that should always run when the workflow completes.
          onExit: "exit-handler",
          templates: [
            {
              name: "e2e",
              steps: [
                [
                  {
                    name: "checkout",
                    template: "checkout",
                  },
                ],
                [
                  {
                    name: "python-tests",
                    template: "python-tests",
                  },
                  {
                    name: "build-katib-controller",
                    template: "build-katib-controller",
                  },
                  {
                    name: "build-db-manager",
                    template: "build-db-manager",
                  },
                  {
                    name: "build-ui",
                    template: "build-ui",
                  },
                  {
                    name: "build-file-metrics-collector",
                    template: "build-file-metrics-collector",
                  },
                  {
                    name: "build-tfevent-metrics-collector",
                    template: "build-tfevent-metrics-collector",
                  },
                  {
                    name: "build-suggestion-hyperopt",
                    template: "build-suggestion-hyperopt",
                  },
                  {
                    name: "build-suggestion-chocolate",
                    template: "build-suggestion-chocolate",
                  },
                  {
                    name: "build-suggestion-skopt",
                    template: "build-suggestion-skopt",
                  },
                  {
                    name: "build-suggestion-hyperband",
                    template: "build-suggestion-hyperband",
                  },
                  {
                    name: "build-suggestion-goptuna",
                    template: "build-suggestion-goptuna",
                  },
                  {
                    name: "build-suggestion-enas",
                    template: "build-suggestion-enas",
                  },
                  {
                    name: "build-suggestion-darts",
                    template: "build-suggestion-darts",
                  },
                  {
                    name: "build-earlystopping-medianstop",
                    template: "build-earlystopping-medianstop",
                  },
                  {
                    name: "create-cluster",
                    template: "create-cluster",
                  },
                ],
                [
                  {
                    name: "setup-katib",
                    template: "setup-katib",
                  },
                ],
                [
                  {
                    name: "run-random-e2e-tests",
                    template: "run-random-e2e-tests",
                  },
                  {
                    name: "run-tpe-e2e-tests",
                    template: "run-tpe-e2e-tests",
                  },
                  {
                    name: "run-grid-e2e-tests",
                    template: "run-grid-e2e-tests",
                  },
                  {
                    name: "run-bayesian-e2e-tests",
                    template: "run-bayesian-e2e-tests",
                  },
                  {
                    name: "run-hyperband-e2e-tests",
                    template: "run-hyperband-e2e-tests",
                  },
                  {
                    name: "run-cmaes-e2e-tests",
                    template: "run-cmaes-e2e-tests",
                  },
                  {
                    name: "run-enas-e2e-tests",
                    template: "run-enas-e2e-tests",
                  },
                  {
                    name: "run-darts-e2e-tests",
                    template: "run-darts-e2e-tests",
                  },
                  {
                    name: "run-tfjob-e2e-tests",
                    template: "run-tfjob-e2e-tests",
                  },
                  {
                    name: "run-pytorchjob-e2e-tests",
                    template: "run-pytorchjob-e2e-tests",
                  },
                  {
                    name: "run-file-metricscollector-e2e-tests",
                    template: "run-file-metricscollector-e2e-tests",
                  },
                  {
                    name: "run-never-resume-e2e-tests",
                    template: "run-never-resume-e2e-tests",
                  },
                  {
                    name: "run-from-volume-e2e-tests",
                    template: "run-from-volume-e2e-tests",
                  },
                  {
                    name: "run-medianstop-e2e-tests",
                    template: "run-medianstop-e2e-tests",
                  },
                ],
              ],
            },
            {
              name: "exit-handler",
              steps: [
                [{
                  name: "delete-cluster",
                  template: "delete-cluster",

                }],
              ],
            },
            {
              name: "checkout",
              container: {
                command: [
                  "/usr/local/bin/checkout.sh",
                  srcRootDir,
                ],
                env: prow_env + [{
                  name: "EXTRA_REPOS",
                  value: "kubeflow/testing@HEAD;kubeflow/manifests@HEAD",
                }],
                image: testWorkerImage,
                volumeMounts: [
                  {
                    name: dataVolume,
                    mountPath: mountPath,
                  },
                ],
              },
            }, // checkout
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("create-cluster", testWorkerImage, [
              "/usr/local/bin/create-eks-cluster.sh",
            ]),  // Create cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("delete-cluster", testWorkerImage, [
              "/usr/local/bin/delete-eks-cluster.sh",
            ]),  // Delete cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("python-tests", pythonImage, [
              "test/scripts/v1beta1/python-tests.sh",
            ]),  // run python tests
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-katib-controller", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/katib-controller/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/katib-controller:$(PULL_BASE_SHA)",
            ]),  // build katib-controller
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-db-manager", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/db-manager/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/katib-db-manager:$(PULL_BASE_SHA)",
            ]),  // build katib-db-manager
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-ui", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/ui/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/katib-ui:$(PULL_BASE_SHA)",
            ]),  // build katib-ui
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-file-metrics-collector", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/metricscollector/v1beta1/file-metricscollector/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/file-metrics-collector:$(PULL_BASE_SHA)",
            ]),  // build file metrics collector
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-tfevent-metrics-collector", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/metricscollector/v1beta1/tfevent-metricscollector/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/tfevent-metrics-collector:$(PULL_BASE_SHA)",
            ]),  // build tfevent metrics collector
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-hyperopt", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/hyperopt/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-hyperopt:$(PULL_BASE_SHA)",
            ]),  // build suggestion hyperopt
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-chocolate", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/chocolate/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-chocolate:$(PULL_BASE_SHA)",
            ]),  // build suggestion chocolate
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-skopt", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/skopt/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-skopt:$(PULL_BASE_SHA)",
            ]),  // build suggestion skopt
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-hyperband", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/hyperband/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-hyperband:$(PULL_BASE_SHA)",
            ]),  // build suggestion hyperband
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-goptuna", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/goptuna/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-goptuna:$(PULL_BASE_SHA)",
            ]),  // build suggestion goptuna
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-enas", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/nas/enas/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-enas:$(PULL_BASE_SHA)",
            ]),  // build suggestion enas
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-darts", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/suggestion/nas/darts/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/suggestion-darts:$(PULL_BASE_SHA)",
            ]),  // build suggestion darts
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-earlystopping-medianstop", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + katibDir + "/cmd/earlystopping/medianstop/v1beta1/Dockerfile",
              "--context=dir://" + katibDir,
              "--destination=" + registry + "/katib/v1beta1/earlystopping-medianstop:$(PULL_BASE_SHA)",
            ]),  // build early stopping median stop
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("setup-katib", testWorkerImage, [
              "test/scripts/v1beta1/setup-katib.sh",
            ]),  // check Katib readiness and deploy it
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-random-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/random-example.yaml",
            ]),  // run random algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tpe-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/tpe-example.yaml",
            ]),  // run TPE algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-grid-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/grid-example.yaml",
            ]),  // run grid algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-bayesian-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/bayesianoptimization-example.yaml",
            ]),  // run BO algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-hyperband-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/hyperband-example.yaml",
            ]),  // run hyperband algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-cmaes-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/cmaes-example.yaml",
            ]),  // run CMA-ES algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-enas-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/nas/enas-example-cpu.yaml",
            ]),  // run ENAS algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-darts-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/nas/darts-example-cpu.yaml",
            ]),  // run DARTS algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tfjob-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/tfjob-example.yaml",
            ]),  // run TFJob example
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-pytorchjob-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/pytorchjob-example.yaml",
            ]),  // run PyTorchJob example
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-file-metricscollector-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/file-metricscollector-example.yaml",
            ]),  // run file metrics collector example
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-never-resume-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/resume-experiment/never-resume.yaml",
            ]),  // run never resume example
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-from-volume-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/resume-experiment/from-volume-resume.yaml",
            ]),  // run from volume resume example
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-medianstop-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh",
              "examples/v1beta1/early-stopping/median-stop.yaml",
            ]),  // run median stopping example
            // TODO (andreyvelich): Temporary disable pr-symlink
            // $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("create-pr-symlink", testWorkerImage, [
            //   "python",
            //   "-m",
            //   "kubeflow.testing.prow_artifacts",
            //   "--artifacts_dir=" + outputDir,
            //   "create_pr_symlink",
            //   "--bucket=" + bucket,
            // ]),  // create-pr-symlink
          ],  // templates
        },
      },  // e2e
  },  // parts
}
