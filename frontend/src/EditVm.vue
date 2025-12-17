<template>
  <div class="vm-cfg-container">
    <div class="vm-cfg">
      <div class="menu-header">
        <h2>{{ props.vmName ? '编辑虚拟机' : '创建虚拟机' }}</h2>
        <button class="install-btn" @click="handleInstall">{{ props.vmName ? '保存修改' : '创建虚拟机' }}</button>
        <button class="cancel-btn" @click="handleCancel">取消</button>
      </div>
      <div class="menu-list">
        <div
          v-for="(item, index) in btnGroup"
          :key="index"
          class="menu-item"
          @click="handleMenuClick(item, index)"
        >
          <span class="menu-icon"></span>
          <span class="menu-text">{{ item.name }}</span>
          <!-- 删除按钮，只对用户添加的硬件显示 -->
          <span
            v-if="isRemovableHardware(index)"
            class="delete-btn"
            @click.stop="deleteHardware(index)"
          >
            ×
          </span>
        </div>
        <!-- 添加硬件按钮 -->
        <div class="add-hardware-btn" @click="showAddHardwareModal = true">
          <span class="menu-icon">+</span>
          <span class="menu-text">添加硬件</span>
        </div>
      </div>
    </div>
    <div class="vm-content">
      <Component
        :is="componentMap[selectedMenu] || componentMap['overview']"
        :hostMsg="hostMsg"
        :cfg="getCurrentCfg()"
        @update:cfg="updateComponentCfg"
        @update:menuName="updateMenuName"
      />
    </div>

    <!-- 添加硬件模态框 -->
    <div v-if="showAddHardwareModal" class="modal-overlay" @click="showAddHardwareModal = false">
      <div class="modal-content" @click.stop>
        <h3>选择硬件类型</h3>
        <div class="hardware-types">
          <div
            v-for="type in availableHardwareTypes"
            :key="type.value"
            class="hardware-type-item"
            @click="addHardware(type)"
          >
            {{ type.label }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, provide, watch, onMounted } from 'vue'
import Overview from './components/Overview.vue'
import CPU from './components/CPU.vue'
import Memory from './components/Memery.vue'
import Disk from './components/Disk.vue'
import Interface from './components/Interface.vue'
import Display from './components/Display.vue'
import Sound from './components/Sound.vue'
import Video from './components/Video.vue'
import { xml } from './utils/xml'
import { vmApi } from './api.js'

// 接收props
const props = defineProps({
  vmName: {
    type: String,
    required: false,
    default: ''
  }
})

// 定义事件
const emit = defineEmits(['vm-updated'])

const selectedMenu = ref('overview')

// 组件映射对象
const componentMap = {
  overview: Overview,
  cpu: CPU,
  memory: Memory,
  disk: Disk,
  interface: Interface,
  display: Display,
  sound: Sound,
  video: Video,
}

// 主机信息
const hostMsg = ref({
  hostCpuCount: 16,
  hostMemory: 4096,
  netNames: ['default'], //  虚拟网卡源列表
})

// 使用ref定义按钮组 - 清空所有组件的默认cfg配置
const btnGroup = ref([
  {
    cfg: {
      name: 'ubuntu25.10',
      uuid: 'f0715790-cde5-4faf-8869-d8d72dfaf7d8',
      osMachine: 'q35',
      osFirmware: 'bios',
    },
    name: '概况',
    type: 'overview',
  },
  {
    cfg: {
      cpuCount: 2,
      cpuMode: 'host-passthrough',
      isNotManualTopology: false,
      manualTopology: {
        sockets: 2,
        cores: 1,
        threads: 1,
      },
    },
    name: 'CPU数',
    type: 'cpu',
  },
  {
    cfg: {
      memory: 4096,
      currentMemory: 4096,
    },
    name: '内存',
    type: 'memory',
  },
  {
    cfg: {
      diskType: 'disk',
      sourcePath: '/var/lib/libvirt/images/ubuntu25.10-1.qcow2',
      diskFormat: 'qcow2',
      targetDev: 'vda',
      targetBus: 'virtio',
      isReadOnly: false,
      bootOrder: 0,
    },
    name: 'VirtIO-磁盘',
    type: 'disk',
  },
  {
    cfg: {
      networkType: 'network',
      netName: 'default',
      bridgeName: '',
      macAddress: '52:54:00:b5:b3:e7',
      model: 'virtio',
    },
    name: '虚拟网络',
    type: 'interface',
  },
  {
    cfg: {
      type: 'vnc',
      port: '-1',
      listen: '0.0.0.0',
      passwd: '',
      imageCompression: 'off',
    },
    name: '显示协议-VNC',
    type: 'display',
  },
  {
    cfg: {
      model: 'ac97',
    },
    name: '声音',
    type: 'sound',
  },
  {
    cfg: {
      model: {
        type: 'virtio',
        acceleration: {
          accel3d: 'yes',
        },
      },
    },
    name: '视频',
    type: 'video',
  },
])

