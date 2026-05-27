package dao

import (
    "context"
    "tool-go/internal/model/do"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

var Permission = permissionDao{}

type permissionDao struct {
    Table   string
    Group   string
    Columns permissionColumns
}

type permissionColumns struct {
    Id     string
    Code   string
    Name   string
    MenuId string
}

func init() {
    Permission = permissionDao{
        Table: "permission",
        Group: "default",
        Columns: permissionColumns{
            Id:     "id",
            Code:   "code",
            Name:   "name",
            MenuId: "menu_id",
        },
    }
}

func (d *permissionDao) Ctx(ctx context.Context) *gdb.Model {
    return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *permissionDao) Data(data *do.Permission) *gdb.Model {
    return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
