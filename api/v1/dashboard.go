// ============================================================
// package v1 — API v1 的仪表盘接口数据结构
// ------------------------------------------------------------
// 仪表盘（Dashboard）通常用于在首页展示系统概览数据，
// 如用户数、角色数、访问量等统计信息。
// ============================================================
package v1

import "github.com/gogf/gf/v2/frame/g"

// DashboardStatsReq — 获取仪表盘统计数据请求（无参数）
type DashboardStatsReq struct {
	g.Meta `path:"/dashboard/stats" method:"get" tags:"Dashboard" summary:"获取仪表盘统计数据"`
}

// ============================================================
// DashboardStatsRes — 仪表盘统计数据响应
// ------------------------------------------------------------
// 字段都是基础 int 类型，因为统计结果一定是数字。
// UserVisits 字段使用了 UserVisitItem 结构体（定义在 pageview.go 中），
// 展示了 Go 中跨文件引用同一包下类型的用法。
// 同一个 package v1 下的所有 .go 文件共享所有类型定义。
// ============================================================
type DashboardStatsRes struct {
	UserCount    int             `json:"user_count"    dc:"用户总数"`
	RoleCount    int             `json:"role_count"    dc:"角色总数"`
	OnlineUser   int             `json:"online_user"   dc:"在线用户"`
	ApiRequest   int             `json:"api_request"   dc:"API请求数"`
	TotalVisits  int             `json:"total_visits"  dc:"总访问量"`
	UserVisits   []UserVisitItem `json:"user_visits"   dc:"用户访问统计"`
}
