/**
 * 生成CPU配置的XML
 * @param {*} cfg
 * @returns
 */
function cpuXml({ cfg }) {
  let topology = ''
  if (!cfg.isNotManualTopology) {
    topology = `<topology sockets="${cfg.manualTopology.sockets}" cores="${cfg.manualTopology.cores}" threads="${cfg.manualTopology.threads}"/>`
  }
  return `
<vcpu current="${cfg.cpuCount}">${cfg.cpuCount}</vcpu>
<cpu mode="${cfg.cpuMode}">
  ${topology}
</cpu>
`
}

/**
 * 生成磁盘配置的XML
 * @param {*} cfg
 * @returns
 */
function diskXml({ cfg }) {
  // 只对CDROM类型添加只读标签
  const readonlyTag = cfg.diskType === 'cdrom' ? '<readonly/>' : ''
  // 只对CDROM类型添加可弹出标签
  const removableTag = cfg.diskType === 'cdrom' ? '<removable/>' : ''
  // 只有当bootOrder大于0时才添加启动顺序
  const bootOrderTag = cfg.bootOrder !== undefined && cfg.bootOrder > 0 ? `<boot order="${cfg.bootOrder}"/>` : ''
  return `
<disk type="file" device="${cfg.diskType}">
  <driver name="qemu" type="${cfg.diskFormat}"${cfg.diskType !== 'cdrom' ? ' discard="unmap"' : ''}/>
  <source file="${cfg.sourcePath}"/>
  <target dev="${cfg.targetDev}" bus="${cfg.targetBus}"/>
  ${readonlyTag}
  ${removableTag}
  ${bootOrderTag}
</disk>
`
}
/**
 * 生成显示配置的XML
 * @param {*} cfg
 * @returns
 */
function displayXml({ cfg }) {
  return `
<graphics type="${cfg.type}" port="${cfg.port}" listen="${cfg.listen}" password="${cfg.passwd}">
  <gl enable="no"/>
  <image compression="${cfg.imageCompression}"/>
</graphics>
  `
}

/**
 * 生成网络接口配置的XML
 * @param {*} cfg
 * @returns
 */
function interfaceXml({ cfg }) {
  let model = ''
  if (cfg.model === 'default') {
    model = `<model type="${cfg.model}"/>`
  }
  // 只有当bootOrder大于0时才添加启动顺序
  const bootOrderTag = cfg.bootOrder !== undefined && cfg.bootOrder > 0 ? `<boot order="${cfg.bootOrder}"/>` : ''
  return `
<interface type="${cfg.networkType}">
  <source network="${cfg.networkType === 'bridge' ? cfg.bridgeName : cfg.netName}"/>
  <mac address="${cfg.macAddress}"/>
  ${model}
  ${bootOrderTag}
</interface>
`
}

/**
 * 生成内存配置的XML
 * @param {*} cfg
 * @returns
 */
function memoryXml({ cfg }) {
  // 确保内存值存在且为正数
  const memory = cfg.memory && cfg.memory > 0 ? cfg.memory : 2048
  const currentMemory = cfg.currentMemory && cfg.currentMemory > 0 ? cfg.currentMemory : 2048

  // XML中的内存值通常是以KB为单位的，而前端的内存配置是以MB为单位的，所以需要进行单位转换
  return `
  <memory unit="MiB">${memory}</memory>
  <currentMemory unit="MiB">${currentMemory}</currentMemory>
`
}

/**
 * 生成虚拟机概览配置的XML
 * @param {*} cfg
 * @returns
 */
function overviewXml({ cfg }) {
  let osFirmwareTag = ''
  if (cfg.osFirmware === 'uefi') {
    osFirmwareTag = `firmware="efi"`
  }
  return `
<name>${cfg.name}</name>
<uuid>${cfg.uuid}</uuid>
<metadata>
  <libosinfo:libosinfo xmlns:libosinfo="http://libosinfo.org/xmlns/libvirt/domain/1.0">
    <libosinfo:os id="http://ubuntu.com/ubuntu/25.10"/>
  </libosinfo:libosinfo>
</metadata>
<os ${osFirmwareTag}>
  <type arch="x86_64" machine="${cfg.osMachine}">hvm</type>
</os>
<features>
  <acpi/>
  <apic/>
  <vmport state="off"/>
</features>
<clock offset="utc">
  <timer name="rtc" tickpolicy="catchup"/>
  <timer name="pit" tickpolicy="delay"/>
  <timer name="hpet" present="no"/>
</clock>
`
}

/**
 * 生成声卡配置的XML
 * @param {*} cfg
 * @returns
 */
function soundXml({ cfg }) {
  return `<sound model="${cfg.model}"/>`
}

/**
 * 生成显卡配置的XML
 * @param {*} cfg
 * @returns
 */
function videoXml({ cfg }) {
  let accel3d = ''
  if (cfg.model.acceleration.accel3d === 'yes' && cfg.model.type === 'virtio') {
    accel3d = `<acceleration accel3d="${cfg.model.acceleration.accel3d}"/>`
  }
  return `
<video>
  <model type="${cfg.model.type}">
    ${accel3d}
  </model>
</video>
`
}

