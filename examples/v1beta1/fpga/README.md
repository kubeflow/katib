# FPGA support for your Katib Experiments

Let's spawn [F1 instances](https://aws.amazon.com/ec2/instance-types/f1) and
accelerate time-consuming Katib experiments on AWS, with zero FPGA knowledge!

If you want to read more about provisioning FPGA resources and deploying
accelerated applications (e.g. Kubeflow Pipelines) on any Kubernetes cluster,
visit the [InAccel](https://docs.inaccel.com) documentation.

## Simplifying FPGA management in EKS\* (Elastic Kubernetes Service)

\*_For development and testing purposes you can still [deploy Kubeflow Katib
using Minikube](https://kubeflow.org/docs/started/workstation/minikube-linux) in
a single AMI instance. In production environments, Amazon's managed Kubernetes
service ([EKS](https://aws.amazon.com/eks)) is recommended._

The InAccel FPGA Operator allows administrators of Kubernetes clusters to manage
FPGA nodes just like CPU nodes in the cluster. Instead of provisioning a special
OS image for FPGA nodes, administrators can rely on a standard OS image for both
CPU and FPGA nodes and then rely on the FPGA Operator to provision the required
software components for FPGAs.

Note that the FPGA Operator is specifically useful for scenarios where the
Kubernetes cluster needs to scale quickly - for example provisioning additional
FPGA nodes on the cloud and managing the lifecycle of the underlying software
components.

## Enabling FPGA based workloads

The following section describes how to run a workload on an FPGA based instance
with the InAccel FPGA Operator.

After your FPGA worker nodes join your cluster, you must apply the [InAccel FPGA
Operator](https://artifacthub.io/packages/helm/inaccel/fpga-operator) for
Kubernetes, as a Helm app on your cluster, with the following command.

```sh
helm repo add inaccel https://setup.inaccel.com/helm

helm install -n kube-system inaccel inaccel/fpga-operator
```

You can verify that your nodes have available FPGAs with the following command:

```sh
kubectl get nodes -o custom-columns=NAME:metadata.name,FPGAS:.status.capacity.xilinx/aws-vu9p-f1,SHELL:.metadata.labels.xilinx/aws-vu9p-f1
```

## Experiment

You can submit a new accelerated Experiment and check your Experiment results
using the Web UI, as usual.

#### XGBoost Parameter Tuning [[source](https://github.com/inaccel/jupyter/blob/master/lab/dot/XGBoost/parameter-tuning.py)]

```sh
kubectl apply -f xgboost-example.yaml
```
