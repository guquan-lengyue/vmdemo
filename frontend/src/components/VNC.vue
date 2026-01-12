<template>
  <div class="vnc-container">
    <div id="vnc-viewer" ref="vncViewer"></div>
    <div class="vnc-controls">
      <button @click="connect" :disabled="isConnected || loading">连接</button>
      <button @click="disconnect" :disabled="!isConnected">断开</button>
      <button @click="sendCtrlAltDel">Ctrl+Alt+Del</button>
    </div>
  </div>
</template>

<script>
// 从npm包导入noVNC
import RFB from '@novnc/novnc/lib/rfb';

const loadNoVNC = async () => {
  return { RFB };
};

export default {
  name: 'VNC',
  props: {
    vmName: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      rfb: null,
      isConnected: false,
      loading: false
    };
  },
  mounted() {
    // 组件挂载后自动连接
    this.connect();
  },
  beforeUnmount() {
    // 组件卸载前断开连接
    this.disconnect();
  },
  methods: {
    async connect() {
      const viewer = this.$refs.vncViewer;
      if (!viewer) return;

      this.loading = true;

      try {
        // 加载noVNC库
        const NoVNC = await loadNoVNC();

        // 清除之前的视图
        viewer.innerHTML = '';

        // 获取WebSocket URL
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        const port = window.location.port ? `:${window.location.port}` : '';
        const wsUrl = `${protocol}//${host}${port}/api/vnc/${this.vmName}`;

        // 创建RFB实例 (注意：CDN版本的API可能略有不同)
        const RFB = window.RFB || (NoVNC && NoVNC.RFB);
        if (!RFB) {
          throw new Error('无法找到RFB类');
        }

        this.rfb = new RFB(viewer, wsUrl, {
          credentials: {
            password: '' // 如果需要密码，可以在这里设置
          },
          // 其他选项
          shared: true,
          viewOnly: false
        });

        // 设置事件处理
        this.rfb.addEventListener('connect', () => {
          console.log('VNC连接成功');
          this.isConnected = true;
        });

        this.rfb.addEventListener('disconnect', () => {
          console.log('VNC连接断开');
          this.isConnected = false;
        });

        this.rfb.addEventListener('credentialsrequired', () => {
          console.log('需要VNC密码');
          // 这里可以添加密码输入框
        });

        this.rfb.addEventListener('error', (e) => {
          console.error('VNC连接错误:', e.detail || e);
          this.isConnected = false;
        });
      } catch (error) {
        console.error('创建VNC连接失败:', error);
      } finally {
        this.loading = false;
      }
    },
    disconnect() {
      if (this.rfb) {
        this.rfb.disconnect();
        this.rfb = null;
        this.isConnected = false;
      }
    },
    sendCtrlAltDel() {
      if (this.rfb && this.isConnected) {
        this.rfb.sendCtrlAltDel();
      }
    }
  }
};
</script>

<style scoped>
.vnc-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
}

#vnc-viewer {
  flex: 1;
  border: 1px solid #ccc;
  overflow: hidden;
  background-color: #000;
}

.vnc-controls {
  padding: 10px;
  background-color: #f5f5f5;
  border-top: 1px solid #ccc;
  display: flex;
  gap: 10px;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  background-color: #42b983;
  color: white;
  cursor: pointer;
  font-size: 14px;
}

button:hover {
  background-color: #3aa876;
}

button:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}
</style>
