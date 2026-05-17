// Package entity 是数据实体层（Entity），定义与数据库表一一对应的 Go 结构体。
// 每个实体实例代表数据库表中的一行记录，主要用于从数据库读取数据后的结果映射（Result Mapping）。
// 实体层的字段类型与实际数据库列类型严格对应，与 do 包（Data Object）不同：
//   - entity：类型严格（uint64、string 等），包含全部字段，用于读取数据
//   - do：类型为 any，仅包含需要写入的字段，用于写入数据
package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User 对应数据库 user 表的一行记录。
// 结构体中的每个字段对应表的一列，字段名通过 orm 标签映射到数据库列名。
// 字段类型与数据库列类型一一对应：
//   - 数值类型 → uint64 / uint / int
//   - 字符串类型 → string
//   - 时间类型 → *gtime.Time（指针，可区分 NULL 和零值）
type User struct {
	Id        uint64      `json:"id"        orm:"id"        description:"User ID"`
	Username  string      `json:"username"  orm:"username"  description:"Username"`
	Password  string      `json:"password"  orm:"password"  description:"Password hash"`
	Salt      string      `json:"salt"      orm:"salt"      description:"Password salt"`
	Nickname  string      `json:"nickname"  orm:"nickname"  description:"Nickname"`
	Email     string      `json:"email"     orm:"email"     description:"Email"`
	Phone     string      `json:"phone"     orm:"phone"     description:"Phone"`
	Status    uint        `json:"status"    orm:"status"    description:"Status: 1=active, 0=disabled"`
	// 时间字段使用 *gtime.Time（指针类型）的原因：
	//   数据库中的 DATETIME 字段可能为 NULL，Go 中使用指针可在数据库值为 NULL 时表示为 nil；
	//   而非指针的 time.Time 在数据库字段为 NULL 时会成为零值（0001-01-01），无法区分 NULL 和零值。
	// gtime.Time 是 GoFrame 框架对 time.Time 的增强封装，提供更灵活的格式化、时区转换等能力。
	//
	// orm 标签（如 orm:"created_at"）是 GoFrame ORM 框架的字段映射标记：
	//   格式为 orm:"column_name"，将结构体字段映射到数据库表的指定列。
	// json 标签（如 json:"created_at"）控制 JSON 序列化/反序列化时的字段名。
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" description:"Created at"`
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" description:"Updated at"`
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" description:"Deleted at"`
}
