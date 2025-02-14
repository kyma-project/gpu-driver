#!/bin/bash

BIN_DIR=${BIN_DIR:-/opt/nvidia-installer}

driver() {
  echo "Driver"
  source "${BIN_DIR}/load_install_gpu_driver.sh"
}

fabricManager() {
  echo "FabricManager"
  source "${BIN_DIR}/install_fabricmanager.sh"
}

case "$1" in
  "--driver" )
    driver
    ;;
  "--fabricManager" )
    fabricManager
    ;;
  "" )
    echo "Sleep..."
    sleep 1000d
    ;;
  * )
    echo "Error: Unknown argument $1"
    exit 1
    ;;
esac
