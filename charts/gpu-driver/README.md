# GPU Driver Helm chart

## Installation

Helm repo is not yet established. Clone the repository and `cd charts/gpu-driver`.

```shell
kubectl create ns gpu-driver
helm upgrade --install -n gpu-driver gpu-driver .
```

## Uninstallation

```shell
helm uninstall -n gpu-driver gpu-driver
```

