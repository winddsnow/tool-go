// Package do 是数据对象层（Data Object），用于数据库的写入操作（INSERT / UPDATE）。
// 与 entity 层的核心区别：
//   - entity：类型严格（uint64、string 等），包含数据库表的所有字段，主要用于读取数据后的结果映射
//   - do：字段类型使用 any（interface{}），可以包含任意类型的值
//
// 为什么 do 使用 any 类型？
//   当需要执行部分更新（partial update）时，do 中未设置的字段（零值）不会被 ORM 拼接到 SQL 中。
//   例如只需要更新 Password 字段时，do 只设置 Password 字段，其他字段保持零值，ORM 不会生成
//   其他列的 SET 子句。这避免了 entity 层零值无法区分"未设置"和"值就是零"的问题。
package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User 是 user 表的数据对象，用于 INSERT / UPDATE 操作。
// g.Meta 是 GoFrame 框架的元数据嵌入结构体，提供 ORM 相关的结构体级别配置。
// orm:"table:user, do:true" 标签含义：
//   table:user  — 指定该结构体映射到数据库中的 "user" 表
//   do:true     — 标记此结构体为 DO 类型。ORM 在处理 DO 类型时，
//                  会忽略零值字段（如空字符串 ""、数字 0、nil），
//                  只将非零值字段生成到 INSERT 或 UPDATE 的 SQL 中，实现精确的部分更新。
type User struct {
	g.Meta    `orm:"table:user, do:true"`
	Id        any
	Username  any
	Password  any
	Salt      any
	Nickname  any
	Email     any
	Phone     any
	Status    any
	CreatedAt *gtime.Time
	UpdatedAt *gtime.Time
	DeletedAt *gtime.Time
}
