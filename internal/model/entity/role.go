// Package entity 是数据实体层，定义与数据库表一一对应的结构体。
// 每个实体实例代表数据库表中的一行记录，用于读取数据后的结果映射。
package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role 对应数据库 role 表的一行记录。
// orm 标签将结构体字段映射到数据库列名，json 标签控制 JSON 序列化字段名。
type Role struct {
	Id        uint64      `json:"id"        orm:"id"        description:"Role ID"`
	Name      string      `json:"name"      orm:"name"      description:"Role name"`
	Code      string      `json:"code"      orm:"code"      description:"Role code"`
	Sort      int         `json:"sort"      orm:"sort"      description:"Sort order"`
	Status    uint        `json:"status"    orm:"status"    description:"Status: 1=active, 0=disabled"`
	Desc      string      `json:"desc"      orm:"desc"      description:"Description"`
	// 时间字段使用 *gtime.Time 指针类型，可区分数据库 NULL 值和零值。
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" description:"Created at"`
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" description:"Updated at"`
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" description:"Deleted at"`
}
