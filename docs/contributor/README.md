# GPU Driver

Based on the Garden Linux NVIDIA Installer with a difference that driver is not pre-compiled and distributed
in a container image due to licensing issues, but compiled in place during runtime.

## Flow

Defaults
* KERNEL_TYPE = cloud
* TARGET_ARCH = amd64

Inputs
* DRIVER_VERSION
  * 550.127.08

Run flow:
* run ??? and detect Garden Linux version from the node
  * outputs
    * GARDENLINUX_VERSION
* run ghcr.io/gardenlinux/gardenlinux/kmodbuild:GARDENLINUX_VERSION
  * run `compile.sh`
  * inputs
    * KERNEL_TYPE
    * DRIVER_VERSION
    * TARGET_ARCH
  * outputs
    * `/out`
* run debian:bookworm-slim, install dependencies and build `/rootfs`
  * install deps
  * copy deps to `/rootfs`
  * inputs
    * `/out`
  * outputs
    * `/rootfs`
* run scratch and copy `/rootfs` to `/`
  * run `load_install_gpu_driver.sh` 
  * label the node with `ai.sap.com/nvidia-driver-version=DRIVER_VERSION`
  * run `install_fabricmanager_no_sleep.sh`
  * run `/usr/bin/nvidia-gpu-device-plugin`
