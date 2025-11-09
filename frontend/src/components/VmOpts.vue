<template>
  <div class="vm-opts-container">
    <h1>虚拟机操作</h1>

    <!-- 通用操作区域 -->
    <div class="section">
      <h2>通用操作</h2>
      <div class="form-group">
        <label for="vmName">虚拟机名称:</label>
        <input type="text" id="vmName" v-model="vmName" placeholder="请输入虚拟机名称" />
      </div>
      <div class="button-group">
        <button @click="performAction('start')">启动</button>
        <button @click="performAction('stop')">停止</button>
        <button @click="performAction('suspend')">挂起</button>
        <button @click="performAction('resume')">恢复</button>
        <button @click="performAction('forceShutdown')">强制关闭</button>
        <button @click="performAction('delete')">删除</button>
        <button @click="performAction('getInfo')">获取信息</button>
      </div>
    </div>

    <!-- 虚拟机列表 -->
    <div class="section">
      <h2>虚拟机列表</h2>
      <div class="form-group">
        <label for="listType">列表类型:</label>
        <select id="listType" v-model="listType">
          <option value="all">所有</option>
          <option value="active">活动</option>
          <option value="inactive">非活动</option>
        </select>
        <button @click="fetchVMList">刷新列表</button>
      </div>
      <div v-if="vmList.length > 0" class="vm-list">
        <ul>
          <li v-for="vm in vmList" :key="vm.name">{{ vm.name }} ({{ vm.state }})</li>
        </ul>
      </div>
      <p v-else>没有找到虚拟机。</p>
    </div>

    <!-- 创建虚拟机 -->
    <div class="section">
      <h2>创建新虚拟机</h2>
      <div class="form-group">
        <label for="newVmName">新虚拟机名称:</label>
        <input type="text" id="newVmName" v-model="newVmName" placeholder="请输入新虚拟机名称" />
      </div>
      <div class="form-group">
        <label for="xmlConfig">XML 配置:</label>
        <textarea
          id="xmlConfig"
          v-model="xmlConfig"
          rows="10"
          placeholder="在此处粘贴 XML 配置"
        ></textarea>
      </div>
      <button @click="createVMFromXML">创建虚拟机</button>
    </div>

    <!-- 结果显示 -->
    <div class="section result-section">
      <h2>操作结果</h2>
      <div v-if="loading" class="loading">处理中...</div>
      <div v-if="error" class="error-message">
        <p>错误:</p>
        <pre>{{ error }}</pre>
      </div>
      <div v-if="result" class="result-message">
        <p>成功:</p>
        <pre>{{ result }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  startVM,
  stopVM,
  suspendVM,
  resumeVM,
  forceShutdownVM,
  deleteVM,
  listVMs,
  getVMInfo,
  createVM,
} from '@/api/vmopts.js'

const vmName = ref('')
const newVmName = ref('')
const xmlConfig = ref('')
const listType = ref('all')
const vmList = ref([])
const result = ref(null)
const error = ref(null)
const loading = ref(false)

const actions = {
  start: startVM,
  stop: stopVM,
  suspend: suspendVM,
  resume: resumeVM,
  forceShutdown: forceShutdownVM,
  delete: deleteVM,
  getInfo: getVMInfo,
}

async function performAction(action) {
  if (!vmName.value && action !== 'list') {
    error.value = '请输入虚拟机名称。'
    result.value = null
    return
  }
  loading.value = true
  error.value = null
  result.value = null
  try {
    const response = await actions[action](vmName.value)
    result.value = response
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
    if (action === 'start' || action === 'stop' || action === 'delete') {
      fetchVMList() // Refresh list after state-changing actions
    }
  }
}

async function fetchVMList() {
  loading.value = true
  error.value = null
  result.value = null
  try {
    const response = await listVMs(listType.value)
    vmList.value = response.vms || []
  } catch (e) {
    error.value = e.message
    vmList.value = []
  } finally {
    loading.value = false
  }
}

async function createVMFromXML() {
  if (!newVmName.value || !xmlConfig.value) {
    error.value = '新虚拟机名称和 XML 配置不能为空。'
    result.value = null
    return
  }
  loading.value = true
  error.value = null
  result.value = null
  try {
    const response = await createVM(newVmName.value, xmlConfig.value)
    result.value = response
    newVmName.value = ''
    xmlConfig.value = ''
    fetchVMList() // Refresh list after creating a new VM
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchVMList()
})
</script>

<style scoped>
.vm-opts-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: sans-serif;
}

.section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 5px;
}


.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
}

input[type='text'],
select,
textarea {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.button-group button {
  margin-right: 10px;
  margin-bottom: 10px;
}

button {
  padding: 10px 15px;
  border: none;
  background-color: #007bff;
  color: white;
  border-radius: 5px;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}

.vm-list ul {
  list-style-type: none;
  padding: 0;
}

.vm-list li {
  padding: 5px 0;
}

.result-section .loading,
.result-section .error-message,
.result-section .result-message {
  margin-top: 15px;
  padding: 10px;
  border-radius: 5px;
}

.loading {
  color: #555;
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.result-message {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