/**
 * 生成USB设备配置的XML
 * @param {*} cfg
 * @returns
 */
function usbXml({ cfg }) {
  if (!cfg.attachedDevices || cfg.attachedDevices.length === 0) {
    return ''
  }

  // 为每个附加的USB设备生成XML配置
  return cfg.attachedDevices.map(deviceId => {
    const [vendorId, productId] = deviceId.split(':')
    return `
<hostdev mode="subsystem" type="usb" managed="yes">
  <source>
    <vendor id="0x${vendorId}"/>
    <product id="0x${productId}"/>
  </source>
</hostdev>
`
  }).join('')
}

/**
 * 生成PCI设备配置的XML
 * @param {*} cfg
 * @returns
 */
function pciXml({ cfg }) {
  if (!cfg.attachedDevices || cfg.attachedDevices.length === 0) {
    return ''
  }

  // 为每个附加的PCI设备生成XML配置
  return cfg.attachedDevices.map((deviceId, index) => {
    // 解析PCI设备ID，格式如：0000:00:02.0
    const parts = deviceId.split(':')
    if (parts.length !== 3) {
      return ''
    }
    const domain = parts[0]
    const bus = parts[1]
    const slotFunc = parts[2]
    const [slot, func] = slotFunc.split('.')
    
    return `
<hostdev mode="subsystem" type="pci" managed="yes">
  <source>
    <address domain="0x${domain}" bus="0x${bus}" slot="0x${slot}" function="0x${func}"/>
  </source>
  <alias name="hostdev${index}"/>
  <address type="pci" domain="0x0000" bus="0x00" slot="0x${(index + 5).toString(16).padStart(2, '0')}" function="0x0" multifunction="on"/>
</hostdev>
`
  }).join('')
}

/**
 * 生成CPU配置的XML
 * @param {List<*>} cfg_list
 * @returns
 */
function xml(cfg_list) {
  // 跟踪已使用的启动顺序
  const usedBootOrders = new Set()

  // 处理磁盘XML，确保启动顺序唯一
  const diskXmls = cfg_list
    .filter((i) => i.type === 'disk')
    .map(item => {
      if (item.cfg.bootOrder > 0 && usedBootOrders.has(item.cfg.bootOrder)) {
        // 如果启动顺序已被使用，将其设置为0（不添加启动顺序标签）
        item.cfg.bootOrder = 0
      } else if (item.cfg.bootOrder > 0) {
        usedBootOrders.add(item.cfg.bootOrder)
      }
      return diskXml(item)
    })
    .join('\n')

  // 处理网络接口XML，确保启动顺序唯一
  const interfaceXmls = cfg_list
    .filter((i) => i.type === 'interface')
    .map(item => {
      if (item.cfg.bootOrder > 0 && usedBootOrders.has(item.cfg.bootOrder)) {
        // 如果启动顺序已被使用，将其设置为0（不添加启动顺序标签）
        item.cfg.bootOrder = 0
      } else if (item.cfg.bootOrder > 0) {
        usedBootOrders.add(item.cfg.bootOrder)
      }
      return interfaceXml(item)
    })
    .join('\n')

  const overviewXmls = cfg_list
    .filter((i) => i.type === 'overview')
    .map(overviewXml)
    .join('\n')
  const cpuXmls = cfg_list
    .filter((i) => i.type === 'cpu')
    .map(cpuXml)
    .join('\n')
  const memoryXmls = cfg_list
    .filter((i) => i.type === 'memory')
    .map(memoryXml)
    .join('\n')
  const displayXmls = cfg_list
    .filter((i) => i.type === 'display')
    .map(displayXml)
    .join('\n')
  const soundXmls = cfg_list
    .filter((i) => i.type === 'sound')
    .map(soundXml)
    .join('\n')
  const videoXmls = cfg_list
    .filter((i) => i.type === 'video')
    .map(videoXml)
    .join('')
  const usbXmls = cfg_list
    .filter((i) => i.type === 'usb')
    .map(usbXml)
    .join('')
  const pciXmls = cfg_list
    .filter((i) => i.type === 'pci')
    .map(pciXml)
    .join('')
  return `
<domain type="kvm">
  ${overviewXmls}
  ${cpuXmls}
  ${memoryXmls}
  <pm>
    <suspend-to-mem enabled="no"/>
    <suspend-to-disk enabled="no"/>
  </pm>
  <devices>
    <emulator>/usr/bin/qemu-system-x86_64</emulator>
    ${diskXmls}
    <controller type="usb" model="qemu-xhci" ports="15"/>
    <controller type="pci" model="pcie-root"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="pci" model="pcie-root-port"/>
    <controller type="scsi" index="0" model="virtio-scsi"/>
    ${interfaceXmls}
    <console type="pty"/>
    <channel type="unix">
      <source mode="bind"/>
      <target type="virtio" name="org.qemu.guest_agent.0"/>
    </channel>
    ${displayXmls}
    ${soundXmls}
    ${videoXmls}
    ${usbXmls}
    ${pciXmls}
    <memballoon model="virtio"/>
    <rng model="virtio">
      <backend model="random">/dev/urandom</backend>
    </rng>
  </devices>
</domain>
`
}

export { xml }