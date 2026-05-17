package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UserRole struct {
	g.Meta `orm:"table:user_role, do:true"`
	Id     any
	UserId any
	RoleId any
}
