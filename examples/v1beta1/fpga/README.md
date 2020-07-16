# FPGA support for your Katib Experiments

Let's spawn [F1 instances](https://aws.amazon.com/ec2/instance-types/f1) and
accelerate time-consuming Katib experiments on AWS, with zero FPGA knowledge!

If you want to read more about provisioning FPGA resources and deploying
accelerated applications (e.g. Kubeflow Pipelines) on any Kubernetes cluster,
visit the [InAccel](https://docs.inaccel.com) documentation.

## EKS*-optimized AMI with FPGA support

**For development and testing purposes you can still [deploy Kubeflow Katib
using Minikube](https://kubeflow.org/docs/started/workstation/minikube-linux) in
a single AMI instance. In production environments, Amazon's managed Kubernetes
service ([EKS](https://aws.amazon.com/eks)) is recommended.*

The InAccel EKS-optimized FPGA AMI is built on top of the standard [Amazon
EKS-optimized Linux AMI](https://aws.amazon.com/marketplace/pp/B07GRMYQR5), and
is configured to serve as an optional image for Amazon EKS worker nodes to
support FPGA based workloads.

In addition to the standard Amazon EKS-optimized AMI configuration, the FPGA AMI
includes the following:

* Xilinx FPGA drivers
* InAccel runtime (as the default runtime)

The AMI IDs for the latest InAccel EKS-optimized FPGA AMI are shown in the
following table.

| AWS Region                          | Kubernetes version 1.17.7 | Kubernetes version 1.16.12 |
| ----------------------------------- | ------------------------- | -------------------------- |
| US East (N. Virginia) (`us-east-1`) | `ami-0c4e0b85781a9dde3`   | `ami-0519d9b242546530a`    |

> **Note**: The EKS-optimized FPGA AMI also supports non-FPGA instance types.

## Enabling FPGA based workloads

The following section describes how to run a workload on an FPGA based instance
with the InAccel EKS-optimized FPGA AMI.

After your FPGA worker nodes join your cluster, you must apply the [InAccel FPGA
Operator](https://hub.helm.sh/charts/inaccel/fpga-operator) for Kubernetes, as a
Helm app on your cluster, with the following command.

```sh
helm repo add inaccel https://setup.inaccel.com/helm

helm install inaccel inaccel/fpga-operator --set license=...
```

Get your free Community Edition [license](https://inaccel.com/license), if it's
the first time that you use InAccel toolset.

You can verify that your nodes have available FPGAs with the following command:

```sh
kubectl get nodes -o custom-columns=NAME:metadata.name,FPGAS:.status.capacity.xilinx/aws-vu9p-f1
```

## Experiment

You can submit a new accelerated Experiment and check your Experiment results
using the Web UI, as usual.

#### XGBoost Parameter Tuning [[source](https://github.com/inaccel/jupyter/blob/master/lab/dot/XGBoost/parameter-tuning.py)]

```sh
kubectl apply -f xgboost-example.yaml
```
