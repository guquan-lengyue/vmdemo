#!/usr/bin/env bash
# 创建 Intel GVT-g mdev 实例（如果主机支持）
# 用法： sudo ./create_mdev.sh <pci-address> <mdev-type>
# 举例： sudo ./create_mdev.sh 0000:00:02.0 1

if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root (use sudo)."
  exit 1
fi

PCI_ADDR=${1}
MDEV_TYPE=${2}

if [[ -z "${PCI_ADDR}" || -z "${MDEV_TYPE}" ]]; then
  echo "Usage: $0 <pci-address> <mdev-type>"
  echo "Example: $0 0000:00:02.0 1"
  exit 2
fi

MDEV_DIR="/sys/bus/pci/devices/${PCI_ADDR}/mdev_supported_types"

if [[ ! -d "${MDEV_DIR}" ]]; then
  echo "mdev_supported_types not found for ${PCI_ADDR}. Is GVT-g enabled and i915 driver loaded?"
  ls -la /sys/bus/pci/devices/${PCI_ADDR} || true
  exit 3
fi

TYPE_PATH="${MDEV_DIR}/${MDEV_TYPE}"
if [[ ! -d "${TYPE_PATH}" ]]; then
  echo "mdev type ${MDEV_TYPE} doesn't exist in ${MDEV_DIR}"
  ls -la ${MDEV_DIR}
  exit 4
fi

UUID=$(echo ${MDEV_TYPE} > ${TYPE_PATH}/create 2>&1 | tr -d '\n')
if [[ -z "${UUID}" ]]; then
  echo "Failed to create mdev instance. Output: $UUID" >&2
  exit 5
fi

echo "Created mdev UUID: ${UUID}"

echo "To attach to a VM, run:
  virsh attach-device <vm-name> - < /tmp/mdev.xml (or use vmdemo REST API /vm/attach-mdev?name=<vm>&uuid=${UUID})"

cat <<EOF > /tmp/mdev_${UUID}.xml
<hostdev mode=\"subsystem\" type=\"mdev\" managed=\"yes\">
  <source>
    <address uuid=\"${UUID}\"/>
  </source>
</hostdev>
EOF

echo "mdev xml prepared in /tmp/mdev_${UUID}.xml"
exit 0