// 系统默认菜单项的数量
const defaultMenuCount = 3

// 可添加的硬件类型
const availableHardwareTypes = [
  { label: '网络接口', value: 'interface' },
  { label: '磁盘', value: 'disk' },
  { label: '显示', value: 'display' },
  { label: '视频', value: 'video' },
  { label: '声音', value: 'sound' },
]

// 跟踪当前选中的菜单项索引
const currentMenuIndex = ref(-1)

// 模态框显示状态
const showAddHardwareModal = ref(false)

// 硬件计数，用于生成唯一名称
const hardwareCounts = ref({
  interface: 1,
  disk: 1,
  display: 1,
  video: 1,
  sound: 1,
})

// 菜单点击事件处理
const handleMenuClick = (item, index) => {
  selectedMenu.value = item.type
  currentMenuIndex.value = index
}

// 获取当前选中组件的配置
const getCurrentCfg = () => {
  if (currentMenuIndex.value === -1) return undefined

  const currentItem = btnGroup.value[currentMenuIndex.value]
  // 如果cfg为空对象，返回undefined，这样组件会使用自身的默认值
  return currentItem && Object.keys(currentItem.cfg).length > 0 ? currentItem.cfg : undefined
}

// 更新组件配置
const updateComponentCfg = (newCfg) => {
  if (currentMenuIndex.value !== -1) {
    btnGroup.value[currentMenuIndex.value].cfg = newCfg
  }
}

// 更新菜单名称
const updateMenuName = (newName) => {
  if (currentMenuIndex.value !== -1) {
    btnGroup.value[currentMenuIndex.value].name = newName
  }
}

// 将btnGroup提供给子组件
provide('btnGroup', btnGroup.value)

// 组件挂载时获取虚拟机配置
onMounted(async () => {
  if (props.vmName && props.vmName !== 'new') {
    try {
      const vmInfo = await vmApi.getVMInfo(props.vmName)
      // 这里需要解析vmInfo中的XML配置并更新btnGroup
      // 由于没有具体的XML解析逻辑，目前保持默认配置
      console.log('获取到虚拟机配置:', vmInfo)
    } catch (error) {
      console.error('获取虚拟机配置失败:', error)
    }
  } else if (props.vmName === 'new') {
    // 如果是新建虚拟机，初始化一个默认名称
    const overviewItem = btnGroup.value.find(item => item.type === 'overview')
    if (overviewItem) {
      overviewItem.cfg.name = `vm-${Date.now()}`
    }
  }
})

function handleInstall() {
  const xmlStr = xml(btnGroup.value)
  console.log('生成的XML配置:', xmlStr)

  // 保存配置到后端
  saveVMConfig(xmlStr)
}

// 取消操作
function handleCancel() {
  emit('vm-updated')
}

// 保存虚拟机配置到后端
async function saveVMConfig(xmlConfig) {
  try {
    if (props.vmName) {
      // 更新现有虚拟机配置
      await vmApi.updateVM(props.vmName, xmlConfig)
      emit('vm-updated')
    } else {
      // 创建新虚拟机
      const vmName = btnGroup.value.find(item => item.type === 'overview')?.cfg?.name || 'new-vm'
      await vmApi.createVM(vmName, xmlConfig)
      emit('vm-updated')
    }
  } catch (error) {
    console.error('保存虚拟机配置失败:', error)
  }
}

