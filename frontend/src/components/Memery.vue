<!-- 内存设置 -->
<template>
  <div class="vm-memory">
    <h2>内存设置</h2>
    <div class="info-item">
      <span class="label">总主机内存:</span>
      <span class="value">{{ props.hostMsg.hostMemory }} MiB</span>
    </div>
    <div class="memory-info">
      <div class="info-item">
        <span class="label">内存大小:</span>
        <input type="number" v-model="memory" min="1" class="input-field" /> MiB
      </div>
      <div class="info-item">
        <span class="label">当前内存:</span>
        <input type="number" v-model="currentMemory" min="1" class="input-field" /> MiB
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
const props = defineProps({
  hostMsg: {
    type: Object,
    required: false,
    default: () => ({
      hostMemory: 1,
    }),
  },
})
const memory = ref(4194304)
const currentMemory = ref(4194304)
</script>
<style scoped>
.vm-memory {
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

.memory-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  align-items: center;
  color: #333;
}

.label {
  width: 120px;
  font-weight: bold;
  color: #555;
}

.input-field {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  margin-right: 10px;
}
</style>
