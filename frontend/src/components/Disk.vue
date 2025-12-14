<!-- 磁盘设置 -->
<template>
  <div class="vm-disk">
    <h2>磁盘设置</h2>
    <div class="disk-info">
      <div class="info-item">
        <span class="label">磁盘类型:</span>
        <select v-model="localCfg.diskType" class="select-field" @change="updateCfg">
          <option v-for="item in disKTypes" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">源文件路径:</span>
        <input
          type="text"
          v-model="localCfg.sourcePath"
          class="input-field"
          placeholder="/path/to/disk.img"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">磁盘格式:</span>
        <select v-model="localCfg.diskFormat" class="select-field" @change="updateCfg">
          <option value="qcow2">qcow2</option>
          <option value="raw">raw</option>
          <option value="vmdk">vmdk</option>
          <option value="qcow">qcow</option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">目标设备:</span>
        <input
          type="text"
          v-model="localCfg.targetDev"
          class="small-input"
          placeholder="vda"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">总线类型:</span>
        <select v-model="localCfg.targetBus" class="select-field" @change="updateCfg">
          <option v-for="item in targetBusTypes" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </div>
      <div class="info-item" v-if="localCfg.diskType === 'cdrom'">
        <span class="label">只读:</span>
        <input type="checkbox" v-model="localCfg.isReadOnly" class="checkbox" @change="updateCfg" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

// 接收父组件传递的配置
const props = defineProps({
  cfg: {
    type: Object,
    default: () => ({
      diskType: 'disk',
      sourcePath: '',
      diskFormat: 'qcow2',
      targetDev: 'vda',
      targetBus: 'virtio',
      isReadOnly: false,
      bootOrder: 0,
    }),
  },
})

// 定义事件
const emit = defineEmits(['update:cfg', 'update:menuName'])

// 本地配置副本
const localCfg = ref({ ...props.cfg })

// 监听配置变化
watch(
  () => props.cfg,
  (newVal) => {
    localCfg.value = { ...newVal }
  },
  { deep: true },
)

const disKTypes = [
  { label: '磁盘', value: 'disk' },
  { label: '光驱', value: 'cdrom' },
]

const targetBusTypes = [
  { label: 'VirtIO', value: 'virtio' },
  { label: 'SATA', value: 'sata' },
  { label: 'IDE', value: 'ide' },
  { label: 'SCSI', value: 'scsi' },
]

// 更新配置
const updateCfg = () => {
  if (localCfg.value.diskType === 'cdrom') {
    localCfg.value.isReadOnly = true
    localCfg.value.diskFormat = 'raw'
  }
  // 计算菜单名称：{总线类型}:{磁盘类型}
  const targetBusLabel = targetBusTypes.find((i) => i.value === localCfg.value.targetBus)?.label
  const diskTypeLabel = disKTypes.find((i) => i.value === localCfg.value.diskType)?.label
  const menuName = `${targetBusLabel}-${diskTypeLabel}`
  emit('update:cfg', { ...localCfg.value })
  emit('update:menuName', menuName)
}

// 生成磁盘配置的XML
const xml = computed(() => {
  const readonlyTag =
    localCfg.value.isReadOnly && localCfg.value.diskType === 'cdrom' ? '<readonly/>' : ''
  return `
<disk type="file" device="${localCfg.value.diskType}">
  <driver name="qemu" type="${localCfg.value.diskFormat}" discard="unmap"/>
  <source file="${localCfg.value.sourcePath}"/>
  <target dev="${localCfg.value.targetDev}" bus="${localCfg.value.targetBus}"/>
  ${readonlyTag}
  <boot order="${localCfg.value.bootOrder}"/>
</disk>
`
})
</script>

<style scoped>
.vm-disk {
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

.disk-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  align-items: center;
}

.label {
  width: 120px;
  font-weight: bold;
  color: #555;
}

.input-field,
.select-field {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.small-input {
  width: 80px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}
</style>
