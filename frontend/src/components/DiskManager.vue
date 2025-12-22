<!-- 磁盘设置 -->
<template>
  <div class="vm-disk">
    <h2>磁盘设置</h2>
    <div class="disk-info">
      <div class="info-item">
        <span class="label">磁盘类型:</span>
        <select v-model="localCfg.diskType" class="select-field" @change="updateCfg">
          <option v-for="item in diskTypes" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">源文件路径:</span>
        <!-- 根据磁盘类型显示不同的输入方式 -->
        <div v-if="localCfg.diskType === 'disk'" class="disk-selector">
          <select v-model="localCfg.sourcePath" class="select-field" @change="updateCfg">
            <option v-for="disk in availableDisks" :key="disk.path" :value="disk.path">
              {{ disk.name }} ({{ disk.size }})
            </option>
          </select>
          <button class="create-disk-btn" @click="showCreateDialog = true">创建新磁盘</button>
        </div>
        <input
          v-else
          type="text"
          v-model="localCfg.sourcePath"
          class="input-field"
          placeholder="/path/to/cdrom.iso"
          @input="updateCfg"
        />
      </div>
      <!-- 创建磁盘对话框 -->
      <div v-if="showCreateDialog" class="dialog-overlay">
        <div class="dialog">
          <h3 style="color: #333;">创建新磁盘</h3>
          <form @submit.prevent="createDisk">
            <div class="form-item">
              <label style="color: #333;" for="disk-name">磁盘名称:</label>
              <input
                id="disk-name"
                type="text"
                v-model="newDisk.name"
                placeholder="disk.qcow2"
                required
              />
            </div>
            <div class="form-item">
              <label style="color: #333;" for="disk-format">磁盘格式:</label>
              <select id="disk-format" v-model="newDisk.format">
                <option value="qcow2">qcow2</option>
                <option value="raw">raw</option>
                <option value="vmdk">vmdk</option>
                <option value="qcow">qcow</option>
              </select>
            </div>
            <div class="form-item">
              <label style="color: #333;" for="disk-size">磁盘大小:</label>
              <input
                id="disk-size"
                type="text"
                v-model="newDisk.size"
                placeholder="10G"
                required
              />
            </div>
            <div class="dialog-actions">
              <button type="button" class="btn btn-secondary" @click="showCreateDialog = false">取消</button>
              <button type="submit" class="btn btn-primary">创建</button>
            </div>
          </form>
        </div>
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
import { ref, watch, onMounted } from 'vue'
import { diskApi } from '../api.js'

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

const diskTypes = [
  { label: '磁盘', value: 'disk' },
  { label: '光驱', value: 'cdrom' },
]

const targetBusTypes = [
  { label: 'VirtIO', value: 'virtio' },
  { label: 'SATA', value: 'sata' },
  { label: 'IDE', value: 'ide' },
  { label: 'SCSI', value: 'scsi' },
]

// 可用磁盘列表
const availableDisks = ref([])

// 创建磁盘对话框状态
const showCreateDialog = ref(false)

// 新磁盘配置
const newDisk = ref({
  name: '',
  format: 'qcow2',
  size: '10G'
})

// 加载可用磁盘
const loadAvailableDisks = async () => {
  try {
    const response = await diskApi.listDisks()
    availableDisks.value = response.disks
  } catch (error) {
    console.error('Failed to load disks:', error)
  }
}

// 创建新磁盘
const createDisk = async () => {
  try {
    await diskApi.addDisk(newDisk.value.name, newDisk.value.format, newDisk.value.size)
    // 重新加载可用磁盘列表
    await loadAvailableDisks()
    // 关闭对话框
    showCreateDialog.value = false
    // 重置表单
    newDisk.value = {
      name: '',
      format: 'qcow2',
      size: '10G'
    }
  } catch (error) {
    console.error('Failed to create disk:', error)
    alert('创建磁盘失败: ' + error.message)
  }
}

// 组件挂载时加载可用磁盘
onMounted(() => {
  loadAvailableDisks()
})

// 更新配置
const updateCfg = () => {
  if (localCfg.value.diskType === 'cdrom') {
    localCfg.value.isReadOnly = true
    localCfg.value.diskFormat = 'raw'
  }
  // 计算菜单名称：{总线类型}:{磁盘类型}
  const targetBusLabel = targetBusTypes.find((i) => i.value === localCfg.value.targetBus)?.label
  const diskTypeLabel = diskTypes.find((i) => i.value === localCfg.value.diskType)?.label
  const menuName = `${targetBusLabel}-${diskTypeLabel}`
  emit('update:cfg', { ...localCfg.value })
  emit('update:menuName', menuName)
}
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
  /* 创建磁盘对话框样式 */
  .dialog-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .dialog {
    background-color: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    width: 400px;
  }

  .dialog h3 {
    margin-top: 0;
    margin-bottom: 20px;
  }

  .dialog form {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .dialog .form-item {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }

  .dialog label {
    font-weight: bold;
    font-size: 14px;
  }

  .dialog input,
  .dialog select {
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
  }

  .dialog .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 10px;
  }

  .dialog .btn {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .dialog .btn-primary {
    background-color: #409eff;
    color: white;
  }

  .dialog .btn-secondary {
    background-color: #909399;
    color: white;
  }

  /* 磁盘选择器样式 */
  .disk-selector {
    display: flex;
    gap: 10px;
    align-items: center;
  }

  .create-disk-btn {
    padding: 8px 16px;
    background-color: #409eff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }

  .create-disk-btn:hover {
    background-color: #66b1ff;
  }
</style>
