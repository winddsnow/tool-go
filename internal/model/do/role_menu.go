package do

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
)

type RoleMenu struct {
    g.Meta    `orm:"table:role_menu, do:true"`
    Id        any
    RoleId    any
    MenuId    any
    CreatedAt *gtime.Time
}
