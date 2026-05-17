package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type Role struct {
	g.Meta    `orm:"table:role, do:true"`
	Id        any
	Name      any
	Code      any
	Sort      any
	Status    any
	Desc      any
	CreatedAt *gtime.Time
	UpdatedAt *gtime.Time
	DeletedAt *gtime.Time
}
