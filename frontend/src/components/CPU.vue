<!-- CPU设置 -->
<template>
  <div class="vm-cpu">
    <h2>CPU设置</h2>
    <div class="cpu-info">
      <div class="info-item">
        <span class="label">CPU数量:</span>
        <input type="number" v-model="cpuCount" min="1" max="16" class="input-field" />
      </div>
      <div class="info-item">
        <span class="label">CPU模式:</span>
        <select v-model="cpuMode" class="select-field">
          <option value="host-passthrough">host-passthrough</option>
          <option value="host-model">host-model</option>
          <option value="custom">custom</option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">CPU拓扑:</span>
        <div class="topology">
          <div class="topology-item">
            <span> sockets: </span>
            <input type="number" v-model="sockets" min="1" class="small-input" />
          </div>
          <div class="topology-item">
            <span> cores: </span>
            <input type="number" v-model="cores" min="1" class="small-input" />
          </div>
          <div class="topology-item">
            <span> threads: </span>
            <input type="number" v-model="threads" min="1" class="small-input" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const cpuCount = ref(2)
const cpuMode = ref('host-passthrough')
const sockets = ref(1)
const cores = ref(2)
const threads = ref(1)

const xml = computed(() => {
  return `
<cpu mode="${cpuMode.value}">
  <topology sockets="${sockets.value}" cores="${cores.value}" threads="${threads.value}"/>
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
  flex: 1;
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
}

.small-input {
  width: 60px;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}
</style>
