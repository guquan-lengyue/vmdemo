<!-- 概况 -->
<template>
  <div class="vm-overview">
    <h2>虚拟机概况</h2>
    <div class="info-item">
      <span class="label">名称:</span>
      <input class="value" v-model="localCfg.name" @input="updateCfg" />
    </div>
    <div class="info-item">
      <span class="label">UUID:</span>
      <span class="value">{{ localCfg.uuid }}</span>
    </div>
    <div class="info-item">
      <span class="label">操作系统:</span>
      <span class="value">Ubuntu 25.10</span>
    </div>
    <div class="info-item">
      <span class="label">架构:</span>
      <span class="value">x86_64</span>
    </div>
    <div class="info-item">
      <span class="label">固件类型:</span>
      <select
        class="value"
        v-model="localCfg.osFirmware"
        name="osFirmware"
        id="osFirmware"
        @change="updateCfg"
      >
        <option value="bios">BIOS</option>
        <option value="uefi">UEFI</option>
      </select>
    </div>
    <div class="boot-order">
      <span class="label">启动顺序:</span>
      <div class="boot-order-list">
        <div
          v-for="(component, index) in bootableComponents"
          :key="component.id"
          class="boot-order-item"
          draggable="true"
          @dragstart="onDragStart(index)"
          @dragover.prevent
          @drop="onDrop(index)"
        >
          <div class="boot-order-content">
            <span class="boot-order-number">{{ index + 1 }}</span>
            <span class="boot-order-label">{{ component.label }}</span>
            <span class="boot-order-type">({{ component.type }})</span>
          </div>
          <div class="boot-order-actions">
            <button v-if="index > 0" class="move-up-btn" @click="moveUp(index)" title="上移">
              ↑
            </button>
            <button
              v-if="index < bootableComponents.length - 1"
              class="move-down-btn"
              @click="moveDown(index)"
              title="下移"
            >
              ↓
            </button>
          </div>
        </div>
        <div v-if="bootableComponents.length === 0" class="no-bootable-devices">没有可启动设备</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, inject } from 'vue'

// 接收父组件传递的配置
const props = defineProps({
  cfg: {
    type: Object,
    default: () => ({
      name: 'ubuntu25.10',
      uuid: 'f0715790-cde5-4faf-8869-d8d72dfaf7d8',
      osMachine: 'pc',
      osFirmware: 'bios',
    }),
  },
})

// 定义事件
const emit = defineEmits(['update:cfg'])

// 本地配置副本
const localCfg = ref({ ...props.cfg })

// 获取父组件的btnGroup数据
const btnGroup = inject('btnGroup')

// 拖拽相关状态
const draggedItemIndex = ref(-1)

// 扫描可启动设备
const bootableComponents = computed(() => {
  if (!btnGroup) return []

  const bootableDevices = []

  // 扫描disk和interface类型的组件
  btnGroup.forEach((item, index) => {
    if (item.type === 'disk' || item.type === 'interface') {
      const cfg = item.cfg || {}

      // 为每个设备生成唯一标识
      const deviceId = `${item.type}-${index}`

      bootableDevices.push({
        id: deviceId,
        type: item.type,
        label: getDeviceLabel(item, cfg),
        cfg: cfg,
        originalIndex: index,
      })
    }
  })

  // 按bootOrder排序
  return bootableDevices.sort((a, b) => {
    const orderA = a.cfg.bootOrder || 0
    const orderB = b.cfg.bootOrder || 0
    return orderA - orderB
  })
})

// 获取设备显示标签
const getDeviceLabel = (item, cfg) => {
  if (item.type === 'disk') {
    if (cfg.diskType === 'cdrom') {
      return 'CD-ROM'
    }
    return `磁盘 ${cfg.targetDev || 'vda'}`
  } else if (item.type === 'interface') {
    return `网络 ${cfg.targetDevice || 'vnet0'}`
  }
  return item.name
}

