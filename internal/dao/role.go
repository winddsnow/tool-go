package dao

import (
	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

var Role = roleDao{}

type roleDao struct {
	table    string
	group    string
	columns  roleColumns
}

type roleColumns struct {
	Id        string
	Name      string
	Code      string
	Sort      string
	Status    string
	Desc      string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

func init() {
	Role = roleDao{
		table:   "role",
		group:   "default",
		columns: roleColumns{
			Id:        "id",
			Name:      "name",
			Code:      "code",
			Sort:      "sort",
			Status:    "status",
			Desc:      "desc",
			CreatedAt: "created_at",
			UpdatedAt: "updated_at",
			DeletedAt: "deleted_at",
		},
	}
}

func (dao *roleDao) DB(ctx context.Context) *gdb.Model {
	return gdb.Ctx(ctx).Safe().Model(dao.table).Safe()
}

func (dao *roleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return gdb.Ctx(ctx).Transaction(ctx, f)
}

func (dao *roleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB(ctx)
}

func (dao *roleDao) Data(data *do.Role) *gdb.Model {
	return dao.DB(gctx.New()).Data(data)
}

func (dao *roleDao) As(as string) *gdb.Model {
	return dao.DB(gctx.New()).As(as)
}
