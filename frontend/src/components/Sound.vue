<!-- 声音设置 -->
<template>
  <div class="sound-info">
    <h2 style="color: #333">声音设置</h2>
    <div class="info-item">
      <span class="label">音频设备:</span>
      <select v-model="localCfg.model" class="select-field" @change="updateCfg">
        <option value="ac97">AC97</option>
        <option value="ich6">HDA(ICH6)</option>
        <option value="ich9">HDA(ICH9)</option>
      </select>
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
      model: 'ac97',
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
.vm-sound {
  margin-bottom: 20px;
}

.sound-info {
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
  width: 120px;
  font-weight: 500;
}

.select-field {
  flex: 1;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
