package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `path:"/login" method:"post" tags:"Auth" summary:"用户登录"`
	Username string `json:"username" v:"required#用户名不能为空" dc:"用户名"`
	Password string `json:"password" v:"required#密码不能为空" dc:"密码"`
}

type LoginRes struct {
	Token    string   `json:"token" dc:"访问令牌"`
	UserId   uint64   `json:"user_id" dc:"用户ID"`
	Username string   `json:"username" dc:"用户名"`
	Nickname string   `json:"nickname" dc:"昵称"`
	Roles    []string `json:"roles" dc:"角色列表"`
}

type GetUserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" tags:"Auth" summary:"获取当前用户信息"`
}

type GetUserInfoRes struct {
	UserId   uint64   `json:"user_id" dc:"用户ID"`
	Username string   `json:"username" dc:"用户名"`
	Nickname string   `json:"nickname" dc:"昵称"`
	Roles    []string `json:"roles" dc:"角色列表"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" method:"post" tags:"Auth" summary:"退出登录"`
}

type LogoutRes struct{}
