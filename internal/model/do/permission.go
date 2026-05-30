package do

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
)

type Permission struct {
    g.Meta    `orm:"table:permission, do:true"`
    Id        any
    Code      any
    Name      any
    MenuId    any
    CreatedAt *gtime.Time
}
