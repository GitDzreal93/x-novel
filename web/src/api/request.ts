import axios from 'axios';

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 60000,
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 添加设备 ID
    let deviceId = localStorage.getItem('x-novel-device-id');
    if (!deviceId) {
      deviceId = generateDeviceId();
      localStorage.setItem('x-novel-device-id', deviceId);
    }
    config.headers['X-Device-ID'] = deviceId;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data;

    // 如果响应中有设备 ID，保存到本地
    if (response.headers['x-device-id']) {
      localStorage.setItem('x-novel-device-id', response.headers['x-device-id']);
    }

    return res;
  },
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

// 生成设备 ID
function generateDeviceId(): string {
  return 'device_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

export default request;
