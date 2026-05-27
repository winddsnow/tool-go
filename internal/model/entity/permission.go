package entity

import "github.com/gogf/gf/v2/os/gtime"

type Permission struct {
    Id        uint64      `orm:"id"`
    Code      string      `orm:"code"`
    Name      string      `orm:"name"`
    MenuId    uint64      `orm:"menu_id"`
    CreatedAt *gtime.Time `orm:"created_at"`
}
