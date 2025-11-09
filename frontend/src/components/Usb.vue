<template>
  <div class="usb-list-container">
    <h2>USB 设备列表</h2>
    <button @click="fetchUsbList" :disabled="loading">
      {{ loading ? '加载中...' : '刷新列表' }}
    </button>
    <div v-if="error" class="error-message">
      <p>错误: {{ error }}</p>
    </div>
    <div v-if="usbList.length > 0" class="usb-list">
      <ul>
        <li v-for="usb in usbList" :key="usb.device">
          <strong>ID:</strong> {{ usb.id }}
          <br />
          <strong>设备:</strong> {{ usb.info }}
        </li>
      </ul>
    </div>
    <p v-else-if="!loading">没有找到USB设备。</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { listUsb } from '@/api/usb.js'

const usbList = ref([])
const loading = ref(false)
const error = ref(null)

async function fetchUsbList() {
  loading.value = true
  error.value = null
  try {
    const response = await listUsb()
    usbList.value = response || []
  } catch (e) {
    error.value = e.message
    usbList.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchUsbList()
})
</script>

<style scoped>
.usb-list-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 20px;
  font-family: sans-serif;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  padding: 10px;
  border-radius: 5px;
  margin-top: 15px;
}

.usb-list ul {
  list-style-type: none;
  padding: 0;
}

.usb-list li {
  padding: 8px 0;
  border-bottom: 1px solid #eee;
}

button {
  padding: 10px 15px;
  border: none;
  background-color: #007bff;
  color: white;
  border-radius: 5px;
  cursor: pointer;
  margin-bottom: 15px;
}

button:hover {
  background-color: #0056b3;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
</style>
