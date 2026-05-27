package entity

import "github.com/gogf/gf/v2/os/gtime"

type RolePermission struct {
    Id           uint64      `orm:"id"`
    RoleId       uint64      `orm:"role_id"`
    PermissionId uint64      `orm:"permission_id"`
    CreatedAt    *gtime.Time `orm:"created_at"`
}
