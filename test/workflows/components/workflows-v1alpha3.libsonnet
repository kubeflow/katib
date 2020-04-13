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
    project:: "kubeflow-ci",
    zone:: "us-east1-d",
    // Default registry to use.
    //registry:: "gcr.io/" + $.defaultParams.project,

    // The image tag to use.
    // Defaults to a value based on the name.
    versionTag:: null,

    // The name of the secret containing GCP credentials.
    gcpCredentialsSecretName:: "kubeflow-testing-credentials",
  },

  // overrides is a dictionary of parameters to provide in addition to defaults.
  parts(namespace, name, overrides):: {
    // Workflow to run the e2e test.
    e2e(prow_env, bucket):
      local params = $.defaultParams + overrides;

      // mountPath is the directory where the volume to store the test data
      // should be mounted.
      local mountPath = "/mnt/" + "test-data-volume";
      // testDir is the root directory for all data for a particular test run.
      local testDir = mountPath + "/" + name;
      // outputDir is the directory to sync to GCS to contain the output for this job.
      local outputDir = testDir + "/output";
      local artifactsDir = outputDir + "/artifacts";
      local goDir = testDir + "/go";
      // Source directory where all repos should be checked out
      local srcRootDir = testDir + "/src";
      // The directory containing the kubeflow/katib repo
      local srcDir = srcRootDir + "/kubeflow/katib";
      // The directory containing the kubeflow/manifests repo;
      local manifestsDir = srcRootDir + "/kubeflow/manifests";
      local testWorkerImage = "gcr.io/kubeflow-ci/test-worker:v20190802-c6f9140-e3b0c4";
      local pythonImage = "python:3.6-jessie";
      // The name of the NFS volume claim to use for test files.
      // local nfsVolumeClaim = "kubeflow-testing";
      local nfsVolumeClaim = "nfs-external";
      // The name to use for the volume to use to contain test data.
      local dataVolume = "kubeflow-test-volume";
      local versionTag = if params.versionTag != null then
        params.versionTag
        else name;

      // The namespace on the cluster we spin up to deploy into.
      local deployNamespace = "kubeflow";
      // The directory within the kubeflow_testing submodule containing
      // py scripts to use.
      local k8sPy = srcDir;
      local kubeflowPy = srcRootDir + "/kubeflow/testing/py";

      local project = params.project;
      // GKE cluster to use
      // We need to truncate the cluster to no more than 40 characters because
      // cluster names can be a max of 40 characters.
      // We expect the suffix of the cluster name to be unique salt.
      // We prepend a z because cluster name must start with an alphanumeric character
      // and if we cut the prefix we might end up starting with "-" or other invalid
      // character for first character.
      local cluster =
        if std.length(name) > 40 then
          "z" + std.substr(name, std.length(name) - 39, 39)
        else
          name;
      local zone = params.zone;
      local registry = params.registry;
      local chart = srcDir + "/katib-chart";
      {
        // Build an Argo template to execute a particular command.
        // step_name: Name for the template
        // command: List to pass as the container command.
        buildTemplate(step_name, image, command):: {
          name: step_name,
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
                name: "GCP_ZONE",
                value: zone,
              },
              {
                name: "GCP_PROJECT",
                value: project,
              },
              {
                name: "GCP_REGISTRY",
                value: registry,
              },
              {
                name: "DEPLOY_NAMESPACE",
                value: deployNamespace,
              },
              {
                name: "GOOGLE_APPLICATION_CREDENTIALS",
                value: "/secret/gcp-credentials/key.json",
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
                name: "gcp-credentials",
                mountPath: "/secret/gcp-credentials",
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
              name: "gcp-credentials",
              secret: {
                secretName: params.gcpCredentialsSecretName,
              },
            },
            {
              name: dataVolume,
              persistentVolumeClaim: {
                claimName: nfsVolumeClaim,
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
                  {
                    name: "build-suggestion-enas",
                    template: "build-suggestion-enas",
                  },
                  {
                    name: "build-manager",
                    template: "build-manager",
                  },
                  {
                    name: "build-katib-controller",
                    template: "build-katib-controller",
                  },
                  {
                    name: "build-suggestion-chocolate",
                    template: "build-suggestion-chocolate",
                  },
                  {
                    name: "build-suggestion-hyperband",
                    template: "build-suggestion-hyperband",
                  },
                  {
                    name: "build-suggestion-hyperopt",
                    template: "build-suggestion-hyperopt",
                  },
                  {
                    name: "build-suggestion-skopt",
                    template: "build-suggestion-skopt",
                  },
                  {
                    name: "build-earlystopping-median",
                    template: "build-earlystopping-median",
                  },
                  {
                    name: "build-ui",
                    template: "build-ui",
                  },
                  {
                    name: "create-pr-symlink",
                    template: "create-pr-symlink",
                  },
                ],
                [
                  {
                    name: "setup-cluster",
                    template: "setup-cluster",
                  },
                ],
                [
                  {
                    name: "check-katib-ready",
                    template: "check-katib-ready",
                  },
                ],
                [
                  {
                    name: "run-random-e2e-tests",
                    template: "run-random-e2e-tests",
                  },
                  {
                    name: "run-grid-e2e-tests",
                    template: "run-grid-e2e-tests",
                  },
                  {
                    name: "run-file-metricscollector-e2e-tests",
                    template: "run-file-metricscollector-e2e-tests",
                  },
                  {
                    name: "run-custom-metricscollector-e2e-tests",
                    template: "run-custom-metricscollector-e2e-tests",
                  },
                  {
                    name: "run-bayesian-e2e-tests",
                    template: "run-bayesian-e2e-tests",
                  },
                  {
                    name: "run-enas-e2e-tests",
                    template: "run-enas-e2e-tests",
                  },
                  {
                    name: "run-hyperband-e2e-tests",
                    template: "run-hyperband-e2e-tests",
                  },
                  {
                    name: "run-tpe-e2e-tests",
                    template: "run-tpe-e2e-tests",
                  },
                  {
                    name: "run-tfjob-e2e-tests",
                    template: "run-tfjob-e2e-tests",
                  },
                  {
                    name: "run-pytorchjob-e2e-tests",
                    template: "run-pytorchjob-e2e-tests",
                  },
                ],
              ],
            },
            {
              name: "exit-handler",
              steps: [
                [{
                  name: "teardown-cluster",
                  template: "teardown-cluster",

                }],
                [{
                  name: "copy-artifacts",
                  template: "copy-artifacts",
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
            },  // checkout
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("setup-cluster",testWorkerImage, [
              "test/scripts/v1alpha3/create-cluster.sh",
            ]),  // setup cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("python-tests", pythonImage, [
              "test/scripts/v1alpha3/python-tests.sh",
            ]),  // run python tests
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("check-katib-ready", testWorkerImage, [
              "test/scripts/v1alpha3/check-katib-ready.sh",
            ]),  // check katib readiness
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-random-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-random.sh",
            ]),  // run random algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tpe-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-tpe.sh",
            ]),  // run tpe algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-tfjob-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-tfjob.sh",
            ]),  // run tfjob
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-pytorchjob-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-pytorchjob.sh",
            ]),  // run pytorchjob
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-hyperband-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-hyperband.sh",
            ]),  // run hyperband algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-grid-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-grid.sh",
            ]),  // run grid algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-enas-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-enas.sh",
            ]),  // run enas algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-bayesian-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-suggestion-bayesian.sh",
            ]),  // run bayesian algorithm
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-file-metricscollector-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-file-metricscollector.sh",
            ]),  // run file metrics collector test
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("run-custom-metricscollector-e2e-tests", testWorkerImage, [
              "test/scripts/v1alpha3/run-custom-metricscollector.sh",
            ]),  // run custom metrics collector test
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("create-pr-symlink", testWorkerImage, [
              "python",
              "-m",
              "kubeflow.testing.prow_artifacts",
              "--artifacts_dir=" + outputDir,
              "create_pr_symlink",
              "--bucket=" + bucket,
            ]),  // create-pr-symlink
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("teardown-cluster",testWorkerImage, [
              "test/scripts/v1alpha3/delete-cluster.sh",
             ]),  // teardown cluster
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("copy-artifacts", testWorkerImage, [
              "python",
              "-m",
              "kubeflow.testing.prow_artifacts",
              "--artifacts_dir=" + outputDir,
              "copy_artifacts",
              "--bucket=" + bucket,
            ]),  // copy-artifacts
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-manager", testWorkerImage, [
              "test/scripts/v1alpha3/build-manager.sh",
            ]),  // build-manager
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-katib-controller", testWorkerImage, [
              "test/scripts/v1alpha3/build-katib-controller.sh",
            ]),  // build-katib-controller
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-hyperband", testWorkerImage, [
              "test/scripts/v1alpha3/build-suggestion-hyperband.sh",
            ]),  // build-suggestion-hyperband
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-hyperopt", testWorkerImage, [
              "test/scripts/v1alpha3/build-suggestion-hyperopt.sh",
            ]),  // build-suggestion-hyperopt
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-skopt", testWorkerImage, [
              "test/scripts/v1alpha3/build-suggestion-skopt.sh",
            ]),  // build-suggestion-skopt
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-chocolate", testWorkerImage, [
              "test/scripts/v1alpha3/build-suggestion-chocolate.sh",
            ]),  // build-suggestion-chocolate
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-suggestion-enas", testWorkerImage, [
              "test/scripts/v1alpha3/build-suggestion-enas.sh",
            ]),  // build-suggestion-enas
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-earlystopping-median", testWorkerImage, [
              "test/scripts/v1alpha3/build-earlystopping-median.sh",
            ]),  // build-earlystopping-median
            $.parts(namespace, name, overrides).e2e(prow_env, bucket).buildTemplate("build-ui", testWorkerImage, [
              "test/scripts/v1alpha3/build-ui.sh",
            ]),  // build-ui
          ],  // templates
        },
      },  // e2e
  },  // parts
}
