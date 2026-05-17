package dao

import (
	"context"

	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var UserRole = userRoleDao{}

type userRoleDao struct {
	Table   string
	Group   string
	Columns userRoleColumns
}

type userRoleColumns struct {
	Id     string
	UserId string
	RoleId string
}

func init() {
	UserRole = userRoleDao{
		Table: "user_role",
		Group: "default",
		Columns: userRoleColumns{
			Id:     "id",
			UserId: "user_id",
			RoleId: "role_id",
		},
	}
}

func (d *userRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *userRoleDao) Data(data *do.UserRole) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
