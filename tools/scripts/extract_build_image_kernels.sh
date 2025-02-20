#!/usr/bin/env bash

IMAGE=ghcr.io/gardenlinux/gardenlinux/kmodbuild

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

TAGS=$(docker image ls $IMAGE | awk '{ print $2 }' | tail -n +2)

while IFS= read -r TAG; do
    KERNEL=$(docker run --platform linux/amd64 --rm -v $SCRIPT_DIR/../../charts/gpu-driver/files/gardenlinux-nvidia-installer:/mnt/scripts $IMAGE:$TAG /mnt/scripts/extract_kernel_name.sh cloud)
    echo "  $KERNEL: $TAG"
done <<< "$TAGS"
