#!/usr/bin/env bash

set -euo pipefail

echo "Node name: $NODENAME"

KERNEL_VERSION=$(kubectl get node $NODENAME -o jsonpath='{.status.nodeInfo.kernelVersion}')

kubectl label node $NODENAME gpu.kyma-project.io/kernel-version=$KERNEL_VERSION --overwrite

echo "Labeled node $NODENAME with gpu.kyma-project.io/kernel-version=$KERNEL_VERSION"

sleep Infinity & wait
