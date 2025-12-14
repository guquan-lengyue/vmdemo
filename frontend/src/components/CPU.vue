<!-- CPU设置 -->
<template>
  <div class="vm-cpu">
    <h2>CPU设置</h2>
    <div class="cpu-info">
      <div class="info-item">
        <span class="label">主机CPU数量:</span>
        <span class="value">{{ props.hostMsg.hostCpuCount }}</span>
      </div>
      <div class="info-item">
        <span class="label">CPU数量:</span>
        <input
          :disabled="localCfg.isNotManualTopology"
          type="number"
          v-model="localCfg.cpuCount"
          min="1"
          class="input-field"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">CPU模式:</span>
        <select v-model="localCfg.cpuMode" class="select-field" @change="updateCfg">
          <option value="host-passthrough">host-passthrough</option>
          <option value="host-model">host-model</option>
          <option value="custom">custom</option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">手动拓扑:</span>
        <input
          type="checkbox"
          v-model="localCfg.isNotManualTopology"
          class="checkbox"
          @change="updateCfg"
        />
      </div>
      <div v-if="localCfg.isNotManualTopology" class="info-item">
        <span class="label">CPU拓扑:</span>
        <div class="topology">
          <div class="topology-item">
            <span> sockets: </span>
            <input
              :disabled="!localCfg.isNotManualTopology"
              type="number"
              v-model="localCfg.manualTopology.sockets"
              min="1"
              class="small-input"
              @input="updateCfg"
            />
          </div>
          <div class="topology-item">
            <span> cores: </span>
            <input
              :disabled="!localCfg.isNotManualTopology"
              type="number"
              v-model="localCfg.manualTopology.cores"
              min="1"
              class="small-input"
              @input="updateCfg"
            />
          </div>
          <div class="topology-item">
            <span> threads: </span>
            <input
              :disabled="!localCfg.isNotManualTopology"
              type="number"
              v-model="localCfg.manualTopology.threads"
              min="1"
              class="small-input"
              @input="updateCfg"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  hostMsg: {
    type: Object,
    required: false,
    default: () => ({
      hostCpuCount: 1,
    }),
  },
  cfg: {
    type: Object,
    default: () => ({
      cpuCount: 2,
      cpuMode: 'host-passthrough',
      isNotManualTopology: false,
      manualTopology: {
        sockets: 2,
        cores: 1,
        threads: 1,
      },
    }),
  },
})

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

// 计算cpuCount和manualTopology的联动
watch(
  () => localCfg.value.cpuCount,
  (newVal) => {
    if (localCfg.value.isNotManualTopology) {
      return
    }
    localCfg.value.manualTopology = {
      sockets: newVal,
      cores: 1,
      threads: 1,
    }
    updateCfg()
  },
)

watch(
  () => localCfg.value.manualTopology,
  (newVal) => {
    if (!localCfg.value.isNotManualTopology) {
      return
    }
    localCfg.value.cpuCount = newVal.sockets * newVal.cores * newVal.threads
    updateCfg()
  },
  { deep: true },
)

const xml = computed(() => {
  let topology = ''
  if (!localCfg.value.isNotManualTopology) {
    topology = `<topology sockets="${localCfg.value.manualTopology.sockets}" cores="${localCfg.value.manualTopology.cores}" threads="${localCfg.value.manualTopology.threads}"/>`
  }
  return `
<vcpu current="6">${localCfg.value.cpuCount}</vcpu>
<cpu mode="${localCfg.value.cpuMode}">
  ${topology}
</cpu>
`
})
</script>

<style scoped>
.vm-cpu {
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.value {
  color: #333;
}
h2 {
  margin-bottom: 20px;
  color: #333;
  font-size: 20px;
}

.cpu-info {
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
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.topology {
  display: flex;
  gap: 20px;
  align-items: center;
}

.topology-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: black;
}

.small-input {
  width: 60px;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}
</style>
