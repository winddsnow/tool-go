package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type PageView struct {
	Id        uint64      `json:"id"        orm:"id"        description:"ID"`
	PagePath  string      `json:"page_path" orm:"page_path" description:"Page path"`
	UserId    uint64      `json:"user_id"   orm:"user_id"   description:"User ID"`
	Username  string      `json:"username"  orm:"username"  description:"Username"`
	IpAddress string      `json:"ip_address" orm:"ip_address" description:"IP address"`
	UserAgent string      `json:"user_agent" orm:"user_agent" description:"User agent"`
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" description:"Created at"`
}
