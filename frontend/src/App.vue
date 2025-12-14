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
      <Component
        :is="componentMap[selectedMenu] || componentMap['overview']"
        :hostMsg="hostMsg"
        :cfg="getCurrentCfg()"
        @update:cfg="updateComponentCfg"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Overview from './components/Overview.vue'
import CPU from './components/CPU.vue'
import Memory from './components/Memery.vue'
import Disk from './components/Disk.vue'

const selectedMenu = ref('overview')

// 组件映射对象
const componentMap = {
  overview: Overview,
  cpu: CPU,
  memory: Memory,
  disk: Disk,
}

// 主机信息
const hostMsg = ref({
  hostCpuCount: 16,
  hostMemory: 4096,
})

// 使用ref定义按钮组
const btnGroup = ref([
  {
    cfg: {
      name: 'ubuntu25.10',
      uuid: 'f0715790-cde5-4faf-8869-d8d72dfaf7d8',
      osMachine: 'pc',
      osFirmware: 'bios',
    },
    name: '概况',
    type: 'overview',
  },
  {
    cfg: {
      cpuCount: 2,
      cpuMode: 'host-passthrough',
      isNotManualTopology: false,
      manualTopology: { sockets: 2, cores: 1, threads: 1 },
    },
    name: 'CPU数',
    type: 'cpu',
  },
  { cfg: { memory: 4194304, currentMemory: 4194304 }, name: '内存', type: 'memory' },
  {
    cfg: {
      diskType: 'disk',
      sourcePath: '/var/lib/libvirt/images/ubuntu25.10.qcow2',
      diskFormat: 'qcow2',
      targetDev: 'vda',
      targetBus: 'virtio',
      isReadOnly: false,
    },
    name: '磁盘',
    type: 'disk',
  },
  { cfg: {}, name: 'CDROM', type: 'cdrom' },
  { cfg: {}, name: '虚拟网络', type: 'network' },
  { cfg: {}, name: '显示协议', type: 'display' },
  { cfg: {}, name: '声卡', type: 'sound' },
  { cfg: {}, name: '控制台', type: 'console' },
  { cfg: {}, name: '通道', type: 'channel-qemu' },
  { cfg: {}, name: '通道', type: 'channel-spice' },
  { cfg: {}, name: '视频', type: 'video' },
  { cfg: {}, name: '控制器USB', type: 'controller-usb' },
  { cfg: {}, name: '控制器PCIe', type: 'controller-pcie' },
  { cfg: {}, name: 'USB转发器1', type: 'usb-redir1' },
  { cfg: {}, name: 'USB转发器2', type: 'usb-redir2' },
  { cfg: {}, name: 'RNG /dev/urandom', type: 'rng' },
])

// 获取当前选中组件的配置
const getCurrentCfg = () => {
  const currentItem = btnGroup.value.find((item) => item.type === selectedMenu.value)
  return currentItem ? currentItem.cfg : {}
}

// 更新组件配置
const updateComponentCfg = (newCfg) => {
  const currentIndex = btnGroup.value.findIndex((item) => item.type === selectedMenu.value)
  if (currentIndex !== -1) {
    btnGroup.value[currentIndex].cfg = newCfg
  }
}
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
