const API_BASE_URL = '/api'; // 假设后端API服务在 /api 路径下

// 统一处理 fetch 请求的响应
async function handleResponse(response) {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: '解析错误响应失败' }));
    throw new Error(errorData.error || `HTTP 错误！状态: ${response.status}`);
  }
  return response.json();
}

// 启动虚拟机
export function startVM(name) {
  return fetch(`${API_BASE_URL}/vm/start?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 停止虚拟机
export function stopVM(name) {
  return fetch(`${API_BASE_URL}/vm/stop?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 挂起虚拟机
export function suspendVM(name) {
  return fetch(`${API_BASE_URL}/vm/suspend?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 恢复虚拟机
export function resumeVM(name) {
  return fetch(`${API_BASE_URL}/vm/resume?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 强制关闭虚拟机
export function forceShutdownVM(name) {
  return fetch(`${API_BASE_URL}/vm/force-shutdown?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 删除虚拟机
export function deleteVM(name) {
  return fetch(`${API_BASE_URL}/vm/delete?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 获取虚拟机列表
// type 可以是 'all', 'active', 或 'inactive'
export function listVMs(type = 'all') {
  return fetch(`${API_BASE_URL}/vm/list?type=${encodeURIComponent(type)}`).then(handleResponse);
}

// 获取虚拟机详细信息
export function getVMInfo(name) {
  return fetch(`${API_BASE_URL}/vm/info?name=${encodeURIComponent(name)}`).then(handleResponse);
}

// 通过 XML 配置创建新虚拟机
export function createVM(name, xmlConfig) {
  const formData = new FormData();
  formData.append('name', name);
  formData.append('xmlConfig', xmlConfig);

  return fetch(`${API_BASE_URL}/vm/create`, {
    method: 'POST',
    body: formData,
  }).then(handleResponse);
}

// 挂载USB设备
export function attachUsb(name, usbId) {
  return fetch(`${API_BASE_URL}/vm/attach-usb?name=${encodeURIComponent(name)}&usbId=${encodeURIComponent(usbId)}`).then(handleResponse);
}
// 卸载USB设备
export function detachUsb(name, usbId) {
  return fetch(`${API_BASE_URL}/vm/detach-usb?name=${encodeURIComponent(name)}&usbId=${encodeURIComponent(usbId)}`).then(handleResponse);
}
