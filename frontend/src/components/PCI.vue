<!-- PCI设备设置 -->
<template>
  <div class="vm-pci">
    <h2>PCI设备设置</h2>
    <div class="pci-info">
      <div class="info-item">
        <span class="label">可用PCI设备:</span>
        <div class="pci-device-list">
          <div v-if="loading" class="loading">加载中...</div>
          <div v-else-if="pciDevices.length === 0" class="no-devices">无可用PCI设备</div>
          <div v-else class="devices">
            <div
              v-for="device in pciDevices"
              :key="device.id"
              class="pci-device-item"
              :class="{ 'selected': isDeviceSelected(device.id) }"
              @click="toggleDeviceSelection(device.id)"
            >
              <span class="device-id">{{ device.id }}</span>
              <span style="color: #000;" class="device-name">{{ device.name }}</span>
              <span class="device-vendor">{{ device.vendor }} {{ device.device }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { pciApi } from '../api.js'

const props = defineProps({
  cfg: {
    type: Object,
    default: () => ({
      attachedDevices: [] // 存储已挂载的PCI设备ID列表
    }),
  },
})

const emit = defineEmits(['update:cfg'])

// 使用 reactive 替代 ref，确保深度响应式
const localCfg = reactive({ ...props.cfg })

// PCI设备列表
const pciDevices = ref([])
const loading = ref(false)

// 监听配置变化 - 使用 deep: true 确保对象内部属性变化也能被监听
watch(
  () => props.cfg,
  (newVal) => {
    // 深拷贝，确保所有属性都被更新
    Object.assign(localCfg, newVal)
  },
  { deep: true },
)

// 更新配置 - 使用对象扩展语法创建新对象，确保Vue能检测到变化
const updateCfg = () => {
  emit('update:cfg', { ...localCfg })
}

// 获取PCI设备列表
const fetchPciDevices = async () => {
  loading.value = true
  try {
    const devices = await pciApi.listPciDevices()
    pciDevices.value = devices
  } catch (error) {
    console.error('获取PCI设备列表失败:', error)
    pciDevices.value = []
  } finally {
    loading.value = false
  }
}

// 检查设备是否已选中
const isDeviceSelected = (deviceId) => {
  return localCfg.attachedDevices.includes(deviceId)
}

// 切换设备选择状态
const toggleDeviceSelection = (deviceId) => {
  const index = localCfg.attachedDevices.indexOf(deviceId)
  if (index === -1) {
    // 添加设备
    localCfg.attachedDevices.push(deviceId)
  } else {
    // 移除设备
    localCfg.attachedDevices.splice(index, 1)
  }
  updateCfg()
}

// 组件挂载时获取PCI设备列表
onMounted(() => {
  fetchPciDevices()
})
</script>

<style scoped>
.vm-pci {
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-bottom: 20px;
  color: #333;
  font-size: 20px;
}

.pci-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.label {
  font-weight: bold;
  color: #333;
}

.pci-device-list {
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  max-height: 300px;
  overflow-y: auto;
}

.loading,
.no-devices {
  text-align: center;
  padding: 20px;
  color: #666;
}

.devices {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pci-device-item {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.pci-device-item:hover {
  background-color: #f5f5f5;
}

.pci-device-item.selected {
  background-color: #e3f2fd;
  border-color: #2196F3;
}

.device-id {
  display: block;
  font-weight: bold;
  color: #2196F3;
  margin-bottom: 5px;
}

.device-name {
  display: block;
  font-size: 14px;
  color: #333;
  margin-bottom: 3px;
}

.device-vendor {
  font-size: 12px;
  color: #666;
}
</style>