<!-- 显示设置 -->
<template>
  <div class="vm-display">
    <h2>显示设置</h2>
    <div class="display-info">
      <div class="info-item">
        <span class="label">显示类型:</span>
        <select v-model="localCfg.type" class="select-field" @change="updateCfg">
          <option v-for="(v, k) in displayTypes" :key="k" :value="k">{{ v }}</option>
        </select>
      </div>
      <div class="info-item">
        <span class="label">端口:</span>
        <input
          type="text"
          v-model="localCfg.port"
          class="input-field"
          placeholder="-1 (自动分配)"
          @input="updateCfg"
        />
      </div>

      <div class="info-item">
        <span class="label">监听地址:</span>
        <input
          type="text"
          v-model="localCfg.listen"
          class="input-field"
          placeholder="0.0.0.0"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">密码:</span>
        <input
          type="password"
          v-model="localCfg.passwd"
          class="input-field"
          placeholder="设置密码"
          @input="updateCfg"
        />
      </div>
      <div class="info-item">
        <span class="label">图像压缩:</span>
        <select v-model="localCfg.imageCompression" class="select-field" @change="updateCfg">
          <option value="off">关闭</option>
          <option value="auto">自动</option>
          <option value="on">开启</option>
        </select>
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
      type: 'vnc',
      port: '-1',
      listen: '0.0.0.0',
      passwd: '',
      imageCompression: 'off',
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

const displayTypes = {
  vnc: 'VNC',
  spice: 'SPICE',
}

// 更新配置
const updateCfg = () => {
  // 计算菜单名称：显示协议:{显示类型}
  const menuName = `显示协议-${displayTypes[localCfg.value.type]}`
  emit('update:cfg', { ...localCfg.value })
  emit('update:menuName', menuName)
}

const xml = computed(() => {
  return `
<graphics type="${localCfg.value.type}" port="${localCfg.value.port}" listen="${localCfg.value.listen}" password="${localCfg.value.passwd}">
  <gl enable="no"/>
  <image compression="${localCfg.value.imageCompression}"/>
</graphics>
  `
})
</script>

<style scoped>
.vm-display {
  margin-bottom: 20px;
}

.display-info {
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
