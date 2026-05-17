import request from '@/utils/request'

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
}

export interface UserInfoRes {
  user_id: number
  username: string
  nickname: string
  roles: string[]
}

export const authApi = {
  login: (data: LoginReq) => request.post('/login', data),
  getUserInfo: () => request.get('/user/info'),
  logout: () => request.post('/logout'),
}
