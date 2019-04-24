{
  global: {
    // User-defined global parameters; accessible to all component and environments, Ex:
    // replicas: 4,
  },
  components: {
    // Component-level parameters, defined initially from 'ks prototype use ...'
    // Each object below should correspond to a component in the components/ directory
    "workflows-v1alpha1": {
      bucket: "kubeflow-ci_temp",
      name: "some-very-very-very-very-very-long-name-katib-v1alpha1-presubmit-test-374-6e32",
      namespace: "kubeflow-test-infra",
      registry: "gcr.io/kubbeflow-ci",
      prow_env: "JOB_NAME=katib-v1alpha1-presubmit-test,JOB_TYPE=presubmit,PULL_NUMBER=374,REPO_NAME=k8s,REPO_OWNER=tensorflow,BUILD_NUMBER=6e32",
      versionTag: null,
    },
    "workflows-v1alpha2": {
      bucket: "kubeflow-ci_temp",
      name: "some-very-very-very-very-very-long-name-katib-v1alpha2-presubmit-test-374-6e32",
      namespace: "kubeflow-test-infra",
      registry: "gcr.io/kubbeflow-ci",
      prow_env: "JOB_NAME=katib-v1alpha2-presubmit-test,JOB_TYPE=presubmit,PULL_NUMBER=374,REPO_NAME=k8s,REPO_OWNER=tensorflow,BUILD_NUMBER=6e32",
      versionTag: null,
    },
  },
}
