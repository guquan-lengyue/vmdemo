<!-- 磁盘设置 -->
<template>
  <div class="vm-disk">
    <h2>磁盘设置</h2>
    <div class="disk-info">
      <div class="info-item">
        <span class="label">磁盘类型:</span>
        <select v-model="localCfg.diskType" class="select-field" @change="updateCfg">
          <option value="disk">磁盘</option>
          <option value="cdrom">光驱</option>
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
          <option value="virtio">virtio</option>
          <option value="sata">sata</option>
          <option value="ide">ide</option>
          <option value="scsi">scsi</option>
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
      sourcePath: '/var/lib/libvirt/images/ubuntu25.10.qcow2',
      diskFormat: 'qcow2',
      targetDev: 'vda',
      targetBus: 'virtio',
      isReadOnly: false,
    }),
  },
})

// 定义事件
const emit = defineEmits(['update:cfg'])

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

// 更新配置
const updateCfg = () => {
  emit('update:cfg', { ...localCfg.value })
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
