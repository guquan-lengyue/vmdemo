<template>
  <div class="vm-opts-container">
    <h1>虚拟机操作</h1>

    <!-- 通用操作区域 -->
    <div class="section">
      <h2>通用操作</h2>
      <div class="form-group">
        <label for="vmName">虚拟机名称:</label>
        <input type="text" id="vmName" v-model="vmName" placeholder="请输入虚拟机名称" />
      </div>
      <div class="button-group">
        <button @click="performAction('start')">启动</button>
        <button @click="performAction('stop')">停止</button>
        <button @click="performAction('suspend')">挂起</button>
        <button @click="performAction('resume')">恢复</button>
        <button @click="performAction('forceShutdown')">强制关闭</button>
        <button @click="performAction('delete')">删除</button>
        <button @click="performAction('getInfo')">获取信息</button>
        <div class="form-group">
          <label for="usbId">USB设备ID:</label>
          <input type="text" id="usbId" v-model="usbId" placeholder="请输入USB设备ID" />
        </div>
        <div>
          <button @click="handleAttachUsbClick">挂载USB设备</button>
          <button @click="handleDetachUsbClick">卸载USB设备</button>
        </div>
      </div>
    </div>

    <!-- 虚拟机列表 -->
    <div class="section">
      <h2>虚拟机列表</h2>
      <div class="form-group">
        <label for="listType">列表类型:</label>
        <select id="listType" v-model="listType">
          <option value="all">所有</option>
          <option value="active">活动</option>
          <option value="inactive">非活动</option>
        </select>
        <button @click="fetchVMList">刷新列表</button>
      </div>
      <div v-if="vmList.length > 0" class="vm-list">
        <ul>
          <li v-for="vm in vmList" :key="vm.name">{{ vm.name }} ({{ vm.state }})</li>
        </ul>
      </div>
      <p v-else>没有找到虚拟机。</p>
    </div>

    <!-- 创建虚拟机 -->
    <div class="section">
      <h2>创建新虚拟机</h2>
      <div class="form-group">
        <label for="newVmName">新虚拟机名称:</label>
        <input type="text" id="newVmName" v-model="newVmName" placeholder="请输入新虚拟机名称" />
      </div>
      <div class="form-group">
        <label for="newVmCPU">CPU 核心数:</label>
        <input
          type="number"
          id="newVmCPU"
          v-model.number="newVmCPU"
          placeholder="请输入 CPU 核心数"
        />
      </div>
      <div class="form-group">
        <label for="newVmMemory">内存 (GB):</label>
        <input
          type="number"
          id="newVmMemory"
          v-model.number="newVmMemory"
          placeholder="请输入内存大小 (GB)"
        />
      </div>
      <div class="form-group">
        <label for="newVmDiskPath">磁盘路径:</label>
        <select id="newVmDiskPath" v-model="newVmDiskPath">
          <option v-for="disk in availableDisks" :key="disk.name" :value="disk.path">
            {{ disk.name }} - {{ disk.size }}
          </option>
        </select>
      </div>
      <div class="form-group">
        <label for="newVmIsoPath">ISO 路径:</label>
        <input type="text" id="newVmIsoPath" v-model="newVmIsoPath" placeholder="请输入 ISO 路径" />
      </div>
      <button @click="createVMFromForm">创建虚拟机</button>
    </div>

    <!-- 结果显示 -->
    <div class="section result-section">
      <h2>操作结果</h2>
      <div v-if="loading" class="loading">处理中...</div>
      <div v-if="error" class="error-message">
        <p>错误:</p>
        <pre>{{ error }}</pre>
      </div>
      <div v-if="result" class="result-message">
        <p>成功:</p>
        <pre>{{ result }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  startVM,
  stopVM,
  suspendVM,
  resumeVM,
  forceShutdownVM,
  deleteVM,
  listVMs,
  getVMInfo,
  createVM,
  attachUsb,
  detachUsb,
} from '@/api/vmopts.js'
import { listDisks } from '@/api/diskopts.js'

const vmName = ref('')
const newVmName = ref('my-new-vm')
const newVmCPU = ref(2)
const newVmMemory = ref(4) // in GiB
const newVmDiskPath = ref('')
const newVmIsoPath = ref('/home/k/下载/ubuntu-25.10-desktop-amd64.iso')
const listType = ref('all')
const vmList = ref([])
const availableDisks = ref([])
const result = ref(null)
const error = ref(null)
const loading = ref(false)
const usbId = ref('')

const actions = {
  start: startVM,
  stop: stopVM,
  suspend: suspendVM,
  resume: resumeVM,
  forceShutdown: forceShutdownVM,
  delete: deleteVM,
  getInfo: getVMInfo,
}

async function handleAttachUsbClick() {
  if (!usbId.value) {
    error.value = '请输入USB设备ID。'
    result.value = null
    return
  }
  try {
    const response = await attachUsb(vmName.value, usbId.value)
    result.value = response
  } catch (e) {
    error.value = e.message
    result.value = null
  }
}

async function handleDetachUsbClick() {
  if (!usbId.value) {
    error.value = '请输入USB设备ID。'
    result.value = null
    return
  }
  try {
    const response = await detachUsb(vmName.value, usbId.value)
    result.value = response
  } catch (e) {
    error.value = e.message
    result.value = null
  }
}

