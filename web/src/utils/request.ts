import axios, { AxiosResponse, type AxiosInstance } from 'axios'
import { ElMessage } from 'element-plus'

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
})

service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

service.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data
    if (code === 0) {
      return data
    }
    ElMessage.error(message || '请求失败')
    return Promise.reject(new Error(message || '请求失败'))
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      } else if (status === 403) {
        window.location.href = '/403'
      } else {
        ElMessage.error(data?.message || error.message || '网络错误')
      }
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  },
)

type Request = {
  get<T = any>(url: string, config?: any): Promise<T>
  post<T = any>(url: string, data?: any, config?: any): Promise<T>
  put<T = any>(url: string, data?: any, config?: any): Promise<T>
  delete<T = any>(url: string, config?: any): Promise<T>
}

const request: Request = {
  get: <T>(url: string, config?: any) => service.get<T>(url, config).then(res => res as unknown as T),
  post: <T>(url: string, data?: any, config?: any) => service.post<T>(url, data, config).then(res => res as unknown as T),
  put: <T>(url: string, data?: any, config?: any) => service.put<T>(url, data, config).then(res => res as unknown as T),
  delete: <T>(url: string, config?: any) => service.delete<T>(url, config).then(res => res as unknown as T),
}

export default request
