// role.ts — 角色管理相关的 API 接口
import request from '@/utils/request'

// RoleCreateReq：创建角色时前端需要发送的请求体
export interface RoleCreateReq {
  name: string         // 角色显示名称，例如"超级管理员"
  code: string         // 角色标识码，例如"super_admin"，后端权限判断使用 code
  sort?: number        // 排序值（可选），数值越小越靠前
  status?: number      // 状态：1=启用，0=禁用（可选）
  desc?: string        // 角色描述（可选）
}

// RoleUpdateReq：更新角色时前端需要发送的请求体，所有字段可选
export interface RoleUpdateReq {
  name?: string
  code?: string
  sort?: number
  status?: number
  desc?: string
}

// RoleItem：后端返回的角色对象结构
export interface RoleItem {
  id: number
  name: string
  code: string
  sort: number
  status: number
  desc: string
  created_at: string
}

// RoleListRes：角色列表查询时后端返回的分页数据
export interface RoleListRes {
  list: RoleItem[]
  total: number
  page: number
}

// roleApi 对象统一管理角色相关的 CRUD 操作
export const roleApi = {
  // 创建角色：POST /api/v1/role
  create: (data: RoleCreateReq) => request.post('/role', data),
  // 删除角色：DELETE /api/v1/role/:id
  delete: (id: number) => request.delete(`/role/${id}`),
  // 更新角色：PUT /api/v1/role/:id
  update: (id: number, data: RoleUpdateReq) => request.put(`/role/${id}`, data),
  // 获取单个角色：GET /api/v1/role/:id
  getOne: (id: number) => request.get(`/role/${id}`),
  // 查询角色列表：GET /api/v1/role?page=1&page_size=10&name=xxx&status=1
  list: (params: { page: number; page_size: number; name?: string; status?: number }) =>
    request.get('/role', { params }),
  // 获取角色的权限ID列表：GET /api/v1/role/:id/permissions
  getPermissions: (id: number) => request.get<{ permission_ids: number[] }>(`/role/${id}/permissions`),
  // 为角色分配权限：PUT /api/v1/role/:id/permissions
  assignPermissions: (id: number, data: { permission_ids: number[] }) =>
    request.put(`/role/${id}/permissions`, data),
}
