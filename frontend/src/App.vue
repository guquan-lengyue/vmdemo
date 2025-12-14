<template>
  <div class="vm-cfg-container">
    <div class="vm-cfg">
      <div class="menu-header">
        <button class="install-btn">开始安装</button>
        <button class="cancel-btn">取消安装</button>
      </div>
      <div class="menu-list">
        <div
          v-for="(item, index) in btnGroup"
          :key="index"
          class="menu-item"
          @click="selectedMenu = item.type"
        >
          <span class="menu-icon"></span>
          <span class="menu-text">{{ item.name }}</span>
        </div>
      </div>
    </div>
    <div class="vm-content">
      <Component :is="componentMap[selectedMenu] || componentMap['overview']" :hostMsg="hostMsg" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Overview from './components/Overview.vue'
import CPU from './components/CPU.vue'
import Memory from './components/Memery.vue'

const selectedMenu = ref('overview')

// 组件映射对象
const componentMap = {
  overview: Overview,
  cpu: CPU,
  memory: Memory,
}
// 主机信息
const hostMsg = ref({
  hostCpuCount: 16,
  hostMemory: 4096,
})

// 使用ref定义按钮组
const btnGroup = ref([
  { name: '概况', type: 'overview' },
  { name: 'CPU数', type: 'cpu' },
  { name: '内存', type: 'memory' },
  { name: '磁盘', type: 'disk' },
  { name: 'CDROM', type: 'cdrom' },
  { name: '虚拟网络', type: 'network' },
  { name: '显示协议', type: 'display' },
  { name: '声卡', type: 'sound' },
  { name: '控制台', type: 'console' },
  { name: '通道', type: 'channel-qemu' },
  { name: '通道', type: 'channel-spice' },
  { name: '视频', type: 'video' },
  { name: '控制器USB', type: 'controller-usb' },
  { name: '控制器PCIe', type: 'controller-pcie' },
  { name: 'USB转发器1', type: 'usb-redir1' },
  { name: 'USB转发器2', type: 'usb-redir2' },
  { name: 'RNG /dev/urandom', type: 'rng' },
])
</script>

<style scoped>
.vm-cfg-container {
  display: flex;
  height: 100vh;
  font-family: Arial, sans-serif;
}

.vm-cfg {
  width: 280px;
  height: 100vh;
  border: 1px solid #ccc;
  display: flex;
  flex-direction: column;
  background-color: #2c3e50;
}

.menu-header {
  background-color: #e74c3c;
  padding: 10px;
  display: flex;
  justify-content: space-between;
}

.install-btn,
.cancel-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.install-btn {
  background-color: #fff;
  color: #e74c3c;
}

.cancel-btn {
  background-color: transparent;
  color: #fff;
  border: 1px solid #fff;
}

.menu-list {
  flex: 1;
  overflow-y: auto;
}

.menu-item {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border-bottom: 1px solid #34495e;
}

.menu-item:hover {
  background-color: #34495e;
}

.menu-icon {
  width: 20px;
  height: 20px;
  margin-right: 12px;
  background-color: #7f8c8d;
  border-radius: 2px;
  display: inline-block;
}

.menu-text {
  font-size: 14px;
  color: white;
}

.vm-content {
  flex: 1;
  padding: 20px;
  background-color: #ecf0f1;
  overflow-y: auto;
}
</style>
