// Package entity 是数据实体层，定义与数据库表一一对应的结构体。
package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PageView 对应数据库 page_view 表的一行记录，用于存储页面访问日志。
type PageView struct {
	Id        uint64      `json:"id"        orm:"id"        description:"ID"`
	PagePath  string      `json:"page_path" orm:"page_path" description:"Page path"`
	UserId    uint64      `json:"user_id"   orm:"user_id"   description:"User ID"`
	Username  string      `json:"username"  orm:"username"  description:"Username"`
	IpAddress string      `json:"ip_address" orm:"ip_address" description:"IP address"`
	UserAgent string      `json:"user_agent" orm:"user_agent" description:"User agent"`
	// *gtime.Time 指针类型用于存储数据库时间列，nil 表示数据库 NULL 值。
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" description:"Created at"`
}
