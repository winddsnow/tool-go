package do

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
)

type RolePermission struct {
    g.Meta       `orm:"table:role_permission, do:true"`
    Id           any
    RoleId       any
    PermissionId any
    CreatedAt    *gtime.Time
}
