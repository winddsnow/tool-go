// auth.ts — 认证相关的 API 接口
// request 是封装好的 Axios 实例（见 src/utils/request.ts），会自动处理 token、错误提示等
import request from '@/utils/request'

// TypeScript 接口：定义请求体 / 响应体的数据结构
// 接口中的每个字段都标注了类型，例如 username 必须是字符串类型
// 这样在代码中写 data.username 时就有自动补全和类型检查，避免运行时错误

// LoginReq：登录请求的参数，包含用户名和密码
export interface LoginReq {
  username: string
  password: string
}

// LoginRes：登录成功后后端返回的数据结构
export interface LoginRes {
  token: string     // JWT 令牌，后续请求通过 Authorization 请求头发送
  user_id: number
  username: string
  nickname: string  // 用户昵称（显示名称）
  roles: string[]   // 角色列表，例如 ['super_admin', 'admin']，用于权限控制
}

// UserInfoRes：获取当前登录用户信息时，后端返回的数据结构
export interface UserInfoRes {
  user_id: number
  username: string
  nickname: string
  roles: string[]
}

// authApi 对象统一管理认证相关的 HTTP 请求
// request.post<T>(url, data) 表示 POST 请求，泛型 T 指定返回数据的类型
export const authApi = {
  // 登录：POST /api/v1/login，请求体为 LoginReq，返回 LoginRes
  login: (data: LoginReq) => request.post<LoginRes>('/login', data),
  // 获取当前用户信息：GET /api/v1/user/info，返回 UserInfoRes
  getUserInfo: () => request.get<UserInfoRes>('/user/info'),
  // 退出登录：POST /api/v1/logout
  logout: () => request.post('/logout'),
}
