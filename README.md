[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/gpu-driver)](https://api.reuse.software/info/github.com/kyma-project/gpu-driver)

# Overview

Kyma GPU driver is the remake of 
[gardenlinux-nvidia-installer](https://github.com/gardenlinux/gardenlinux-nvidia-installer)
that does in-cluster node kernel detection and driver compilation, not requiring 
you to maintain a container repository with all possible images built for different 
kernel versions.

# Prerequisites

Kyma GPU driver operator requires a Kyma cluster with GPU machine types.

# Installation

You can install Kyma GPU driver operator by Helm chart and with plain manifest. 
Helm chart is recommended since it simplifies the removal of unnecessary resources
and uninstallation.

## Install with Helm chart

Requirements:
* Helm CLI - for details check [Installing Helm](https://helm.sh/docs/intro/install/)

Once Helm has been set up correctly, add the repo as follows:

```shell
helm repo add kyma-gpu-driver https://kyma-project.github.io/gpu-driver
```

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages. You can then run
`helm search repo kyma-gpu-driver` to see the charts.

To install the `gpu-driver-operator` chart you should run:

```shell
helm upgrade --install gpu-driver kyma-gpu-driver/gpu-driver-operator
```

## Install with kubectl 

Requirements:
* kubectl - for details check [Kubernetes tools](https://kubernetes.io/docs/tasks/tools/#kubectl)

To install the Kyma GPU driver operator you should run

```shell
kubectl apply -f https://raw.githubusercontent.com/kyma-project/gpu-driver/refs/heads/main/config/dist/all.yaml
```

> [!NOTE]  
> If you are not familiar with the Kubernetes platform and details for the 
> Kyma GPU driver operator manifests, it is recommended to use Helm installation
> procedure. Installing the plain manifest with kubectl does not delete the old
> resources, previously installed in some older version but removed
> in the newer release. 

# Usage and API

To instruct Kyma GPU driver operator on which nodes you would like to have
GPU driver installed you should create a GpuDriver custom resource.

```yaml
apiVersion: gpu.kyma-project.io/v1beta1
kind: GpuDriver
metadata:
  name: gpu1
spec:
  nodeSelector:
    worker.gardener.cloud/pool: gpu-worker-pool
```

The resource above specifies that all nodes from the node pool `gpu-worker-pool`, will be instrumented
with the GPU driver. You can use any other set of labels as node selector. If node selector is empty, it 
will match all nodes. 


# Troubleshooting

List help repositories with expectation that `kyma-gpu-driver https://kyma-project.github.io/gpu-driver` is defined.

```shell
helm repo list | grep gpu 
```

List charts in the kyma-gpu-driver repo with expectations that latest version of the `kyma-gpu-driver/gpu-driver-operator` chart is present.

```shell
helm search repo kyma-gpu-driver
```

List all k8s nodes with expectation to have at least one node pool of machine types with GPU device, and that nodes in that pool have `.status.capacity["nvidia.com/gpu"]` > 0.

```shell
kubectl get nodes -o yaml
```

List k8s API resources with expectation that `gpudrivers gpu.kyma-project.io/v1beta1` exists.

```shell
kubectl api-resources | grep gpu
```

List k8s namespaces with expectation that `gpu-driver-system` namespace exists.

```shell
kubectl get ns
```

List `GpuDriver` CR resources with expectation to have one for each node pool having GPU capable machine type, with the appropriate node selector matching that node pool.

```shell
kubectl get GpuDriver -A -o yaml
```

List operator pods in the `gpu-driver-system` namespace.

```shell
kubectl get pods -n gpu-driver-system -o yaml
```


# Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing Rules](CONTRIBUTING.md).

# Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

# Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
