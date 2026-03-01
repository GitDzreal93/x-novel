import axios from 'axios';
import type { AxiosRequestConfig } from 'axios';

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 60000,
});

instance.interceptors.request.use(
  (config) => {
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

instance.interceptors.response.use(
  (response) => {
    if (response.headers['x-device-id']) {
      localStorage.setItem('x-novel-device-id', response.headers['x-device-id']);
    }
    return response.data;
  },
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

function generateDeviceId(): string {
  return 'device_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

/**
 * 封装 axios 实例，使返回类型与拦截器一致（拦截器返回 response.data）。
 */
const request = {
  get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.get(url, config) as unknown as Promise<T>;
  },
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.post(url, data, config) as unknown as Promise<T>;
  },
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.put(url, data, config) as unknown as Promise<T>;
  },
  delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.delete(url, config) as unknown as Promise<T>;
  },
};

export default request;
