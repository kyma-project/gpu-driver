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


## Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing Rules](CONTRIBUTING.md).

## Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
