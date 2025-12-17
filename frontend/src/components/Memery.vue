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
        <input
          type="number"
          v-model="localCfg.memory"
          min="1"
          class="input-field"
          @input="updateCfg"
        />
        MiB
      </div>
      <div class="info-item">
        <span class="label">当前内存:</span>
        <input
          type="number"
          v-model="localCfg.currentMemory"
          min="1"
          class="input-field"
          @input="updateCfg"
        />
        MiB
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue'

const props = defineProps({
  hostMsg: {
    type: Object,
    required: false,
    default: () => ({
      hostMemory: 1,
    }),
  },
  cfg: {
    type: Object,
    default: () => ({
      memory: 4194304,
      currentMemory: 4194304,
    }),
  },
})

const emit = defineEmits(['update:cfg'])

// 使用 reactive 替代 ref，确保深度响应式
const localCfg = reactive({ ...props.cfg })

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
