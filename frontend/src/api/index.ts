import axios from 'axios';
import { useAuthStore } from '@/store/auth';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(config => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// 响应拦截器
apiClient.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    // 统一错误处理
    if (error.response) {
      const status = error.response.status;
      
      if (status === 401) {
        // 未授权，跳转到登录页
        useAuthStore.getState().logout();
        window.location.href = '/login';
      }
      
      return Promise.reject({
        status,
        message: error.response.data?.message || '请求失败',
      });
    }
    
    return Promise.reject({
      status: 500,
      message: '网络错误，请检查网络连接',
    });
  }
);

export default apiClient;
