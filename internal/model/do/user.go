package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type User struct {
	g.Meta    `orm:"table:user, do:true"`
	Id        any
	Username  any
	Password  any
	Nickname  any
	Email     any
	Phone     any
	Status    any
	CreatedAt *gtime.Time
	UpdatedAt *gtime.Time
	DeletedAt *gtime.Time
}
