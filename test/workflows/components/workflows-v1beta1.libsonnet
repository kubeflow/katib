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
    // project:: "automl-ci",
    // zone:: "us-east1-d",
    // // Default registry to use.
    // //registry:: "gcr.io/" + $.defaultParams.project,

    // // The image tag to use.
    // // Defaults to a value based on the name.
    // versionTag:: null,

    // // The name of the secret containing GCP credentials.
    // gcpCredentialsSecretName:: "kubeflow-testing-credentials",
  },

  // overrides is a dictionary of parameters to provide in addition to defaults.
  parts(namespace, name, overrides):: {
    // Workflow to run the e2e test.
    e2e(prow_env, bucket):
      local params = $.defaultParams + overrides;

      // mountPath is the directory where the volume to store the test data
      // should be mounted.
      local mountPath = "/mnt/test-data-volume";
      // testDir is the root directory for all data for a particular test run.
      // local testDir = mountPath + "/" + name;
      // outputDir is the directory to sync to GCS to contain the output for this job.
      // local outputDir = testDir + "/output";
      // local artifactsDir = outputDir + "/artifacts";
      // local goDir = testDir + "/go";
      // Source directory where all repos should be checked out

      local goDir = mountPath + "/" + name + "/go";
      local srcRootDir = goDir + "/src";
      // The directory containing the kubeflow/katib repo
      local srcDir = srcRootDir + "/kubeflow/katib";
      // The directory containing the kubeflow/manifests repo;
      local manifestsDir = srcRootDir + "/kubeflow/manifests";
      local testWorkerImage = "348134392524.dkr.ecr.us-west-2.amazonaws.com/aws-kubeflow-ci/test-worker:0.1";
      local kanikoExecutorImage = "gcr.io/kaniko-project/executor:v1.0.0";
      local pythonImage = "python:3.6-jessie";
      // The name of the NFS volume claim to use for test files.
      local nfsVolumeClaim = "nfs-external";
      // The name to use for the volume to use to contain test data.
      local dataVolume = "kubeflow-test-volume";

      // TODO (andreyvelich): Do we need it ?
      // The directory within the kubeflow_testing submodule containing
      // py scripts to use.
      local k8sPy = srcDir;
      local kubeflowPy = srcRootDir + "/kubeflow/testing/py";

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
      local registry = params.registry;
      {
        // Build an Argo template to execute a particular command.
        // step_name: Name for the template
        // command: List to pass as the container command.
        buildTemplate(step_name, image, command):: {
          name: step_name,
          retryStrategy: {
            limit: 3,
            retryPolicy: "Always",
          },
          container: {
            command: command,
            image: image,
            workingDir: srcDir,
            env: [
              {
                // Add the source directories to the python path.
                name: "PYTHONPATH",
                value: k8sPy + ":" + kubeflowPy,
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
        // TODO(jlewi): Use OnExit to run cleanup steps.
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
                [{
                  name: "checkout",
                  template: "checkout",
                }],
                [
                  {
                    name: "python-tests",
                    template: "python-tests",
                  },
                  // {
                  //   name: "build-suggestion-enas",
                  //   template: "build-suggestion-enas",
                  // },
                  // {
                  //   name: "build-manager",
                  //   template: "build-manager",
                  // },
                  {
                    name: "build-katib-controller",
                    template: "build-katib-controller",
                  },
                  // {
                  //   name: "build-file-metrics-collector",
                  //   template: "build-file-metrics-collector",
                  // },
                  // {
                  //   name: "build-tfevent-metrics-collector",
                  //   template: "build-tfevent-metrics-collector",
                  // },
                  // {
                  //   name: "build-suggestion-chocolate",
                  //   template: "build-suggestion-chocolate",
                  // },
                  // {
                  //   name: "build-suggestion-hyperband",
                  //   template: "build-suggestion-hyperband",
                  // },
                  // {
                  //   name: "build-suggestion-hyperopt",
                  //   template: "build-suggestion-hyperopt",
                  // },
                  // {
                  //   name: "build-suggestion-skopt",
                  //   template: "build-suggestion-skopt",
                  // },
                  // {
                  //   name: "build-suggestion-goptuna",
                  //   template: "build-suggestion-goptuna",
                  // },
                  // {
                  //   name: "build-suggestion-darts",
                  //   template: "build-suggestion-darts",
                  // },
                  // {
                  //   name: "build-earlystopping-median",
                  //   template: "build-earlystopping-median",
                  // },
                  // {
                  //   name: "build-ui",
                  //   template: "build-ui",
                  // },
                  // Temporarily disable py symplink
                  // {
                  //   name: "create-pr-symlink",
                  //   template: "create-pr-symlink",
                  // },
                ],
                [
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
                  // {
                  //   name: "run-grid-e2e-tests",
                  //   template: "run-grid-e2e-tests",
                  // },
                  // {
                  //   name: "run-file-metricscollector-e2e-tests",
                  //   template: "run-file-metricscollector-e2e-tests",
                  // },
                  // {
                  //   name: "run-custom-metricscollector-e2e-tests",
                  //   template: "run-custom-metricscollector-e2e-tests",
                  // },
                  // {
                  //   name: "run-bayesian-e2e-tests",
                  //   template: "run-bayesian-e2e-tests",
                  // },
                  // {
                  //   name: "run-enas-e2e-tests",
                  //   template: "run-enas-e2e-tests",
                  // },
                  // {
                  //   name: "run-hyperband-e2e-tests",
                  //   template: "run-hyperband-e2e-tests",
                  // },
                  // {
                  //   name: "run-tpe-e2e-tests",
                  //   template: "run-tpe-e2e-tests",
                  // },
                  // {
                  //   name: "run-tfjob-e2e-tests",
                  //   template: "run-tfjob-e2e-tests",
                  // },
                  // {
                  //   name: "run-pytorchjob-e2e-tests",
                  //   template: "run-pytorchjob-e2e-tests",
                  // },
                  // {
                  //   name: "run-cmaes-e2e-tests",
                  //   template: "run-cmaes-e2e-tests",
                  // },
                  // {
                  //   name: "run-never-resume-e2e-tests",
                  //   template: "run-never-resume-e2e-tests",
                  // },
                  // {
                  //   name: "run-darts-e2e-tests",
                  //   template: "run-darts-e2e-tests",
                  // },
                  // {
                  //   name: "run-from-volume-e2e-tests",
                  //   template: "run-from-volume-e2e-tests",
                  // },
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
                // [{
                //   name: "copy-artifacts",
                //   template: "copy-artifacts",
                // }],
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
            },  // checkout
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("create-cluster",testWorkerImage, [
              "test/scripts/v1beta1/create-cluster.sh",
            ]),  // setup cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("python-tests", pythonImage, [
              "test/scripts/v1beta1/python-tests.sh",
            ]),  // run python tests
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("setup-katib", testWorkerImage, [
              "test/scripts/v1beta1/setup-katib.sh",
            ]),  // check katib readiness and deploy it
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-random-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-e2e-experiment.sh examples/v1beta1/random-example.yaml",
            ]),  // run random algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tpe-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-tpe.sh",
            ]),  // run tpe algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tfjob-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-tfjob.sh",
            ]),  // run tfjob
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-pytorchjob-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-pytorchjob.sh",
            ]),  // run pytorchjob
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-hyperband-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-hyperband.sh",
            ]),  // run hyperband algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-grid-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-grid.sh",
            ]),  // run grid algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-enas-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-enas.sh",
            ]),  // run enas algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-bayesian-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-bayesian.sh",
            ]),  // run bayesian algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-file-metricscollector-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-file-metricscollector.sh",
            ]),  // run file metrics collector test
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-custom-metricscollector-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-custom-metricscollector.sh",
            ]),  // run custom metrics collector test
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-cmaes-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-cmaes.sh",
            ]),  // run cmaes algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-never-resume-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-never-resume.sh",
            ]),  // run never resume suggestion test
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-darts-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-suggestion-darts.sh",
            ]),  // run darts algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-from-volume-e2e-tests", testWorkerImage, [
              "test/scripts/v1beta1/run-from-volume.sh",
            ]),  // run resume from volume suggestion test
            // $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("create-pr-symlink", testWorkerImage, [
            //   "python",
            //   "-m",
            //   "kubeflow.testing.prow_artifacts",
            //   "--artifacts_dir=" + outputDir,
            //   "create_pr_symlink",
            //   "--bucket=" + bucket,
            // ]),  // create-pr-symlink
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("delete-cluster",testWorkerImage, [
              "test/scripts/v1beta1/delete-cluster.sh",
             ]),  // teardown cluster
            // $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("copy-artifacts", testWorkerImage, [
            //   "python",
            //   "-m",
            //   "kubeflow.testing.prow_artifacts",
            //   "--artifacts_dir=" + outputDir,
            //   "copy_artifacts",
            //   "--bucket=" + bucket,
            // ]),  // copy-artifacts
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-katib-controller", kanikoExecutorImage, [
              "/kaniko/executor",
              "--dockerfile=" + srcDir + "/cmd/katib-controller/v1beta1/Dockerfile",
              "--context=dir://" + srcDir,
              "--destination=" + "527798164940.dkr.ecr.us-west-2.amazonaws.com/katib/v1beta1/katib-controller:$(PULL_BASE_SHA)",
              // need to add volume mounts and extra env.
            ]),  // build-katib-controller
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-manager", testWorkerImage, [
              "test/scripts/v1beta1/build-manager.sh",
            ]),  // build-manager
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-file-metrics-collector", testWorkerImage, [
              "test/scripts/v1beta1/build-file-metrics-collector.sh",
            ]),  // build-file-metrics-collector
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-tfevent-metrics-collector", testWorkerImage, [
              "test/scripts/v1beta1/build-tfevent-metrics-collector.sh",
            ]),  // build-tfevent-metrics-collector
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-hyperband", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-hyperband.sh",
            ]),  // build-suggestion-hyperband
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-hyperopt", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-hyperopt.sh",
            ]),  // build-suggestion-hyperopt
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-skopt", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-skopt.sh",
            ]),  // build-suggestion-skopt
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-chocolate", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-chocolate.sh",
            ]),  // build-suggestion-chocolate
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-enas", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-enas.sh",
            ]),  // build-suggestion-enas
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-goptuna", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-goptuna.sh",
            ]),  // build-suggestion-goptuna
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-darts", testWorkerImage, [
              "test/scripts/v1beta1/build-suggestion-darts.sh",
            ]),  // build-suggestion-darts
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-earlystopping-median", testWorkerImage, [
              "test/scripts/v1beta1/build-earlystopping-median.sh",
            ]),  // build-earlystopping-median
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-ui", testWorkerImage, [
              "test/scripts/v1beta1/build-ui.sh",
            ]),  // build-ui
          ],  // templates
        },
      },  // e2e
  },  // parts
}
