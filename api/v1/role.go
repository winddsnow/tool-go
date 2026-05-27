// ============================================================
// package v1 — API v1 的角色管理接口数据结构
// ------------------------------------------------------------
// 本文件定义了"角色管理"的 CRUD 请求/响应结构体。
// 角色（Role）通常代表一组权限的集合，如"超级管理员"、"普通用户"等。
// 在 RBAC（基于角色的访问控制）模型中，权限分配给角色，
// 角色再分配给用户，这样管理权限更加灵活。
// ============================================================
package v1

import "github.com/gogf/gf/v2/frame/g"

// ============================================================
// RoleCreateReq — 创建角色请求参数
// ------------------------------------------------------------
//   Code — 角色编码，如 "super_admin"、"admin"、"user"。
//          编码用于在代码中判断角色（比用 ID 更可读、更稳定）。
//   Sort — 排序值，值越小排序越靠前，用于角色列表的顺序控制。
//   d:"0" — 默认值 0，如果请求不传 sort 字段，框架自动设为 0。
// ============================================================
type RoleCreateReq struct {
	g.Meta `path:"/role" method:"post" tags:"Role" summary:"创建角色"`
	Name   string `json:"name" v:"required#角色名称不能为空" dc:"角色名称"`
	Code   string `json:"code" v:"required#角色编码不能为空" dc:"角色编码"`
	Sort   int    `json:"sort" d:"0" dc:"排序"`
	Status uint   `json:"status" d:"1" dc:"状态: 1=启用, 0=禁用"`
	Desc   string `json:"desc" dc:"描述"`
}

// RoleCreateRes — 创建角色成功响应，返回新角色的 ID
type RoleCreateRes struct {
	Id uint64 `json:"id" dc:"角色ID"`
}

// RoleDeleteReq — 删除角色，路径参数 {id} 指定角色 ID
type RoleDeleteReq struct {
	g.Meta `path:"/role/{id}" method:"delete" tags:"Role" summary:"删除角色"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

// RoleDeleteRes — 删除成功响应（空）
type RoleDeleteRes struct{}

// RoleUpdateReq — 更新角色信息
// 使用 PUT 方法全量更新，未传的字段不会被修改（由业务逻辑决定）
type RoleUpdateReq struct {
	g.Meta `path:"/role/{id}" method:"put" tags:"Role" summary:"更新角色"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
	Name   string `json:"name" dc:"角色名称"`
	Code   string `json:"code" dc:"角色编码"`
	Sort   int    `json:"sort" dc:"排序"`
	Status uint   `json:"status" dc:"状态: 1=启用, 0=禁用"`
	Desc   string `json:"desc" dc:"描述"`
}

// RoleUpdateRes — 更新成功响应（空）
type RoleUpdateRes struct{}

// RoleGetOneReq — 获取单个角色详情
type RoleGetOneReq struct {
	g.Meta `path:"/role/{id}" method:"get" tags:"Role" summary:"获取角色详情"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

// RoleGetOneRes — 角色详情响应
type RoleGetOneRes struct {
	Id        uint64 `json:"id" dc:"角色ID"`
	Name      string `json:"name" dc:"角色名称"`
	Code      string `json:"code" dc:"角色编码"`
	Sort      int    `json:"sort" dc:"排序"`
	Status    uint   `json:"status" dc:"状态"`
	Desc      string `json:"desc" dc:"描述"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
	UpdatedAt string `json:"updated_at" dc:"更新时间"`
}

// ============================================================
// RoleListReq — 角色列表查询（分页 + 筛选）
// ------------------------------------------------------------
// Status *int — 同样使用指针类型解决"0 值"歧义问题，
// 具体解释见 user.go 中的 UserListReq.Status 注释。
// ============================================================
type RoleListReq struct {
	g.Meta   `path:"/role" method:"get" tags:"Role" summary:"获取角色列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Name     string `json:"name" dc:"按名称筛选"`
	Status   *int   `json:"status" dc:"按状态筛选"`
}

// RoleListRes — 角色列表分页响应
type RoleListRes struct {
	List  []RoleItem `json:"list" dc:"角色列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
}

// RoleItem — 列表中的单个角色项
type RoleItem struct {
	Id        uint64 `json:"id" dc:"角色ID"`
	Name      string `json:"name" dc:"角色名称"`
	Code      string `json:"code" dc:"角色编码"`
	Sort      int    `json:"sort" dc:"排序"`
	Status    uint   `json:"status" dc:"状态"`
	Desc      string `json:"desc" dc:"描述"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
}

type RoleGetPermissionsReq struct {
	g.Meta `path:"/role/{id}/permissions" method:"get" tags:"Role" summary:"获取角色权限"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

type RoleGetPermissionsRes struct {
	PermissionIds []uint64 `json:"permission_ids" dc:"权限ID列表"`
}

type RoleAssignPermissionsReq struct {
	g.Meta        `path:"/role/{id}/permissions" method:"put" tags:"Role" summary:"分配角色权限"`
	Id            uint64   `path:"id" v:"required#ID不能为空" dc:"角色ID"`
	PermissionIds []uint64 `json:"permission_ids" v:"required#权限列表不能为空" dc:"权限ID列表"`
}

type RoleAssignPermissionsRes struct{}
