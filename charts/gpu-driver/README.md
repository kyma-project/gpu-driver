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

To install the `gpu-driver` chart:

```shell
helm install my-gpu-driver kyma-gpu-driver/gpu-driver
```

To uninstall the chart:

```shell
helm uninstall my-gpu-driver
```