// 监听配置变化
watch(
  () => props.cfg,
  (newVal) => {
    localCfg.value = { ...newVal }
  },
  { deep: true },
)

// 更新配置
const updateCfg = () => {
  emit('update:cfg', { ...localCfg.value })
}

// 拖拽开始
const onDragStart = (index) => {
  draggedItemIndex.value = index
}

// 拖拽放置
const onDrop = (targetIndex) => {
  if (draggedItemIndex.value === -1) return

  const fromIndex = draggedItemIndex.value
  if (fromIndex === targetIndex) return

  // 重新排序设备
  reorderBootDevices(fromIndex, targetIndex)
  draggedItemIndex.value = -1
}

// 重新排序启动设备
const reorderBootDevices = (fromIndex, toIndex) => {
  if (!btnGroup) return

  // 获取当前排序的设备列表
  const currentOrder = [...bootableComponents.value]

  // 移动设备
  const [movedItem] = currentOrder.splice(fromIndex, 1)
  currentOrder.splice(toIndex, 0, movedItem)

  // 更新所有设备的bootOrder
  currentOrder.forEach((device, index) => {
    const bootOrder = index + 1

    // 找到对应的btnGroup项并更新bootOrder
    const btnGroupItem = btnGroup.find((item, idx) => `${item.type}-${idx}` === device.id)

    if (btnGroupItem && btnGroupItem.cfg) {
      btnGroupItem.cfg.bootOrder = bootOrder
    }
  })

  // 触发配置更新
  updateCfg()
}

// 上移设备
const moveUp = (index) => {
  if (index <= 0) return
  reorderBootDevices(index, index - 1)
}

// 下移设备
const moveDown = (index) => {
  if (index >= bootableComponents.value.length - 1) return
  reorderBootDevices(index, index + 1)
}

// 计算xml
const xml = computed(() => {
  return `
<name>${localCfg.value.name}</name>
<uuid>${localCfg.value.uuid}</uuid>
<metadata>
  <libosinfo:libosinfo xmlns:libosinfo="http://libosinfo.org/xmlns/libvirt/domain/1.0">
    <libosinfo:os id="http://ubuntu.com/ubuntu/25.10"/>
  </libosinfo:libosinfo>
</metadata>
<os firmware="${localCfg.value.osFirmware}">
  <type arch="x86_64" machine="${localCfg.value.osMachine}">hvm</type>
</os>
<features>
  <acpi/>
  <apic/>
  <vmport state="off"/>
</features>
`
})
</script>

<style scoped>
.vm-overview {
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

.info-item {
  display: flex;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.info-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.label {
  width: 100px;
  font-weight: bold;
  color: #555;
}

.value {
  flex: 1;
  color: #333;
}

.boot-order {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.boot-order-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 8px;
  background-color: #fafafa;
}

.boot-order-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: move;
  transition: all 0.2s ease;
  user-select: none;
}

.boot-order-item:hover {
  border-color: #007bff;
  box-shadow: 0 2px 4px rgba(0, 123, 255, 0.1);
}

.boot-order-item.dragging {
  opacity: 0.5;
  background-color: #f8f9fa;
}

.boot-order-content {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.boot-order-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background-color: #007bff;
  color: white;
  border-radius: 50%;
  font-size: 12px;
  font-weight: bold;
}

.boot-order-label {
  font-weight: bold;
  color: #333;
}

.boot-order-type {
  color: #666;
  font-size: 0.9em;
}

.boot-order-actions {
  display: flex;
  gap: 4px;
}

.move-up-btn,
.move-down-btn {
  width: 28px;
  height: 28px;
  border: 1px solid #ddd;
  background-color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.move-up-btn:hover,
.move-down-btn:hover {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}

.move-up-btn:disabled,
.move-down-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.no-bootable-devices {
  text-align: center;
  color: #999;
  font-style: italic;
  padding: 20px;
}

/* 拖拽时的视觉效果 */
.boot-order-item.drag-over {
  border: 2px dashed #007bff;
  background-color: #f0f8ff;
}
</style>
