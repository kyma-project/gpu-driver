# GPU Driver Operator Helm chart

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

To install the `gpu-driver-operator` chart:

```shell
helm upgrade --install gpu-driver kyma-gpu-driver/gpu-driver-operator
```

## Upgrade

When new version of gpu-driver-operator is released, to upgrade the operator you should first update
the gpu-driver helm repository and fetch its new version. 

```shell
helm repo update
```

Now, you can upgrade existing installation with

```shell
helm upgrade --install gpu-driver kyma-gpu-driver/gpu-driver-operator
```

## Uninstall

To uninstall the gpu-driver-operator you should run:

```shell
helm uninstall gpu-driver
```
