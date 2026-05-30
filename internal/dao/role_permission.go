package dao

import (
    "context"
    "tool-go/internal/model/do"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

var RolePermission = rolePermissionDao{}

type rolePermissionDao struct {
    Table   string
    Group   string
    Columns rolePermissionColumns
}

type rolePermissionColumns struct {
    Id           string
    RoleId       string
    PermissionId string
}

func init() {
    RolePermission = rolePermissionDao{
        Table: "role_permission",
        Group: "default",
        Columns: rolePermissionColumns{
            Id:           "id",
            RoleId:       "role_id",
            PermissionId: "permission_id",
        },
    }
}

func (d *rolePermissionDao) Ctx(ctx context.Context) *gdb.Model {
    return g.Model(d.Table).Safe().Ctx(ctx)
}

func (d *rolePermissionDao) Data(data *do.RolePermission) *gdb.Model {
    return g.Model(d.Table).Ctx(gctx.New()).Data(data)
}
