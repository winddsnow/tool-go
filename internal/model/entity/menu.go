package entity

import "github.com/gogf/gf/v2/os/gtime"

type Menu struct {
    Id        uint64      `orm:"id"`
    ParentId  uint64      `orm:"parent_id"`
    Name      string      `orm:"name"`
    Path      string      `orm:"path"`
    Component string      `orm:"component"`
    Icon      string      `orm:"icon"`
    Sort      int         `orm:"sort"`
    Visible   uint        `orm:"visible"`
    Status    uint        `orm:"status"`
    Type      uint        `orm:"type"`
    CreatedAt *gtime.Time `orm:"created_at"`
    UpdatedAt *gtime.Time `orm:"updated_at"`
    DeletedAt *gtime.Time `orm:"deleted_at"`
}
