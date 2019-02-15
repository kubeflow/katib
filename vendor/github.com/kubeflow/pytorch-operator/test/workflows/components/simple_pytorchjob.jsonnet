local params = std.extVar("__ksonnet/params").components.simple_pytorchjob;

local k = import 'k.libsonnet';

local defaultTestImage = "pytorch/pytorch:v0.2";
local parts(namespace, name, image) = {
  local actualImage = if image != "" then
    image
  else defaultTestImage,
  job:: {
    apiVersion: "kubeflow.org/v1alpha1",
    kind: "PyTorchJob",
    metadata: {
      name: name,
      namespace: namespace,
    },
    spec: {
      replicaSpecs: [
        {
          replicas: 1,
          template: {
            spec: {
              containers: [
                {
                  image: actualImage,
                  name: "pytorch",
                },
              ],
              restartPolicy: "OnFailure",
            },
          },
          replicaType: "MASTER",
        },
        {
          replicas: 1,
          template: {
            spec: {
              containers: [
                {
                  image: actualImage,
                  name: "pytorch",
                },
              ],
              restartPolicy: "OnFailure",
            },
          },
          replicaType: "WORKER",
        },
      ],
    },
  },
};

std.prune(k.core.v1.list.new([parts(params.namespace, params.name, params.image).job]))
