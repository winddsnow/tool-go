import request from '@/utils/request'
import { MenuTree } from '@/api/menu'

export interface LoginReq {
  username: string
  password: string
}

export interface LoginRes {
  token: string
  user_id: number
  username: string
  nickname: string
  roles: string[]
  menus: MenuTree[]
}

export interface UserInfoRes {
  user_id: number
  username: string
  nickname: string
  roles: string[]
  menus: MenuTree[]
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
