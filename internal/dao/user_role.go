package dao

import (
	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

var UserRole = userRoleDao{}

type userRoleDao struct {
	table   string
	group   string
	columns userRoleColumns
}

type userRoleColumns struct {
	Id     string
	UserId string
	RoleId string
}

func init() {
	UserRole = userRoleDao{
		table:   "user_role",
		group:   "default",
		columns: userRoleColumns{
			Id:     "id",
			UserId: "user_id",
			RoleId: "role_id",
		},
	}
}

func (dao *userRoleDao) DB(ctx context.Context) *gdb.Model {
	return gdb.Ctx(ctx).Safe().Model(dao.table).Safe()
}

func (dao *userRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB(ctx)
}

func (dao *userRoleDao) Data(data *do.UserRole) *gdb.Model {
	return dao.DB(gctx.New()).Data(data)
}
