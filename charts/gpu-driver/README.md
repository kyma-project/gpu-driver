# GPU Driver Helm chart

## Installation

```shell
kubectl create ns gpu-driver
helm upgrade --install -n gpu-driver gpu-driver .
```

## Uninstallation

```shell
helm uninstall -n gpu-driver gpu-driver
```
