package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
)

var Dashboard = cDashboard{}

type cDashboard struct{}

func (c *cDashboard) GetStats(ctx context.Context, req *v1.DashboardStatsReq) (*v1.DashboardStatsRes, error) {
	userCount, _ := dao.User.Ctx(ctx).WhereNull(dao.User.Columns.DeletedAt).Count()
	roleCount, _ := dao.Role.Ctx(ctx).WhereNull(dao.Role.Columns.DeletedAt).Count()

	return &v1.DashboardStatsRes{
		UserCount:  userCount,
		RoleCount:  roleCount,
		OnlineUser: 0,
		ApiRequest: 0,
	}, nil
}
