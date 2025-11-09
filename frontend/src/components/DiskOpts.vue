<template>
  <div class="disk-opts-container">
    <h1>虚拟硬盘管理</h1>

    <!-- 硬盘列表 -->
    <div class="section">
      <h2>硬盘列表</h2>
      <button @click="fetchDiskList">刷新列表</button>
      <div v-if="diskList.length > 0" class="disk-list">
        <ul>
          <li v-for="disk in diskList" :key="disk.name">
            <span>{{ disk.name }} - {{ disk.size }}</span>
            <div class="button-group">
              <button @click="showDiskInfo(disk.name)">信息</button>
              <button @click="prepareResize(disk.name)">调整大小</button>
              <button @click="confirmDelete(disk.name)">删除</button>
            </div>
          </li>
        </ul>
      </div>
      <p v-else>没有找到虚拟硬盘。</p>
    </div>

    <!-- 添加硬盘 -->
    <div class="section">
      <h2>添加新硬盘</h2>
      <div class="form-group">
        <label for="newDiskName">名称:</label>
        <input type="text" id="newDiskName" v-model="newDisk.name" placeholder="例如: my-disk" />
      </div>
      <div class="form-group">
        <label for="newDiskFormat">格式:</label>
        <input type="text" id="newDiskFormat" v-model="newDisk.format" placeholder="例如: qcow2" />
      </div>
      <div class="form-group">
        <label for="newDiskSize">大小:</label>
        <input type="text" id="newDiskSize" v-model="newDisk.size" placeholder="例如: 10G" />
      </div>
      <button @click="addNewDisk">添加硬盘</button>
    </div>

    <!-- 调整大小 -->
    <div v-if="resizingDisk.name" class="section">
      <h2>调整硬盘大小: {{ resizingDisk.name }}</h2>
      <div class="form-group">
        <label for="newSize">新大小:</label>
        <input type="text" id="newSize" v-model="resizingDisk.newSize" placeholder="例如: 20G" />
      </div>
      <button @click="executeResize">确认调整</button>
      <button @click="cancelResize">取消</button>
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
import { ref, reactive, onMounted } from 'vue';
import {
  listDisks,
  addDisk,
  getDiskInfo,
  resizeDisk,
  deleteDisk,
} from '@/api/diskopts.js';

const diskList = ref([]);
const newDisk = reactive({
  name: '',
  format: 'qcow2',
  size: '10G',
});
const resizingDisk = reactive({
  name: '',
  newSize: '',
});
const result = ref(null);
const error = ref(null);
const loading = ref(false);

async function fetchDiskList() {
  loading.value = true;
  error.value = null;
  try {
    const response = await listDisks();
    diskList.value = response.disks || [];
  } catch (e) {
    error.value = e.message;
    diskList.value = [];
  } finally {
    loading.value = false;
  }
}

async function addNewDisk() {
  if (!newDisk.name || !newDisk.format || !newDisk.size) {
    error.value = '名称、格式和大小不能为空。';
    return;
  }
  loading.value = true;
  error.value = null;
  result.value = null;
  try {
    const response = await addDisk(newDisk.name, newDisk.format, newDisk.size);
    result.value = response;
    newDisk.name = ''; // Reset form
    fetchDiskList();
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}

async function showDiskInfo(name) {
  loading.value = true;
  error.value = null;
  result.value = null;
  try {
    const response = await getDiskInfo(name);
    result.value = response.diskInfo;
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}

function prepareResize(name) {
  resizingDisk.name = name;
  resizingDisk.newSize = '';
}

function cancelResize() {
  resizingDisk.name = '';
  resizingDisk.newSize = '';
}

async function executeResize() {
  if (!resizingDisk.newSize) {
    error.value = '请输入新的大小。';
    return;
  }
  loading.value = true;
  error.value = null;
  result.value = null;
  try {
    const response = await resizeDisk(resizingDisk.name, resizingDisk.newSize);
    result.value = response;
    cancelResize();
    fetchDiskList();
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}

function confirmDelete(name) {
  if (window.confirm(`您确定要删除虚拟硬盘 "${name}" 吗？此操作不可撤销。`)) {
    executeDelete(name);
  }
}

async function executeDelete(name) {
  loading.value = true;
  error.value = null;
  result.value = null;
  try {
    const response = await deleteDisk(name);
    result.value = response;
    fetchDiskList();
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchDiskList();
});
</script>

<style scoped>
.disk-opts-container {
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

input[type="text"] {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.button-group {
  display: inline-block;
  margin-left: 20px;
}

.button-group button {
  margin-left: 10px;
}

button {
  padding: 8px 12px;
  border: none;
  background-color: #007bff;
  color: white;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 10px;
}

button:hover {
  background-color: #0056b3;
}

.disk-list ul {
  list-style-type: none;
  padding: 0;
}

.disk-list li {
  padding: 10px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