async function performAction(action) {
  if (!vmName.value && action !== 'list') {
    error.value = '请输入虚拟机名称。'
    result.value = null
    return
  }
  loading.value = true
  error.value = null
  result.value = null
  try {
    const response = await actions[action](vmName.value)
    result.value = response
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
    if (action === 'start' || action === 'stop' || action === 'delete') {
      fetchVMList() // Refresh list after state-changing actions
    }
  }
}

async function fetchVMList() {
  loading.value = true
  error.value = null
  result.value = null
  try {
    const response = await listVMs(listType.value)
    vmList.value = response.vms || []
  } catch (e) {
    error.value = e.message
    vmList.value = []
  } finally {
    loading.value = false
  }
}

async function fetchAvailableDisks() {
  try {
    const response = await listDisks()
    availableDisks.value = response.disks || []
    if (availableDisks.value.length > 0) {
      newVmDiskPath.value = availableDisks.value[0].path
    }
  } catch (e) {
    error.value = '获取可用磁盘列表失败: ' + e.message
  }
}

async function createVMFromForm() {
  if (
    !newVmName.value ||
    !newVmCPU.value ||
    !newVmMemory.value ||
    !newVmIsoPath.value ||
    !newVmDiskPath.value
  ) {
    error.value = '请填写所有虚拟机创建字段。'
    result.value = null
    return
  }

  const memoryInKiB = newVmMemory.value * 1024 * 1024

  const xml = `
<domain type='kvm'>
  <name>${newVmName.value}</name>
  <memory unit='KiB'>${memoryInKiB}</memory>
  <vcpu placement='static'>${newVmCPU.value}</vcpu>
  <os firmware='efi'>
    <type arch='x86_64' machine='pc-q35-8.2'>hvm</type>
    <firmware>
      <feature enabled='yes' name='enrolled-keys'/>
      <feature enabled='yes' name='secure-boot'/>
    </firmware>
    <loader readonly='yes' secure='yes' type='pflash'>/usr/share/OVMF/OVMF_CODE_4M.ms.fd</loader>
    <nvram template='/usr/share/OVMF/OVMF_VARS_4M.ms.fd'/>
  </os>
  <features>
    <acpi/>
    <apic/>
    <smm state='on'/>
  </features>
  <cpu mode='host-passthrough' check='none' migratable='on'/>
  <clock offset='utc'>
    <timer name='rtc' tickpolicy='catchup'/>
    <timer name='pit' tickpolicy='delay'/>
    <timer name='hpet' present='no'/>
  </clock>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <pm>
    <suspend-to-mem enabled='no'/>
    <suspend-to-disk enabled='no'/>
  </pm>
  <devices>
    <emulator>/usr/bin/qemu-system-x86_64</emulator>
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2' discard='unmap'/>
      <source file='${newVmDiskPath.value}'/>
      <target dev='vda' bus='virtio'/>
      <boot order='2'/>
    </disk>
    <disk type='file' device='cdrom'>
      <driver name='qemu' type='raw'/>
      <source file='${newVmIsoPath.value}'/>
      <target dev='sda' bus='sata'/>
      <readonly/>
      <boot order='1'/>
    </disk>
    <controller type='usb' index='0' model='qemu-xhci' ports='15'/>
    <controller type='pci' index='0' model='pcie-root'/>
    <interface type='network'>
      <source network='default'/>
      <model type='virtio'/>
    </interface>
    <serial type='pty'>
      <target type='isa-serial' port='0'>
        <model name='isa-serial'/>
      </target>
    </serial>
    <console type='pty'>
      <target type='serial' port='0'/>
    </console>
    <input type='mouse' bus='ps2'/>
    <input type='keyboard' bus='ps2'/>
    <graphics type='vnc' port='-1' autoport='yes' listen='0.0.0.0'>
      <listen type='address' address='0.0.0.0'/>
    </graphics>
    <video>
      <model type='virtio' heads='1' primary='yes'/>
    </video>
    <memballoon model='virtio'/>
    <rng model='virtio'>
      <backend model='random'>/dev/urandom</backend>
    </rng>
  </devices>
</domain>
  `

  loading.value = true
  error.value = null
  result.value = null
  try {
    const response = await createVM(newVmName.value, xml)
    result.value = response
    fetchVMList() // Refresh list after creating a new VM
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchVMList()
  fetchAvailableDisks()
})
</script>

<style scoped>
.vm-opts-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: sans-serif;
}

.section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
}

input[type='text'],
select,
textarea {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.button-group button {
  margin-right: 10px;
  margin-bottom: 10px;
}

button {
  padding: 10px 15px;
  border: none;
  background-color: #007bff;
  color: white;
  border-radius: 5px;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.vm-list ul {
  list-style-type: none;
  padding: 0;
}

.vm-list li {
  padding: 5px 0;
}

.result-section .loading,
.result-section .error-message,
.result-section .result-message {
  margin-top: 15px;
  padding: 10px;
  border-radius: 5px;
}

.loading {
  color: #555;
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.result-message {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
