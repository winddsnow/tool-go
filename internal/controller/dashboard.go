package controller

import (
	"context"
	"strings"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
)

var Dashboard = cDashboard{}

type cDashboard struct{}

func (c *cDashboard) GetStats(ctx context.Context, req *v1.DashboardStatsReq) (*v1.DashboardStatsRes, error) {
	userCount, _ := dao.User.Ctx(ctx).WhereNull(dao.User.Columns.DeletedAt).Count()
	roleCount, _ := dao.Role.Ctx(ctx).WhereNull(dao.Role.Columns.DeletedAt).Count()

	totalVisits, _ := dao.PageView.Ctx(ctx).Count()

	type visitCount struct {
		Username string `json:"username"`
		Count    int    `json:"count"`
	}
	var userCounts []visitCount
	err := dao.PageView.Ctx(ctx).
		Fields("username, COUNT(*) as count").
		Where("username != ?", "").
		Group("username").
		OrderDesc("count").
		Limit(10).
		Scan(&userCounts)
	if err != nil {
		userCounts = nil
	}

	items := make([]v1.UserVisitItem, 0, len(userCounts))
	for _, v := range userCounts {
		if strings.TrimSpace(v.Username) == "" {
			continue
		}
		items = append(items, v1.UserVisitItem{
			Username: v.Username,
			Count:    v.Count,
		})
	}

	return &v1.DashboardStatsRes{
		UserCount:   userCount,
		RoleCount:   roleCount,
		TotalVisits: totalVisits,
		UserVisits:  items,
	}, nil
}
