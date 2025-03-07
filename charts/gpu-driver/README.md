# GPU Driver Helm chart


## Usage

[Helm](https://helm.sh) must be installed to use the charts. Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

```shell
helm repo add kyma-gpu-driver https://kyma-project.github.io/gpu-driver
```

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages. You can then run 
`helm search repo kyma-gpu-driver` to see the charts.

## Installation

To install the `gpu-driver` chart:

### 1. Find the name of the node pool where you want to install the GPU driver

```shell
NODE_POOL="my-node-pool" # this must match you node pool name
```

You can find the node pool name in the node label `worker.gardener.cloud/pool`. 
To get the list of node pool names you can run:

```shell
kubectl get nodes -o jsonpath="{range .items[*]}{.metadata.labels.worker\.gardener\.cloud/pool}{'\n'}{end}" | uniq
```


### 2. Find the kernel version of the nodes in the chosen node pool

```shell
KERNEL_VERSION=$(kubectl get nodes -l worker.gardener.cloud/pool=$NODE_POOL \
  -o jsonpath='{range .items[*]}{.status.nodeInfo.kernelVersion}{end}' | head)
```


### 3. Upgrade/install the helm chart for specified node pool and kernel version

```shell
helm upgrade --install my-gpu-driver kyma-gpu-driver/gpu-driver \
  --set kernelVersion=$KERNEL_VERSION --set nodePool=$NODE_POOL
```

## Uninstall 

To uninstall the chart:

```shell
helm uninstall my-gpu-driver
```
