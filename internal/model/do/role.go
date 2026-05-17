// Package do 是数据对象层（Data Object），用于数据库的写入操作（INSERT / UPDATE）。
// 字段使用 any 类型，零值字段不参与 SQL（实现部分更新）。
// 与 entity 层的区别：
//   - entity：类型严格，包含全部字段，用于读取数据
//   - do：类型为 any，用于写入数据
package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role 是 role 表的数据对象，用于 INSERT / UPDATE 操作。
// g.Meta 内嵌结构体通过 orm:"table:role, do:true" 标签配置 ORM 行为：
//   table:role — 指定数据库表名
//   do:true    — 标记为 DO 类型，ORM 忽略零值字段
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
