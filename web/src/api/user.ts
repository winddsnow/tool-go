// user.ts — 用户管理相关的 API 接口
import request from '@/utils/request'

// UserCreateReq：创建用户时前端需要发送的请求体
// ? 表示可选字段，可以不传
export interface UserCreateReq {
  username: string      // 用户名（必填）
  password: string      // 密码（必填）
  nickname?: string     // 昵称（可选）
  email?: string        // 邮箱（可选）
  phone?: string        // 手机号（可选）
  status?: number       // 状态：1=启用，0=禁用（可选，不传则使用后端默认值）
}

// UserUpdateReq：更新用户时前端需要发送的请求体
// 所有字段都是可选的，更新时只传需要修改的字段即可
export interface UserUpdateReq {
  username?: string
  nickname?: string
  email?: string
  phone?: string
  status?: number
}

// UserItem：后端返回的用户对象结构
export interface UserItem {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  status: number       // 1=启用，0=禁用
  created_at: string   // 创建时间，ISO 8601 格式字符串，例如 "2025-01-15T10:30:00Z"
}

// UserListRes：用户列表查询时后端返回的分页数据
export interface UserListRes {
  list: UserItem[]     // 当前页的用户数据数组
  total: number        // 符合条件的用户总数（用于分页组件计算总页数）
  page: number         // 当前页码
}

// userApi 对象统一管理用户相关的 CRUD（增删改查）操作
// 路径中的 :id 会被替换为实际的用户 ID，例如 DELETE /api/v1/user/3
export const userApi = {
  // 创建用户：POST /api/v1/user
  create: (data: UserCreateReq) => request.post('/user', data),
  // 删除用户：DELETE /api/v1/user/:id
  delete: (id: number) => request.delete(`/user/${id}`),
  // 更新用户：PUT /api/v1/user/:id
  update: (id: number, data: UserUpdateReq) => request.put(`/user/${id}`, data),
  // 获取单个用户：GET /api/v1/user/:id
  getOne: (id: number) => request.get(`/user/${id}`),
  // 查询用户列表：GET /api/v1/user?page=1&page_size=10&username=xxx&status=1
  // params 中的 username 和 status 是可选筛选条件
  list: (params: { page: number; page_size: number; username?: string; status?: number }) =>
    request.get('/user', { params }),
  // 获取用户的角色 ID 列表：GET /api/v1/user/:id/roles
  // 返回 { role_ids: [1, 2, 3] } 格式的数据
  getRoles: (id: number) => request.get<{ role_ids: number[] }>(`/user/${id}/roles`),
  // 为用户分配角色：PUT /api/v1/user/:id/roles，请求体为 { role_ids: [1, 2, 3] }
  assignRoles: (id: number, role_ids: number[]) => request.put(`/user/${id}/roles`, { role_ids }),
}
