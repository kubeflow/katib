local params = import '../../components/params.libsonnet';

params {
  components+: {
    workflows+: {
      namespace: 'kubeflow-test-infra',
      name: 'katib-release-47d4d8f7-kunming',
      prow_env: 'JOB_NAME=katib-release,JOB_TYPE=katib-release,REPO_NAME=katib,REPO_OWNER=kubeflow,BUILD_NUMBER=2649,PULL_BASE_SHA=47d4d8f7',
      versionTag: 'v20190702-47d4d8f7',
      registry: 'gcr.io/kubeflow-images-public',
      bucket: 'kubeflow-releasing-artifacts',
    },
    "workflows-v1alpha2"+: {
      namespace: 'kubeflow-test-infra',
      name: 'katib-release-47d4d8f7-ffa4-kunming',
      prow_env: 'JOB_NAME=katib-release,JOB_TYPE=katib-release,REPO_NAME=katib,REPO_OWNER=kubeflow,BUILD_NUMBER=FFA4,PULL_BASE_SHA=47d4d8f7',
      versionTag: 'v20190702-47d4d8f7',
      registry: 'gcr.io/kubeflow-images-public',
      bucket: 'kubeflow-releasing-artifacts',
    },
  },
}