// 添加硬件的方法
function addHardware(type) {
  // 增加对应硬件类型的计数
  hardwareCounts.value[type.value]++

  // 创建新的硬件配置
  const newHardware = {
    cfg: {}, // 空配置，让组件使用默认值
    name: `${getHardwareLabel(type.value)} ${hardwareCounts.value[type.value]}`,
    type: type.value,
  }

  // 添加到按钮组
  btnGroup.value.push(newHardware)

  // 关闭模态框
  showAddHardwareModal.value = false

  // 选中新添加的硬件
  const newIndex = btnGroup.value.length - 1
  handleMenuClick(newHardware, newIndex)
}

// 获取硬件类型的中文标签
function getHardwareLabel(type) {
  const hardwareType = availableHardwareTypes.find((h) => h.value === type)
  return hardwareType ? hardwareType.label : type
}

// 判断是否为可删除的硬件（即用户添加的硬件）
function isRemovableHardware(index) {
  // 只有通过"添加硬件"功能添加的菜单项才能删除
  // 默认菜单项（前defaultMenuCount个）不能删除
  return index >= defaultMenuCount
}

// 删除硬件的方法
function deleteHardware(index) {
  if (!isRemovableHardware(index)) return

  // 如果删除的是当前选中的项，需要重置选中状态
  if (currentMenuIndex.value === index) {
    currentMenuIndex.value = -1
    selectedMenu.value = 'overview'
  }

  // 调整当前选中项的索引（如果删除的项在当前选中项之前）
  if (currentMenuIndex.value > index) {
    currentMenuIndex.value--
  }

  // 从按钮组中删除该项
  btnGroup.value.splice(index, 1)
}
</script>

<style scoped>
/* 菜单头部样式 */
.menu-header {
  background-color: #e74c3c;
  padding: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.menu-header h2 {
  color: white;
  margin: 0;
  font-size: 18px;
}

.install-btn,
.cancel-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.install-btn {
  background-color: #fff;
  color: #e74c3c;
}

.cancel-btn {
  background-color: transparent;
  color: #fff;
  border: 1px solid #fff;
  margin-left: 8px;
}

/* 虚拟机配置容器 - 适应模态框 */
.vm-cfg-container {
  display: flex;
  width: 100%;
  height: 80vh;
  max-height: 800px;
  font-family: Arial, sans-serif;
  overflow: hidden;
}

/* 左侧菜单 */
.vm-cfg {
  width: 280px;
  height: 100%;
  border: 1px solid #ccc;
  display: flex;
  flex-direction: column;
  background-color: #2c3e50;
  overflow: hidden;
}

.menu-list {
  flex: 1;
  overflow-y: auto;
}

.menu-item {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border-bottom: 1px solid #34495e;
  position: relative;
}

.menu-item:hover {
  background-color: #34495e;
}

.add-hardware-btn {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border-bottom: 1px solid #34495e;
  background-color: #3498db;
  color: white;
}

.add-hardware-btn:hover {
  background-color: #2980b9;
}

.menu-icon {
  width: 20px;
  height: 20px;
  margin-right: 12px;
  background-color: #7f8c8d;
  border-radius: 2px;
  display: inline-block;
  text-align: center;
  line-height: 20px;
  font-size: 16px;
}

.add-hardware-btn .menu-icon {
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 50%;
}

.menu-text {
  font-size: 14px;
  color: white;
  flex: 1;
}

/* 删除按钮样式 */
.delete-btn {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: #e74c3c;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  cursor: pointer;
  margin-left: 8px;
  transition: background-color 0.3s;
}

.delete-btn:hover {
  background-color: #c0392b;
}

/* 右侧内容区域 */
.vm-content {
  flex: 1;
  padding: 20px;
  background-color: #ecf0f1;
  overflow-y: auto;
}

/* 添加硬件模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  width: 65vw;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 硬件类型选择 */
.hardware-types {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 10px;
  margin-top: 15px;
}

.hardware-type-item {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  text-align: center;
  transition: background-color 0.3s;
}

.hardware-type-item:hover {
  background-color: #f0f0f0;
  border-color: #3498db;
}

.modal-content h3 {
  margin-top: 0;
  color: #2c3e50;
}

.hardware-types {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.hardware-type-item {
  color: #333;
  padding: 12px;
  background-color: #f5f5f5;
  border-radius: 4px;
  cursor: pointer;
  text-align: center;
  transition: background-color 0.3s;
}

.hardware-type-item:hover {
  background-color: #3498db;
  color: white;
}
</style>
