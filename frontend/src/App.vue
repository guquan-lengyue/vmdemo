<template>
  <div class="vm-management">
    <h1>虚拟机管理</h1>

    <!-- 操作按钮区域 -->
    <div class="action-buttons">
      <button @click="showCreateModal = true" class="create-btn">创建虚拟机</button>
    </div>

    <!-- 虚拟机列表 -->
    <div class="vm-list">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="vms.length === 0" class="empty-state">暂无虚拟机</div>
      <div v-else class="vm-cards">
        <div v-for="vm in vms" :key="vm.name" class="vm-card">
          <div class="vm-header">
            <h3>{{ vm.name }}</h3>
            <span :class="['status-indicator', (vm.state || 'unknown').toLowerCase()]">
              {{ vm.state || '未知状态' }}
            </span>
          </div>
          <div class="vm-actions">
            <button @click="handleStart(vm.name)" :disabled="vm.state === 'running'" class="btn start-btn">启动</button>
            <button @click="handleStop(vm.name)" :disabled="vm.state === 'shut off'" class="btn stop-btn">停止</button>
            <button @click="handleForceShutdown(vm.name)" :disabled="vm.state !== 'running'" class="btn force-shutdown-btn">强制关机</button>
            <button @click="handleSuspend(vm.name)" :disabled="vm.state !== 'running'" class="btn suspend-btn">挂起</button>
            <button @click="handleResume(vm.name)" :disabled="vm.state !== 'paused'" class="btn resume-btn">恢复</button>
            <button @click="handleEdit(vm.name)" class="btn edit-btn">编辑</button>
            <button @click="handleDelete(vm.name)" class="btn delete-btn">删除</button>
            <button @click="handleVncConnect(vm.name)" :disabled="vm.state !== 'running'" class="btn vnc-btn">连接VNC</button>
          </div>
        </div>
      </div>
    </div>



    <!-- 创建虚拟机模态框 -->
    <div v-if="showCreateModal" class="modal-overlay create-vm-modal" @click="showCreateModal = false">
      <div class="modal-content" @click.stop>
        <EditVm :vm-name="'new'" @vm-updated="handleCreateModalUpdated" />
      </div>
    </div>

    <!-- 编辑虚拟机模态框 -->
    <div v-if="showEditModal" class="modal-overlay create-vm-modal" @click="showEditModal = false">
      <div class="modal-content" @click.stop>
        <EditVm :vm-name="editingVm" @vm-updated="handleVmUpdated" />
      </div>
    </div>

    <!-- VNC连接模态框 -->
    <div v-if="showVncModal" class="modal-overlay vnc-modal" @click="showVncModal = false">
      <div class="modal-content vnc-content" @click.stop>
        <div class="vnc-header">
          <h3>VNC - {{ connectedVmName }}</h3>
          <button @click="showVncModal = false" class="close-btn">×</button>
        </div>
        <VNC :vm-name="connectedVmName" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { vmApi } from './api.js'
import EditVm from './EditVm.vue'
import VNC from './components/VNC.vue'

// 数据状态
const vms = ref([])
const loading = ref(false)
const editingVm = ref('')
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showVncModal = ref(false)
const connectedVmName = ref('')

// 获取虚拟机列表
const fetchVMs = async () => {
  loading.value = true
  try {
    const data = await vmApi.listVMs()
    // 后端返回的数据结构是{"vms": [虚拟机列表]}，需要正确提取
    vms.value = data.vms || []
  } catch (error) {
    console.error('获取虚拟机列表失败:', error)
    vms.value = [] // 出错时确保vms是数组
  } finally {
    loading.value = false
  }
}

// 生命周期钩子 - 页面加载时获取虚拟机列表
onMounted(() => {
  fetchVMs()
})

// 操作处理函数
const handleStart = async (vmName) => {
  try {
    await vmApi.startVM(vmName)
    await fetchVMs() // 刷新列表
  } catch (error) {
    console.error('启动虚拟机失败:', error)
  }
}

