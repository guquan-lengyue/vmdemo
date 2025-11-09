const API_BASE_URL = '/api'; // 假设后端API服务在 /api 路径下

// 统一处理 fetch 请求的响应
async function handleResponse(response) {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: '解析错误响应失败' }));
    throw new Error(errorData.error || `HTTP 错误！状态: ${response.status}`);
  }
  // 检查响应是否为空
  const text = await response.text();
  return text ? JSON.parse(text) : {};
}

// 列出所有虚拟硬盘
export function listDisks() {
  return fetch(`${API_BASE_URL}/disk/list`).then(handleResponse);
}

// 添加虚拟硬盘
export function addDisk(name, format, size) {
  return fetch(`${API_BASE_URL}/disk/add`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, format, size }),
  }).then(handleResponse);
}

// 获取虚拟硬盘信息
export function getDiskInfo(name) {
  return fetch(`${API_BASE_URL}/disk/info/${encodeURIComponent(name)}`).then(handleResponse);
}

// 调整虚拟硬盘大小
export function resizeDisk(name, newSize) {
  return fetch(`${API_BASE_URL}/disk/resize`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, newSize }),
  }).then(handleResponse);
}

// 删除虚拟硬盘
export function deleteDisk(name) {
  return fetch(`${API_BASE_URL}/disk/delete/${encodeURIComponent(name)}`, {
    method: 'DELETE',
  }).then(handleResponse);
}
