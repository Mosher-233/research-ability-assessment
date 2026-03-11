import axios from 'axios'
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '../types/user'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const authApi = {
  // 登录
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/login', data)
    return response.data
  },
  
  // 注册
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/register', data)
    return response.data
  },
  
  // 获取用户信息
  getUserInfo: async (): Promise<{ code: number; message: string; data: User }> => {
    const response = await api.get('/user/info')
    return response.data
  }
}

export default api