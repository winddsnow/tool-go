package do

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gtime"
)

type Menu struct {
    g.Meta    `orm:"table:menu, do:true"`
    Id        any
    ParentId  any
    Name      any
    Path      any
    Component any
    Icon      any
    Sort      any
    Visible   any
    Status    any
    Type      any
    CreatedAt *gtime.Time
    UpdatedAt *gtime.Time
    DeletedAt *gtime.Time
}
