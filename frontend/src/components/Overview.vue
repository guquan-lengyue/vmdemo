<!-- 概况 -->
<template>
  <div class="vm-overview">
    <h2>虚拟机概况</h2>
    <div class="info-item">
      <span class="label">名称:</span>
      <input class="value" v-model="name" />
    </div>
    <div class="info-item">
      <span class="label">UUID:</span>
      <span class="value">{{ uuid }}</span>
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
      <select class="value" v-model="osFirmware" name="osFirmware" id="osFirmware">
        <option value="bios">BIOS</option>
        <option value="uefi">UEFI</option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const name = ref('ubuntu25.10')
const uuid = ref('f0715790-cde5-4faf-8869-d8d72dfaf7d8')
const osMachine = ref('pc')
const osFirmware = ref('bios')

const xml = computed(() => {
  return ```
<name>${name.value}</name>
<uuid>${uuid.value}</uuid>
<metadata>
  <libosinfo:libosinfo xmlns:libosinfo="http://libosinfo.org/xmlns/libvirt/domain/1.0">
    <libosinfo:os id="http://ubuntu.com/ubuntu/25.10"/>
  </libosinfo:libosinfo>
</metadata>
<os firmware="${osFirmware.value}">
  <type arch="x86_64" machine="${osMachine.value}">hvm</type>
</os>
<features>
  <acpi/>
  <apic/>
  <vmport state="off"/>
</features>
```
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
