// Package dao 是数据访问层（Data Access Object），封装对数据库表的基本操作。
package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PageView 是全局导出的 DAO 实例。
var PageView = pageViewDao{}

// pageViewDao 是 page_view（页面访问记录）表的 DAO 结构体。
// 用于记录用户访问页面的日志信息，包括访问路径、用户信息、IP 地址等。
type pageViewDao struct {
	Table   string
	Group   string
	Columns pageViewColumns
}

// pageViewColumns 定义 page_view 表的列名常量。
type pageViewColumns struct {
	Id        string
	PagePath  string
	UserId    string
	Username  string
	IpAddress string
	UserAgent string
	CreatedAt string
}

func init() {
	PageView = pageViewDao{
		Table: "page_view",
		Group: "default",
		Columns: pageViewColumns{
			Id:        "id",
			PagePath:  "page_path",
			UserId:    "user_id",
			Username:  "username",
			IpAddress: "ip_address",
			UserAgent: "user_agent",
			CreatedAt: "created_at",
		},
	}
}

// Ctx 创建带上下文的 ORM 查询模型，适用于 SELECT 查询。
// page_view 表为只追加写入的日志表，通常不需要 Data 和 As 方法。
func (d *pageViewDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}
