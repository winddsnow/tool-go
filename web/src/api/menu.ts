import request from '@/utils/request'

export interface MenuTree {
  id: number
  parent_id: number
  name: string
  path: string
  component: string
  icon: string
  sort: number
  visible: number
  type: number
  children?: MenuTree[]
}

export interface MenuCreateReq {
  parent_id?: number
  name: string
  path?: string
  component?: string
  icon?: string
  sort?: number
  visible?: number
  status?: number
  type?: number
}

export interface MenuItem {
  id: number
  parent_id: number
  name: string
  path: string
  component: string
  icon: string
  sort: number
  visible: number
  status: number
  type: number
  created_at: string
}

export interface MenuListRes {
  list: MenuItem[]
  total: number
  page: number
}

export const menuApi = {
  create: (data: MenuCreateReq) => request.post('/menu', data),
  delete: (id: number) => request.delete(`/menu/${id}`),
  update: (id: number, data: Partial<MenuCreateReq>) => request.put(`/menu/${id}`, data),
  getOne: (id: number) => request.get(`/menu/${id}`),
  list: (params: { page: number; page_size: number; name?: string; status?: number }) =>
    request.get('/menu', { params }),
  getUserMenus: () => request.get<{ menus: MenuTree[] }>('/menu/user'),
}
