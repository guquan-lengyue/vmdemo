<!-- 视频设置 -->
<template>
  <div class="video-info">
    <h2 style="color: #333">视频设置</h2>
    <div class="info-item">
      <span class="label">视频设备:</span>
      <select v-model="localCfg.model.type" class="select-field" @change="updateCfg">
        <option value="virtio">VirtIO</option>
        <option value="vga">VGA</option>
        <option value="ramfb">Ramfb</option>
        <option value="qxl">QXL</option>
        <option value="bochs">Bochs</option>
        <option value="none">None</option>
      </select>
    </div>
    <div class="info-item" v-if="localCfg.model.type === 'virtio'">
      <span class="label">3D加速:</span>
      <select
        v-model="localCfg.model.acceleration.accel3d"
        class="select-field"
        @change="updateCfg"
      >
        <option value="yes">开启</option>
        <option value="no">关闭</option>
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
      model: {
        type: 'virtio',
        acceleration: {
          accel3d: 'yes',
        },
      },
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
  // 如果视频设备不是VirtIO，确保3D加速配置被移除或设置为关闭
  if (localCfg.value.model.type !== 'virtio') {
    delete localCfg.value.model.acceleration
  }
  emit('update:cfg', { ...localCfg.value })
}
</script>

<style scoped>
.vm-video {
  margin-bottom: 20px;
}

.video-info {
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
