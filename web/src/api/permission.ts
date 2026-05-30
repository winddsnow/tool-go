import request from '@/utils/request'

export interface PermissionItem {
  id: number
  code: string
  name: string
  menu_id: number
}

export interface PermissionListRes {
  list: PermissionItem[]
  total: number
  page: number
}

export const permissionApi = {
  list: (params: { page: number; page_size: number; code?: string }) =>
    request.get<PermissionListRes>('/permission', { params }),
}
