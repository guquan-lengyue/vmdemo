<!-- 虚拟网络 -->
<template>
  <div class="vm-interface">
    <h2>网络设置</h2>
    <div class="interface-info">
      <div class="info-item">
        <span class="label">网络类型:</span>
        <select v-model="localCfg.networkType" class="select-field" @change="updateCfg">
          <option value="bridge">桥接</option>
          <option value="nat">NAT</option>
          <option value="internal">内部网络</option>
          <option value="isolated">隔离网络</option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">源设备:</span>
        <input
          type="text"
          v-model="localCfg.sourceDevice"
          class="input-field"
          placeholder="eth0"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">网络模型:</span>
        <select v-model="localCfg.networkModel" class="select-field" @change="updateCfg">
          <option value="virtio">virtio</option>
          <option value="e1000">e1000</option>
          <option value="rtl8139">rtl8139</option>
          <option value="ne2k_pci">ne2k_pci</option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">目标设备:</span>
        <input
          type="text"
          v-model="localCfg.targetDevice"
          class="small-input"
          placeholder="vnet0"
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
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

// 接收父组件传递的配置
const props = defineProps({
  cfg: {
    type: Object,
    default: () => ({
      networkType: 'bridge',
      sourceDevice: 'eth0',
      networkModel: 'virtio',
      targetDevice: 'vnet0',
      macAddress: '',
      bootOrder: 0,
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
</script>

<style scoped>
.vm-interface {
  padding: 1rem;
  background-color: #f5f5f5;
  border-radius: 8px;
  margin-bottom: 1rem;
}

h2 {
  margin-top: 0;
  color: #333;
  font-size: 1.2rem;
  margin-bottom: 1rem;
}

.interface-info {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.label {
  width: 120px;
  font-weight: bold;
  color: #555;
}

.input-field,
.select-field {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.small-input {
  width: 150px;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.checkbox {
  width: 18px;
  height: 18px;
}
</style>
