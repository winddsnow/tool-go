package dao

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

var PageView = pageViewDao{}

type pageViewDao struct {
	Table   string
	Group   string
	Columns pageViewColumns
}

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

func (d *pageViewDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}
