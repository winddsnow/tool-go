package dao

import (
	"context"

	"tool-go/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var Role = roleDao{}

type roleDao struct {
	Table   string
	Group   string
	Columns roleColumns
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
		Table: "role",
		Group: "default",
		Columns: roleColumns{
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

func (d *roleDao) Ctx(ctx context.Context) *gdb.Model {
	return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *roleDao) Data(data *do.Role) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}

func (d *roleDao) As(as string) *gdb.Model {
	return g.Model(d.Table).Ctx(gctx.New()).As(as)
}
