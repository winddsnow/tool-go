package v1

import "github.com/gogf/gf/v2/frame/g"

type DashboardStatsReq struct {
	g.Meta `path:"/dashboard/stats" method:"get" tags:"Dashboard" summary:"获取仪表盘统计数据"`
}

type DashboardStatsRes struct {
	UserCount    int             `json:"user_count"    dc:"用户总数"`
	RoleCount    int             `json:"role_count"    dc:"角色总数"`
	OnlineUser   int             `json:"online_user"   dc:"在线用户"`
	ApiRequest   int             `json:"api_request"   dc:"API请求数"`
	TotalVisits  int             `json:"total_visits"  dc:"总访问量"`
	UserVisits   []UserVisitItem `json:"user_visits"   dc:"用户访问统计"`
}
