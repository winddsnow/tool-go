package v1

import "github.com/gogf/gf/v2/frame/g"

type UserCreateReq struct {
	g.Meta   `path:"/user" method:"post" tags:"User" summary:"创建用户"`
	Username string `json:"username" v:"required#用户名不能为空" dc:"用户名"`
	Password string `json:"password" v:"required#密码不能为空" dc:"密码"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`
	Phone    string `json:"phone" dc:"手机号"`
	Status   uint   `json:"status" d:"1" dc:"状态: 1=启用, 0=禁用"`
}

type UserCreateRes struct {
	Id uint64 `json:"id" dc:"用户ID"`
}

type UserDeleteReq struct {
	g.Meta `path:"/user/{id}" method:"delete" tags:"User" summary:"删除用户"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
}

type UserDeleteRes struct{}

type UserUpdateReq struct {
	g.Meta   `path:"/user/{id}" method:"put" tags:"User" summary:"更新用户"`
	Id       uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
	Username string `json:"username" dc:"用户名"`
	Nickname string `json:"nickname" dc:"昵称"`
	Email    string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`
	Phone    string `json:"phone" dc:"手机号"`
	Status   uint   `json:"status" dc:"状态: 1=启用, 0=禁用"`
}

type UserUpdateRes struct{}

type UserGetOneReq struct {
	g.Meta `path:"/user/{id}" method:"get" tags:"User" summary:"获取用户详情"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
}

type UserGetOneRes struct {
	Id        uint64 `json:"id" dc:"用户ID"`
	Username  string `json:"username" dc:"用户名"`
	Nickname  string `json:"nickname" dc:"昵称"`
	Email     string `json:"email" dc:"邮箱"`
	Phone     string `json:"phone" dc:"手机号"`
	Status    uint   `json:"status" dc:"状态"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
	UpdatedAt string `json:"updated_at" dc:"更新时间"`
}

type UserListReq struct {
	g.Meta   `path:"/user" method:"get" tags:"User" summary:"获取用户列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Username string `json:"username" dc:"按用户名筛选"`
	Status   uint   `json:"status" dc:"按状态筛选"`
	UserId   uint64 `json:"-" dc:"当前用户ID"`
}

type UserListRes struct {
	List  []UserItem `json:"list" dc:"用户列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
}

type UserItem struct {
	Id        uint64 `json:"id" dc:"用户ID"`
	Username  string `json:"username" dc:"用户名"`
	Nickname  string `json:"nickname" dc:"昵称"`
	Email     string `json:"email" dc:"邮箱"`
	Phone     string `json:"phone" dc:"手机号"`
	Status    uint   `json:"status" dc:"状态"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
}

type UserGetRolesReq struct {
	g.Meta `path:"/user/{id}/roles" method:"get" tags:"User" summary:"获取用户角色"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"用户ID"`
}

type UserGetRolesRes struct {
	RoleIds []uint64 `json:"role_ids" dc:"角色ID列表"`
}

type UserAssignRolesReq struct {
	g.Meta  `path:"/user/{id}/roles" method:"put" tags:"User" summary:"分配用户角色"`
	Id      uint64   `path:"id" v:"required#ID不能为空" dc:"用户ID"`
	RoleIds []uint64 `json:"role_ids" v:"required#角色列表不能为空" dc:"角色ID列表"`
}

type UserAssignRolesRes struct{}
