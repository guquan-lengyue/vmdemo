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

// 列出所有usb设备
export function listUsb() {
  return fetch(`${API_BASE_URL}/usb/list`).then(handleResponse);
}
