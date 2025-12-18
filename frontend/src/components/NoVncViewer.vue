<template>
  <div class="novnc-container">
    <!-- 1. noVNC 渲染容器（必须指定宽高） -->
    <div ref="vncContainer" class="vnc-viewer"></div>

    <!-- 2. 连接控制区域（输入参数 + 按钮） -->
    <div class="vnc-controls">
      <input
        v-model="vncConfig.host"
        placeholder="VNC 服务器地址（如 192.168.1.100）"
        class="input"
      />
      <input
        v-model="vncConfig.wsPort"
        type="number"
        placeholder="WebSocket 端口（如 5700）"
        class="input"
      />
      <input
        v-model="vncConfig.password"
        type="password"
        placeholder="VNC 密码（可选）"
        class="input"
      />
      <button @click="connectVNC" :disabled="isConnecting">
        {{ isConnecting ? '连接中...' : '连接远程桌面' }}
      </button>
      <button @click="disconnectVNC" :disabled="!isConnected">断开连接</button>
    </div>

    <!-- 3. 连接状态提示 -->
    <div class="vnc-status" :class="statusClass">
      {{ statusText }}
    </div>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
// 引入 noVNC 核心类和样式
import RFB from '@novnc/novnc/lib/rfb'

// 1. 状态管理
const vncContainer = ref(null) // 渲染容器的 DOM 引用
const rfbInstance = ref(null) // RFB 实例（管理 VNC 连接）
const isConnecting = ref(false) // 是否正在连接
const isConnected = ref(false) // 是否已连接
const statusText = ref('未连接') // 状态提示文本
const statusClass = ref('status-idle') // 状态样式（用于颜色区分）

// 2. VNC 连接配置（可根据需求扩展）
const vncConfig = ref({
  host: 'localhost', // VNC 服务器 IP（如宿主机 IP）
  wsPort: 5900, // WebSocket 端口（非原生 VNC 端口）
  password: '', // VNC 密码（若服务器设置了密码）
  path: '/websockify', // WebSocket 路径（默认 /，若用 Nginx 转发需匹配）
})

// 3. 初始化 RFB 实例，建立连接
const connectVNC = async () => {
  if (!vncContainer.value) {
    statusText.value = '渲染容器不存在'
    statusClass.value = 'status-error'
    return
  }

  // 拼接 WebSocket 地址（格式：ws://host:port/path）
  const wsUrl = `ws://${vncConfig.value.host}:${vncConfig.value.wsPort}${vncConfig.value.path}`

  try {
    isConnecting.value = true
    statusText.value = '正在连接远程桌面...'
    statusClass.value = 'status-connecting'

    // 初始化 RFB 实例（核心步骤）
    rfbInstance.value = new RFB(vncContainer.value, wsUrl, {
      // noVNC 配置选项（按需调整）
      credentials: { password: vncConfig.value.password }, // 密码（可选）
      scaleViewport: true, // 自动缩放以适应容器
      showToolbar: true, // 显示默认控制栏（全屏、快捷键、断开等）
      focusOnClick: true, // 点击容器时获取焦点（用于键盘输入）
      dragViewport: true, // 支持拖拽视图（当画面大于容器时）
      renderDuration: 16, // 渲染间隔（毫秒，默认 16，接近 60fps）
    })

    // 4. 绑定 RFB 事件（处理连接状态）
    bindRfbEvents()
  } catch (error) {
    isConnecting.value = false
    statusText.value = `连接失败：${error.message}`
    statusClass.value = 'status-error'
  }
}

// 5. 绑定 RFB 事件（监听连接、断开、错误）
const bindRfbEvents = () => {
  const rfb = rfbInstance.value

  // 连接成功事件
  rfb.addEventListener('connect', () => {
    isConnecting.value = false
    isConnected.value = true
    statusText.value = '已成功连接远程桌面'
    statusClass.value = 'status-success'
    console.log('VNC 连接成功')
  })

  // 连接断开事件（主动/被动断开）
  rfb.addEventListener('disconnect', (event) => {
    isConnecting.value = false
    isConnected.value = false
    // 区分主动断开和错误断开
    const reason = event.detail.clean ? '主动断开连接' : '连接被意外断开'
    statusText.value = `已断开：${reason}`
    statusClass.value = 'status-idle'
    console.log('VNC 连接断开：', reason)
  })

  // 连接错误事件（如密码错误、服务器不可达）
  rfb.addEventListener('error', (event) => {
    isConnecting.value = false
    isConnected.value = false
    statusText.value = `连接错误：${event.detail.message}`
    statusClass.value = 'status-error'
    console.error('VNC 连接错误：', event.detail.message)
  })

  // 密码请求事件（若初始化时未传密码，服务器会主动请求）
  rfb.addEventListener('credentialsrequired', () => {
    // 动态输入密码（示例：用 prompt，实际可替换为组件内输入框）
    const password = prompt('请输入 VNC 密码：')
    if (password) {
      rfb.sendCredentials({ password })
    } else {
      rfb.disconnect() // 取消密码输入，断开连接
    }
  })
}

// 6. 断开 VNC 连接
const disconnectVNC = () => {
  if (rfbInstance.value && isConnected.value) {
    rfbInstance.value.disconnect() // 主动断开
    rfbInstance.value = null // 销毁实例
  }
}

// 7. 组件卸载时清理连接（避免内存泄漏）
onBeforeUnmount(() => {
  disconnectVNC()
})
</script>

<style scoped>
.novnc-container {
  width: 100%;
  max-width: 1200px;
  margin: 20px auto;
  padding: 20px;
  box-sizing: border-box;
}

/* 渲染容器：必须指定固定宽高或百分比，否则画面无法显示 */
.vnc-viewer {
  width: 100%;
  height: 600px; /* 高度可根据需求调整 */
  border: 1px solid #eee;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
}

/* 控制区域样式 */
.vnc-controls {
  display: flex;
  gap: 10px;
  margin: 15px 0;
  flex-wrap: wrap;
}

.input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  flex: 1;
  min-width: 150px;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  background: #42b983;
  color: white;
  cursor: pointer;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

/* 状态提示样式 */
.vnc-status {
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
}

.status-idle {
  background: #f5f5f5;
  color: #666;
}

.status-connecting {
  background: #fff3cd;
  color: #856404;
}

.status-success {
  background: #d4edda;
  color: #155724;
}

.status-error {
  background: #f8d7da;
  color: #721c24;
}
</style>
