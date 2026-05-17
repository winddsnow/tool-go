import request from '@/utils/request'

export interface UserCreateReq {
  username: string
  password: string
  nickname?: string
  email?: string
  phone?: string
  status?: number
}

export interface UserUpdateReq {
  username?: string
  nickname?: string
  email?: string
  phone?: string
  status?: number
}

export interface UserItem {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  status: number
  created_at: string
}

export interface UserListRes {
  list: UserItem[]
  total: number
  page: number
}

export const userApi = {
  create: (data: UserCreateReq) => request.post('/user', data),
  delete: (id: number) => request.delete(`/user/${id}`),
  update: (id: number, data: UserUpdateReq) => request.put(`/user/${id}`, data),
  getOne: (id: number) => request.get(`/user/${id}`),
  list: (params: { page: number; page_size: number; username?: string; status?: number }) =>
    request.get('/user', { params }),
  getRoles: (id: number) => request.get<{ role_ids: number[] }>(`/user/${id}/roles`),
  assignRoles: (id: number, role_ids: number[]) => request.put(`/user/${id}/roles`, { role_ids }),
}
