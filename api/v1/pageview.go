// ============================================================
// package v1 — API v1 的页面访问追踪接口数据结构
// ------------------------------------------------------------
// 页面访问追踪功能用于记录用户浏览了哪些页面，以及统计访问量。
// 这类数据通常用于分析用户行为、生成访问报告或显示在仪表盘上。
// ============================================================
package v1

import "github.com/gogf/gf/v2/frame/g"

// ============================================================
// PageViewTrackReq — 记录页面访问请求
// ------------------------------------------------------------
// POST 请求——客户端每次页面跳转时调用此接口上报页面路径。
// PagePath 是页面 URL，如 "/dashboard"、"/user/list"。
// 不需要认证，因为登录页面也需要被记录。
// ============================================================
type PageViewTrackReq struct {
	g.Meta  `path:"/pageview/track" method:"post" tags:"PageView" summary:"记录页面访问"`
	PagePath string `json:"page_path" dc:"页面路径"`
}

// PageViewTrackRes — 记录成功响应（空）
type PageViewTrackRes struct{}

// ============================================================
// UserVisitItem — 单个用户的访问统计项
// ------------------------------------------------------------
// 这个结构体在 pageview.go 和 dashboard.go 中都被使用。
// 因为它在 package v1 中定义，所以同一包下的所有文件都可以引用它。
//
//   Username string — 用户名（用于显示谁访问了系统）
//   Count    int    — 该用户的访问次数
// ============================================================
type UserVisitItem struct {
	Username string `json:"username" dc:"用户名"`
	Count    int    `json:"count"    dc:"访问次数"`
}

// PageViewStatsReq — 获取页面访问统计请求（无参数）
type PageViewStatsReq struct {
	g.Meta `path:"/pageview/stats" method:"get" tags:"PageView" summary:"获取页面访问统计"`
}

// PageViewStatsRes — 页面访问统计响应
type PageViewStatsRes struct {
	TotalVisits int             `json:"total_visits" dc:"总访问量"`
	UserVisits  []UserVisitItem `json:"user_visits"  dc:"用户访问统计"`
}
