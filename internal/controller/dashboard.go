// Package controller 包含所有 HTTP 请求处理器。
// 本文件是仪表盘（Dashboard）相关的 Controller，提供首页统计数据。
package controller

import (
	"context"
	"strings"

	v1 "tool-go/api/v1"
	"tool-go/internal/dao"
)

// var Dashboard = cDashboard{} 是 package-level 单例，路由注册时通过 controller.Dashboard.GetStats 引用。
var Dashboard = cDashboard{}

// cDashboard 是仪表盘 controller，提供聚合统计接口。
type cDashboard struct{}

// GetStats 获取仪表盘统计数据，包括用户数、角色数、总访问量和用户访问排行。
//
// 函数签名按 GoFrame 控制器规范：
//   - (c *cDashboard) 是 Go 的方法接收器（method receiver），类似其他语言的 this/self。
//     使用指针 *cDashboard 而非值 cDashboard，确保方法调用时不会复制接收器（虽然空结构体无区别）。
//   - ctx context.Context：请求上下文，贯穿整个请求生命周期。
//   - req *v1.DashboardStatsReq：请求参数指针。*T 表示"指向 T 的指针"。
//   - *v1.DashboardStatsRes：响应结构体指针。
//   - error：Go 内置错误接口。nil 表示成功。
//
// 注意：这里所有 Count 调用的错误都被忽略（使用 _ 丢弃），即使查询失败也返回 0。
// 这种做法适用于仪表盘这种"尽力而为"的场景，用户看到的只是一个数字而非错误页面。
func (c *cDashboard) GetStats(ctx context.Context, req *v1.DashboardStatsReq) (*v1.DashboardStatsRes, error) {
	// Count() 是 GoFrame ORM 的聚合函数，生成 SELECT COUNT(*) 并返回 int64。
	// WhereNull 确保只统计未被软删除的记录。
	// Go 的多重赋值：第一个返回值赋给 userCount，第二个错误值赋给 _（下划线表示忽略）。
	userCount, _ := dao.User.Ctx(ctx).WhereNull(dao.User.Columns.DeletedAt).Count()
	roleCount, _ := dao.Role.Ctx(ctx).WhereNull(dao.Role.Columns.DeletedAt).Count()

	totalVisits, _ := dao.PageView.Ctx(ctx).Count()

	// type visitCount struct 是 Go 的局部类型声明。
	// 在函数内部定义类型，外部不可见。这是 Go 的常见模式：
	// 当某个结构体只在当前函数内使用时，就近定义以提高可读性。
	// struct 字段后的反引号 `json:"username"` 是 struct tag（结构体标签），
	// 告诉 Go 的 JSON 序列化器在编码时使用这个字段名。
	//
	// Fields("username, COUNT(*) as count") 指定 SELECT 子句，支持原生 SQL 表达式。
	// Where("username != ?", "") 使用占位符 ? 防止 SQL 注入，Go 会自动转义。
	// Group("username") 生成 GROUP BY username。
	// OrderDesc("count") 生成 ORDER BY count DESC（降序排列）。
	// Limit(10) 只返回前 10 条。
	// Scan(&userCounts) 将查询结果映射到 visitCount 切片（slice）中。
	// Scan 根据字段名的 json tag 自动匹配（不匹配的列抛出错误或静默忽略）。
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

	// make([]v1.UserVisitItem, 0, len(userCounts)) 预分配切片：
	//   - 类型：[]v1.UserVisitItem（结构体切片）
	//   - 初始长度：0（空切片，无元素）
	//   - 容量（cap）：len(userCounts)（底层数组大小，避免 append 时多次扩容）
	// for _, v := range userCounts 遍历切片：
	//   - _ 表示忽略索引值（Go 不允许声明了但不使用的变量）
	//   - v 是每次遍历的元素副本（值拷贝）
	// append 是 Go 内置函数，向切片追加元素，容量不足时自动扩容。
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
