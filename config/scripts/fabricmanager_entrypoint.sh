#!/usr/bin/env bash

set -euo pipefail

: "${DRIVER_VERSION:?Argument DRIVER_VERSION needs to be set and non-empty.}"
: "${TARGET_ARCH:?Argument TARGET_ARCH needs to be set and non-empty.}"

DRIVER_BRANCH=$(echo "$DRIVER_VERSION" | grep -oE '^[0-9]+')

declare -A arch_translation
arch_translation=(["amd64"]="x86_64" ["arm64"]="aarch64")

if [[ ! ${arch_translation[$TARGET_ARCH]+_} ]]; then
    echo "Error: Unsupported TARGET_ARCH value."
    exit 2
fi


mkdir -p /tmp/nvidia

# shellcheck disable=SC2164
pushd /tmp/nvidia

install_packages ca-certificates wget

# Download Fabric Manager tarball
wget -O /tmp/keyring.deb https://developer.download.nvidia.com/compute/cuda/repos/debian12/x86_64/cuda-keyring_1.1-1_all.deb && dpkg -i /tmp/keyring.deb
apt-get update -qq
apt-get install -y --no-install-recommends -V nvidia-fabricmanager-"$DRIVER_BRANCH"="$DRIVER_VERSION"-1

rm -r /var/lib/apt/lists /var/cache/apt/archives
