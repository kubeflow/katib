local params = std.extVar("__ksonnet/params").components.gpu_pytorchjob;

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
          template: {
            spec: {
              containers: [
                {
                  image: actualImage,
                  name: "pytorch",
                  resources: {
                    limits: {
                      "nvidia.com/gpu": 1,
                    },
                  },
                },
              ],
              restartPolicy: "OnFailure",
            },
          },
          replicaType: "MASTER",
        },
      ],
    },
  },  // job
};

std.prune(k.core.v1.list.new([parts(params.namespace, params.name, params.image).job]))
