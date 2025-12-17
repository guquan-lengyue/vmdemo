// API 对接文件，使用 fetch 实现所有后端接口调用

const API_BASE_URL = 'http://localhost:5173/api';

// 封装 fetch 请求
async function request(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const mergedOptions = { ...defaultOptions, ...options };

  try {
    const response = await fetch(url, mergedOptions);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
}

// 虚拟机操作相关 API
export const vmApi = {
  // 启动虚拟机
  startVM: (vmName) => request(`/vm/start?name=${vmName}`),

  // 停止虚拟机
  stopVM: (vmName) => request(`/vm/stop?name=${vmName}`),

  // 挂起虚拟机
  suspendVM: (vmName) => request(`/vm/suspend?name=${vmName}`),

  // 恢复虚拟机
  resumeVM: (vmName) => request(`/vm/resume?name=${vmName}`),

  // 强制关闭虚拟机
  forceShutdownVM: (vmName) => request(`/vm/force-shutdown?name=${vmName}`),

  // 删除虚拟机
  deleteVM: (vmName) => request(`/vm/delete?name=${vmName}`),

  // 获取虚拟机列表
  listVMs: (type = 'all') => request(`/vm/list?type=${type}`),

  // 获取虚拟机详细信息
  getVMInfo: (vmName) => request(`/vm/info?name=${vmName}`),

  // 创建虚拟机
  createVM: (vmName, xmlConfig) => {
    const formData = new FormData();
    formData.append('name', vmName);
    formData.append('xmlConfig', xmlConfig);

    return fetch(`${API_BASE_URL}/vm/create`, {
      method: 'POST',
      body: formData,
    }).then(response => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
    });
  },

  // 更新虚拟机配置
  updateVM: (vmName, xmlConfig) => {
    const formData = new FormData();
    formData.append('name', vmName);
    formData.append('xmlConfig', xmlConfig);

    return fetch(`${API_BASE_URL}/vm/update`, {
      method: 'POST',
      body: formData,
    }).then(response => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
    });
  },

  // 为虚拟机添加USB设备
  attachUsbDevice: (vmName, usbId) => request(`/vm/attach-usb?name=${vmName}&usbId=${usbId}`),

  // 为虚拟机移除USB设备
  detachUsbDevice: (vmName, usbId) => request(`/vm/detach-usb?name=${vmName}&usbId=${usbId}`),
};

// 磁盘操作相关 API
export const diskApi = {
  // 列出硬盘池中的所有虚拟硬盘
  listDisks: () => request(`/disk/list`),

  // 添加虚拟硬盘到硬盘池
  addDisk: (name, format, size) => request(`/disk/add`, {
    method: 'POST',
    body: JSON.stringify({ name, format, size }),
  }),

  // 获取硬盘池中某个虚拟硬盘的信息
  getDiskInfo: (name) => request(`/disk/info/${name}`),

  // 调整硬盘池中某个虚拟硬盘的大小
  resizeDisk: (name, newSize) => request(`/disk/resize`, {
    method: 'POST',
    body: JSON.stringify({ name, newSize }),
  }),

  // 从硬盘池中删除虚拟硬盘
  deleteDisk: (name) => request(`/disk/delete/${name}`, {
    method: 'DELETE',
  }),
};

// VNC 相关 API
export const vncApi = {
  // 获取虚拟机的 VNC WebSocket 连接
  getVncWs: (vmName) => `${API_BASE_URL}/vnc/${vmName}`,
};

// USB 操作相关 API
export const usbApi = {
  // 获取 USB 设备列表
  getUsbList: () => request(`/usb/list`),
};

// UUID 相关 API
export const uuidApi = {
  // 生成 UUID
  generateUUID: () => request(`/uuid/generate`),
};

export default {
  vm: vmApi,
  disk: diskApi,
  vnc: vncApi,
  usb: usbApi,
  uuid: uuidApi,
};
