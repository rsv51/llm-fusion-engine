// 统一的 API 客户端

import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '../types'

// 创建 axios 实例
const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加认证 token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 统一错误处理
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error: AxiosError) => {
    // 处理 401 未授权错误
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
      return Promise.reject(new Error('未授权,请重新登录'))
    }

    // 处理 403 禁止访问错误
    if (error.response?.status === 403) {
      return Promise.reject(new Error('权限不足'))
    }

    // 处理 404 未找到错误
    if (error.response?.status === 404) {
      return Promise.reject(new Error('资源未找到'))
    }

    // 处理 500 服务器错误
    if (error.response?.status === 500) {
      return Promise.reject(new Error('服务器内部错误'))
    }

    // 处理网络错误
    if (!error.response) {
      return Promise.reject(new Error('网络连接失败,请检查网络'))
    }

    // 其他错误
    const message = (error.response?.data as ApiResponse)?.message || error.message || '请求失败'
    return Promise.reject(new Error(message))
  }
)

export const excelApi = {
  // 导出 Excel
  async exportToExcel(): Promise<Blob> {
    return api.get('/admin/export/excel', { responseType: 'blob' })
  },

  // 导入 Excel
  async importFromExcel(file: File): Promise<{ message: string }> {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/admin/import/excel', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export default api