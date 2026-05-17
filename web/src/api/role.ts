import request from '@/utils/request'

export interface RoleCreateReq {
  name: string
  code: string
  sort?: number
  status?: number
  desc?: string
}

export interface RoleUpdateReq {
  name?: string
  code?: string
  sort?: number
  status?: number
  desc?: string
}

export interface RoleItem {
  id: number
  name: string
  code: string
  sort: number
  status: number
  desc: string
  created_at: string
}

export interface RoleListRes {
  list: RoleItem[]
  total: number
  page: number
}

export const roleApi = {
  create: (data: RoleCreateReq) => request.post('/role', data),
  delete: (id: number) => request.delete(`/role/${id}`),
  update: (id: number, data: RoleUpdateReq) => request.put(`/role/${id}`, data),
  getOne: (id: number) => request.get(`/role/${id}`),
  list: (params: { page: number; page_size: number; name?: string; status?: number }) =>
    request.get('/role', { params }),
}
