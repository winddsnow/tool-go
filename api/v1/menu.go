package v1

import "github.com/gogf/gf/v2/frame/g"

type MenuCreateReq struct {
	g.Meta    `path:"/menu" method:"post" tags:"Menu" summary:"创建菜单"`
	ParentId  uint64 `json:"parent_id" d:"0" dc:"父菜单ID"`
	Name      string `json:"name" v:"required#菜单名称不能为空" dc:"菜单名称"`
	Path      string `json:"path" dc:"路由路径"`
	Component string `json:"component" dc:"前端组件路径"`
	Icon      string `json:"icon" dc:"图标名"`
	Sort      int    `json:"sort" d:"0" dc:"排序"`
	Visible   uint   `json:"visible" d:"1" dc:"是否显示"`
	Status    uint   `json:"status" d:"1" dc:"状态"`
	Type      uint   `json:"type" d:"1" dc:"类型: 1=目录, 2=菜单, 3=按钮"`
}

type MenuCreateRes struct {
	Id uint64 `json:"id" dc:"菜单ID"`
}

type MenuDeleteReq struct {
	g.Meta `path:"/menu/{id}" method:"delete" tags:"Menu" summary:"删除菜单"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"菜单ID"`
}

type MenuDeleteRes struct{}

type MenuUpdateReq struct {
	g.Meta    `path:"/menu/{id}" method:"put" tags:"Menu" summary:"更新菜单"`
	Id        uint64  `path:"id" v:"required#ID不能为空" dc:"菜单ID"`
	ParentId  *uint64 `json:"parent_id" dc:"父菜单ID"`
	Name      string  `json:"name" dc:"菜单名称"`
	Path      string  `json:"path" dc:"路由路径"`
	Component string  `json:"component" dc:"前端组件路径"`
	Icon      string  `json:"icon" dc:"图标名"`
	Sort      *int    `json:"sort" dc:"排序"`
	Visible   *uint   `json:"visible" dc:"是否显示"`
	Status    *uint   `json:"status" dc:"状态"`
	Type      *uint   `json:"type" dc:"类型"`
}

type MenuUpdateRes struct{}

type MenuGetOneReq struct {
	g.Meta `path:"/menu/{id}" method:"get" tags:"Menu" summary:"获取菜单详情"`
	Id     uint64 `path:"id" v:"required#ID不能为空" dc:"菜单ID"`
}

type MenuGetOneRes struct {
	Id        uint64 `json:"id" dc:"菜单ID"`
	ParentId  uint64 `json:"parent_id" dc:"父菜单ID"`
	Name      string `json:"name" dc:"菜单名称"`
	Path      string `json:"path" dc:"路由路径"`
	Component string `json:"component" dc:"前端组件路径"`
	Icon      string `json:"icon" dc:"图标名"`
	Sort      int    `json:"sort" dc:"排序"`
	Visible   uint   `json:"visible" dc:"是否显示"`
	Status    uint   `json:"status" dc:"状态"`
	Type      uint   `json:"type" dc:"类型"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
	UpdatedAt string `json:"updated_at" dc:"更新时间"`
}

type MenuListReq struct {
	g.Meta   `path:"/menu" method:"get" tags:"Menu" summary:"获取菜单列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Name     string `json:"name" dc:"按名称筛选"`
	Status   *int   `json:"status" dc:"按状态筛选"`
}

type MenuListRes struct {
	List  []MenuItem `json:"list" dc:"菜单列表"`
	Total int        `json:"total" dc:"总数"`
	Page  int        `json:"page" dc:"当前页"`
}

type MenuItem struct {
	Id        uint64 `json:"id" dc:"菜单ID"`
	ParentId  uint64 `json:"parent_id" dc:"父菜单ID"`
	Name      string `json:"name" dc:"菜单名称"`
	Path      string `json:"path" dc:"路由路径"`
	Component string `json:"component" dc:"前端组件路径"`
	Icon      string `json:"icon" dc:"图标名"`
	Sort      int    `json:"sort" dc:"排序"`
	Visible   uint   `json:"visible" dc:"是否显示"`
	Status    uint   `json:"status" dc:"状态"`
	Type      uint   `json:"type" dc:"类型"`
	CreatedAt string `json:"created_at" dc:"创建时间"`
}

type MenuGetUserMenusReq struct {
	g.Meta `path:"/menu/user" method:"get" tags:"Menu" summary:"获取当前用户菜单"`
}

type MenuTree struct {
	Id        uint64     `json:"id" dc:"菜单ID"`
	ParentId  uint64     `json:"parent_id" dc:"父菜单ID"`
	Name      string     `json:"name" dc:"菜单名称"`
	Path      string     `json:"path" dc:"路由路径"`
	Component string     `json:"component" dc:"前端组件路径"`
	Icon      string     `json:"icon" dc:"图标名"`
	Sort      int        `json:"sort" dc:"排序"`
	Visible   uint       `json:"visible" dc:"是否显示"`
	Type      uint       `json:"type" dc:"类型"`
	Children  []MenuTree `json:"children" dc:"子菜单"`
}

type MenuGetUserMenusRes struct {
	Menus []MenuTree `json:"menus" dc:"用户菜单树"`
}
