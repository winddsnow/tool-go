package v1

import "github.com/gogf/gf/v2/frame/g"

type PageViewTrackReq struct {
	g.Meta  `path:"/pageview/track" method:"post" tags:"PageView" summary:"记录页面访问"`
	PagePath string `json:"page_path" dc:"页面路径"`
}

type PageViewTrackRes struct{}

type UserVisitItem struct {
	Username string `json:"username" dc:"用户名"`
	Count    int    `json:"count"    dc:"访问次数"`
}

type PageViewStatsReq struct {
	g.Meta `path:"/pageview/stats" method:"get" tags:"PageView" summary:"获取页面访问统计"`
}

type PageViewStatsRes struct {
	TotalVisits int             `json:"total_visits" dc:"总访问量"`
	UserVisits  []UserVisitItem `json:"user_visits"  dc:"用户访问统计"`
}
