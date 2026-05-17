package v1

import "github.com/gogf/gf/v2/frame/g"

type RoleCreateReq struct {
	g.Meta `path:"/role" method:"post" tags:"Role" summary:"创建角色"`
	Name   string `json:"name" v:"required#角色名称不能为空" dc:"角色名称"`
	Code   string `json:"code" v:"required#角色编码不能为空" dc:"角色编码"`
	Sort   int    `json:"sort" d:"0" dc:"排序"`
	Status uint   `json:"status" d:"1" dc:"状态: 1=启用, 0=禁用"`
	Desc   string `json:"desc" dc:"描述"`
}

type RoleCreateRes struct {
	Id uint64 `json:"id" dc:"角色ID"`
}

type RoleDeleteReq struct {
	g.Meta `path:"/role/{id}" method:"delete" tags:"Role" summary:"删除角色"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

type RoleDeleteRes struct{}

type RoleUpdateReq struct {
	g.Meta `path:"/role/{id}" method:"put" tags:"Role" summary:"更新角色"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
	Name   string `json:"name" dc:"角色名称"`
	Code   string `json:"code" dc:"角色编码"`
	Sort   int    `json:"sort" dc:"排序"`
	Status uint   `json:"status" dc:"状态: 1=启用, 0=禁用"`
	Desc   string `json:"desc" dc:"描述"`
}

type RoleUpdateRes struct{}

type RoleGetOneReq struct {
	g.Meta `path:"/role/{id}" method:"get" tags:"Role" summary:"获取角色详情"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"角色ID"`
}

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

type RoleListReq struct {
	g.Meta   `path:"/role" method:"get" tags:"Role" summary:"获取角色列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Name     string `json:"name" dc:"按名称筛选"`
	Status   *int   `json:"status" dc:"按状态筛选"`
}

type RoleListRes struct {
	List  []RoleItem `json:"list" dc:"角色列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
}

type RoleItem struct {
	Id        uint64 `json:"id" dc:"角色ID"`
	Name      string `json:"name" dc:"角色名称"`
	Code      string `json:"code" dc:"角色编码"`
	Sort      int    `json:"sort" dc:"排序"`
	Status    uint   `json:"status" dc:"状态"`
	Desc      string `json:"desc" dc:"描述"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
}
