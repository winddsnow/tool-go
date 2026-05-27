package dao

import (
    "context"
    "tool-go/internal/model/do"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

var RoleMenu = roleMenuDao{}

type roleMenuDao struct {
    Table   string
    Group   string
    Columns roleMenuColumns
}

type roleMenuColumns struct {
    Id     string
    RoleId string
    MenuId string
}

func init() {
    RoleMenu = roleMenuDao{
        Table: "role_menu",
        Group: "default",
        Columns: roleMenuColumns{
            Id:     "id",
            RoleId: "role_id",
            MenuId: "menu_id",
        },
    }
}

func (d *roleMenuDao) Ctx(ctx context.Context) *gdb.Model {
    return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *roleMenuDao) Data(data *do.RoleMenu) *gdb.Model {
    return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
