#!/usr/bin/env bash

echo "GPU driver installer entrypoint"

set -euo pipefail

if ${DEBUG}; then
  set -x
  env
fi

: "${KERNEL_TYPE:?Build argument KERNEL_TYPE needs to be set and non-empty.}"
: "${TARGET_ARCH:?Build argument TARGET_ARCH needs to be set and non-empty.}"
: "${DRIVER_VERSION:?Build argument DRIVER_VERSION needs to be set and non-empty.}"
: "${HOST_DRIVER_PATH:?Build argument HOST_DRIVER_PATH needs to be set and non-empty.}"

log() {
  local NAME="kyma-gpu-driver"
  local MSG=$1
  local TS=$(date '+%Y-%m-%d %H:%M:%S')
  echo "${TS} ${NAME}: $MSG"
}

export KERNEL_NAME=$(./extract_kernel_name.sh ${KERNEL_TYPE} ${TARGET_ARCH})

COMPILED_FILENAME="${DRIVER_VERSION}"

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
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

log "Downloading Fabric Manager"

./download_fabricmanager.sh


log "Removing dependency packages"

apt-get remove -y --autoremove --allow-remove-essential --ignore-hold \
    libgnutls30 apt openssl wget ncurses-base ncurses-bin

rm -rf /var/lib/apt/lists/* /usr/bin/dpkg /sbin/start-stop-daemon /usr/lib/x86_64-linux-gnu/libsystemd.so* \
    /var/lib/dpkg/info/libdb5.3* /usr/lib/x86_64-linux-gnu/libdb-5.3.so* /usr/share/doc/libdb5.3 \
    /usr/bin/chfn /usr/bin/gpasswd

log "Installing GPU driver"

mkdir -p /opt/nvidia-installer
cp *.sh /opt/nvidia-installer

/opt/nvidia-installer/load_install_gpu_driver.sh


log "Run fabric manager"

/opt/nvidia-installer/install_fabricmanager.sh

echo "Sleep infinity"
sleep infinity
