import request from '@/utils/request'

export interface RoleCreateReq {
  name: string
  code: string
  sort?: number
  status?: number
  desc?: string
}

export interface RoleUpdateReq {
  id: number
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
  create: (data: RoleCreateReq) => request.post('/v1/role', data),
  delete: (id: number) => request.delete(`/v1/role/${id}`),
  update: (id: number, data: RoleUpdateReq) => request.put(`/v1/role/${id}`, data),
  getOne: (id: number) => request.get(`/v1/role/${id}`),
  list: (params: { page: number; page_size: number; name?: string; status?: number }) =>
    request.get('/v1/role', { params }),
}
