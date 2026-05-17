// request.ts — Axios HTTP 请求封装
// Axios 是一个基于 Promise 的 HTTP 客户端，用于浏览器和 Node.js
// 封装后统一处理：baseURL 配置、请求自动携带 token、响应错误提示、401/403 处理等
// 各 API 模块（auth.ts、user.ts 等）都通过此文件发送请求

import axios, { AxiosResponse, type AxiosInstance } from 'axios'
// ElMessage 是 Element Plus 的轻量级消息提示组件，用于在屏幕顶部显示错误提示
import { ElMessage } from 'element-plus'

// 创建 Axios 实例（service），所有请求共享此配置
const service: AxiosInstance = axios.create({
  // baseURL：所有请求 URL 的前缀
  // import.meta.env.VITE_API_BASE_URL 来自项目根目录的 .env 文件
  // 开发环境通常设置为 "/api/v1"（通过 Vite 代理转发到后端）
  // 生产环境可能设置为完整的后端地址
  // 如果环境变量未定义，则使用默认值 "/api/v1"
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  // 请求超时时间：15 秒，超过此时间未收到响应则报错
  timeout: 15000,
})

// ----- 请求拦截器（request interceptor）-----
// 在每次请求发送之前执行，用于在请求头中添加认证令牌
service.interceptors.request.use(
  (config) => {
    // 从 localStorage 读取 JWT 令牌
    const token = localStorage.getItem('token')
    if (token) {
      // 将令牌以 Bearer 格式放入 Authorization 请求头
      // 后端通过解析此请求头来验证用户身份
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// ----- 响应拦截器（response interceptor）-----
// 在收到后端响应之后、返回给业务代码之前执行
service.interceptors.response.use(
  // 第一个回调：处理成功的响应（HTTP 状态码 2xx）
  (response: AxiosResponse) => {
    // 后端统一响应格式为 { code: number, message: string, data: any }
    const { code, message, data } = response.data
    if (code === 0) {
      // code === 0 表示业务成功，只返回 data 部分给调用方
      // 这样业务代码中不需要再处理 response.data.data，直接拿到结果
      return data
    }
    // code !== 0 表示业务失败，用 ElMessage 弹出错误提示
    ElMessage.error(message || '请求失败')
    return Promise.reject(new Error(message || '请求失败'))
  },
  // 第二个回调：处理失败的响应（HTTP 状态码 4xx/5xx）或网络错误
  (error) => {
    if (error.response) {
      // 服务器返回了响应（有 HTTP 状态码）
      const { status, data } = error.response
      if (status === 401 || status === 403) {
        // 401=未认证（未登录或 token 过期），403=无权限
        // 此处不弹出错误提示，也不做页面跳转
        // 因为路由守卫（src/router/index.ts）会检测到无 token 或无角色，自动跳转到登录页或 403 页
        return Promise.reject(error)
      }
      // 其他错误状态码（如 500 服务器错误），显示后端返回的错误消息
      ElMessage.error(data?.message || error.message || '网络错误')
    } else {
      // 网络错误（服务器未响应），例如跨域问题、网络断开等
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  },
)

// ----- 自定义请求类型 -----
// Request 类型定义了四个方法（get/post/put/delete），每个都支持泛型 <T>
// Promise<T> 表示异步操作完成后会返回类型为 T 的数据
// T = any 表示不指定泛型时默认类型为 any
type Request = {
  get<T = any>(url: string, config?: any): Promise<T>
  post<T = any>(url: string, data?: any, config?: any): Promise<T>
  put<T = any>(url: string, data?: any, config?: any): Promise<T>
  delete<T = any>(url: string, config?: any): Promise<T>
}

// 封装后的 request 对象，在 API 模块中使用，例如 request.get<UserListRes>('/user', { params })
// 每个方法都调用 service 的对应方法，然后用 .then(res => res as unknown as T) 提取数据
// service 的响应拦截器已经返回了 response.data.data，这里再做一层类型转换
const request: Request = {
  get: <T>(url: string, config?: any) => service.get<T>(url, config).then(res => res as unknown as T),
  post: <T>(url: string, data?: any, config?: any) => service.post<T>(url, data, config).then(res => res as unknown as T),
  put: <T>(url: string, data?: any, config?: any) => service.put<T>(url, data, config).then(res => res as unknown as T),
  delete: <T>(url: string, config?: any) => service.delete<T>(url, config).then(res => res as unknown as T),
}

// 默认导出，各 API 文件通过 import request from '@/utils/request' 使用
export default request