const handleStop = async (vmName) => {
  try {
    await vmApi.stopVM(vmName)
    await fetchVMs() // 刷新列表
  } catch (error) {
    console.error('停止虚拟机失败:', error)
  }
}

const handleForceShutdown = async (vmName) => {
  if (confirm(`确定要强制关闭虚拟机 ${vmName} 吗？此操作可能会导致数据丢失！`)) {
    try {
      await vmApi.forceShutdownVM(vmName)
      await fetchVMs() // 刷新列表
    } catch (error) {
      console.error('强制关闭虚拟机失败:', error)
    }
  }
}

const handleSuspend = async (vmName) => {
  try {
    await vmApi.suspendVM(vmName)
    await fetchVMs() // 刷新列表
  } catch (error) {
    console.error('挂起虚拟机失败:', error)
  }
}

const handleResume = async (vmName) => {
  try {
    await vmApi.resumeVM(vmName)
    await fetchVMs() // 刷新列表
  } catch (error) {
    console.error('恢复虚拟机失败:', error)
  }
}

const handleDelete = async (vmName) => {
  if (confirm(`确定要删除虚拟机 ${vmName} 吗？`)) {
    try {
      await vmApi.deleteVM(vmName)
      await fetchVMs() // 刷新列表
    } catch (error) {
      console.error('删除虚拟机失败:', error)
    }
  }
}

const handleEdit = (vmName) => {
  editingVm.value = vmName
  showEditModal.value = true
}

const handleVmUpdated = () => {
  editingVm.value = ''
  showEditModal.value = false
  fetchVMs() // 刷新列表
}

const handleCreateModalUpdated = () => {
  showCreateModal.value = false
  fetchVMs() // 刷新列表
}

const handleVncConnect = (vmName) => {
  connectedVmName.value = vmName
  showVncModal.value = true
}
</script>

<style scoped>
.vm-management {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.action-buttons {
  margin-bottom: 20px;
}

.create-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

.create-btn:hover {
  background-color: #45a049;
}

.vm-list {
  margin-top: 20px;
}

.loading, .empty-state {
  text-align: center;
  padding: 40px;
  color: #666;
}

.vm-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.vm-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.vm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.vm-header h3 {
  margin: 0;
  font-size: 18px;
}

.status-indicator {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}

.status-indicator.running {
  background-color: #4CAF50;
  color: white;
}

.status-indicator.shutoff {
  background-color: #f44336;
  color: white;
}

.status-indicator.paused {
  background-color: #ff9800;
  color: white;
}

.status-indicator.suspended {
  background-color: #9c27b0;
  color: white;
}

.vm-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  flex: 1;
  min-width: 60px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.start-btn {
  background-color: #4CAF50;
  color: white;
}

.stop-btn {
  background-color: #f44336;
  color: white;
}

.force-shutdown-btn {
  background-color: #d32f2f;
  color: white;
  border: 1px solid #b71c1c;
}

.force-shutdown-btn:hover {
  background-color: #b71c1c;
}

.suspend-btn {
  background-color: #ff9800;
  color: white;
}

.resume-btn {
  background-color: #2196F3;
  color: white;
}

.edit-btn {
  background-color: #9E9E9E;
  color: white;
}

.delete-btn {
  background-color: #757575;
  color: white;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

/* 基础模态框内容样式 - 会被特定模态框样式覆盖 */
.modal-content {
  background-color: white;
  padding: 30px;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
}

/* 创建虚拟机模态框特定样式 */
.create-vm-modal .modal-content {
  width: 65vw;
  max-width: 65vw;
  max-height: 80vh;
  padding: 20px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* VNC模态框特定样式 */
.vnc-modal .modal-content {
  width: 95vw;
  max-width: 95vw;
  height: 95vh;
  max-height: 95vh;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.vnc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
}

.vnc-header h3 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.close-btn:hover {
  color: #333;
}

.vnc-btn {
  background-color: #2196F3;
  color: white;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 30px;
}

.cancel-btn {
  background-color: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
}
</style>
