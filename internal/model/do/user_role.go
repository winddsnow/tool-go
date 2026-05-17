// Package do 是数据对象层（Data Object），用于数据库的写入操作。
package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserRole 是 user_role 关联表的数据对象，用于 INSERT / UPDATE。
// g.Meta 通过 orm 标签指定表名和 DO 行为，any 类型字段支持部分更新。
type UserRole struct {
	g.Meta `orm:"table:user_role, do:true"`
	Id     any
	UserId any
	RoleId any
}
