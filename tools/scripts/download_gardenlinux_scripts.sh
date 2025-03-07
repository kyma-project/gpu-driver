#!/usr/bin/env bash

set -ex

GARDENLINUX_NVIDIA_INSTALLER_VERSION="${GARDENLINUX_NVIDIA_INSTALLER_VERSION:-0.0.8}"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT_DIR=$( realpath $SCRIPT_DIR/../.. )

ZIPFILE="$(mktemp).zip"
curl -L "https://github.com/gardenlinux/gardenlinux-nvidia-installer/archive/refs/tags/$GARDENLINUX_NVIDIA_INSTALLER_VERSION.zip" -o "$ZIPFILE"

DIR="$(mktemp -d)"

unzip $ZIPFILE -d $DIR

mkdir -p $ROOT_DIR/config/scripts/gardenlinux/

cp -r "$DIR/gardenlinux-nvidia-installer-$GARDENLINUX_NVIDIA_INSTALLER_VERSION/resources/." $ROOT_DIR/charts/gpu-driver-operator/files/gardenlinux-nvidia-installer/
