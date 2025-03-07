#!/usr/bin/env bash

set -euo pipefail

if ${DEBUG}; then
  set -x
  env
fi

: "${KERNEL_TYPE:?Argument KERNEL_TYPE needs to be set and non-empty.}"
: "${TARGET_ARCH:?Argument TARGET_ARCH needs to be set and non-empty.}"
: "${KERNEL_NAME:?Argument KERNEL_NAME needs to be set and non-empty.}"
: "${DRIVER_VERSION:?Argument DRIVER_VERSION needs to be set and non-empty.}"
: "${HOST_DRIVER_PATH:?Argument HOST_DRIVER_PATH needs to be set and non-empty.}"
: "${NODE_NAME:?Argument NODE_NAME needs to be set and non-empty.}"

LABEL_GPU_NAME="${LABEL_GPU_NAME:-gpu.kyma-project.io/gpu-name}"

log() {
  set +x
  local NAME="kyma-gpu-driver"
  local MSG=$1
  local TS
  TS=$(date '+%Y-%m-%d %H:%M:%S')
  echo "${TS} ${NAME}: $MSG"
  if ${DEBUG}; then
    set -x
  fi
}

log "GPU driver installer entrypoint"

find /dev -name "nvidia*"

SHOULD_INSTALL=true

DEVICE_EXISTS=false
if [ -e /dev/nvidia0 ] && [ -e /dev/nvidiactl ] && [ -e /dev/nvidia-uvm ]; then
  DEVICE_EXISTS=true
fi

if [ "$DEVICE_EXISTS" = true ]; then
  SHOULD_INSTALL=false
  NVIDIA_ROOT="${NVIDIA_ROOT:-/opt/nvidia-installer/cache/nvidia/$DRIVER_VERSION}"
  export NVIDIA_ROOT
  LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$NVIDIA_ROOT/lib}"
  export LD_LIBRARY_PATH

  echo $NVIDIA_ROOT
  echo $LD_LIBRARY_PATH

  GPU_NAME=$(/opt/nvidia-installer/cache/nvidia/550.127.08/bin/nvidia-smi -i 0 --query-gpu=name --format=csv,noheader | tr -d '\n' || echo '')

  # NVIDIA-SMI has failed
  if [[ $GPU_NAME == *"NVIDIA-SMI"* ]]; then
    echo "$GPU_NAME"
    echo "Reinstalling the driver"
    find /dev -name "nvidia*"
    find /dev -name "nvidia*" -exec rm -rf {} \;
    SHOULD_INSTALL=true
  else
    echo "Detected GPU device name: $GPU_NAME"
  fi
fi

if [ "$SHOULD_INSTALL" = false ]; then

  log "GPU driver already installed"

else

  log "Compiling the GPU driver"

  ./compile.sh

  ls -la /out
  ls -la /out/nvidia


  log "Installing dependency packages"

  apt-get update && apt-get install --no-install-recommends -y \
      kmod \
      pciutils \
      ca-certificates \
      wget \
      curl \
      xz-utils \
      && rm -rf /var/lib/apt/lists/*

  log "Installing GPU driver"

  mkdir -p /opt/nvidia-installer
  cp ./*.sh /opt/nvidia-installer

  source /opt/nvidia-installer/load_install_gpu_driver.sh

fi



NVIDIA_ROOT="${NVIDIA_ROOT:-/opt/nvidia-installer/cache/nvidia/$DRIVER_VERSION}"
LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$NVIDIA_ROOT/lib}"

echo "NVIDIA_ROOT: $NVIDIA_ROOT"
echo "LD_LIBRARY_PATH: $LD_LIBRARY_PATH"

GPU_NAME=$(/opt/nvidia-installer/cache/nvidia/550.127.08/bin/nvidia-smi -i 0 --query-gpu=name --format=csv,noheader | tr -d '\n' || echo '')

echo "GPU_NAME=$GPU_NAME"

cd /opt/nvidia-installer

echo "Downloading kubectl"

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

chmod +x ./kubectl

GPU_NAME_LABEL="${GPU_NAME// /-}"

echo "Labeling node $NODE_NAME with $LABEL_GPU_NAME = $GPU_NAME_LABEL"

./kubectl label node $NODE_NAME $LABEL_GPU_NAME="$GPU_NAME_LABEL" --overwrite


if ${DEBUG}; then
  echo "Entrypoint sleep infinity"
  sleep infinity
fi

