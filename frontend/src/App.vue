<template>
  <div>
    <div style="display: flex">
      <button class="btn" @click="selected = 'vnc'">VNC</button>
      <button class="btn" @click="selected = 'vmopts'">虚拟机</button>
      <button class="btn" @click="selected = 'diskopts'">虚拟硬盘</button>
      <button class="btn" @click="selected = 'usb'">usb设备</button>
    </div>
    <NoVncViewer v-show="selected === 'vnc'" />
    <VmOpts v-show="selected === 'vmopts'" />
    <DiskOpts v-show="selected === 'diskopts'" />
    <Usb v-show="selected === 'usb'" />
  </div>
</template>

<script setup>
import NoVncViewer from '@/components/NoVncViewer.vue'
import VmOpts from '@/components/VmOpts.vue'
import DiskOpts from '@/components/DiskOpts.vue'
import Usb from '@/components/Usb.vue'
import { ref } from 'vue'
const selected = ref('vnc')
</script>

<style scoped>
.btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-right: 4px;
}
</style>
