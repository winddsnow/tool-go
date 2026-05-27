package entity

import "github.com/gogf/gf/v2/os/gtime"

type RoleMenu struct {
    Id        uint64      `orm:"id"`
    RoleId    uint64      `orm:"role_id"`
    MenuId    uint64      `orm:"menu_id"`
    CreatedAt *gtime.Time `orm:"created_at"`
}
