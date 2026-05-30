package v1

import "github.com/gogf/gf/v2/frame/g"

type PermissionListReq struct {
	g.Meta   `path:"/permission" method:"get" tags:"Permission" summary:"获取权限列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Code     string `json:"code" dc:"按权限码筛选"`
}

type PermissionListRes struct {
	List  []PermissionItem `json:"list" dc:"权限列表"`
	Total int              `json:"total" dc:"总数"`
	Page  int              `json:"page" dc:"当前页"`
}

type PermissionItem struct {
	Id     uint64 `json:"id" dc:"权限ID"`
	Code   string `json:"code" dc:"权限码"`
	Name   string `json:"name" dc:"权限名称"`
	MenuId uint64 `json:"menu_id" dc:"关联菜单ID"`
}
