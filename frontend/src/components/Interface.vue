<!-- 网络接口设置 -->
<template>
  <div class="vm-interface">
    <h2>网络接口设置</h2>
    <div class="interface-info">
      <div class="info-item">
        <span class="label">网络类型:</span>
        <select v-model="localCfg.networkType" class="select-field" @change="updateCfg">
          <option v-for="source in props.sources" :key="source" :value="source">
            {{ source }}
          </option>
          <option value="bridge">桥接模式</option>
        </select>
      </div>
      <div v-if="localCfg.networkType === 'bridge'" class="info-item">
        <span class="label">桥接设备名称:</span>
        <input
          type="text"
          v-model="localCfg.bridgeName"
          class="input-field"
          placeholder="例如: br0"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">MAC地址:</span>
        <input
          type="text"
          v-model="localCfg.macAddress"
          class="input-field"
          placeholder="自动生成"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">设备型号:</span>
        <select v-model="localCfg.model" class="select-field" @change="updateCfg">
          <option value="virtio">virtio</option>
          <option value="e1000">e1000</option>
          <option value="default">虚拟机管理默认</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

// 接收父组件传递的配置和参数
const props = defineProps({
  cfg: {
    type: Object,
    default: () => ({
      networkType: 'bridge',
      bridgeName: '',
      macAddress: '',
      model: 'virtio',
    }),
  },
  sources: {
    type: Array,
    default: () => ['default'],
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
</script>

<style scoped>
.vm-interface {
  margin-bottom: 20px;
}

.interface-info {
  display: flex;
  flex-direction: column;
  gap: 15px;
  padding: 15px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.label {
  color: #333;
  width: 150px;
  font-weight: 500;
}

.select-field,
.input-field {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.checkbox {
  width: 18px;
  height: 18px;
}
</style>
