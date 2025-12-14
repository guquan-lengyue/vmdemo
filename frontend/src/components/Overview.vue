<!-- 概况 -->
<template>
  <div class="vm-overview">
    <h2>虚拟机概况</h2>
    <div class="info-item">
      <span class="label">名称:</span>
      <input class="value" v-model="localCfg.name" @input="updateCfg" />
    </div>
    <div class="info-item">
      <span class="label">UUID:</span>
      <span class="value">{{ localCfg.uuid }}</span>
    </div>
    <div class="info-item">
      <span class="label">操作系统:</span>
      <span class="value">Ubuntu 25.10</span>
    </div>
    <div class="info-item">
      <span class="label">架构:</span>
      <span class="value">x86_64</span>
    </div>
    <div class="info-item">
      <span class="label">固件类型:</span>
      <select
        class="value"
        v-model="localCfg.osFirmware"
        name="osFirmware"
        id="osFirmware"
        @change="updateCfg"
      >
        <option value="bios">BIOS</option>
        <option value="uefi">UEFI</option>
      </select>
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
      name: 'ubuntu25.10',
      uuid: 'f0715790-cde5-4faf-8869-d8d72dfaf7d8',
      osMachine: 'pc',
      osFirmware: 'bios',
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

// 计算xml
const xml = computed(() => {
  return `
<name>${localCfg.value.name}</name>
<uuid>${localCfg.value.uuid}</uuid>
<metadata>
  <libosinfo:libosinfo xmlns:libosinfo="http://libosinfo.org/xmlns/libvirt/domain/1.0">
    <libosinfo:os id="http://ubuntu.com/ubuntu/25.10"/>
  </libosinfo:libosinfo>
</metadata>
<os firmware="${localCfg.value.osFirmware}">
  <type arch="x86_64" machine="${localCfg.value.osMachine}">hvm</type>
</os>
<features>
  <acpi/>
  <apic/>
  <vmport state="off"/>
</features>
`
})
</script>

<style scoped>
.vm-overview {
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

.info-item {
  display: flex;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.info-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.label {
  width: 100px;
  font-weight: bold;
  color: #555;
}

.value {
  flex: 1;
  color: #333;
}
</style>
