/**
 * 生成CPU配置的XML
 * @param {*} cfg
 * @returns
 */
function cpuXml(cfg) {
  let topology = ''
  if (!localCfg.isNotManualTopology) {
    topology = `<topology sockets="${localCfg.manualTopology.sockets}" cores="${localCfg.manualTopology.cores}" threads="${localCfg.manualTopology.threads}"/>`
  }
  return `
<vcpu current="${localCfg.cpuCount}">${localCfg.cpuCount}</vcpu>
<cpu mode="${localCfg.cpuMode}">
  ${topology}
</cpu>
`
}

/**
 * 生成内存配置的XML
 * @param {*} cfg
 * @returns
 */
function memoryXml(cfg) {
  return `
  <memory>${cfg.memory}</memory>
  <currentMemory>${cfg.currentMemory}</currentMemory>
`
}

/**
 * 生成磁盘配置的XML
 * @param {*} cfg
 * @returns
 */
function diskXml(cfg) {
  const readonlyTag = cfg.isReadOnly && cfg.diskType === 'cdrom' ? '<readonly/>' : ''
  return `
<disk type="file" device="${cfg.diskType}">
  <driver name="qemu" type="${cfg.diskFormat}" discard="unmap"/>
  <source file="${cfg.sourcePath}"/>
  <target dev="${cfg.targetDev}" bus="${cfg.targetBus}"/>
  ${readonlyTag}
  <boot order="${cfg.bootOrder}"/>
</disk>
`
}
/**
 * 生成显示配置的XML
 * @param {*} cfg
 * @returns
 */
function displayXml(cfg) {
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
function interfaceXml(cfg) {
  let model = ''
  if (cfg.model === 'default') {
    model = `<model type="${cfg.model}"/>`
  }
  return `
<interface type="${cfg.networkType}">
  <source network="${cfg.networkType === 'bridge' ? cfg.bridgeName : cfg.sourceDevice}"/>
  <mac address="${cfg.macAddress}"/>
  ${model}
  <boot order="${cfg.bootOrder}"/>
</interface>
`
}

/**
 * 生成内存配置的XML
 * @param {*} cfg
 * @returns
 */
function memoryXml(cfg) {
  return `
  <memory>${cfg.memory}</memory>
  <currentMemory>${cfg.currentMemory}</currentMemory>
`
}

/**
 * 生成虚拟机概览配置的XML
 * @param {*} cfg
 * @returns
 */
function overviewXml(cfg) {
  return `
<name>${cfg.name}</name>
<uuid>${cfg.uuid}</uuid>
<metadata>
  <libosinfo:libosinfo xmlns:libosinfo="http://libosinfo.org/xmlns/libvirt/domain/1.0">
    <libosinfo:os id="http://ubuntu.com/ubuntu/25.10"/>
  </libosinfo:libosinfo>
</metadata>
<os firmware="${cfg.osFirmware}">
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
function soundXml(cfg) {
  return `<sound model="${localCfg.value.model}"/>`
}

/**
 * 生成显卡配置的XML
 * @param {*} cfg
 * @returns
 */
function videoXml(cfg) {
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
 * 生成CPU配置的XML
 * @param {List<*>} cfg_list
 * @returns
 */
function xml(cfg_list) {
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
  const diskXmls = cfg_list
    .filter((i) => i.type === 'disk')
    .map(diskXml)
    .join('\n')
  const interfaceXmls = cfg_list
    .filter((i) => i.type === 'interface')
    .map(interfaceXml)
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
    .join('\n')
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
    ${interfaceXmls}
    <console type="pty"/>
    <channel type="unix">
      <source mode="bind"/>
      <target type="virtio" name="org.qemu.guest_agent.0"/>
    </channel>
    <channel type="spicevmc">
      <target type="virtio" name="com.redhat.spice.0"/>
    </channel>
    ${displayXmls}
    ${soundXmls}
    ${videoXmls}
    <redirdev bus="usb" type="spicevmc"/>
    <redirdev bus="usb" type="spicevmc"/>
    <memballoon model="virtio"/>
    <rng model="virtio">
      <backend model="random">/dev/urandom</backend>
    </rng>
  </devices>
</domain>
`
}

module.exports = {
  xml,
}
